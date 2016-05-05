package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

const scriptCreateRepo = `#!/bin/sh
repos=$1
for repo in $repos;
do
	mkdir $repo
	(cd $repo && git init)
	for file in ` + "`seq 1 100`" + `;
	do
		file=$file.bin
		(cd $repo && dd if=/dev/urandom of=$file count=1 && git add $file && git commit -m 'Add $file' $file)
	done
	(cd $repo && git checkout -b release)
	for file in ` + "`seq 1 100`" + `;
	do
		file=$file.release
		(cd $repo && dd if=/dev/urandom of=$file count=1 && git add $file && git commit -m 'Add $file' $file)
	done
done
`

func runScript(t *testing.T, script string, runAtDirectory string, arg ...string) {
	scriptFile := path.Join(runAtDirectory, "createRepos.sh")
	if err := ioutil.WriteFile(scriptFile, []byte(script), 0755); err != nil {
		t.Error(err)
	}

	cmd := exec.Command(scriptFile, arg...)
	cmd.Dir = runAtDirectory
	err := cmd.Start()
	if err != nil {
		t.Error()
	}

	err = cmd.Wait()
	if err != nil {
		t.Error()
	}
}

func createManifestFile(t *testing.T, manifestFile string, reposDirectory string, reposName []string) {
	manifestContent := "manifest:\n"
	for _, repoName := range reposName {
		repoPath := path.Join(reposDirectory, repoName)
		manifestContent += fmt.Sprintf("- repo: file://%s\n", repoPath)
	}

	if err := ioutil.WriteFile(manifestFile, []byte(manifestContent), 0644); err != nil {
		t.Error(err)
	}
}

func TestSync(t *testing.T) {
	tmpReposDirectory, err := ioutil.TempDir("/tmp", "multirepo")

	if err != nil {
		t.Error(err)
	}
	reposDirectory = tmpReposDirectory
	defer os.RemoveAll(reposDirectory)

	reposName := []string{"repoA", "repoB", "repoC"}

	reposOriginDir := path.Join(reposDirectory, "origin")
	if err = os.MkdirAll(reposOriginDir, 0755); err != nil {
		t.Error(err)
	}
	runScript(t, scriptCreateRepo, reposOriginDir, strings.Join(reposName, " "))

	manifestFile = path.Join(reposDirectory, DEFAULT_MANIFEST_FILE)
	createManifestFile(t, manifestFile, reposOriginDir, reposName)

	var noArgs []string
	runSync(nil, noArgs)

	for _, repoName := range reposName {
		repoAPath := path.Join(reposDirectory, repoName)

		if _, err = os.Stat(repoAPath); err != nil {
			t.Errorf("Clone operation failed repository not exist on local, %s", err)
		}
	}

	branchName = "release"
	runSync(nil, noArgs)
}
