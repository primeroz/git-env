package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var featureBranchName string
	var envBranchName string
	var helpFlag bool
	var dryFlag bool

	// Setup SubCommands
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initCmd.BoolVar(&helpFlag, "help", false, "show help")

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	//startBranchName := startCmd.String("branch", "", "Start a new feature branch")
	startCmd.StringVar(&featureBranchName, "branch", "", "Feature Branch Name")
	startCmd.StringVar(&featureBranchName, "b", "", "Feature Branch Name")
	startCmd.BoolVar(&helpFlag, "help", false, "show help")
	startCmd.BoolVar(&dryFlag, "dry", false, "dry-run - only print commands to stdout without running them")

	deployCmd := flag.NewFlagSet("deploy", flag.ExitOnError)
	deployCmd.StringVar(&featureBranchName, "branch", "", "Feature Branch Name")
	deployCmd.StringVar(&featureBranchName, "b", "", "Feature Branch Name")
	deployCmd.StringVar(&envBranchName, "env", "", "Env Branch Name")
	deployCmd.StringVar(&envBranchName, "e", "", "Env Branch Name")
	deployCmd.BoolVar(&helpFlag, "help", false, "show help")
	deployCmd.BoolVar(&dryFlag, "dry", false, "dry-run - only print commands to stdout without running them")

	if len(os.Args) < 2 {
		fmt.Println("Commands:")
		fmt.Println("  git env init                               - configure which ENV branches are being used")
		fmt.Println("  git env start BRANCH_NAME                  - start a new feature branch")
		fmt.Println("  git env deploy ENV_BRANCH [FEATURE_BRANCH] - deploy a feature branch to an ENV branch (FEATURE_BRANCH defaults to current branch)")
		os.Exit(1)
	}

	var err error
	config, err = loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])

		if helpFlag {
			fmt.Println("configure the local git configs with the env branches settings")
			initCmd.PrintDefaults()
			os.Exit(0)
		} else {
			cmdInit()
			return
		}
	case "start":
		startCmd.Parse(os.Args[2:])

		if helpFlag {
			fmt.Println("start a new feature branch")
			startCmd.PrintDefaults()
			os.Exit(0)
		} else if featureBranchName == "" {
			fmt.Println("Invalid Flags")
			startCmd.PrintDefaults()
			os.Exit(1)
		} else {
			cmdStart(featureBranchName, dryFlag)
			return
		}
	case "deploy":
		deployCmd.Parse(os.Args[2:])

		if helpFlag {
			fmt.Println("deploy a feature branch to an ENV branch")
			deployCmd.PrintDefaults()
			os.Exit(0)
		} else if envBranchName == "" {
			fmt.Println("Invalid Flags")
			deployCmd.PrintDefaults()
			os.Exit(1)
		} else {
			cmdDeploy(envBranchName, featureBranchName, dryFlag)
			return
		}
	}
}
