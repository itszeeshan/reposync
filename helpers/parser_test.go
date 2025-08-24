package helpers

import (
	"testing"
)

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"empty token", "", true},
		{"short token", "short", true},
		{"valid token", "glpat_abcdefghijklmnopqrstuvwxyz1234567890", false},
		{"github token", "ghp_abcdefghijklmnopqrstuvwxyz1234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOrganizationName(t *testing.T) {
	tests := []struct {
		name    string
		org     string
		wantErr bool
	}{
		{"empty org", "", true},
		{"valid org", "my-org", false},
		{"valid org with numbers", "my-org-123", false},
		{"invalid org with spaces", "my org", true},
		{"invalid org with special chars", "my@org", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrganizationName(tt.org)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrganizationName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGroupID(t *testing.T) {
	tests := []struct {
		name    string
		groupID string
		wantErr bool
	}{
		{"empty group ID", "", true},
		{"valid group ID", "123456", false},
		{"invalid group ID", "abc", true},
		{"invalid group ID with letters", "123abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGroupID(tt.groupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGroupID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetGitLabAPIURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		endpoint string
		want     string
	}{
		{"default URL", "", "/groups/123", "https://gitlab.com/api/v4/groups/123"},
		{"custom URL", "https://gitlab.company.com", "/groups/123", "https://gitlab.company.com/api/v4/groups/123"},
		{"URL with trailing slash", "https://gitlab.company.com/", "/groups/123", "https://gitlab.company.com/api/v4/groups/123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGitLabAPIURL(tt.baseURL, tt.endpoint)
			if got != tt.want {
				t.Errorf("GetGitLabAPIURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGitHubAPIURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		endpoint string
		want     string
	}{
		{"default URL", "", "/orgs/my-org", "https://api.github.com/orgs/my-org"},
		{"custom URL", "https://github.company.com/api/v3", "/orgs/my-org", "https://github.company.com/api/v3/orgs/my-org"},
		{"URL with trailing slash", "https://github.company.com/api/v3/", "/orgs/my-org", "https://github.company.com/api/v3/orgs/my-org"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGitHubAPIURL(tt.baseURL, tt.endpoint)
			if got != tt.want {
				t.Errorf("GetGitHubAPIURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
