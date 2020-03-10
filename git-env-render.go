package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"text/template"
)

func cmdRender(dryRun bool) {
	var commit string
	var tagName string

	// Get the current branch , Render always work on current branch
	branch, err := getCurrentBranch()
	if err != nil {
		log.Fatal(err)
	}

	if config.isEnv(branch) {
		// If branch is one of the ENVs then render ensure branch is in sync with upstream

		// Check Repo status
		hasChanges, err := gitHasAnyChange()
		if err != nil {
			log.Fatal(err)
		}
		if hasChanges {
			log.Fatalf("There are uncommited changes in the repo, cannot proceed")
		}

		// Check that current repo HEAD is in sync with upstream
		inSync, err := gitBranchesInSync(branch, config.getProdRemote()+"/"+branch)
		if err != nil {
			log.Fatalf("Failed to compare commit Ids for %s and %s", branch, config.getProdRemote()+"/"+branch)
		}

		if !inSync {
			log.Fatalf("Your branch of %s is not in sync with upstream", branch)
		}

		// Create a Signed tag and push it upstream
		commit, _ = getGitBranchCommitId(branch, true)
		tagName = fmt.Sprintf("render/%s-%s", branch, commit)

		gitCommand(dryRun, "tag", "--sign", "--annotate", tagName, "--message", strconv.Quote(fmt.Sprintf("Rendering Manifests for branch %s @ commit %s", branch, commit)))
		gitCommand(dryRun, "push", config.getProdRemote(), tagName)

	} else {
		// If branch is a development branch then render from current HEAD
		inSync, _ := gitBranchesInSync(branch, config.getProdRemote()+"/"+branch)
		if !inSync {
			log.Printf("WARNING: Rendering from an unclean repo with Changes. Be Careful")
		}

		commit, _ = getGitBranchCommitId(branch, true)
		tagName = fmt.Sprintf("render/%s-%s", branch, commit)
		gitCommand(dryRun, "tag", "--annotate", tagName, "--message", strconv.Quote(fmt.Sprintf("Rendering Manifests for branch %s @ commit %s", branch, commit)))
		gitCommand(dryRun, "push", config.getProdRemote(), tagName)
	}

	// Create Temp Dir for rendered manifests repo
	tmpDir, err := ioutil.TempDir("/tmp", "manifests-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone, shallow, the manifests repo into tempdir
	gitRepoClone(dryRun, true, config.RenderedManifestsRepo, tmpDir)

	// Render and Run RenderCmd
	s := bytes.NewBufferString("")
	err = template.Must(template.New("").Parse(config.RenderCmd)).Execute(s, map[string]string{"OUTPUTDIR": tmpDir})
	if err != nil {
		panic(err)
	}

	// cd into root dir of repo and run command
	repoRootDir, err := getGitRepoRootDir()
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(repoRootDir)
	runCommand(dryRun, "sh", "-c", s.String())

	if !dryRun {
		repoGitUrl, err := getGitRemoteUrl()
		if err != nil {
			log.Fatal("Failed to get URL for repo")
		}

		// If there are Changes push
		defer os.Chdir(repoRootDir)
		os.Chdir(tmpDir)

		hasChanges, err := gitHasAnyChange()
		if err != nil {
			fmt.Printf("I am in %s\n", tmpDir)
			log.Fatal(err)
		}
		if hasChanges {
			// TODO Should we do an interactive commit here ?
			gitCommand(dryRun, "checkout", "-b", tagName, "--message", strconv.Quote(fmt.Sprintf("Rendered Manifest from repo %s at tag %s", repoGitUrl, tagName)))

		}
		// getGitlabMRUrl(dryRun, pushBranch, pushEnv)
	}
}
