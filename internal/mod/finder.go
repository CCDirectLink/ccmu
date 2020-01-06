package mod

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

var errInvalidName = errors.New("mod: Invalid mod name")

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
		}
	}

	return unknown, pkg.ErrNotFound
}

func parseMod(path, gamePath string, game game) (Mod, error) {
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
