package db

import (
    "database/sql"
    "fmt"
)

type Task struct {
    ID      string `json:"id,omitempty"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

func GetTask(id string) (*Task, error) {
    var t Task
    t.ID = id
    err := DB.QueryRow(`
        SELECT date, title, comment, repeat
        FROM scheduler
        WHERE id = ?
    `, id).Scan(&t.Date, &t.Title, &t.Comment, &t.Repeat)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("Задача не найдена")
    }
    if err != nil {
        return nil, err
    }
    return &t, nil
}

func UpdateTask(task *Task) error {
    res, err := DB.Exec(`
        UPDATE scheduler
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?
    `, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("Задача не найдена")
    }
    return nil
}