package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/CCDirectLink/ccmu/cmd/internal"
)

func main() {

	flag.Usage = printHelp
	flag.String("game", "", "if set it overrides the path of the game")
	url := flag.String("url", "", "the url that executed ccmu")
	flag.Parse()

	if len(os.Args) == 1 {
		printHelp()
		return
	}

	op := flag.Arg(0)
	args := flag.Args()[1:]

	if url != nil && *url != "" {
		raw := *url
		args = strings.Split(raw[7:len(raw)-1], " ")
	}

	switch op {
	case "install",
		"i":
		internal.Install(args)
	case "installed":
		internal.Installed()
	case "remove",
		"delete",
		"uninstall":
		internal.Uninstall(args)
	case "update":
		internal.Update(args)
	case "list":
		internal.List()
	case "outdated":
		internal.Outdated()
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
