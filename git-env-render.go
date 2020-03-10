package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"
)

func cmdRender(dryRun bool) {
	branch, err := getCurrentBranch()
	if err != nil {
		log.Fatal(err)
	}

	tagName := fmt.Sprintf("render/%s-%d", branch, time.Now().Unix())

	gitCommand(dryRun, "pull")
	commit, err := getGitRevparseBranch(branch)
	if err != nil {
		log.Fatalf("Failed to fetch Current HEAD commit id for branch %s", branch)
	}

	if config.isEnv(branch) {
		// IF branch is one of the ENVs then render from current upstream

		commitUpstream, err := getGitRevparseBranch(config.getProdRemote() + "/" + branch)
		if err != nil {
			log.Fatalf("Failed to fetch upstream commit id for branch %s", branch)
		}

		if commit != commitUpstream {
			log.Fatalf("Your branch of %s is not in sync with upstream", branch)
		}

		gitCommand(dryRun, "tag", "--sign", "--annotate", tagName, "--message", strconv.Quote(fmt.Sprintf("Rendering Manifests for branch %s @ commit %s", branch, commit)))
	} else {
		// IF branch is a development branch then render from current HEAD

		gitCommand(dryRun, "tag", "--annotate", tagName, "--message", strconv.Quote(fmt.Sprintf("Rendering Manifests for branch %s @ commit %s", branch, commit)))
	}

	gitCommand(dryRun, "push", config.getProdRemote(), tagName)

	// Create Temp Dir for rendered manifests repo
	dir, err := ioutil.TempDir("/tmp", "manifests-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Clone, shallow, the manifests repo into tempdir
	gitRepoClone(dryRun, true, config.RenderedManifestsRepo, dir)

	// Render and Run RenderCmd
	s := bytes.NewBufferString("")
	err = template.Must(template.New("").Parse(config.RenderCmd)).Execute(s, map[string]string{"OUTPUTDIR": dir})
	if err != nil {
		panic(err)
	}

	// cd into root dir of repo and run command
	rootDir, err := getGitRepoRootDir()
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(rootDir)
	runCommand(dryRun, "sh", "-c", s.String())

	// getGitlabMRUrl(dryRun, pushBranch, pushEnv)
}
