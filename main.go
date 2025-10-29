package main

import (
	"log"
	"path/filepath"

	"final_project/pkg/db"
	"final_project/pkg/server"
)

func main() {

	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal(err)
	}

	webDir := filepath.Join(".", "web")
	if err := server.Start(webDir); err != nil {
		log.Fatal(err)
	}
}
