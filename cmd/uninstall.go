package cmd

import (
	"fmt"
	"os"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
)

//Uninstall removes a mod from a directory
func Uninstall(args []string) {
	for _, name := range args {
		mod, err := internal.GetLocalMod(name)
		if err != nil {
			fmt.Printf("Could not find mod '%s'\n", name)
			continue
		}

		err = os.RemoveAll(mod.BasePath)
		if err != nil {
			fmt.Printf("Could not remove mod '%s' because of an error in %s\n", name, err.Error())
		}
	}
}
