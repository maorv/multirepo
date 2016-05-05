package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const DEFAULT_MANIFEST_FILE = "multirepo.manifest"

type ManifestEntry struct {
	RepoURL string `yaml:"repo"`
}

type ManifestFile struct {
	Entries []ManifestEntry `yaml:"manifest"`
}

func LoadManifest(manifestPath string) *ManifestFile {
	content, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		log.Fatal(err)
	}

	var manifest ManifestFile
	yaml.Unmarshal(content, &manifest)
	return &manifest
}
