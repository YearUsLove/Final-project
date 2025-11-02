package api

import (
	"encoding/json"
	"net/http"

	"final_project/pkg/db"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]string{"error": "method not allowed"})
		return
	}

	rows, err := db.DB.Query(`
        SELECT id, date, title, comment, repeat 
        FROM scheduler
        ORDER BY date
    `)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var tasks []db.Task
	for rows.Next() {
		var t db.Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
