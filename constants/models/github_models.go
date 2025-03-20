package models

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
