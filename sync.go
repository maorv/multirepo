package main

import (
	"flag"
	"github.com/maorv/git-module"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var CmdSync = Command{
	Run:   runSync,
	Name:  "sync",
	Usage: "sync: sync repositories dpecified to latest version of requested branch",
	Flag:  flag.NewFlagSet("sync", flag.ExitOnError),
}

func init() {
	CmdSync.Flag.StringVar(&manifestFile, "manifest", DEFAULT_MANIFEST_FILE, "Path to repositories manifest file")
	cwd, _ := os.Getwd()
	CmdSync.Flag.StringVar(&reposDirectory, "repos-directory", cwd, "Path to directory contains all repositories")
	CmdSync.Flag.StringVar(&branchName, "branch", "origin/master", "Branch to checkout")
}

func updateRepo(repoPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Starting fetch repository %s\n", repoPath)
	err := git.Fetch(repoPath, git.FetchRemoteOptions{false, -1})
	if err != nil {
		log.Fatal("Failed to open repo %s reason: %s\n", repoPath, err)
	}
	log.Printf("Finished fetch repository %s\n", repoPath)
	git.Rebase(repoPath)
}

func cloneRepo(repoPath string, repoUrl string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Starting cloning %s\n", repoUrl)
	err := git.Clone(repoUrl, repoPath, git.CloneRepoOptions{false, false, false, -1})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Finished cloning %s\n", repoUrl)
}

func runSync(cmd *Command, args []string) {
	manifest := LoadManifest(manifestFile)
	var wg sync.WaitGroup
	reposPath := make([]string, 0, len(manifest.Entries))
	for _, repoUrl := range manifest.Entries {
		fileName := path.Base(repoUrl.RepoURL)
		fileName = strings.TrimSuffix(fileName, ".git")
		clonedPath := filepath.Join(reposDirectory, fileName)
		if _, err := os.Stat(clonedPath); os.IsNotExist(err) {
			wg.Add(1)
			go cloneRepo(clonedPath, repoUrl.RepoURL, &wg)
		} else {
			wg.Add(1)
			go updateRepo(clonedPath, &wg)
		}
		reposPath = append(reposPath, clonedPath)
	}
	wg.Wait()

	for _, repoPath := range reposPath {
		err := git.Checkout(repoPath, branchName)
		if err != nil {
			log.Fatalf("Failed to checkout branch %s in repo %s reason: %s\n", branchName, repoPath, err)
		}
	}

	log.Println("Finished sync manifest")
}
