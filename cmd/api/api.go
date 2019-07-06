package api

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/api"
)

//Start api server
func Start() {
	var port int
	flag.IntVar(&port, "port", 9392, "the port which the api server listens on")

	url := fmt.Sprintf(":%d", port)
	fmt.Printf("API server listening on %s\n", url)

	http.HandleFunc("/api/v1/install", api.Install)
	http.HandleFunc("/api/v1/uninstall", api.Uninstall)
	http.HandleFunc("/api/v1/update", api.Update)
	http.HandleFunc("/api/v1/get/local", api.GetLocalMods)
	http.HandleFunc("/api/v1/get/global", api.GetGlobalMods)
	http.HandleFunc("/api/v1/get/outdated", api.Outdated)

	http.ListenAndServe(url, nil)
}
