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
	http.ListenAndServe(url, nil)
}
