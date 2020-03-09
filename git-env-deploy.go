package main

import (
	"bytes"
	"log"
	"text/template"
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

	// Push Vars
	pushBranch := featureBranch + "-TO-" + deployEnv
	pushFlags := []string{}

	// Rebase from origin master
	// TODO: Any conflict will need to be resolved here , this will stop the git env deploy command ? how to resume ?
	gitCommand(dryRun, "checkout", featureBranch)
	gitCommand(dryRun, "pull", "--rebase", config.getProdRemote(), config.ProdBranch)

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

	if config.Mode == "push" {
		// All Merging is done through PUSH And MR in gitlab / github

		// Build the push command with dynamic flags
		push_args := []string{"push"}
		push_args = append(push_args, pushFlags...)
		push_args = append(push_args, config.getProdRemote())
		push_args = append(push_args, pushBranch)
		gitCommand(dryRun, push_args...)
		getGitlabMRUrl(dryRun, pushBranch, deployEnv)

	} else if config.Mode == "merge" {
		// TODO: REVIEW THIS
		if config.isProd(deployEnv) {
			// In a production merge use --no-ff so the branch names are preserved
			gitCommand(dryRun, "checkout", featureBranch)

			s := bytes.NewBufferString("")
			err := template.Must(template.New("").Parse(config.ProdDeployCmd)).Execute(s, map[string]string{"env": deployEnv, "feature": featureBranch})
			if err != nil {
				panic(err)
			}

			runCommand(dryRun, "sh", "-c", s.String())
		} else {
			gitCommand(dryRun, "checkout", deployEnv)
			gitCommand(dryRun, "pull", "--rebase", config.getProdRemote(), deployEnv)
			gitCommand(dryRun, "merge", "--no-ff", deployEnv)
		}
	}

	gitCommand(dryRun, "checkout", featureBranch)
}
