package main

import ()

func cmdStart(args []string) {
	if len(args) < 1 {
		help("start")
	}

	newBranch := args[0]

	gitCommand("checkout", config.ProdBranch)
	gitCommand("pull", "--rebase", config.getProdRemote(), config.ProdBranch)
	gitCommand("checkout", "-b", newBranch)
}
