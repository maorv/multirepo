package main

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	t.Log("Running  Apply test")
	reposName := []string{"repoA", "repoB", "repoC"}
	setupRepos(t, reposName)
	defer os.RemoveAll(reposDirectory)

	var noArgs []string
	runSync(nil, noArgs)

	for _, repoName := range reposName {
		repoAPath := path.Join(reposDirectory, repoName)

		if _, err := os.Stat(repoAPath); err != nil {
			t.Errorf("Clone operation failed repository not exist on local, %s", err)
		}
	}

	branchName = ""
	remoteName = "origin"
	runSync(nil, noArgs)

	branchName = "mytest"
	reposOriginDir := path.Join(reposDirectory, "origin")
	runScript(t, scriptCreateBranchWithCommits, reposOriginDir, strings.Join(reposName, " "), branchName)
	runApply(nil, noArgs)
	for _, repoName := range reposName {
		repoPath := path.Join(reposDirectory, repoName)
		filePath := path.Join(repoPath, "100."+branchName)
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("Expected file %s in directory %s\n", filePath, repoPath)
		}

	}
}
