package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//LocalMod contains the data of the installed mod
type LocalMod struct {
	Name     string
	BasePath string
	Version  string
}

//GetLocalMods finds all local mods
func GetLocalMods() ([]LocalMod, error) {
	game, err := findGame()
	if err != nil {
		return nil, err
	}

	mods := filepath.Join(game, "assets/mods")
	if exists, _ := exists(mods); !exists {
		return []LocalMod{}, nil
	}

	dirs, err := ioutil.ReadDir(mods)
	if err != nil {
		return nil, err
	}

	var result []LocalMod
	for _, dir := range dirs {
		if dir.IsDir() {
			mod, err := parseMod(filepath.Join(mods, dir.Name(), "package.json"))
			if err == nil {
				result = append(result, mod)
			}
		}
	}

	return result, nil
}

//GetLocalMod finds the installed mod by name
func GetLocalMod(name string) (LocalMod, error) {
	mods, err := GetLocalMods()
	if err != nil {
		return LocalMod{}, err
	}

	for _, mod := range mods {
		if name == mod.Name {
			return mod, nil
		}
	}

	return LocalMod{}, fmt.Errorf("cmd/internal: Could not find mod '%s'", name)
}

func parseMod(path string) (LocalMod, error) {
	file, err := os.Open(path)
	if err != nil {
		return LocalMod{}, nil
	}
	defer file.Close()

	var data struct {
		Name    string  `json:"name"`
		Version *string `json:"version"`
	}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return LocalMod{}, nil
	}

	var version string
	if data.Version != nil {
		version = *data.Version
	} else {
		version = "0.0.0"
	}

	return LocalMod{
		data.Name,
		filepath.Dir(path),
		version,
	}, nil
}
