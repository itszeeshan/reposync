package client

import (
	"fmt"
	"net/http"
)

/*
Request executes authenticated API requests to GitLab/GitHub.
Adds Bearer token authentication header and handles HTTP errors:
- 401 Unauthorized: Returns permission denied error
- 429 Too Many Requests: Returns rate limit error
- Other errors: Returns appropriate error with status code
*/
func Request(method, url, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", "RepoSync/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("permission denied - check if your token is valid")
	} else if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limit exceeded - please wait and try again")
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return resp, nil
}
