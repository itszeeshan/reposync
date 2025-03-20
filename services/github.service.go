package services

import (
	"encoding/json"
	"fmt"

	client "github.com/itszeeshan/reposync/client"
	colors "github.com/itszeeshan/reposync/constants/colors"
	models "github.com/itszeeshan/reposync/constants/models"
	helpers "github.com/itszeeshan/reposync/helpers"
)

/*
CloneGitHubRepositories clones all repositories in a GitHub organization.
Handles pagination implicitly through GitHub's API response structure,
cloning all repositories in flat structure under specified base directory.
*/
func CloneGitHubRepositories(token string, org string, cloneMethod string, baseDir string) {
	fmt.Println(colors.Cyan + "Fetching GitHub repositories..." + colors.Reset)
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=100", org)
	resp := client.Request("GET", url, token)
	defer resp.Body.Close()

	var repositories []models.GitHubRepository
	json.NewDecoder(resp.Body).Decode(&repositories)

	for _, repository := range repositories {
		repoURL := helpers.GetPreferredRepositoryURL(repository.HTTPSURL, repository.SSHURL, cloneMethod)
		helpers.CloneRepository(repoURL, baseDir, repository.Name)
	}
}
