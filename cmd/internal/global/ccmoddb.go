package global

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const link = "https://raw.githubusercontent.com/CCDirectLink/CCModDB/master/mods.json"

//CCModDb contains data about mods
type CCModDb struct {
	Mods map[string]Mod `json:"mods"`
}

//Mod defines the CCModDb mod structure
type Mod struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	License     *string `josn:"license"`
	Page        []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	ArchiveLink string `json:"archive_link"`
	Hash        struct {
		Sha256 string `json:"sha256"`
	} `json:"hash"`
	Version string `json:"version"`
	Dir     *struct {
		Any string `json:"any"`
	} `json:"dir"`
}

var data *CCModDb

//FetchModData from CCModDb
func FetchModData() (*CCModDb, error) {
	if data != nil {
		return data, nil
	}

	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	data = &CCModDb{}
	err = json.NewDecoder(res.Body).Decode(data)
	return data, err
}

//GetMod returns the ccmoddb mod by name
func GetMod(name string) (Mod, error) {
	_, err := FetchModData()
	if err != nil {
		return Mod{}, err
	}

	for _, mod := range data.Mods {
		if mod.Name == name {
			return mod, nil
		}
	}
	return Mod{}, fmt.Errorf("cmd/internal: Could not find mod '%s'", name)
}

func modKnown(name string) (bool, error) {
	_, err := FetchModData()
	if err != nil {
		return false, err
	}

	for _, mod := range data.Mods {
		if mod.Name == name {
			return true, nil
		}
	}
	return false, nil
}
