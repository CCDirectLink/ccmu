package cmd

import (
	"fmt"
	"os"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
)

//List prints a list of all available mods
func List() {
	data, err := internal.FetchModData()
	if err != nil {
		fmt.Printf("Could not list mods because of an error in %s\n", err.Error())
		os.Exit(1)
	}

	for _, mod := range data.Mods {
		fmt.Printf("%s %s\n", mod.Version, mod.Name)
	}
}
