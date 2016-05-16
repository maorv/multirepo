package main

import (
	"flag"
	"github.com/maorv/git-module"
	"log"
	"os"
	"path"
	"path/filepath"
)

var CmdSave = Command{
	Run:   runSave,
	Name:  "save",
	Usage: "save: save the current state of all files into specified dep file",
	Flag:  flag.NewFlagSet("save", flag.ExitOnError),
}

var (
	fileName string
)

func init() {
	CmdSave.Flag.StringVar(&fileName, "output", "", "Path to new dependecy file")
	CmdSave.Flag.StringVar(&manifestFile, "manifest", DEFAULT_MANIFEST_FILE, "Path to repositories manifest file")
	cwd, _ := os.Getwd()
	CmdSave.Flag.StringVar(&reposDirectory, "repos-directory", cwd, "Path to directory contains all repositories")
}

func runSave(cmd *Command, args []string) {
	var depFile DepFile
	manifest := LoadManifest(manifestFile)
	for _, repoUrl := range manifest.Entries {
		repoDir := path.Base(repoUrl.RepoURL)
		clonedPath := filepath.Join(reposDirectory, repoDir)
		repo, err := git.OpenRepository(clonedPath)

		if err != nil {
			log.Fatal("Error open repo %s %s\n", clonedPath, err)
		}

		branch, err := repo.GetHEADBranch()

		if err != nil {
			log.Fatal("Error get repo HEAD branch ", clonedPath, " ", err)
		}

		hash, err := repo.GetBranchCommitID(branch.Name)

		if err != nil {
			log.Fatal("Error get repo commit id ", clonedPath, " ", err)
		}

		depFile.Entries = append(depFile.Entries, DepFileEntry{RepoURL: repoUrl.RepoURL, Version: hash})
	}
	depFile.SaveDepFile(fileName)
}
