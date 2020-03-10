package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func cmdInit() {
	values := map[string]string{}
	gitSettings := map[string]string{}
	reader := bufio.NewReader(os.Stdin)

	// Manage git env options
	for _, opt := range options {
		fmt.Printf("%s [%s] ", opt.Question, opt.Default)
		value, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if value == "\n" {
			values[opt.Name] = opt.Default
		} else {
			values[opt.Name] = value[:len(value)-1]
		}
	}

	for k, v := range values {
		_, err := setGitEnvOption(k, v)
		if err != nil {
			panic(err)
		}
	}

	// Manage git generic options
	for _, s := range settings {
		fmt.Printf("%s [%s] ", s.Question, s.Default)
		value, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if value == "\n" {
			gitSettings[s.Name] = s.Default
		} else {
			gitSettings[s.Name] = value[:len(value)-1]
		}
	}

	for k, v := range gitSettings {
		_, err := setGitOption(k, v)
		if err != nil {
			panic(err)
		}
	}

	// Set Version
	_, err := setGitEnvOption("version", version)
	if err != nil {
		panic(err)
	}

	// Install pre-commit
	err = exec.Command("pre-commit", "install").Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("PreCommit installed")

	// Install other hooks
	commit_msg_hook := `#!/usr/bin/env bash

# set this to your active development branch
current_branch="$(git rev-parse --abbrev-ref HEAD)"
check_branch_regex='^(f/[0-9]+.*)'

# Only check branches matching regex
[[ ! ${current_branch} =~ ${check_branch_regex} ]] && exit 0

# regex to validate in commit msg
commit_regex='(f/[0-9]+|merge)'
error_msg="Aborting commit. Your commit message is missing either a Feature Issue ('f/1111') or 'Merge'"

if ! grep -iqE "$commit_regex" "$1"; then
    echo "$error_msg" >&2
    exit 1
fi
`
	f, err := os.Create(".git/hooks/commit-msg")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	_, err = f.WriteString(commit_msg_hook)
	if err != nil {
		panic(err)
	}
	f.Sync()
	f.Chmod(0750)

	fmt.Println("Hooks Created")

	fmt.Println("You're ready to go.")
}
