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
		fmt.Printf("Could not find game folder. Make sure you executed the command inside the game folder.")
		return
	}

	_, err := global.FetchModData()
	if err != nil {
		fmt.Printf("Could not download mod data because an error occured in %s", err.Error())
		return
	}

	for _, name := range args {
		if _, err := local.GetMod(name); err != nil {
			fmt.Printf("Could not update '%s' because it was not installed", name)
			continue
		}

		if _, err := global.GetMod(name); err != nil {
			fmt.Printf("Could find '%s'", name)
			continue
		}

		err := install.Install(name, true)
		if err != nil {
			fmt.Printf("Could not update '%s' because an error occured in %s", name, err.Error())
		}
	}
}
