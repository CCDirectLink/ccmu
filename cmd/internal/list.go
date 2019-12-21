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
		fmt.Printf("%s %s %s\n", info.NewestVersion, info.NiceName, err)
	}
}
