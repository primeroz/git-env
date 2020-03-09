package main

import (
	"errors"
	"os/exec"
	"strings"
)

type Option struct {
	Name     string
	Question string
	Default  string
}

type Config struct {
	Mode          string
	DeployHook    string
	ProdBranch    string
	OtherBranches []string
	ProdDeployCmd string
}

var (
	config   *Config
	settings = []Option{
		{
			Name:     "rerere.enabled",
			Question: "Enable git.rerere in your local config ?",
			Default:  "true",
		},
		{
			Name:     "merge.conflictstyle",
			Question: "Set conflictstyle in your global config",
			Default:  "diff3",
		},
	}
	options = []Option{
		{
			Name:     "mode",
			Question: "What type of workflow to use ? ( push / merge )",
			Default:  "push",
		},
		{
			Name:     "deploy-hook",
			Question: "Run hook command , relative to root of repo, on deploy",
			Default:  "exit 0",
			// Default:  "",
			// Default: "make && git add -A && git diff-index --quiet HEAD || git commit -m \"deploy hook commit\"",
		},
		{
			Name:     "prod",
			Question: "What is your production branch?",
			Default:  "master",
		},
		{
			Name:     "other",
			Question: "What other environment branches do you have?",
			Default:  "stage dev",
		},
		{
			Name:     "prod-deploy",
			Question: "What command should be run to deploy to the production branch? - Default push to origin for MR",
			Default:  "git push origin {{.feature}}",
			//Default:  "git checkout {{.env}} && git merge --no-ff {{.feature}}",
		},
	}
)

const configPrefix = "env-branch"

func loadConfig_(getOption func(string) (string, error)) (*Config, error) {
	// Load Configuration values from git config
	config := Config{}

	cfg := map[string]string{}

	for _, opt := range options {
		s, err := getOption(opt.Name)
		if err != nil {
			return nil, err
		}
		cfg[opt.Name] = s
	}
	config.Mode = cfg["mode"]
	config.DeployHook = cfg["deploy-hook"]
	config.ProdBranch = cfg["prod"]
	config.OtherBranches = strings.Split(cfg["other"], " ")
	config.ProdDeployCmd = cfg["prod-deploy"]

	return &config, nil
}

func loadConfig() (*Config, error) {
	return loadConfig_(getGitOption)
}

func getGitOption(opt string) (string, error) {
	// Get single option from the stored git config
	stdout, err := exec.Command("git", "config", configPrefix+"."+opt).Output()
	if err != nil {
		return "", errors.New("This repo isn't git env enabled. Run 'git env init' first.")
	}
	return string(stdout)[:len(stdout)-1], nil
}

func setGitOption(opt string, value string) (string, error) {
	// Set a key=value option in the git config
	err := exec.Command("git", "config", "--local", "--replace-all", opt, value).Run()
	if err != nil {
		return "", err
	}
	return "Option successfully set", nil
}

func setGitEnvOption(opt string, value string) (string, error) {
	out, err := setGitOption(configPrefix+"."+opt, value)
	if err != nil {
		return "", err
	}
	return out, nil
}
