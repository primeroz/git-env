package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"text/template"
	"time"
)

func cmdRender(dryRun bool) {
	var commit string
	var tagName string
	var renderedRepoBranchName string

	// Get the current branch , Render always work on current branch
	branch, err := getCurrentBranch()
	if err != nil {
		log.Printf("Failed to get current Branch")
		panic(err)
	}
	// Check if we are actually on a tag detached
	tag, err := gitIsTag()
	if err != nil {
		log.Printf("Failed to check if current HEAD is a tag")
		panic(err)
	}

	// Fetch Refs from remote so we can compare against local branches
	gitCommand(dryRun, "fetch", "--all")

	if config.isEnv(branch) {
		// If branch is one of the ENVs then render ensure branch is in sync with upstream

		// Check Repo status
		hasChanges, err := gitHasAnyChange()
		if err != nil {
			log.Fatal(err)
		}
		if hasChanges {
			log.Fatalf("There are uncommited changes in the repo, cannot proceed. Please push all your changes and retry")
		}

		// Check that current repo HEAD is in sync with upstream
		inSync, err := gitBranchesInSync(branch, config.getProdRemote()+"/"+branch)
		if err != nil {
			log.Fatalf("Failed to compare commit Ids for %s and %s", branch, config.getProdRemote()+"/"+branch)
		}

		if !inSync {
			log.Fatalf("Your branch of %s is not in sync with upstream. Please run 'git env pull' and retry", branch)
		}

		fmt.Printf("+ %s branch is clean and in sync with origin\n", branch)

		// Create a Signed tag and push it upstream
		commit, _ = getGitBranchCommitId(branch, true)
		tagName = fmt.Sprintf("render/%s-%s-%d", branch, commit, time.Now().Unix())
		renderedRepoBranchName = tagName

		gitCommand(dryRun, "tag", "--sign", "--annotate", tagName, "-m", strconv.Quote(fmt.Sprintf("Rendering Manifests for branch %s @ commit %s", branch, commit)))
		gitCommand(dryRun, "push", config.getProdRemote(), tagName)

	} else if tag != "" {
		hasChanges, err := gitHasAnyChange()
		if err != nil {
			log.Fatal(err)
		}
		if hasChanges {
			log.Printf("There are uncommited changes in the repo, cannot proceed. This is a TAG checkout so you should do no changes here")
			//log.Fatalf("There are uncommited changes in the repo, cannot proceed. This is a TAG checkout so you should do no changes here")
		}
		// If we are re-rendering a tag , augment renderedRepoBranchName with epoch
		fmt.Printf("+ rendering from tag %s\n", tag)

		re_renderedTag := regexp.MustCompile("^rendered/")
		if re_renderedTag.Match([]byte(tag)) {
			tagName = tag
		} else {
			tagName = fmt.Sprintf("rendered/%s", tag)
		}

		renderedRepoBranchName = fmt.Sprintf("%s-%d", tagName, time.Now().Unix())
	} else {
		// If branch is a development branch then render from current HEAD
		inSync, _ := gitBranchesInSync(branch, config.getProdRemote()+"/"+branch)
		if !inSync {
			log.Printf("WARNING: Rendering from a repo not in sync with upstream. Be Careful")
		}
		hasChanges, _ := gitHasAnyChange()
		if hasChanges {
			log.Printf("WARNING: Rendering from an unclean repo with Changes. Be Careful")
		}

		commit, _ = getGitBranchCommitId(branch, true)
		tagName = fmt.Sprintf("render/%s-%s-%d", branch, commit, time.Now().Unix())
		renderedRepoBranchName = tagName

		gitCommand(dryRun, "tag", "--annotate", tagName, "-m", strconv.Quote(fmt.Sprintf("Rendering Manifests for branch %s @ commit %s", branch, commit)))
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

	//repoGitUrl, err := getGitRemoteUrl()
	//if err != nil {
	//	log.Fatal("Failed to get URL for repo")
	//}

	// If there are Changes push
	defer os.Chdir(repoRootDir)
	os.Chdir(tmpDir)

	var renderedRepoHasChanges bool
	if !dryRun {
		renderedRepoHasChanges, err = gitHasAnyChange()
		if err != nil {
			fmt.Printf("Error Checking for changes in %s\n", tmpDir)
			log.Fatal(err)
		}
	}
	if dryRun || renderedRepoHasChanges {
		// TODO Should we do an interactive commit here ?
		gitCommand(dryRun, "checkout", "-b", renderedRepoBranchName)
		gitCommand(dryRun, "add", "-A")
		gitCommand(dryRun, "diff", "--cached")
		gitCommand(dryRun, "status")
		gitCommand(dryRun, "push", "origin", renderedRepoBranchName, "-m", strconv.Quote(fmt.Sprintf("Rendered Manifest from repo %s at tag %s", repoGitUrl, tagName)))
		//getGitlabMRUrl(dryRun, pushBranch, pushEnv)
	}
}
