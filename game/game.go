package game

import (
	"net/url"
	"os"
	"strings"

	"github.com/CCDirectLink/ccmu/internal/mod"
	"github.com/CCDirectLink/ccmu/internal/moddb"
	"github.com/CCDirectLink/ccmu/pkg"
	"github.com/CCDirectLink/ccmu/tool"
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

//BasePath returns the resolved path of the game.
func (g Game) BasePath() (string, error) {
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
		var p pkg.Package

		switch strings.ToLower(info.NiceName) {
		case "ccloader":
			p, err = tool.Get(g, "ccloader")
			if err != nil {
				return nil, err
			}
		default:
			p = mod.Mod{
				Name: info.Name,
				Path: g.Path,
				Game: &g,
			}
		}

		result = append(result, p)
	}

	return result, nil
}

//FromSource returns a mod that represents the source object
func (g Game) FromSource(source string) (pkg.Package, error) {
	result, err := mod.FromSource(source, g)
	if err != nil {
		return result, pkg.NewError(pkg.ModeUnknown, result, err)
	}
	return result, nil
}

//Get mod by name.
func (g Game) Get(name string) (pkg.Package, error) {
	switch strings.ToLower(name) {
	case "ccloader":
		return tool.Get(g, "ccloader")
	case "crosscode":
		return tool.Get(g, "crosscode")
	case "browser":
		return tool.Get(g, "browser")
	}

	path, err := g.BasePath()

	result := mod.Mod{
		Name: name,
		Path: path,
		Game: &g,
	}

	if err != nil {
		return result, pkg.NewError(pkg.ModeUnknown, result, err)
	}

	if !result.Installed() && !result.Available() {
		return result, pkg.NewError(pkg.ModeUnknown, result, pkg.ErrNotFound)
	}

	return result, nil
}

//Find matching mods with (part of) the given name.
func (g Game) Find(name string) []pkg.Package {
	avail, _ := g.Available()
	inst, _ := g.Installed()

	src := g.findSource(name)
	if src != nil {
		return []pkg.Package{src}
	}

	exact := findExact(name, avail, inst)
	if exact != nil {
		return []pkg.Package{exact}
	}

	toolPkg, err := tool.Get(g, name)
	if err == nil && toolPkg != nil {
		return []pkg.Package{toolPkg}
	}

	return findAll(name, avail, inst)
}

func (g Game) findSource(name string) pkg.Package {
	if ok, _ := isValidURL(name); ok {
		result, err := mod.FromSource(name, g)
		if err != nil {
			return nil
		}
		return result
	}

	stat, err := os.Stat(name)
	if err == nil && !stat.IsDir() {
		result, err := mod.FromSource(name, g)
		if err != nil {
			return nil
		}
		return result
	}
	return nil
}

func isValidURL(toTest string) (bool, string) {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false, ""
	}

	u, err := url.Parse(toTest)
	return err == nil && u.Scheme != "" && u.Host != "", u.Path
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
