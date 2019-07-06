package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd"
)

//UpdateRequest for incoming update requests
type UpdateRequest struct {
	Game  *string  `json:"game"`
	Names []string `json:"names"`
}

//UpdateResponse for update requests
type UpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

//Update a mod via api request
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := update(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&UpdateResponse{
			Success: true,
		})
	} else {
		encoder.Encode(&UpdateResponse{
			Success: false,
			Message: err.Error(),
		})
	}
}

func update(decoder *json.Decoder) error {
	var req UpdateRequest
	if err := decoder.Decode(&req); err != nil {
		return fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
	}

	if req.Game != nil {
		if err := flag.Set("game", *req.Game); err != nil {
			return fmt.Errorf("cmd/internal/api: Could set game flag: %s", err.Error())
		}
	}

	if result := cmd.Update(req.Names); !result {
		return fmt.Errorf("cmd/internal/api: Could not update mods")
	}

	return nil
}
