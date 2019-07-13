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
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Stats   *cmd.Stats `json:"stats,omitempty"`
}

//Update a mod via api request
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	stats, err := update(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&UpdateResponse{
			Success: true,
			Stats:   stats,
		})
	} else {
		encoder.Encode(&UpdateResponse{
			Success: false,
			Message: err.Error(),
			Stats:   stats,
		})
	}
}

func update(decoder *json.Decoder) (*cmd.Stats, error) {
	var req UpdateRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
	}

	if req.Game != nil {
		if err := flag.Set("game", *req.Game); err != nil {
			return nil, fmt.Errorf("cmd/internal/api: Could set game flag: %s", err.Error())
		}
	}

	return cmd.Update(req.Names)
}
