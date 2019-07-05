package cmd

import (
	"fmt"
	"os"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
)

//Outdated displays old mods and their new version
func Outdated() {
	if _, err := local.GetGame(); err != nil {
		fmt.Printf("Could not find game folder. Make sure you executed the command inside the game folder.\n")
		return
	}

	mods, err := local.GetMods()
	if err != nil {
		fmt.Printf("Could not list mods because of an error in %s\n", err.Error())
		os.Exit(1)
	}

	outdated := false
	for _, mod := range mods {
		if out, _ := mod.Outdated(); out {
			new, err := global.GetMod(mod.Name)
			if err != nil {
				fmt.Printf("An error occured in %s\n", err.Error())
				continue
			}

			if !outdated {
				outdated = true
				fmt.Println("New     Current Name")
			}

			fmt.Printf("%s   %s   %s\n", new.Version, mod.Version, mod.Name)
		}
	}
}
