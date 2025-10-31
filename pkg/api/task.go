package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
    "encoding/json"
    "net/http"

	"final_project/pkg/db"
    "final_project/pkg/db"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodPut:
		handleUpdateTask(w, r)
	default:
		writeJSON(w, map[string]string{"error": "method not allowed"})
	}
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	task, err := getTaskByID(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, map[string]string{"error": "Задача не найдена"})
		} else {
			writeJSON(w, map[string]string{"error": err.Error()})
		}
		return
	}

	writeJSON(w, map[string]string{
		"id":      strconv.FormatInt(task.ID, 10),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	})
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Проверка ID
	if t.ID == 0 {
		writeJSON(w, map[string]string{"error": "Некорректный идентификатор"})
		return
	}

	// Проверка даты
	today := time.Now().Format("20060102")
	if len(t.Date) != 8 || t.Date < today {
		writeJSON(w, map[string]string{"error": "Дата не может быть меньше сегодняшней"})
		return
	}

	// Проверка title
	if t.Title == "" {
		writeJSON(w, map[string]string{"error": "Пустой заголовок"})
		return
	}

	// Проверка существования задачи
	if _, err := getTaskByID(t.ID); err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	// Обновление
	if err := updateTask(t.ID, &t); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Успешно — пустой JSON {}
	writeJSON(w, map[string]string{})
}

func getTaskByID(id int64) (*db.Task, error) {
	row := db.DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
	var t db.Task
	if err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		return nil, err
	}
	return &t, nil
}

func updateTask(id int64, t *db.Task) error {
	res, err := db.DB.Exec(`UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`,
		t.Date, t.Title, t.Comment, t.Repeat, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}