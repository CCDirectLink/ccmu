package internal

import (
	"flag"

	ccmu "github.com/CCDirectLink/CCUpdaterCLI"
)

func getGame() ccmu.Game {
	return ccmu.At(flag.Lookup("game").Value.String())
}
