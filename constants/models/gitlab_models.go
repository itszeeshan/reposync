package models

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
