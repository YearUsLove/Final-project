package api

import (
	"encoding/json"
	"final_project/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// taskHandler — единый обработчик /api/task для POST, GET, PUT, DELETE
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
			return
		}
		task, err := db.GetTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, task)

	case http.MethodPut:
		var task db.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJSON(w, map[string]string{"error": "invalid json"})
			return
		}
		if strings.TrimSpace(task.ID) == "" {
			writeJSON(w, map[string]string{"error": "id is required"})
			return
		}
		if _, err := strconv.Atoi(task.ID); err != nil {
			writeJSON(w, map[string]string{"error": "invalid id"})
			return
		}
		if strings.TrimSpace(task.Title) == "" {
			writeJSON(w, map[string]string{"error": "task title is required"})
			return
		}
		if strings.TrimSpace(task.Repeat) != "" {
			if err := validateRepeatFormat(task.Repeat); err != nil {
				writeJSON(w, map[string]string{"error": "invalid repeat format"})
				return
			}
		}
		if err := processTaskDates(&task); err != nil {
			writeJSON(w, map[string]string{"error": err.Error(), "flag": "1"})
			return
		}
		if err := db.UpdateTask(&task); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if strings.TrimSpace(id) == "" {
			writeJSON(w, map[string]string{"error": "id is required"})
			return
		}
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})

	default:
		writeJSON(w, map[string]string{"error": "method is not supported"})
	}
}

// listTasksHandler — шаг 5 (список задач)
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]string{"error": "method is not supported"})
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

// taskDoneHandler — шаг 7 (POST /api/task/done)
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "method is not supported"})
		return
	}

	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	// Правильное вычисление следующей даты от текущей даты задачи
	start, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "invalid task date"})
		return
	}

	next, err := NextDate(start, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}
