package main

import (
	"flag"
	"github.com/maorv/git-module"
	"log"
	"os"
	"sync"
)

var CmdApply = Command{
	Run:   runApply,
	Name:  "apply",
	Usage: "apply: put local changes on top of requested remote branch",
	Flag:  flag.NewFlagSet("apply", flag.ExitOnError),
}

func init() {
	CmdApply.Flag.StringVar(&manifestFile, "manifest", DEFAULT_MANIFEST_FILE, "Path to repositories manifest file")
	cwd, _ := os.Getwd()
	CmdApply.Flag.StringVar(&reposDirectory, "repos-directory", cwd, "Path to directory contains all repositories")
	CmdApply.Flag.StringVar(&branchName, "branch", "", "Branch name")
	CmdApply.Flag.StringVar(&remoteName, "remote", "origin", "Remote name")
}

func runApply(cmd *Command, args []string) {
	manifest := LoadManifest(manifestFile)
	var wg sync.WaitGroup
	for _, repoUrl := range manifest.Entries {
		clonedRepo := repoDirectory(repoUrl.RepoURL)
		if _, err := os.Stat(clonedRepo); os.IsNotExist(err) {
			log.Fatalf("Error local repo not exist location %s\n", repoUrl)
		}
		wg.Add(1)
		fetchRepo(clonedRepo, &wg)
	}

	wg.Wait()

	for _, repoUrl := range manifest.Entries {
		clonedRepo := repoDirectory(repoUrl.RepoURL)

		branch := ""
		if branchName != "" {
			branch = remoteName + "/" + branchName
		}

		err := git.Rebase(clonedRepo, git.RebaseOptions{branch})
		if err != nil {
			log.Fatalf("Error rebase repo not exist location %s\n", repoUrl)
		}
	}
}
