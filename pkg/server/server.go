package server

import (
	"fmt"
	"net/http"
	"os"

	"final_project/pkg/api"
)

func Start(webDir string) error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	mux := http.NewServeMux()

	api.Init(mux)

	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fs)

	
	fmt.Printf("Server running on port %s\n", port)
	
	return http.ListenAndServe(":"+port, mux)
}