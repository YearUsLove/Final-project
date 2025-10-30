package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
    dateStr := r.URL.Query().Get("date")
    repeat := r.URL.Query().Get("repeat")

    // проверка формата даты
    if len(dateStr) != 8 {
        writeJSON(w, map[string]string{"error": "invalid date format"})
        return
    }
    year, err1 := strconv.Atoi(dateStr[:4])
    month, err2 := strconv.Atoi(dateStr[4:6])
    day, err3 := strconv.Atoi(dateStr[6:])
    if err1 != nil || err2 != nil || err3 != nil {
        writeJSON(w, map[string]string{"error": "invalid date format"})
        return
    }

    currentDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
    var nextDate time.Time

    switch strings.ToLower(repeat) {
    case "daily":
        nextDate = currentDate.AddDate(0, 0, 1)
    case "weekly":
        nextDate = currentDate.AddDate(0, 0, 7)
    case "monthly":
        nextDate = currentDate.AddDate(0, 1, 0)
    case "yearly":
        nextDate = currentDate.AddDate(1, 0, 0)
    default:
        writeJSON(w, map[string]string{"error": "invalid repeat format"})
        return
    }

    writeJSON(w, map[string]string{"nextdate": nextDate.Format("20060102")})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(v)
}