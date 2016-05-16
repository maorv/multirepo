package main

import (
	"github.com/maorv/git-module"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func fetchRepo(repoPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Starting fetch repository %s\n", repoPath)
	err := git.Fetch(repoPath, git.FetchRemoteOptions{false, -1})
	if err != nil {
		log.Fatal("Failed to open repo %s reason: %s\n", repoPath, err)
	}
	log.Printf("Finished fetch repository %s\n", repoPath)
}

func repoDirectory(repoUrl string) string {
	fileName := path.Base(repoUrl)
	fileName = strings.TrimSuffix(fileName, ".git")
	clonedPath := filepath.Join(reposDirectory, fileName)
	return clonedPath
}
