package db

import (
	"database/sql"
	"log"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask add new task to database
func AddTask(task *Task) (int64, error) {
	var id int64
	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	var rows []*Task
	query := `SELECT * FROM scheduler ORDER BY date LIMIT :limit `
	res, err := DB.Query(query, sql.Named("limit", limit))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		row := Task{}
		err := res.Scan(&row.ID, &row.Date, &row.Title, &row.Comment, &row.Repeat)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		rows = append(rows, &row)
	}

	return rows, nil
}
