package api

import (
	"net/http"

	"final_project/pkg/db"
)

func ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	out := make([]map[string]string, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, map[string]string{
			"id":      int64ToStr(t.ID),
			"date":    t.Date,
			"title":   t.Title,
			"comment": t.Comment,
			"repeat":  t.Repeat,
		})
	}

	writeJSON(w, map[string]any{"tasks": out})
}