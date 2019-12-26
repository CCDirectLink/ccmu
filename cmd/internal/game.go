package internal

import (
	"flag"
	"fmt"

	ccmu "github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

func getGame() ccmu.Game {
	return ccmu.At(flag.Lookup("game").Value.String())
}

func getAll(names []string) []pkg.Package {
	game := getGame()
	result := make([]pkg.Package, 0, len(names))
	for _, name := range names {
		pkg, err := game.Get(name)
		if err != nil {
			//TODO:
			fmt.Printf("Error occured in %s\n", err)
			continue
		}

		result = append(result, pkg)
	}
	return result
}
