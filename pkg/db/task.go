package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
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

// Tasks return all tasks from database
func Tasks(limit int, search string) ([]*Task, error) {
	var rows []*Task
	var res *sql.Rows
	var err error

	if search != "" {
		date, err := time.Parse("02.01.2006", search)
		if err == nil {
			res, err = DB.Query(
				`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :search ORDER BY date LIMIT :limit`,
				sql.Named("search", date.Format("20060102")), sql.Named("limit", limit),
			)
		} else {
			res, err = DB.Query(
				`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE %:search% OR comment LIKE %:search% ORDER BY date LIMIT :limit`,
				sql.Named("search", search), sql.Named("limit", limit),
			)
		}
	} else {
		res, err = DB.Query(
			`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`,
			sql.Named("limit", limit),
		)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		row := Task{}
		err = res.Scan(&row.ID, &row.Date, &row.Title, &row.Comment, &row.Repeat)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		rows = append(rows, &row)
	}

	return rows, nil
}

// GetTask return task by id
func GetTask(id string) (*Task, error) {
	var task = &Task{}
	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`,
		sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return task, nil
}

// UpdateTask update task by id
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// DeleteTask delete task by id
func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}
