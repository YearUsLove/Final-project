package api

import (
	"net/http"
)

func InitRoutes() {
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/tasks", listTasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}
