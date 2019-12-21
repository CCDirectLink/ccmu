package internal

import (
	"fmt"
	"os"
)

//Outdated displays old mods and their new version
func Outdated() {
	pkgs, err := getGame().Installed()
	if err != nil {
		fmt.Printf("Could not list mods because of an error in %s\n", err.Error())
		os.Exit(1)
	}

	outdated := false
	for _, pkg := range pkgs {
		info, _ := pkg.Info()
		if out, _ := info.Outdated(); out {
			if !outdated {
				outdated = true
				fmt.Println("New     Current Name")
			}

			fmt.Printf("%s   %s   %s\n", info.NewestVersion, info.CurrentVersion, info.NiceName)
		}
	}
}
