package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

type GitLabProject struct {
	HTTPURLToRepo string `json:"http_url_to_repo"`
	SSHURLToRepo  string `json:"ssh_url_to_repo"`
	Name          string `json:"name"`
}

type GitLabGroup struct {
	ID       int    `json:"id"`
	FullPath string `json:"full_path"`
}

type GitHubRepo struct {
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
	Name     string `json:"name"`
}

func fetchAndCloneGitLab(token string, groupID int, method string, baseDir string) {
	fmt.Println(Cyan + "Fetching GitLab repositories..." + Reset)
	subgroupsURL := fmt.Sprintf("https://gitlab.com/api/v4/groups/%d/subgroups", groupID)
	req, err := http.NewRequest("GET", subgroupsURL, nil)
	if err != nil {
		log.Fatalf(Red+"Failed to create request: %v"+Reset, err)
	}
	req.Header.Set("Private-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf(Red+"Failed to fetch subgroups: %v"+Reset, err)
	}
	defer resp.Body.Close()

	var subgroups []GitLabGroup
	json.NewDecoder(resp.Body).Decode(&subgroups)

	for _, subgroup := range subgroups {
		newPath := filepath.Join(baseDir, subgroup.FullPath)
		os.MkdirAll(newPath, os.ModePerm)
		fmt.Println(Yellow + "Processing subgroup: " + subgroup.FullPath + Reset)
		fetchAndCloneGitLab(token, subgroup.ID, method, newPath)
	}

	projectsURL := fmt.Sprintf("https://gitlab.com/api/v4/groups/%d/projects", groupID)
	req, err = http.NewRequest("GET", projectsURL, nil)
	if err != nil {
		log.Fatalf(Red+"Failed to create request: %v"+Reset, err)
	}
	req.Header.Set("Private-Token", token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf(Red+"Failed to fetch projects: %v"+Reset, err)
	}
	defer resp.Body.Close()

	var projects []GitLabProject
	json.NewDecoder(resp.Body).Decode(&projects)

	for _, project := range projects {
		repoURL := project.HTTPURLToRepo
		if method == "ssh" {
			repoURL = project.SSHURLToRepo
		}
		cloneRepo(repoURL, baseDir, project.Name)
	}
}

func fetchAndCloneGitHub(token string, org string, method string, baseDir string) {
	fmt.Println(Cyan + "Fetching GitHub repositories..." + Reset)
	reposURL := fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=100", org)
	req, err := http.NewRequest("GET", reposURL, nil)
	if err != nil {
		log.Fatalf(Red+"Failed to create request: %v"+Reset, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf(Red+"Failed to fetch repositories: %v"+Reset, err)
	}
	defer resp.Body.Close()

	var repos []GitHubRepo
	json.NewDecoder(resp.Body).Decode(&repos)

	for _, repo := range repos {
		repoURL := repo.CloneURL
		if method == "ssh" {
			repoURL = repo.SSHURL
		}
		cloneRepo(repoURL, baseDir, repo.Name)
	}
}

func cloneRepo(repoURL, baseDir, name string) {
	path := filepath.Join(baseDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(Green + "Cloning: " + name + Reset)
		cmd := exec.Command("git", "clone", repoURL, path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else {
		fmt.Println(Yellow + "Skipping: " + name + " (Already cloned)" + Reset)
	}
}

func main() {
	provider := flag.String("p", "", "Provider: gitlab or github")
	token := flag.String("t", "", "Personal Access Token")
	groupID := flag.String("g", "", "Group/Organization ID")
	method := flag.String("m", "https", "Clone method: https or ssh")
	help := flag.Bool("h", false, "Show help message")
	// Override default help message
	flag.Usage = func() {
		fmt.Println(`reposync - Sync repositories from GitHub or GitLab
	
Usage:
  reposync -p <gitlab|github> -t <TOKEN> -g <GROUP_ID> [-m <https|ssh>]

Flags:
  -p  Provider: gitlab or github
  -t  Personal Access Token
  -g  Group/Organization ID
  -m  Clone method: https or ssh (default: https)
  -h  Show help message`)
	}
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println("Usage: reposync -p <gitlab|github> -t <TOKEN> -g <GROUP_ID> [-m <https|ssh>]")
		os.Exit(1)
	}
	if len(flag.Args()) > 0 {
		fmt.Println(Red+"Unknown argument(s):", flag.Args(), Reset)
		fmt.Println("Use -h or --help to see the available options.")
		os.Exit(1)
	}
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Println(Blue + "Starting repository cloning process..." + Reset)

	if *provider == "gitlab" {
		fetchAndCloneGitLab(*token, atoi(*groupID), *method, ".")
	} else if *provider == "github" {
		fetchAndCloneGitHub(*token, *groupID, *method, ".")
	} else {
		fmt.Println(Red + "Unsupported provider. Use 'gitlab' or 'github'." + Reset)
	}
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf(Red+"Invalid group ID: %s"+Reset, s)
	}
	return n
}
