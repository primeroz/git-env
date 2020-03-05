package main

import ()

func cmdStart(newBranch string, dryRun bool) {
	gitCommand(dryRun, "checkout", config.ProdBranch)
	if config.Mode == "merge" {
		gitCommand(dryRun, "pull")
	} else {
		gitCommand(dryRun, "pull", "--rebase", config.getProdRemote(), config.ProdBranch)
	}
	gitCommand(dryRun, "checkout", "-b", newBranch)
}
