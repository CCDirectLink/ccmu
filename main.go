package main

import (
	"fmt"
	"os"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd"
)

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	op := os.Args[1]
	args := os.Args[2:]

	switch op {
	case "install",
		"i":
		cmd.Install(args)
	case "remove",
		"delete",
		"uninstall":
		cmd.Uninstall(args)
	case "update":
		cmd.Update(args)
	case "list":
		cmd.List()
	case "outdated":
		cmd.Outdated()
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown operation %s\n", op)
		printHelp()
		os.Exit(1)
	}
}
