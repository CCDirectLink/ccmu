package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd"
)

func main() {
	flag.String("game", "", "if set it overrides the path of the game")
	flag.Parse()

	if len(os.Args) == 1 {
		printHelp()
		return
	}

	op := flag.Arg(0)
	args := flag.Args()[1:]

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
	case "version":
		printVersion()
	case "help":
		printHelp()
	default:
		fmt.Printf("%s\n is not a command", op)
		printHelp()
		os.Exit(1)
	}
}
