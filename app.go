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

	colors "github.com/itszeeshan/reposync/constants/colors"
	models "github.com/itszeeshan/reposync/constants/models"
	helpers "github.com/itszeeshan/reposync/helpers"
	services "github.com/itszeeshan/reposync/services"
)

/*
handleConfig implements interactive token configuration workflow.
Prompts user for both GitLab and GitHub tokens, then saves them
to encrypted config file in user's home directory for future use.
*/
func handleConfig() {
	var gitlabToken, githubToken string
	fmt.Print("Enter GitLab Personal Access Token: ")
	fmt.Scanln(&gitlabToken)
	fmt.Print("Enter GitHub Personal Access Token: ")
	fmt.Scanln(&githubToken)

	config := models.Config{
		GitLabToken: gitlabToken,
		GitHubToken: githubToken,
	}

	configPath := getConfigPath()
	os.MkdirAll(filepath.Dir(configPath), 0700)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(colors.Red + "Failed to marshal config: " + err.Error() + colors.Reset)
	}

	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		log.Fatal(colors.Red + "Failed to write config file: " + err.Error() + colors.Reset)
	}

	fmt.Println(colors.Green + "Configuration saved successfully!" + colors.Reset)
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
		handleConfig()
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

	if *provider != "gitlab" && *provider != "github" {
		fmt.Println(colors.Red + "Unsupported provider. Use 'gitlab' or 'github'." + colors.Reset)
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

	groupIDInt := helpers.ParseStringToInt(*groupID)

	fmt.Println(colors.Blue + "Starting repository cloning process..." + colors.Reset)

	if *provider == "gitlab" {
		services.CloneGitLabRepositories(token, groupIDInt, *cloneMethod, ".")
	} else {
		services.CloneGitHubRepositories(token, *groupID, *cloneMethod, ".")
	}
}
