package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	client "github.com/itszeeshan/reposync/client"
	colors "github.com/itszeeshan/reposync/constants/colors"
	models "github.com/itszeeshan/reposync/constants/models"
	helpers "github.com/itszeeshan/reposync/helpers"
)

/*
getGitLabSubgroups fetches subgroup hierarchy from GitLab API.
Uses paginated API to retrieve all subgroups within specified parent group,
enabling complete group structure analysis for directory creation.
Supports both cloud GitLab and self-hosted instances.
*/
func getGitLabSubgroups(token string, groupID int, baseURL string) ([]models.GitLabSubgroup, error) {
	url := helpers.GetGitLabAPIURL(baseURL, fmt.Sprintf("/groups/%d/subgroups", groupID))
	resp, err := client.Request("GET", url, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subgroups: %w", err)
	}
	defer resp.Body.Close()

	var subgroups []models.GitLabSubgroup
	if err := json.NewDecoder(resp.Body).Decode(&subgroups); err != nil {
		return nil, fmt.Errorf("failed to decode subgroups: %w", err)
	}
	return subgroups, nil
}

/*
getGitLabRepositories retrieves project list from GitLab group.
Fetches all repositories in specified group, including those shared
from parent groups, using GitLab's projects API endpoint.
Supports both cloud GitLab and self-hosted instances.
*/
func getGitLabRepositories(token string, groupID int, baseURL string) ([]models.GitLabRepository, error) {
	url := helpers.GetGitLabAPIURL(baseURL, fmt.Sprintf("/groups/%d/projects", groupID))
	resp, err := client.Request("GET", url, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}
	defer resp.Body.Close()

	var repositories []models.GitLabRepository
	if err := json.NewDecoder(resp.Body).Decode(&repositories); err != nil {
		return nil, fmt.Errorf("failed to decode repositories: %w", err)
	}
	return repositories, nil
}

/*
getGitLabGroupInfo fetches basic information about a GitLab group.
Returns the group name and path for directory structure creation.
*/
func getGitLabGroupInfo(token string, groupID int, baseURL string) (string, string, error) {
	url := helpers.GetGitLabAPIURL(baseURL, fmt.Sprintf("/groups/%d", groupID))
	resp, err := client.Request("GET", url, token)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch group info: %w", err)
	}
	defer resp.Body.Close()

	var groupInfo struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&groupInfo); err != nil {
		return "", "", fmt.Errorf("failed to decode group info: %w", err)
	}
	return groupInfo.Name, groupInfo.Path, nil
}

/*
CloneGitLabRepositories recursively clones all repositories in a GitLab group.
Handles both direct repositories and nested subgroups by:
1. Processing subgroups first to create directory structure
2. Cloning repositories in current group level
3. Using depth-first recursion for subgroup processing
Supports both cloud GitLab and self-hosted instances.
*/
func CloneGitLabRepositories(token string, groupID int, cloneMethod string, baseDir string) error {
	return CloneGitLabRepositoriesWithURL(token, groupID, cloneMethod, baseDir, "")
}

/*
CloneGitLabRepositoriesWithURL recursively clones all repositories in a GitLab group with custom URL.
Allows specifying custom GitLab instance URL for self-hosted installations.
*/
func CloneGitLabRepositoriesWithURL(token string, groupID int, cloneMethod string, baseDir string, baseURL string) error {
	fmt.Println(colors.Cyan + "Fetching GitLab repositories..." + colors.Reset)

	// Get group info to create proper root directory
	groupName, groupPath, err := getGitLabGroupInfo(token, groupID, baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch group info: %w", err)
	}

	// Create root directory with group path
	rootDir := filepath.Join(baseDir, groupPath)
	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create root directory %s: %w", rootDir, err)
	}

	fmt.Printf("Creating directory structure for group: %s (%s)\n", groupName, groupPath)

	// Process all subgroups first to create directory structure
	subgroups, err := getGitLabSubgroups(token, groupID, baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch subgroups: %w", err)
	}

	for _, subgroup := range subgroups {
		fmt.Println(colors.Yellow + "Processing subgroup: " + subgroup.FullPath + colors.Reset)

		// Recursively process the subgroup - pass the root directory
		if err := CloneGitLabRepositoriesWithURL(token, subgroup.ID, cloneMethod, rootDir, baseURL); err != nil {
			fmt.Printf(colors.Red+"Failed to process subgroup %s: %v\n"+colors.Reset, subgroup.FullPath, err)
			continue // Continue with other subgroups
		}
	}

	// Process repositories in current group
	repositories, err := getGitLabRepositories(token, groupID, baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch repositories: %w", err)
	}

	fmt.Printf("Found %d repositories in current group\n", len(repositories))

	for i, repository := range repositories {
		fmt.Printf("Progress: %d/%d (%.1f%%)\n", i+1, len(repositories), float64(i+1)/float64(len(repositories))*100)

		repoURL := helpers.GetPreferredRepositoryURL(repository.HTTPSURL, repository.SSHURL, cloneMethod)
		if err := helpers.CloneRepository(repoURL, rootDir, repository.Path, token); err != nil {
			fmt.Printf(colors.Red+"Failed to clone %s: %v\n"+colors.Reset, repository.Name, err)
			continue // Continue with other repos
		}
	}

	// Add rate limiting to avoid hitting GitLab's rate limits
	time.Sleep(100 * time.Millisecond)

	return nil
}
