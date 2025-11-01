package api

import (
	"final_project/pkg/db"
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		http.Error(w, "method is not supported", http.StatusMethodNotAllowed)
	}
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method is not supported", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, TasksResp{Tasks: tasks})
}
