package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var commands = []Command{
	CmdSync,
	CmdSave,
	CmdPush,
	CmdApply,
}

func PrintDefaults() {
	for _, cmd := range commands {
		fmt.Fprintln(os.Stderr, cmd.Usage)
		cmd.Flag.PrintDefaults()
	}
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("multirepo: ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: of %s:\n", os.Args[0])
		PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	if args[0] == "help" {
		flag.Usage()
		os.Exit(0)
	}

	for _, cmd := range commands {
		if args[0] == cmd.Name {
			cmd.Flag.Parse(args[1:])
			cmd.Run(&cmd, cmd.Flag.Args())
		}
	}
}
