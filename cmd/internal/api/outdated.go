package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
)

//OutdatedRequest for incoming outdated requests
type OutdatedRequest struct {
	Game *string `json:"game"`
}

//OutdatedResponse contains a list of outdated mods
type OutdatedResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message,omitempty"`
	Mods    []OutdatedDescription `json:"mods"`
}

//OutdatedDescription contains basic information about an outdated mod
type OutdatedDescription struct {
	Current string `json:"current"`
	Newest  string `json:"newest"`
	Name    string `json:"name"`
}

//Outdated returns all available mods
func Outdated(w http.ResponseWriter, r *http.Request) {
	var decoder *json.Decoder
	if r.Method == "POST" {
		decoder = json.NewDecoder(r.Body)
	}

	mods, err := outdated(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&OutdatedResponse{
			Success: true,
			Mods:    mods,
		})
	} else {
		encoder.Encode(&OutdatedResponse{
			Success: false,
			Message: err.Error(),
		})
	}
}

func outdated(decoder *json.Decoder) ([]OutdatedDescription, error) {
	if decoder != nil {
		var req OutdatedRequest
		if err := decoder.Decode(&req); err != nil {
			return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
		}

		if req.Game != nil {
			if err := flag.Set("game", *req.Game); err != nil {
				return nil, fmt.Errorf("cmd/internal/api: Could set game flag: %s", err.Error())
			}
		}
	}

	mods, err := local.GetMods()
	if err != nil {
		return nil, fmt.Errorf("cmd/internal/api: Could not list mods because of an error in %s", err.Error())
	}

	var res []OutdatedDescription
	for _, mod := range mods {
		if out, _ := mod.Outdated(); out {
			new, err := global.GetMod(mod.Name)
			if err != nil {
				fmt.Printf("An error occured in %s\n", err.Error())
				continue
			}

			res = append(res, OutdatedDescription{
				Current: mod.Version,
				Newest:  new.Version,
				Name:    mod.Name,
			})
		}
	}
	return res, nil
}
