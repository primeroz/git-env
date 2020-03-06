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

	if config.Mode == "local-rebase" {

		// Rebase feature and env against upstream
		gitCommand(dryRun, "checkout", featureBranch)
		gitCommand(dryRun, "pull", "--rebase", config.getProdRemote(), config.ProdBranch)
		gitCommand(dryRun, "checkout", deployEnv)
		gitCommand(dryRun, "pull", "--rebase", config.getProdRemote(), deployEnv)

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
			// In a non-production merge rebase against the remote env branch and merge
			gitCommand(dryRun, "merge", featureBranch)
		}
	} else if config.Mode == "push-rebase" || config.Mode == "merge" {
		pushBranch := featureBranch + "-" + deployEnv
		pushFlags := []string{}

		gitCommand(dryRun, "fetch")
		gitCommand(dryRun, "checkout", deployEnv)
		gitCommand(dryRun, "pull")

		_, exist := gitRefsExists(pushBranch)
		if exist == nil {
			gitCommand(dryRun, "checkout", pushBranch)
			gitCommand(dryRun, "reset", "--hard", featureBranch)
			pushFlags = append(pushFlags, "--force")
		} else {
			gitCommand(dryRun, "checkout", featureBranch)
			gitCommand(dryRun, "checkout", "-b", pushBranch)
		}

		if config.Mode == "push-rebase" {
			gitCommand(dryRun, "pull", "--rebase", config.getProdRemote(), deployEnv)
		} else if config.Mode == "merge" {
			// The --no-ff flag causes the merge to always create a new commit object, even if the merge could be performed with a fast-forward. This avoids losing information about the historical existence of a feature branch
			gitCommand(dryRun, "merge", "--no-ff", deployEnv)
		}

		if config.DeployHook != "" {
			s := bytes.NewBufferString("")
			err := template.Must(template.New("").Parse(config.DeployHook)).Execute(s, map[string]string{"env": deployEnv, "feature": featureBranch})
			if err != nil {
				panic(err)
			}

			runCommand(dryRun, "sh", "-c", s.String())
		}

		// Build the push command with dynamic flags
		push_args := []string{"push"}
		push_args = append(push_args, pushFlags...)
		push_args = append(push_args, config.getProdRemote())
		push_args = append(push_args, pushBranch)
		gitCommand(dryRun, push_args...)
		gitCommand(dryRun, "checkout", featureBranch)
		getGitlabMRUrl(dryRun, pushBranch, deployEnv)
	}
}
