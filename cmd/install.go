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

	if _, err := internal.GetGame(); err != nil {
		fmt.Printf("Could not find game folder. Make sure you executed the command inside the game folder.")
		return
	}

	if _, err := internal.FetchModData(); err != nil {
		fmt.Printf("Could not download mod data because an error occured in %s", err.Error())
		return
	}

	for _, name := range args {
		if _, err := internal.GetLocalMod(name); err == nil {
			fmt.Printf("Could not install '%s' because it was already installed", name)
			continue
		}

		if _, err := internal.GetGlobalMod(name); err != nil {
			fmt.Printf("Could find '%s'", name)
			continue
		}

		if err := internal.Install(name, false); err != nil {
			fmt.Printf("Could not install '%s' because an error occured in %s", name, err.Error())
		}
	}
}
