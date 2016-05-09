package main

import (
	"flag"
	"github.com/maorv/git-module"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var CmdPush = Command{
	Run:   runPush,
	Name:  "push",
	Usage: "push: push repositories specified to requested remote branch",
	Flag:  flag.NewFlagSet("push", flag.ExitOnError),
}

var (
	remoteName string
	force      bool
)

func init() {
	CmdPush.Flag.StringVar(&manifestFile, "manifest", DEFAULT_MANIFEST_FILE, "Path to repositories manifest file")
	cwd, _ := os.Getwd()
	CmdPush.Flag.StringVar(&reposDirectory, "repos-directory", cwd, "Path to directory contains all repositories")
	CmdPush.Flag.StringVar(&branchName, "branch", "", "Branch name")
	CmdPush.Flag.StringVar(&remoteName, "remote", "origin", "Remote name")
	CmdPush.Flag.BoolVar(&force, "force", false, "Force push to remote branch")
}

func pushAll(manifest *ManifestFile) {
	var wg sync.WaitGroup
	pushOpt := git.PushOptions{Force: force}
	for _, repoUrl := range manifest.Entries {
		repoDir := path.Base(repoUrl.RepoURL)
		clonedPath := filepath.Join(reposDirectory, repoDir)
		headToPush := "HEAD:" + branchName
		wg.Add(1)
		go func() {
			err := git.Push(clonedPath, remoteName, headToPush, pushOpt)

			if err != nil {
				log.Fatalf("Error push repo %s %s\n", clonedPath, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func checkoutAll(manifest *ManifestFile) {
	for _, repoUrl := range manifest.Entries {
		repoDir := path.Base(repoUrl.RepoURL)
		clonedPath := filepath.Join(reposDirectory, repoDir)
		err := git.Checkout(clonedPath, branchName)

		if err != nil {
			log.Fatalf("Checkout branch %s in %s failed err: %s\n", branchName, clonedPath, err)
		}
	}
}

func runPush(cmd *Command, args []string) {
	manifest := LoadManifest(manifestFile)
	pushAll(manifest)
	checkoutAll(manifest)
}
