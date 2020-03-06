package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func (c Config) isEnv(branch string) bool {
	// Check if the branch string passed as argument is one of the managed branches
	if branch == c.ProdBranch {
		return true
	}
	for _, b := range c.OtherBranches {
		if branch == b {
			return true
		}
	}
	return false
}

func (c Config) isProd(branch string) bool {
	// Check if the branch string passed as argument is the production branch
	return branch == c.ProdBranch
}

func (c Config) getProdRemote() string {
	// Get the remote for the production branch
	stdout, err := exec.Command("git", "config", "branch."+c.ProdBranch+".remote").Output()
	if err != nil {
		log.Fatalf("Failed to get remote of %s branch.", c.ProdBranch)
	}
	return string(stdout)[:len(stdout)-1]
}

func gitCommand(dryRun bool, args ...string) {
	// Run a generic Git Command
	runCommand(dryRun, "git", args...)
}

func runCommand(dryRun bool, cmd string, args ...string) {
	// Run a generic command
	fmt.Printf("+ %s %s\n", cmd, strings.Join(args, " "))
	if !dryRun {
		c := exec.Command(cmd, args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		err := c.Run()

		if err != nil {
			log.Fatalf("Failed executing command: %s %s", cmd, strings.Join(args, " "))
		}
	}
}

func gitBranch() (string, error) {
	// Wrapper for git branch command
	stdout, err := exec.Command("git", "branch").Output()
	return string(stdout), err
}

func gitRefsExists(ref string) (string, error) {
	stdout, err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+ref).Output()
	return string(stdout), err
}

func getCurrentBranch_(gitBranch func() (string, error)) (string, error) {
	// Get the current git branch
	stdout, err := gitBranch()
	if err != nil {
		return "", err
	}

	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "* ") {
			items := strings.Split(line, " ")

			return items[1], nil
		}
	}

	return "", errors.New("could not detect current branch")
}

func getCurrentBranch() (string, error) {
	return getCurrentBranch_(gitBranch)
}

func getGitRemoteUrl() (string, error) {
	stdout, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()

	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

func getGitlabMRUrl(dryRun bool, pushBranch string, deployEnv string) {
	if !dryRun {
		git_url, err := getGitRemoteUrl()

		if err != nil {
			fmt.Printf("+ Failed to generate MR Url, failed to fetch remote url\n")
			return
		}

		re_git := regexp.MustCompile(`^git.+`)
		if re_git.Match([]byte(git_url)) {
			git_url = strings.TrimSuffix(strings.Replace(git_url, "git@gitlab.com:", "https://gitlab.com/", 1), ".git\n")
		}

		Url, err := url.Parse(git_url)
		Url.Path += "/-/merge_requests/new"
		params := url.Values{}
		params.Add("merge_request[source_branch]", pushBranch)
		params.Add("merge_request[target_branch]", deployEnv)
		params.Add("merge_request[title]", "Merge "+pushBranch+" into "+deployEnv)
		params.Add("merge_request[squash]", "false")
		params.Add("merge_request[remove_source_branch]", "true")
		Url.RawQuery = params.Encode()

		fmt.Printf("+ Create a Giltab Merge Request for branch %s to environment %s\n+\n", pushBranch, deployEnv)
		fmt.Printf("+ %s\n", Url.String())
	}
}
