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
	raw, err := mod.GetMods(g.Path, &g)

	result := make([]pkg.Package, len(raw))
	for i, mod := range raw {
		result[i] = mod
	}

	if err != nil {
		return result, pkg.NewError(pkg.ModeUnknown, nil, err)
	}
	return result, nil
}

//Available mods.
func (g Game) Available() ([]pkg.Package, error) {
	infos, err := moddb.PkgInfos()
	if err != nil {
		return nil, err
	}

	result := make([]pkg.Package, 0, len(infos))
	for _, info := range infos {
		if info.NiceName == "CCLoader" {
			continue
		}

		result = append(result, mod.Mod{
			Name: info.Name,
			Path: g.Path,
			Game: &g,
		})
	}

	return result, nil
}

//Get mod by name.
func (g Game) Get(name string) (pkg.Package, error) {
	path, err := g.path()

	result := mod.Mod{
		Name: name,
		Path: path,
		Game: &g,
	}

	if err != nil {
		return result, pkg.NewError(pkg.ModeUnknown, result, err)
	}
	return result, nil
}
