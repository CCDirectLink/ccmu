package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
)

//LocalModsRequest for incoming installed mod list requests
type LocalModsRequest struct {
	Game *string `json:"game"`
}

//LocalModsResponse contains a list of installed mods
type LocalModsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Mods    []local.Mod `json:"mods"`
}

//GetLocalMods returns all installed mods
func GetLocalMods(w http.ResponseWriter, r *http.Request) {
	var decoder *json.Decoder
	if r.Method == "POST" {
		decoder = json.NewDecoder(r.Body)
	}

	setHeaders(w)

	mods, err := getLocalMods(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&LocalModsResponse{
			Success: true,
			Mods:    mods,
		})
	} else {
		encoder.Encode(&LocalModsResponse{
			Success: false,
			Message: err.Error(),
		})
	}
}

func getLocalMods(decoder *json.Decoder) ([]local.Mod, error) {
	if decoder != nil {
		var req LocalModsRequest
		if err := decoder.Decode(&req); err != nil {
			return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
		}

		if req.Game != nil {
			if err := flag.Set("game", *req.Game); err != nil {
				return nil, fmt.Errorf("cmd/internal/api: Could set game flag: %s", err.Error())
			}
		}
	}

	return local.GetMods()
}
