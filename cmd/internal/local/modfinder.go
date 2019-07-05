package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Mod contains the data of the installed mod
type Mod struct {
	Name     string
	BasePath string
	Version  string
}

//GetMods finds all local mods
func GetMods() ([]Mod, error) {
	game, err := GetGame()
	if err != nil {
		return nil, err
	}

	mods := filepath.Join(game, "assets/mods")
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
			mod, err := parseMod(filepath.Join(mods, dir.Name(), "package.json"))
			if err == nil {
				result = append(result, mod)
			}
		}
	}

	return result, nil
}

//GetMod finds the installed mod by name
func GetMod(name string) (Mod, error) {
	mods, err := GetMods()
	if err != nil {
		return Mod{}, err
	}

	for _, mod := range mods {
		if name == mod.Name {
			return mod, nil
		}
	}

	return Mod{}, fmt.Errorf("cmd/internal: Could not find mod '%s'", name)
}

func parseMod(path string) (Mod, error) {
	file, err := os.Open(path)
	if err != nil {
		return Mod{}, nil
	}
	defer file.Close()

	var data struct {
		Name    string  `json:"name"`
		Version *string `json:"version"`
	}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return Mod{}, nil
	}

	var version string
	if data.Version != nil {
		version = *data.Version
	} else {
		version = "0.0.0"
	}

	return Mod{
		data.Name,
		filepath.Dir(path),
		version,
	}, nil
}
