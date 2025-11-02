package api

import (
	"encoding/json"
	"final_project/pkg/db"
	"net/http"
	"strconv"
	"strings"
)

// taskHandler — единый обработчик /api/task для POST (добавление), GET (получение по id), PUT (редактирование)
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// шаг 3 — добавление задачи
		addTaskHandler(w, r)

	case http.MethodGet:
		// шаг 6 — получение задачи по id
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
		// шаг 6 — редактирование задачи
		var task db.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJSON(w, map[string]string{"error": "invalid json"})
			return
		}

		// Проверка id
		if strings.TrimSpace(task.ID) == "" {
			writeJSON(w, map[string]string{"error": "id is required"})
			return
		}
		if _, err := strconv.Atoi(task.ID); err != nil {
			writeJSON(w, map[string]string{"error": "invalid id"})
			return
		}

		// Проверка заголовка
		if strings.TrimSpace(task.Title) == "" {
			writeJSON(w, map[string]string{"error": "task title is required"})
			return
		}

		// Проверка формата repeat
		if strings.TrimSpace(task.Repeat) != "" {
			if err := validateRepeatFormat(task.Repeat); err != nil {
				writeJSON(w, map[string]string{"error": "invalid repeat format"})
				return
			}
		}

		// Обработка даты
		if err := processTaskDates(&task); err != nil {
			writeJSON(w, map[string]string{"error": err.Error(), "flag": "1"})
			return
		}

		// Сохранение
		if err := db.UpdateTask(&task); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		// Успех — возвращаем пустой объект {}
		writeJSON(w, map[string]string{})

	default:
		http.Error(w, "method is not supported", http.StatusMethodNotAllowed)
	}
}

// TasksResp — структура ответа для /api/tasks
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// listTasksHandler — шаг 5 (получение списка задач)
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
