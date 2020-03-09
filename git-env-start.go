package main

import (
	"fmt"
	"os"
	"regexp"
)

func cmdStart(newBranch string, dryRun bool) {
	re_branch := regexp.MustCompile("^(f|h|feature|hotfix)/[0-9]+.*$")

	if !re_branch.Match([]byte(newBranch)) {
		fmt.Printf("Branch %s does not match required regexp '^(f|h|feature|hotfix)/[0-9]+.*$'", newBranch)
		os.Exit(1)
	}

	gitCommand(dryRun, "checkout", config.ProdBranch)
	gitCommand(dryRun, "pull")
	gitCommand(dryRun, "checkout", "-b", newBranch)
}
