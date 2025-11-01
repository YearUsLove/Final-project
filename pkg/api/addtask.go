package api

import (
	"encoding/json"
	"final_project/pkg/db"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем только POST
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task

	// Читаем JSON из запроса
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "invalid json"})
		return
	}

	// Проверка заголовка задачи
	if strings.TrimSpace(task.Title) == "" {
		writeJSON(w, map[string]string{"error": "task title is required"})
		return
	}

	// Новый способ проверки repeat — без изменения даты
	if strings.TrimSpace(task.Repeat) != "" {
		if err := validateRepeatFormat(task.Repeat); err != nil {
			writeJSON(w, map[string]string{"error": "invalid repeat format"})
			return
		}
	}

	// Обработка даты по правилам ТЗ
	if err := processTaskDates(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error(), "flag": "1"})
		return
	}

	//Сохранение задачи в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	//Ответ с ID сохранённой задачи
	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

// processTaskDates обрабатывает дату задачи по требованиям ТЗ
func processTaskDates(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)

	// Если дата пустая — ставим сегодняшнюю
	if strings.TrimSpace(task.Date) == "" {
		task.Date = today
	}

	// Проверка корректности формата даты
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}

	// Проверка на "сегодня"
	if t.Year() == now.Year() &&
		t.Month() == now.Month() &&
		t.Day() == now.Day() {
		return nil
	}

	// Если дата в прошлом
	if isBeforeToday(t, now) {
		if strings.TrimSpace(task.Repeat) == "" {
			// Без повторения — ставим сегодняшнюю дату
			task.Date = today
		} else {
			// С повторением — считаем следующую дату
			next, err := NextDate(now, today, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	// Если дата в будущем — оставляем
	return nil
}

// Сравнение дат без времени (только yyyyMMdd)
func isBeforeToday(date, now time.Time) bool {
	return date.Format(dateFormat) < now.Format(dateFormat)
}

// Запись JSON-ответа
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// validateRepeatFormat проверяет синтаксис repeat, не меняя дату
func validateRepeatFormat(repeat string) error {
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return fmt.Errorf("invalid format")
	}

	rule := strings.ToLower(strings.TrimSpace(parts[0]))
	switch rule {
	case "d":
		if len(parts) != 2 {
			return fmt.Errorf("invalid format")
		}
		n, err := strconv.Atoi(strings.ReplaceAll(parts[1], "+", ""))
		if err != nil || n < 1 || n > 400 {
			return fmt.Errorf("invalid days number")
		}
	case "y":
		if len(parts) != 1 {
			return fmt.Errorf("invalid format")
		}
	case "w":
		if len(parts) != 2 {
			return fmt.Errorf("invalid format")
		}
		days := strings.Split(parts[1], ",")
		for _, d := range days {
			v, err := strconv.Atoi(strings.TrimSpace(d))
			if err != nil || v < 1 || v > 7 {
				return fmt.Errorf("invalid weekday")
			}
		}
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return fmt.Errorf("invalid format")
		}
		days := strings.Split(parts[1], ",")
		for _, d := range days {
			v, err := strconv.Atoi(strings.TrimSpace(d))
			if err != nil || !(v >= 1 && v <= 31 || v == -1 || v == -2) {
				return fmt.Errorf("invalid month day")
			}
		}
		if len(parts) == 3 {
			months := strings.Split(parts[2], ",")
			for _, m := range months {
				v, err := strconv.Atoi(strings.TrimSpace(m))
				if err != nil || v < 1 || v > 12 {
					return fmt.Errorf("invalid month")
				}
			}
		}
	default:
		return fmt.Errorf("invalid format")
	}
	return nil
}
