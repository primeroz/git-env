package main

import (
	//"log"
	"fmt"
	"strconv"
	"time"
)

func cmdRender(branch string, manifestsRepo string, dryRun bool) {

	if config.isEnv(branch) {
		gitCommand(dryRun, "checkout", branch)
		gitCommand(dryRun, "pull")

		commit, err := getGitRevparseBranch(branch)
		if err != nil {
			panic(err)
		}

		gitCommand(dryRun, "tag", "--sign", "--annotate", fmt.Sprintf("render/%s-%d", branch, time.Now().Unix()), "--message", strconv.Quote(fmt.Sprintf("Rendering Manifests for branch %s @ commit %s", branch, commit)))
		gitCommand(dryRun, "push", "--tags")
	} else {
		fmt.Printf("Placeholder for ELSE")
	}

	// if branch is one of the ENV then pull the latest from remote
	// else use the local HEAD for the branch but WARN ( err ? ) if not pushed
	// Create a TAG out of it so we have an immutable placeholer for the rendering
	// push the tag ( can i push if i did not push head yet ? )
	// Clone Shallow Rendered repo in temp
	// Render ( only the required env ? ) into the temp dir for the cloned repo
	// should use an hook for this rather than hardcode
	// push and create MR with description including TAG reference

	// var err error

	// if featureBranch == "" {
	// 	featureBranch, err = getCurrentBranch()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// if !config.isEnv(pushEnv) {
	// 	log.Fatalf("Branch %s is not an env branch. Can't merge a feature into it.", pushEnv)
	// }

	// if config.isEnv(featureBranch) {
	// 	log.Fatalf("Branch %s is an env branch. Can't merge an env branch into another env branch.", featureBranch)
	// }

	// // TODO - STOP / Break if any uncommited change
	// // TODO - push branch
	// pushBranch := featureBranch + "-TO-" + pushEnv
	// pushFlags := []string{}

	// gitCommand(dryRun, "checkout", config.ProdBranch)
	// gitCommand(dryRun, "pull")
	// gitCommand(dryRun, "checkout", featureBranch)
	// gitCommand(dryRun, "merge", "--no-ff", config.ProdBranch)
	// gitCommand(dryRun, "push", config.getProdRemote(), featureBranch)

	// // Create the push branch
	// _, exist := gitRefsExists(pushBranch)
	// if exist == nil {
	// 	gitCommand(dryRun, "checkout", pushBranch)
	// 	gitCommand(dryRun, "reset", "--hard", featureBranch)
	// 	pushFlags = append(pushFlags, "--force")
	// } else {
	// 	gitCommand(dryRun, "checkout", featureBranch)
	// 	gitCommand(dryRun, "checkout", "-b", pushBranch)
	// }

	// push_args := []string{"push"}
	// push_args = append(push_args, pushFlags...)
	// push_args = append(push_args, config.getProdRemote())
	// push_args = append(push_args, pushBranch)
	// gitCommand(dryRun, push_args...)
	// getGitlabMRUrl(dryRun, pushBranch, pushEnv)

	// gitCommand(dryRun, "checkout", featureBranch)
}
