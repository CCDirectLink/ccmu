package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
)

//Update a mod
func Update(args []string) {
	if len(args) == 0 {
		fmt.Println("No mods updated since no mods were specified")
		return
	}

	if _, err := internal.GetGame(); err != nil {
		fmt.Printf("Could not find game folder. Make sure you executed the command inside the game folder.")
		return
	}

	_, err := internal.FetchModData()
	if err != nil {
		fmt.Printf("Could not download mod data because an error occured in %s", err.Error())
		return
	}

	for _, name := range args {
		if _, err := internal.GetLocalMod(name); err != nil {
			fmt.Printf("Could not update '%s' because it was not installed", name)
			continue
		}

		if _, err := internal.GetGlobalMod(name); err != nil {
			fmt.Printf("Could find '%s'", name)
			continue
		}

		err := internal.Install(name, true)
		if err != nil {
			fmt.Printf("Could not update '%s' because an error occured in %s", name, err.Error())
		}
	}
}
