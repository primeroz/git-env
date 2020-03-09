package main

import (
	"log"
)

func cmdDeploy(deployEnv string, featureBranch string, dryRun bool) {
	var err error

	if featureBranch == "" {
		featureBranch, err = getCurrentBranch()
		if err != nil {
			panic(err)
		}
	}

	if !config.isEnv(deployEnv) {
		log.Fatalf("Branch %s is not an env branch. Can't merge a feature into it.", deployEnv)
	}

	if config.isEnv(featureBranch) {
		log.Fatalf("Branch %s is an env branch. Can't merge an env branch into another env branch.", featureBranch)
	}

	// TODO - STOP / Break if any uncommited change
	// TODO - push branch
	pushBranch := featureBranch + "-TO-" + deployEnv
	pushFlags := []string{}

	gitCommand(dryRun, "checkout", config.ProdBranch)
	gitCommand(dryRun, "pull")
	gitCommand(dryRun, "checkout", featureBranch)
	gitCommand(dryRun, "merge", "--no-ff", config.ProdBranch)
	gitCommand(dryRun, "push", config.getProdRemote(), featureBranch)

	// Create the push branch
	_, exist := gitRefsExists(pushBranch)
	if exist == nil {
		gitCommand(dryRun, "checkout", pushBranch)
		gitCommand(dryRun, "reset", "--hard", featureBranch)
		pushFlags = append(pushFlags, "--force")
	} else {
		gitCommand(dryRun, "checkout", featureBranch)
		gitCommand(dryRun, "checkout", "-b", pushBranch)
	}

	push_args := []string{"push"}
	push_args = append(push_args, pushFlags...)
	push_args = append(push_args, config.getProdRemote())
	push_args = append(push_args, pushBranch)
	gitCommand(dryRun, push_args...)
	getGitlabMRUrl(dryRun, pushBranch, deployEnv)

	gitCommand(dryRun, "checkout", featureBranch)
}
