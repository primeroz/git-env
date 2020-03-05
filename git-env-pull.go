package main

import ()

func cmdPull(dryRun bool) {
	gitCommand(dryRun, "fetch", "--all", "--prune")

	branches := config.OtherBranches
	branches = append(branches, config.ProdBranch)

	for _, v := range branches {
		gitCommand(dryRun, "checkout", v)
		if config.Mode == "merge" {
			gitCommand(dryRun, "pull")
		} else {
			gitCommand(dryRun, "pull", "--rebase", config.getProdRemote(), v)
		}
	}
}
