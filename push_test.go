package main

import (
	"os"
	"path"
	"testing"
)

func TestPush(t *testing.T) {
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

	branchName = "master"
	remoteName = "origin"
	runSync(nil, noArgs)
	branchName = "devel"
	runPush(nil, noArgs)
}
