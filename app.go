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
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// ANSI escape codes for terminal text coloring
// These constants provide consistent color formatting for different message types:
// - Reset returns to default terminal colors
// - Colors are used for success (Green), warnings (Yellow), errors (Red), and information (Blue/Cyan)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

/*
GitLabRepository represents a GitLab project with its clone URLs.
Contains both HTTPS and SSH URLs for cloning, and the repository name
to maintain directory structure during cloning operations.
*/

type GitLabRepository struct {
	HTTPSURL string `json:"http_url_to_repo"`
	SSHURL   string `json:"ssh_url_to_repo"`
	Name     string `json:"name"`
}

/*
GitLabSubgroup represents a nested group structure within GitLab.
Stores the subgroup ID for API navigation and full path for directory structure
replication during repository cloning.
*/
type GitLabSubgroup struct {
	ID       int    `json:"id"`
	FullPath string `json:"full_path"`
}

/*
GitHubRepository represents a GitHub repository with clone information.
Similar to GitLabRepository but matches GitHub's API response structure,
providing both clone URLs and repository name for organization.
*/
type GitHubRepository struct {
	HTTPSURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
	Name     string `json:"name"`
}

/*
Config stores persisted authentication tokens for GitLab and GitHub.
Saved in JSON format in the user's home directory to avoid requiring
tokens in CLI parameters for subsequent runs.
*/
type Config struct {
	GitLabToken string `json:"gitlab"`
	GitHubToken string `json:"github"`
}

/*
cloneGitLabRepositories recursively clones all repositories in a GitLab group.
Handles both direct repositories and nested subgroups by:
1. Processing subgroups first to create directory structure
2. Cloning repositories in current group level
3. Using depth-first recursion for subgroup processing
*/
func cloneGitLabRepositories(token string, groupID int, cloneMethod string, baseDir string) {
	fmt.Println(Cyan + "Fetching GitLab repositories..." + Reset)

	subgroups := getGitLabSubgroups(token, groupID)
	for _, subgroup := range subgroups {
		newPath := filepath.Join(baseDir, subgroup.FullPath)
		os.MkdirAll(newPath, os.ModePerm)
		fmt.Println(Yellow + "Processing subgroup: " + subgroup.FullPath + Reset)
		cloneGitLabRepositories(token, subgroup.ID, cloneMethod, newPath)
	}

	repositories := getGitLabRepositories(token, groupID)
	for _, repository := range repositories {
		repoURL := getPreferredRepositoryURL(repository.HTTPSURL, repository.SSHURL, cloneMethod)
		cloneRepository(repoURL, baseDir, repository.Name)
	}
}

/*
getGitLabSubgroups fetches subgroup hierarchy from GitLab API.
Uses paginated API to retrieve all subgroups within specified parent group,
enabling complete group structure analysis for directory creation.
*/
func getGitLabSubgroups(token string, groupID int) []GitLabSubgroup {
	url := fmt.Sprintf("https://gitlab.com/api/v4/groups/%d/subgroups", groupID)
	resp := makeHTTPRequest(url, token)
	defer resp.Body.Close()

	var subgroups []GitLabSubgroup
	json.NewDecoder(resp.Body).Decode(&subgroups)
	return subgroups
}

/*
getGitLabRepositories retrieves project list from GitLab group.
Fetches all repositories in specified group, including those shared
from parent groups, using GitLab's projects API endpoint.
*/
func getGitLabRepositories(token string, groupID int) []GitLabRepository {
	url := fmt.Sprintf("https://gitlab.com/api/v4/groups/%d/projects", groupID)
	resp := makeHTTPRequest(url, token)
	defer resp.Body.Close()

	var repositories []GitLabRepository
	json.NewDecoder(resp.Body).Decode(&repositories)
	return repositories
}

/*
cloneGitHubRepositories clones all repositories in a GitHub organization.
Handles pagination implicitly through GitHub's API response structure,
cloning all repositories in flat structure under specified base directory.
*/
func cloneGitHubRepositories(token string, org string, cloneMethod string, baseDir string) {
	fmt.Println(Cyan + "Fetching GitHub repositories..." + Reset)
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=100", org)
	resp := makeHTTPRequest(url, token)
	defer resp.Body.Close()

	var repositories []GitHubRepository
	json.NewDecoder(resp.Body).Decode(&repositories)

	for _, repository := range repositories {
		repoURL := getPreferredRepositoryURL(repository.HTTPSURL, repository.SSHURL, cloneMethod)
		cloneRepository(repoURL, baseDir, repository.Name)
	}
}

/*
makeHTTPRequest executes authenticated API requests to GitLab/GitHub.
Adds Bearer token authentication header and handles HTTP errors:
- 401 Unauthorized: Guides user to reconfigure tokens
- Other errors: Fatal exit with status code
*/
func makeHTTPRequest(url, token string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf(Red+"Failed to create request: %v"+Reset, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf(Red+"Failed to fetch data: %v"+Reset, err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println(Red + "Permission denied. Check if your token is valid. Otherwise, run 'reposync config' to configure again." + Reset)
		os.Exit(1)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalf(Red+"Request failed with status code: %d"+Reset, resp.StatusCode)
	}

	return resp
}

/*
getPreferredRepositoryURL determines clone URL based on user preference.
Selects between HTTPS and SSH URLs based on -m flag value,
enabling flexible cloning methods depending on user's authentication setup.
*/
func getPreferredRepositoryURL(httpsURL, sshURL, method string) string {
	if method == "ssh" {
		return sshURL
	}
	return httpsURL
}

/*
cloneRepository executes git clone command for a single repository.
Checks local filesystem first to avoid duplicate cloning,
maintaining existing repositories while synchronizing new ones.
*/
func cloneRepository(repoURL, baseDir, name string) {
	path := filepath.Join(baseDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(Green + "Cloning: " + name + Reset)
		cmd := exec.Command("git", "clone", repoURL, path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else {
		fmt.Println(Yellow + "Skipping: " + name + " (Already cloned)" + Reset)
	}
}

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

	config := Config{
		GitLabToken: gitlabToken,
		GitHubToken: githubToken,
	}

	configPath := getConfigPath()
	os.MkdirAll(filepath.Dir(configPath), 0700)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(Red + "Failed to marshal config: " + err.Error() + Reset)
	}

	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		log.Fatal(Red + "Failed to write config file: " + err.Error() + Reset)
	}

	fmt.Println(Green + "Configuration saved successfully!" + Reset)
}

/*
getConfigPath determines OS-appropriate location for config file.
Uses platform-independent path construction to store configuration
in ~/.reposync/config.json while ensuring proper permissions.
*/
func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(Red + "Failed to get user home directory: " + err.Error() + Reset)
	}
	return filepath.Join(home, ".reposync", "config.json")
}

/*
readConfig loads persisted authentication tokens from disk.
Handles both file existence checks and JSON parsing errors,
providing clear guidance if configuration is missing or corrupted.
*/
func readConfig() (*Config, error) {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
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
		fmt.Println(Red + "Unsupported provider. Use 'gitlab' or 'github'." + Reset)
		os.Exit(1)
	}

	config, err := readConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(Red + "No configuration found. Please run 'reposync config' to configure your tokens." + Reset)
		} else {
			fmt.Println(Red + "Failed to read configuration: " + err.Error() + Reset)
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
		fmt.Printf(Red+"No token found for provider %s. Please run 'reposync config' to configure your tokens.\n"+Reset, *provider)
		os.Exit(1)
	}

	groupIDInt := parseStringToInt(*groupID)

	fmt.Println(Blue + "Starting repository cloning process..." + Reset)

	if *provider == "gitlab" {
		cloneGitLabRepositories(token, groupIDInt, *cloneMethod, ".")
	} else {
		cloneGitHubRepositories(token, *groupID, *cloneMethod, ".")
	}
}

/*
parseStringToInt safely converts group ID string to integer.
Provides user-friendly error handling for invalid numeric inputs,
ensuring valid API requests with properly formatted group IDs.
*/
func parseStringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf(Red+"Invalid group ID: %s"+Reset, s)
	}
	return n
}
