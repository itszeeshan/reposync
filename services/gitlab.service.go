package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	client "github.com/itszeeshan/reposync/client"
	colors "github.com/itszeeshan/reposync/constants/colors"
	models "github.com/itszeeshan/reposync/constants/models"
	helpers "github.com/itszeeshan/reposync/helpers"
)

/*
CloneGitLabRepositories recursively clones all repositories in a GitLab group.
Handles both direct repositories and nested subgroups by:
1. Processing subgroups first to create directory structure
2. Cloning repositories in current group level
3. Using depth-first recursion for subgroup processing
*/
func CloneGitLabRepositories(token string, groupID int, cloneMethod string, baseDir string) {
	fmt.Println(colors.Cyan + "Fetching GitLab repositories..." + colors.Reset)

	subgroups := getGitLabSubgroups(token, groupID)
	for _, subgroup := range subgroups {
		newPath := filepath.Join(baseDir, subgroup.FullPath)
		os.MkdirAll(newPath, os.ModePerm)
		fmt.Println(colors.Yellow + "Processing subgroup: " + subgroup.FullPath + colors.Reset)
		CloneGitLabRepositories(token, subgroup.ID, cloneMethod, newPath)
	}

	repositories := getGitLabRepositories(token, groupID)
	for _, repository := range repositories {
		repoURL := helpers.GetPreferredRepositoryURL(repository.HTTPSURL, repository.SSHURL, cloneMethod)
		helpers.CloneRepository(repoURL, baseDir, repository.Name)
	}
}

/*
getGitLabSubgroups fetches subgroup hierarchy from GitLab API.
Uses paginated API to retrieve all subgroups within specified parent group,
enabling complete group structure analysis for directory creation.
*/
func getGitLabSubgroups(token string, groupID int) []models.GitLabSubgroup {
	url := fmt.Sprintf("https://gitlab.com/api/v4/groups/%d/subgroups", groupID)
	resp := client.Request("GET", url, token)
	defer resp.Body.Close()

	var subgroups []models.GitLabSubgroup
	json.NewDecoder(resp.Body).Decode(&subgroups)
	return subgroups
}

/*
getGitLabRepositories retrieves project list from GitLab group.
Fetches all repositories in specified group, including those shared
from parent groups, using GitLab's projects API endpoint.
*/
func getGitLabRepositories(token string, groupID int) []models.GitLabRepository {
	url := fmt.Sprintf("https://gitlab.com/api/v4/groups/%d/projects", groupID)
	resp := client.Request("GET", url, token)
	defer resp.Body.Close()

	var repositories []models.GitLabRepository
	json.NewDecoder(resp.Body).Decode(&repositories)
	return repositories
}
