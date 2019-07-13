package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd"
)

//InstallRequest for incoming installation requests
type InstallRequest struct {
	Game  *string  `json:"game"`
	Names []string `json:"names"`
}

//InstallResponse for installation requests
type InstallResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
	Stats   cmd.Stats `json:"stats,omitempty"`
}

//Install a mod via api request
func Install(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	decoder := json.NewDecoder(r.Body)
	stats, err := install(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&InstallResponse{
			Success: true,
			Stats:   *stats,
		})
	} else {
		encoder.Encode(&InstallResponse{
			Success: false,
			Message: err.Error(),
		})
	}
}

func install(decoder *json.Decoder) (*cmd.Stats, error) {
	var req InstallRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
	}

	if req.Game != nil {
		if err := flag.Set("game", *req.Game); err != nil {
			return nil, fmt.Errorf("cmd/internal/api: Could set game flag: %s", err.Error())
		}
	}

	return cmd.Install(req.Names)
}
