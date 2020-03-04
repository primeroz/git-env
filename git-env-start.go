package main

import ()

func cmdStart(newBranch string) {
	//	if len(args) < 1 {
	//		help("start")
	//	}

	gitCommand("checkout", config.ProdBranch)
	gitCommand("pull", "--rebase", config.getProdRemote(), config.ProdBranch)
	gitCommand("checkout", "-b", newBranch)
}
