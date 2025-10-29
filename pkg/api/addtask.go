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

const dateFormat = "20060102"

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	if r. Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "invalid json"})
		return 
	}

	if strings.TrimSpace(task.Title) == "" {
		if err := validateRepeatFormat(task.Repeat); err != nil {
		writeJSON(w, map[string]string{"error": "task titile is required"})
		return
	}
} 

if err := processTaskDates(&task); err != nil {
	writeJSON(w, map[string]string{"error": err.Error()})
	return
}

writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})

}

func processTaskDates(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)

	if strings.TrimSpace(task.Date) == "" {
		task.Date = today
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}

	if t.Year() == now.Year() &&
	t.Month() == now.Month() &&
	t.Day() == now.Day() {
		return nil
	}

	if isBeforeToday(t, now) {
		if strings.TrimSpace(task.Repeat) == "" {

			task.Date = today
		} else {
			next, err := NextDate(now, today, task.Repeat)
			if err != nil {
				
			}

		}
	}
}


