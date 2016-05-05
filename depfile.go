package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type DepFileEntry struct {
	RepoURL string `yaml:"repo"`
	Version string `yaml:"ver"`
}

type DepFile struct {
	Entries []DepFileEntry `yaml:"deps"`
}

func LoadDepFile(fileName string) *DepFile {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var depFile DepFile
	yaml.Unmarshal(content, &depFile)
	return &depFile
}

func (dep *DepFile) SaveDepFile(fileName string) {
	content, err := yaml.Marshal(dep)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fileName, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
