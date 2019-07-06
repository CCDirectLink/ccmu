package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/install"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
)

//Update a mod
func Update(args []string) {
	if len(args) == 0 {
		fmt.Println("No mods updated since no mods were specified")
		return
	}

	if _, err := local.GetGame(); err != nil {
		fmt.Printf("Could not find game folder. Make sure you executed the command inside the game folder.\n")
		return
	}

	_, err := global.FetchModData()
	if err != nil {
		fmt.Printf("Could not download mod data because an error occured in %s\n", err.Error())
		return
	}

	count := 0
	for _, name := range args {
		if _, err := local.GetMod(name); err != nil {
			fmt.Printf("Could not update '%s' because it was not installed\n", name)
			continue
		}

		if _, err := global.GetMod(name); err != nil {
			fmt.Printf("Could find '%s'\n", name)
			continue
		}

		err := install.Install(name, true)
		if err != nil {
			fmt.Printf("Could not update '%s' because an error occured in %s\n", name, err.Error())
		}

		count++
	}

	fmt.Printf("Updated %d mods\n", count)
}
