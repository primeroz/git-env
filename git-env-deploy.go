package main

import (
	"bytes"
	"log"
	"text/template"
)

func cmdDeploy(deployEnv string, featureBranch string) {
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

	// Rebase feature and env against upstream
	gitCommand("checkout", featureBranch)
	gitCommand("pull", "--rebase", config.getProdRemote(), config.ProdBranch)
	gitCommand("checkout", deployEnv)
	gitCommand("pull", "--rebase", config.getProdRemote(), deployEnv)

	if config.isProd(deployEnv) {
		// In a production merge use --no-ff so the branch names are preserved
		gitCommand("checkout", featureBranch)

		s := bytes.NewBufferString("")
		err := template.Must(template.New("").Parse(config.ProdDeployCmd)).Execute(s, map[string]string{"env": deployEnv, "feature": featureBranch})
		if err != nil {
			panic(err)
		}

		runCommand("sh", "-c", s.String())
	} else {
		// In a non-production merge rebase against the remote env branch and merge
		gitCommand("merge", featureBranch)
	}
}
