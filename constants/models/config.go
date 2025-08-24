package models

/*
Config stores persisted authentication tokens and configuration for GitLab and GitHub.
Saved in JSON format in the user's home directory to avoid requiring
tokens in CLI parameters for subsequent runs.
Supports both cloud and self-hosted instances.
*/
type Config struct {
	GitLabToken string `json:"gitlab"`
	GitHubToken string `json:"github"`
	GitLabURL   string `json:"gitlab_url,omitempty"` // Support self-hosted GitLab
	GitHubURL   string `json:"github_url,omitempty"` // Support GitHub Enterprise
	CloneMethod string `json:"clone_method,omitempty"`
	MaxRetries  int    `json:"max_retries,omitempty"`
}
