package mod

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/CCDirectLink/ccmu/internal/game"
	"github.com/CCDirectLink/ccmu/pkg"
)

var errInvalidName = errors.New("mod: Invalid mod name")

//GetMods finds all local mods
func GetMods(path string, game game.Game) ([]Mod, error) {
	mods := filepath.Join(path, "assets/mods")
	if exists, _ := exists(mods); !exists {
		return []Mod{}, nil
	}

	dirs, err := ioutil.ReadDir(mods)
	if err != nil {
		return nil, err
	}

	var result []Mod
	for _, dir := range dirs {
		if dir.IsDir() {
			mod, err := parseMod(filepath.Join(mods, dir.Name(), "package.json"), path, game)
			if err == nil && mod.Name != "" {
				result = append(result, mod)
			}
		}
	}

	return result, nil
}

func (m Mod) local() (string, error) {
	mods := filepath.Join(m.Path, "assets", "mods")
	if exists, _ := exists(mods); !exists {
		return "", pkg.ErrNotFound
	}

	unknown := filepath.Join(mods, m.Name)

	dirs, err := ioutil.ReadDir(mods)
	if err != nil {
		return unknown, err
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			result := filepath.Join(mods, dir.Name())

			mod, err := parseMod(filepath.Join(result, "package.json"), m.Path, m.Game)
			if mod.Name == m.Name {
				return result, err
			}
		} else if filepath.Ext(dir.Name()) == ".ccmod" {
			result := filepath.Join(mods, dir.Name())
			pkg, _ := readPackageFromCCMod(result)

			if pkg.Name == m.Name {
				return result, nil
			}
		}
	}

	return unknown, pkg.ErrNotFound
}

func parseMod(path, gamePath string, game game.Game) (Mod, error) {
	file, err := os.Open(path)
	if err != nil {
		return Mod{}, err
	}
	defer file.Close()

	data, err := readPackage(file)
	if err != nil {
		return Mod{}, err
	}

	if data.Name == "" {
		return Mod{}, errInvalidName
	}

	return Mod{
		data.Name,
		gamePath,
		game,
		"",
		"",
	}, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if !os.IsNotExist(err) {
		return true, err
	}
	return false, nil
}
