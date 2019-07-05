package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/install"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
)

//Install a mod
func Install(args []string) {
	if len(args) == 0 {
		fmt.Println("No mods installed since no mods were specified")
		return
	}

	if _, err := local.GetGame(); err != nil {
		fmt.Printf("Could not find game folder. Make sure you executed the command inside the game folder.")
		return
	}

	if _, err := global.FetchModData(); err != nil {
		fmt.Printf("Could not download mod data because an error occured in %s", err.Error())
		return
	}

	for _, name := range args {
		if _, err := local.GetMod(name); err == nil {
			fmt.Printf("Could not install '%s' because it was already installed", name)
			continue
		}

		installMod(name)
	}
}

func installMod(name string) {
	if _, err := global.GetMod(name); err != nil {
		fmt.Printf("Could find '%s'", name)
		return
	}

	if err := install.Install(name, false); err != nil {
		fmt.Printf("Could not install '%s' because an error occured in %s", name, err.Error())
		return
	}

	mod, err := local.GetMod(name)
	if err != nil {
		fmt.Printf("Installed '%s' but it seems to be an invalid mod", name)
		return
	}

	for name := range mod.Dependencies {
		installMod(name)
	}
}
