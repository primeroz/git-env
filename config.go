package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Option struct {
	Name     string
	Question string
	Default  string
}

type Config struct {
	Version               string
	Mode                  string
	RenderedManifestsRepo string
	RenderCmd             string
	// DeployHook            string
	ProdBranch    string
	OtherBranches []string
}

var (
	version  = "0.1"
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
		/* {
			Name:     "deploy-hook",
			Question: "Run hook command , relative to root of repo, on deploy",
			Default:  "exit 0",
			// Default:  "",
			// Default: "make && git add -A && git diff-index --quiet HEAD || git commit -m \"deploy hook commit\"",
		}, */
		{
			Name:     "rendered-repo",
			Question: "Repo where to store the rendered files",
			Default:  "",
		},
		{
			Name:     "render-cmd",
			Question: "Define the render command to run in the root of the repo",
			// Default:  "make OUTPUTDIR={{.OUTPUTDIR}}",
			Default: "",
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
	}
)

const configPrefix = "env-branch"

func loadConfig_(getOption func(string) (string, error)) (*Config, error) {
	// Load Configuration values from git config
	config := Config{}

	cfg := map[string]string{}

	// Check git env version against saved settings
	v, err := getOption("version")
	if err != nil {
		return nil, err
	}
	if v != version {
		return nil, errors.New(fmt.Sprintf("Configuration version does not match current version %s, please rerun 'git env init'", version))
	}
	config.Version = version

	for _, opt := range options {
		s, err := getOption(opt.Name)
		if err != nil {
			return nil, err
		}
		cfg[opt.Name] = s
	}
	config.Mode = "merge"
	// config.RenderedManifestsRepo = "git@gitlab.com:fciocchetti/git-env-rendered-kustomize.git"
	config.RenderedManifestsRepo = cfg["rendered-repo"]
	config.RenderCmd = cfg["render-cmd"]
	// config.DeployHook = cfg["deploy-hook"]
	config.ProdBranch = cfg["prod"]
	config.OtherBranches = strings.Split(cfg["other"], " ")

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
