package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/CCDirectLink/ccmu/cmd/internal"
)

func main() {
	var err error

	flag.Usage = printHelp
	flag.String("game", "", "if set it overrides the path of the game")
	uri := flag.String("url", "", "the url that executed ccmu")
	flag.Parse()

	if len(os.Args) == 1 {
		printHelp()
		return
	}

	op := flag.Arg(0)
	args := flag.Args()[1:]

	if uri != nil && *uri != "" {
		raw := *uri
		args = strings.Split(raw[7:len(raw)-1], "/")
		for i, arg := range args {
			args[i], err = url.PathUnescape(arg)
			if err != nil {
				fmt.Printf("An error occured in %s\n", err)
				os.Exit(1)
			}
		}
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
