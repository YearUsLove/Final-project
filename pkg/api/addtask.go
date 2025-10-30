package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"final_project/pkg/db"
	"strconv"
)

// Заглушка — проверяет формат repeat, например "d 5", "daily", "weekly"
func validateRepeatFormat(repeat string) bool {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return true // пустое значение — допустимо
	}
	// Простейшая проверка
	if repeat == "daily" || repeat == "weekly" {
		return true
	}
	if strings.HasPrefix(repeat, "d ") && len(repeat) > 2 {
		return true
	}
	return false
}

// Возвращает true, если дата раньше сегодняшней (в формате YYYYMMDD)
func isBeforeToday(dateStr string) bool {
	if len(dateStr) != 8 {
		return false
	}
	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return false
	}
	today := time.Now().Truncate(24 * time.Hour)
	return t.Before(today)
}

// Заглушка для NextDate — возвращает ту же дату или модифицированную по repeat
func NextDate(startDate, repeat string) string {
	if repeat == "" {
		return startDate
	}
	// Для примера — daily = +1 день
	d, err := time.Parse("20060102", startDate)
	if err != nil {
		return startDate
	}
	switch repeat {
	case "daily":
		return d.AddDate(0, 0, 1).Format("20060102")
	case "weekly":
		return d.AddDate(0, 0, 7).Format("20060102")
	}
	// пример для "d N"
	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Split(repeat, " ")
		if len(parts) == 2 {
			n, _ := strconv.Atoi(parts[1])
			return d.AddDate(0, 0, n).Format("20060102")
		}
	}
	return startDate
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "method not allowed"})
		return
	}

	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, map[string]string{"error": "invalid json"})
		return
	}

	// Проверка repeat
	if !validateRepeatFormat(t.Repeat) {
		writeJSON(w, map[string]string{"error": "invalid repeat format"})
		return
	}

	// Проверка на дату в прошлом
	if isBeforeToday(t.Date) {
		writeJSON(w, map[string]string{"error": "date is before today"})
		return
	}

	// Вычислим дату следующей задачи (если нужна)
	next := NextDate(t.Date, t.Repeat)
	_ = next // чтобы не было "declared and not used"

	// Вставка в БД
	_, err := db.DB.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)`,
		t.Date, t.Title, t.Comment, t.Repeat,
	)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{}) // успех
}
