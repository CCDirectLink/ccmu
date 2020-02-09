package internal

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/CCDirectLink/ccmu/game"
	"github.com/CCDirectLink/ccmu/internal/mod"
	"github.com/CCDirectLink/ccmu/pkg"
)

func getGame() game.Game {
	return game.At(flag.Lookup("game").Value.String())
}

func getAll(names []string) []pkg.Package {
	reader := bufio.NewReader(os.Stdin)

	result := make([]pkg.Package, 0, len(names))
	game := getGame()
	for _, name := range names {
		matches := game.Find(name)

		if len(matches) == 0 {
			fmt.Printf("Could not find %s\n", name)
		} else if len(matches) > 1 {
			info, _ := matches[0].Info()
			fmt.Printf("'%s' is ambiguous. Did you mean: %s", name, info.NiceName)
			for i := 1; i < len(matches); i++ {
				info, _ = matches[i].Info()
				fmt.Printf(", %s", info.NiceName)
			}
			fmt.Print("\n")
		} else {
			mod, ok := matches[0].(mod.Mod)
			if !ok || mod.Source != name {
				info, _ := matches[0].Info()

				if strings.ToLower(info.Name) != strings.ToLower(name) && strings.ToLower(info.NiceName) != strings.ToLower(name) {
					fmt.Printf("Entered '%s' but only found %s (%s).\n", name, info.NiceName, info.Name)
					response := "invalid"
					for len(response) != 0 && response[0] != 'y' && response[0] != 'Y' && response[0] != 'n' && response[0] != 'N' {
						fmt.Printf("Pick %s instead [Y/n]? ", info.NiceName)
						response, _ = reader.ReadString('\n')
						response = strings.ReplaceAll(strings.ReplaceAll(response, "\r", ""), "\n", "")
					}

					if len(response) != 0 && response[0] != 'y' && response[0] != 'Y' {
						continue
					}
				}
			}

			result = append(result, matches[0])
		}
	}
	return result
}
