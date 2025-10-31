package api

import (
	"net/http"
)

func InitRoutes() {
	// http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)
}
