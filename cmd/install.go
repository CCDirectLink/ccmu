package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
)

//Install a mod
func Install(args []string) {
	if len(args) == 0 {
		fmt.Println("No mods installed since no mods were specified")
		return
	}

	_, err := internal.FetchModData()
	if err != nil {
		fmt.Printf("Could not download mod data because an error occured in %s", err.Error())
		return
	}

	for _, name := range args {
		err := internal.Install(name, false)
		if err != nil {
			fmt.Printf("Could not install '%s' because an error occured in %s", name, err.Error())
		}
	}
}
