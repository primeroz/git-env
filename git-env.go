package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		help("")
		return
	}

	// Do help and init commands before loading config

	switch os.Args[1] {
	case "init":
		cmdInit()
		return
	case "help":
		if len(os.Args) > 2 {
			help(os.Args[2])
		} else {
			help("")
		}
		return
	}

	var err error
	config, err = loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Do start and deploy commands only after loading config

	switch os.Args[1] {
	case "start":
		cmdStart(os.Args[2:])
	case "deploy":
		cmdDeploy(os.Args[2:])
	default:
		help("")
	}
}

func help(arg string) {
	switch arg {
	//TODO: add in-depth documentation for each command
	// case "init":
	// case "start":
	// case "deploy":
	default:
		fmt.Println("Commands:")
		fmt.Println("  git env help                               - show this help")
		fmt.Println("  git env init                               - configure which ENV branches are being used")
		fmt.Println("  git env start BRANCH_NAME                  - start a new feature branch")
		fmt.Println("  git env deploy ENV_BRANCH [FEATURE_BRANCH] - deploy a feature branch to an ENV branch (FEATURE_BRANCH defaults to current branch)")
	}
}
