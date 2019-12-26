package internal

import (
	"fmt"
	"os"
)

//Installed prints a list of all intalled mods
func Installed() {
	pkgs, err := getGame().Installed()
	if err != nil {
		fmt.Printf("Could not list mods because of an error in %s\n", err.Error())
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		info, err := pkg.Info()
		if err == nil {
			version := info.CurrentVersion
			if version == "" {
				version = "x.x.x"
			}
			fmt.Printf("%s %s\n", version, info.NiceName)
		} else {
			fmt.Println(err)
		}
	}
}
