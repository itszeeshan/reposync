package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
cloneRepository executes git clone command for a single repository.
Checks local filesystem first to avoid duplicate cloning,
maintaining existing repositories while synchronizing new ones.
*/
func CloneRepository(repoURL, baseDir, name string) {
	path := filepath.Join(baseDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(colors.Green + "Cloning: " + name + colors.Reset)
		cmd := exec.Command("git", "clone", repoURL, path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else {
		fmt.Println(colors.Yellow + "Skipping: " + name + " (Already cloned)" + colors.Reset)
	}
}
