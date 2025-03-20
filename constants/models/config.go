package models

/*
Config stores persisted authentication tokens for GitLab and GitHub.
Saved in JSON format in the user's home directory to avoid requiring
tokens in CLI parameters for subsequent runs.
*/
type Config struct {
	GitLabToken string `json:"gitlab"`
	GitHubToken string `json:"github"`
}
