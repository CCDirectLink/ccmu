package moddb

import (
	"encoding/json"
	"net/http"
)

const modURL = "https://raw.githubusercontent.com/CCDirectLink/CCModDB/master/mods.json"

//Mods represents the root structure of ccmoddb.
type Mods struct {
	Mods map[string]Mod `json:"mods"`
}

//Mod represents the mod structure of ccmoddb.
type Mod struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ArchiveLink string            `json:"archive_link"`
	Hash        map[string]string `json:"hash"`
	Version     string            `json:"version"`
	Licence     string            `json:"licence"`
}

func getMods() (Mods, error) {
	var result Mods

	resp, err := http.Get(modURL)
	if err != nil {
		return result, err
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
