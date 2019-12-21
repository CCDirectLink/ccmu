package ccmu

import (
	"github.com/CCDirectLink/CCUpdaterCLI/internal/mod"
	"github.com/CCDirectLink/CCUpdaterCLI/internal/moddb"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

//Game represents a game instance at a given path.
type Game struct {
	Path string
}

//Default game instance at the current working directory.
var Default = Game{}

//At returns the game instance at the given path.
func At(path string) Game {
	return Game{path}
}

func (g Game) path() (string, error) {
	if g.Path != "" {
		return g.Path, nil
	}
	return getGame()
}

//Installed mods.
func (g Game) Installed() ([]pkg.Package, error) {
	return nil, nil
}

//Available mods.
func (g Game) Available() ([]pkg.Package, error) {
	infos, err := moddb.ModInfos()
	if err != nil {
		return nil, err
	}

	result := make([]pkg.Package, len(infos))
	for i, info := range infos {
		result[i] = mod.Mod{
			Name: info.Name,
			Game: &g,
		}
	}

	return result, nil
}

//Get mod by name.
func (g Game) Get(name string) (pkg.Package, error) {
	return nil, nil
}
