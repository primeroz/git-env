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

	pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
	pullCmd.BoolVar(&helpFlag, "help", false, "show help")
	pullCmd.BoolVar(&dryFlag, "dry", false, "dry-run - only print commands to stdout without running them")

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	//startBranchName := startCmd.String("branch", "", "Start a new feature branch")
	startCmd.StringVar(&featureBranchName, "branch", "", "Feature Branch Name")
	startCmd.StringVar(&featureBranchName, "b", "", "Feature Branch Name")
	startCmd.BoolVar(&helpFlag, "help", false, "show help")
	startCmd.BoolVar(&dryFlag, "dry", false, "dry-run - only print commands to stdout without running them")

	pushCmd := flag.NewFlagSet("push", flag.ExitOnError)
	pushCmd.StringVar(&featureBranchName, "branch", "", "Feature Branch Name - defaults to current branch")
	pushCmd.StringVar(&featureBranchName, "b", "", "Feature Branch Name - defaults to current branch")
	pushCmd.StringVar(&envBranchName, "env", "", "Env Branch Name")
	pushCmd.StringVar(&envBranchName, "e", "", "Env Branch Name")
	pushCmd.BoolVar(&helpFlag, "help", false, "show help")
	pushCmd.BoolVar(&dryFlag, "dry", false, "dry-run - only print commands to stdout without running them")

	if len(os.Args) < 2 {
		fmt.Println("Commands:")
		fmt.Println("  git env init                                     - configure which ENV branches are being used")
		fmt.Println("  git env pull                                     - pull all the ENV branches")
		fmt.Println("  git env start -b BRANCH_NAME                     - start a new feature branch ( it must match the regex (f|h|feature|hotfix)/[0-9]+.*")
		fmt.Println("  git env push -e ENV_BRANCH -b [FEATURE_BRANCH] - push a feature branch to an ENV branch (FEATURE_BRANCH defaults to current branch)")
		os.Exit(1)
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
	}

	var err error
	config, err = loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	switch os.Args[1] {
	case "pull":
		pullCmd.Parse(os.Args[2:])

		if helpFlag {
			fmt.Println("pull all the ENV branches")
			pullCmd.PrintDefaults()
			os.Exit(0)
		} else {
			cmdPull(dryFlag)
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
	case "push":
		pushCmd.Parse(os.Args[2:])

		if helpFlag {
			fmt.Println("push a feature branch to an ENV branch")
			pushCmd.PrintDefaults()
			os.Exit(0)
		} else if envBranchName == "" {
			fmt.Println("Invalid Flags")
			pushCmd.PrintDefaults()
			os.Exit(1)
		} else {
			cmdPush(envBranchName, featureBranchName, dryFlag)
			return
		}
	}
}
