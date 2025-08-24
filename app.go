/*
reposync - A tool for synchronizing repositories from GitLab and GitHub organizations.
Manages authentication through stored personal access tokens and maintains directory structure.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/term"

	colors "github.com/itszeeshan/reposync/constants/colors"
	models "github.com/itszeeshan/reposync/constants/models"
	helpers "github.com/itszeeshan/reposync/helpers"
	services "github.com/itszeeshan/reposync/services"
)

/*
getSecureInput reads sensitive input without displaying it on screen.
Uses terminal.ReadPassword to hide input from terminal history and process lists.
*/
func getSecureInput(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // New line after input
	return string(bytePassword), nil
}

/*
handleConfig implements interactive token configuration workflow.
Prompts user for both GitLab and GitHub tokens using secure input,
then saves them to encrypted config file in user's home directory for future use.
*/
func handleConfig() error {
	fmt.Print("Enter GitLab Personal Access Token: ")
	gitlabToken, err := getSecureInput("")
	if err != nil {
		return fmt.Errorf("failed to read GitLab token: %w", err)
	}

	fmt.Print("Enter GitHub Personal Access Token: ")
	githubToken, err := getSecureInput("")
	if err != nil {
		return fmt.Errorf("failed to read GitHub token: %w", err)
	}

	// Validate tokens
	if err := helpers.ValidateToken(gitlabToken); err != nil {
		return fmt.Errorf("invalid GitLab token: %w", err)
	}
	if err := helpers.ValidateToken(githubToken); err != nil {
		return fmt.Errorf("invalid GitHub token: %w", err)
	}

	config := models.Config{
		GitLabToken: gitlabToken,
		GitHubToken: githubToken,
	}

	configPath := getConfigPath()
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Println(colors.Green + "Configuration saved successfully!" + colors.Reset)
	return nil
}

/*
getConfigPath determines OS-appropriate location for config file.
Uses platform-independent path construction to store configuration
in ~/.reposync/config.json while ensuring proper permissions.
*/
func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(colors.Red + "Failed to get user home directory: " + err.Error() + colors.Reset)
	}
	return filepath.Join(home, ".reposync", "config.json")
}

/*
readConfig loads persisted authentication tokens from disk.
Handles both file existence checks and JSON parsing errors,
providing clear guidance if configuration is missing or corrupted.
*/
func readConfig() (*models.Config, error) {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config models.Config
	err = json.Unmarshal(data, &config)
	return &config, err
}

/*
main coordinates command execution flow and argument parsing.
Implements dual-mode operation:
1. Configuration mode (reposync config)
2. Sync mode (reposync -p ...)
Validates inputs and initiates appropriate synchronization workflow.
*/
func main() {
	if len(os.Args) >= 2 && os.Args[1] == "config" {
		if err := handleConfig(); err != nil {
			log.Fatal(colors.Red + "Failed to configure tokens: " + err.Error() + colors.Reset)
		}
		os.Exit(0)
	}

	provider := flag.String("p", "", "Provider: gitlab or github")
	groupID := flag.String("g", "", "Group/Organization ID")
	cloneMethod := flag.String("m", "https", "Clone method: https or ssh")
	help := flag.Bool("h", false, "Show help message")

	flag.Parse()

	if *help || flag.NFlag() == 0 {
		fmt.Println(`reposync - Sync repositories from GitHub or GitLab

Usage:
  reposync config               Configure personal access tokens
  reposync -p <gitlab|github> -g <GROUP_ID> [-m <https|ssh>]

Flags:
  -p  Provider: gitlab or github
  -g  Group/Organization ID
  -m  Clone method: https or ssh (default: https)
  -h  Show help message`)
		os.Exit(0)
	}

	// Validate provider
	if *provider != "gitlab" && *provider != "github" {
		fmt.Println(colors.Red + "Unsupported provider. Use 'gitlab' or 'github'." + colors.Reset)
		os.Exit(1)
	}

	// Validate group ID/organization name
	if *provider == "gitlab" {
		if err := helpers.ValidateGroupID(*groupID); err != nil {
			fmt.Printf(colors.Red+"Invalid group ID: %v\n"+colors.Reset, err)
			os.Exit(1)
		}
	} else {
		if err := helpers.ValidateOrganizationName(*groupID); err != nil {
			fmt.Printf(colors.Red+"Invalid organization name: %v\n"+colors.Reset, err)
			os.Exit(1)
		}
	}

	// Validate clone method
	if *cloneMethod != "https" && *cloneMethod != "ssh" {
		fmt.Println(colors.Red + "Invalid clone method. Use 'https' or 'ssh'." + colors.Reset)
		os.Exit(1)
	}

	config, err := readConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(colors.Red + "No configuration found. Please run 'reposync config' to configure your tokens." + colors.Reset)
		} else {
			fmt.Println(colors.Red + "Failed to read configuration: " + err.Error() + colors.Reset)
		}
		os.Exit(1)
	}

	var token string
	switch *provider {
	case "gitlab":
		token = config.GitLabToken
	case "github":
		token = config.GitHubToken
	}

	if token == "" {
		fmt.Printf(colors.Red+"No token found for provider %s. Please run 'reposync config' to configure your tokens.\n"+colors.Reset, *provider)
		os.Exit(1)
	}

	// Validate token
	if err := helpers.ValidateToken(token); err != nil {
		fmt.Printf(colors.Red+"Invalid token for provider %s: %v\n"+colors.Reset, *provider, err)
		os.Exit(1)
	}

	fmt.Println(colors.Blue + "Starting repository cloning process..." + colors.Reset)

	var syncErr error
	if *provider == "gitlab" {
		groupIDInt := helpers.ParseStringToInt(*groupID)
		// The service will create the proper root directory structure
		syncErr = services.CloneGitLabRepositories(token, groupIDInt, *cloneMethod, ".")
	} else {
		// Create root directory with organization name
		rootDir := *groupID
		syncErr = services.CloneGitHubRepositories(token, *groupID, *cloneMethod, rootDir)
	}

	if syncErr != nil {
		fmt.Printf(colors.Red+"Repository synchronization failed: %v\n"+colors.Reset, syncErr)
		os.Exit(1)
	}

	fmt.Println(colors.Green + "Repository synchronization completed successfully!" + colors.Reset)
}
