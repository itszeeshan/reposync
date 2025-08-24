package helpers

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	colors "github.com/itszeeshan/reposync/constants/colors"
)

/*
parseStringToInt safely converts group ID string to integer.
Provides user-friendly error handling for invalid numeric inputs,
ensuring valid API requests with properly formatted group IDs.
*/
func ParseStringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf(colors.Red+"Invalid group ID: %s"+colors.Reset, s)
	}
	return n
}

/*
ValidateToken validates token format and length.
Ensures tokens meet minimum security requirements.
*/
func ValidateToken(token string) error {
	if token == "" {
		return errors.New("token cannot be empty")
	}
	if len(token) < 10 {
		return errors.New("token appears to be too short")
	}
	return nil
}

/*
ValidateOrganizationName validates organization name format.
Ensures organization names contain only valid characters.
*/
func ValidateOrganizationName(org string) error {
	if org == "" {
		return errors.New("organization name cannot be empty")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(org) {
		return errors.New("invalid organization name format")
	}
	return nil
}

/*
ValidateGroupID validates group ID format.
Ensures group ID is a valid integer.
*/
func ValidateGroupID(groupID string) error {
	if groupID == "" {
		return errors.New("group ID cannot be empty")
	}
	if _, err := strconv.Atoi(groupID); err != nil {
		return errors.New("group ID must be a valid integer")
	}
	return nil
}

/*
GetGitLabAPIURL constructs the GitLab API URL for a given endpoint.
Supports both cloud GitLab and self-hosted instances.
*/
func GetGitLabAPIURL(baseURL, endpoint string) string {
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	// Ensure baseURL doesn't end with slash
	baseURL = strings.TrimSuffix(baseURL, "/")
	return fmt.Sprintf("%s/api/v4%s", baseURL, endpoint)
}

/*
GetGitHubAPIURL constructs the GitHub API URL for a given endpoint.
Supports both cloud GitHub and GitHub Enterprise.
*/
func GetGitHubAPIURL(baseURL, endpoint string) string {
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	// Ensure baseURL doesn't end with slash
	baseURL = strings.TrimSuffix(baseURL, "/")
	return fmt.Sprintf("%s%s", baseURL, endpoint)
}
