package client

import (
	"fmt"
	"log"
	"net/http"
	"os"

	colors "github.com/itszeeshan/reposync/constants/colors"
)

/*
MakeHTTPRequest executes authenticated API requests to GitLab/GitHub.
Adds Bearer token authentication header and handles HTTP errors:
- 401 Unauthorized: Guides user to reconfigure tokens
- Other errors: Fatal exit with status code
*/
func Request(method, url, token string) *http.Response {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf(colors.Red+"Failed to create request: %v"+colors.Reset, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf(colors.Red+"Failed to fetch data: %v"+colors.Reset, err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println(colors.Red + "Permission denied. Check if your token is valid. Otherwise, run 'reposync config' to configure again." + colors.Reset)
		os.Exit(1)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalf(colors.Red+"Request failed with status code: %d"+colors.Reset, resp.StatusCode)
	}

	return resp
}
