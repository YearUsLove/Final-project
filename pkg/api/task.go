package api

import (
    "encoding/json"
    "net/http"

    "final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        id := r.URL.Query().Get("id")
        if id == "" {
            writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
            return
        }
        t, err := db.GetTask(id)
        if err != nil {
            writeJSON(w, map[string]string{"error": err.Error()})
            return
        }
        writeJSON(w, t)

    case http.MethodPost:
        var t db.Task
        if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
            writeJSON(w, map[string]string{"error": "invalid json"})
            return
        }
        _, err := db.DB.Exec(`
            INSERT INTO scheduler (date, title, comment, repeat)
            VALUES (?, ?, ?, ?)`,
            t.Date, t.Title, t.Comment, t.Repeat,
        )
        if err != nil {
            writeJSON(w, map[string]string{"error": err.Error()})
            return
        }
        writeJSON(w, map[string]string{})

    case http.MethodPut:
        var t db.Task
        if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
            writeJSON(w, map[string]string{"error": "invalid json"})
            return
        }
        if t.ID == "" {
            writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
            return
        }
        err := db.UpdateTask(&t)
        if err != nil {
            writeJSON(w, map[string]string{"error": err.Error()})
            return
        }
        writeJSON(w, map[string]string{})

    default:
        writeJSON(w, map[string]string{"error": "method not allowed"})
    }
}