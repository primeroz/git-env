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

	// Install pre-commit
	err := exec.Command("pre-commit", "install").Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("PreCommit installed")

	fmt.Println("You're ready to go.")
}
