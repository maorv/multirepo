package main

import "flag"

type Command struct {
	Run   func(cmd *Command, repos []string)
	Name  string
	Usage string
	Flag  *flag.FlagSet
}

var (
	manifestFile   string
	reposDirectory string
)
