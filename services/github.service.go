package services

import (
	"encoding/json"
	"fmt"
	"time"

	client "github.com/itszeeshan/reposync/client"
	colors "github.com/itszeeshan/reposync/constants/colors"
	models "github.com/itszeeshan/reposync/constants/models"
	helpers "github.com/itszeeshan/reposync/helpers"
)

/*
fetchAllGitHubRepositories fetches all repositories from a GitHub organization with pagination.
Handles GitHub's pagination by making multiple API calls until all repositories are retrieved.
Supports both cloud GitHub and GitHub Enterprise.
*/
func fetchAllGitHubRepositories(token, org, baseURL string) ([]models.GitHubRepository, error) {
	var allRepos []models.GitHubRepository
	page := 1

	for {
		url := helpers.GetGitHubAPIURL(baseURL, fmt.Sprintf("/orgs/%s/repos?per_page=100&page=%d", org, page))
		resp, err := client.Request("GET", url, token)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}
		defer resp.Body.Close()

		var repos []models.GitHubRepository
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			return nil, fmt.Errorf("failed to decode page %d: %w", page, err)
		}

		if len(repos) == 0 {
			break // No more repositories
		}

		allRepos = append(allRepos, repos...)
		page++

		// Add rate limiting to avoid hitting GitHub's rate limits
		time.Sleep(100 * time.Millisecond)
	}

	return allRepos, nil
}

/*
CloneGitHubRepositories clones all repositories in a GitHub organization.
Handles pagination through fetchAllGitHubRepositories,
cloning all repositories in flat structure under specified base directory.
Supports both cloud GitHub and GitHub Enterprise.
*/
func CloneGitHubRepositories(token string, org string, cloneMethod string, baseDir string) error {
	return CloneGitHubRepositoriesWithURL(token, org, cloneMethod, baseDir, "")
}

/*
CloneGitHubRepositoriesWithURL clones all repositories in a GitHub organization with custom URL.
Allows specifying custom GitHub instance URL for self-hosted installations.
*/
func CloneGitHubRepositoriesWithURL(token string, org string, cloneMethod string, baseDir string, baseURL string) error {
	// Validate inputs
	if err := helpers.ValidateOrganizationName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}

	fmt.Println(colors.Cyan + "Fetching GitHub repositories..." + colors.Reset)

	repositories, err := fetchAllGitHubRepositories(token, org, baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch repositories: %w", err)
	}

	fmt.Printf("Found %d repositories\n", len(repositories))

	for i, repository := range repositories {
		fmt.Printf("Progress: %d/%d (%.1f%%)\n", i+1, len(repositories), float64(i+1)/float64(len(repositories))*100)

		repoURL := helpers.GetPreferredRepositoryURL(repository.HTTPSURL, repository.SSHURL, cloneMethod)
		if err := helpers.CloneRepository(repoURL, baseDir, repository.Name, token); err != nil {
			fmt.Printf(colors.Red+"Failed to clone %s: %v\n"+colors.Reset, repository.Name, err)
			continue // Continue with other repos
		}
	}

	return nil
}
