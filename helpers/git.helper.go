package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	colors "github.com/itszeeshan/reposync/constants/colors"
)

/*
GetPreferredRepositoryURL determines clone URL based on user preference.
Selects between HTTPS and SSH URLs based on -m flag value,
enabling flexible cloning methods depending on user's authentication setup.
*/
func GetPreferredRepositoryURL(httpsURL, sshURL, method string) string {
	if method == "ssh" {
		return sshURL
	}
	return httpsURL
}

/*
CloneRepository executes git clone command for a single repository.
Checks local filesystem first to avoid duplicate cloning,
maintaining existing repositories while synchronizing new ones.
Includes retry logic for better reliability and token-based authentication as fallback.
*/
func CloneRepository(repoURL, baseDir, name, token string) error {
	path := filepath.Join(baseDir, name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(colors.Green + "Cloning: " + name + colors.Reset)

		// Add retry logic for better reliability
		maxRetries := 3
		for attempt := 1; attempt <= maxRetries; attempt++ {
			var cmd *exec.Cmd

			// First try without authentication (works for public repos and configured credentials)
			if attempt == 1 {
				cmd = exec.Command("git", "clone", repoURL, path)
			} else {
				// On retry, use token authentication as fallback
				authenticatedURL := repoURL
				if token != "" && isHTTPSURL(repoURL) {
					authenticatedURL = constructAuthenticatedURL(repoURL, token)
				}
				cmd = exec.Command("git", "clone", authenticatedURL, path)
			}

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				if attempt == maxRetries {
					return fmt.Errorf("git clone failed for %s after %d attempts: %w", name, maxRetries, err)
				}
				fmt.Printf(colors.Yellow+"Attempt %d failed, retrying with authentication in %d seconds...\n"+colors.Reset, attempt, attempt)
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			break
		}
	} else {
		fmt.Println(colors.Yellow + "Skipping: " + name + " (Already cloned)" + colors.Reset)
	}
	return nil
}

/*
isHTTPSURL checks if the given URL is an HTTPS URL.
*/
func isHTTPSURL(url string) bool {
	return len(url) >= 8 && url[:8] == "https://"
}

/*
constructAuthenticatedURL constructs an authenticated URL for HTTPS cloning.
Inserts the token into the URL for GitLab/GitHub authentication.
*/
func constructAuthenticatedURL(originalURL, token string) string {
	// Replace https:// with https://oauth2:token@
	return "https://oauth2:" + token + "@" + originalURL[8:]
}
