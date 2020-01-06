package game

import (
	"github.com/CCDirectLink/CCUpdaterCLI/internal/mod"
	"github.com/CCDirectLink/CCUpdaterCLI/internal/moddb"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
	"strings"
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
		return searchForGame(g.Path)
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

//Find matching mods with (part of) the given name.
func (g Game) Find(name string) []pkg.Package {
	avail, _ := g.Available()
	inst, _ := g.Installed()

	exact := findExact(name, avail, inst)
	if exact != nil {
		return []pkg.Package{exact}
	}

	return findAll(name, avail, inst)
}

func findExact(name string, avail, inst []pkg.Package) pkg.Package {
	name = strings.ToLower(name)
	for _, pkg := range avail {
		info, _ := pkg.Info()
		if strings.ToLower(info.Name) == name || strings.ToLower(info.NiceName) == name {
			return pkg
		}
	}
	for _, pkg := range inst {
		info, _ := pkg.Info()
		if strings.ToLower(info.Name) == name || strings.ToLower(info.NiceName) == name {
			return pkg
		}
	}
	return nil
}

func findAll(name string, avail, inst []pkg.Package) []pkg.Package {
	name = strings.ToLower(name)
	result := make([]pkg.Package, 0, len(avail)+len(inst))
	for _, pkg := range avail {
		info, _ := pkg.Info()
		if strings.Contains(strings.ToLower(info.Name), name) || strings.Contains(strings.ToLower(info.NiceName), name) {
			result = append(result, pkg)
		}
	}
	for _, pkg := range inst {
		info, _ := pkg.Info()
		if strings.Contains(strings.ToLower(info.Name), name) || strings.Contains(strings.ToLower(info.NiceName), name) {
			new := true
			for _, available := range avail {
				availInfo, _ := available.Info()
				if availInfo.Name == info.Name {
					new = false
					break
				}
			}
			if new {
				result = append(result, pkg)
			}
		}
	}
	return result
}
