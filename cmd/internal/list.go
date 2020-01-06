package internal

import (
	"fmt"
	"os"
)

//List prints a list of all available mods
func List() {
	pkgs, err := getGame().Available()
	if err != nil {
		fmt.Printf("Could not list mods because of an error in %s\n", err.Error())
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		info, err := pkg.Info()
		if err == nil {
			if info.Hidden {
				fmt.Printf("%s %s [hidden]\n", info.NewestVersion, info.NiceName)
			} else {
				fmt.Printf("%s %s\n", info.NewestVersion, info.NiceName)
			}
		} else {
			fmt.Println(err)
		}
	}
}
