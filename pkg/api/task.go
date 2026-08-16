package api

import (
	"encoding/json"
	"errors"
	"log"
	"main/pkg/db"
	"net/http"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]any{"error": "не указан заголовок задачи"})
		return
	}
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, map[string]any{"id": id})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	if !t.After(now) && t.Format(DateFormat) != now.Format(DateFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("writeJson error: %v", err)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := getTask(w, r)
	if err != nil {
		return
	}
	writeJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	if task.ID == "" {
		writeJson(w, map[string]any{"error": "не указан номер задачи"})
		return
	}
	if task.Date == "" {
		writeJson(w, map[string]any{"error": "не указан дата выполнения задачи"})
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]any{"error": "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	err := db.UpdateTask(&task)

	if err != nil {
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, map[string]any{"id": task.ID})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := getTask(w, r)
	if err != nil {
		return
	}
	err = db.DeleteTask(task.ID)
	if err != nil {
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, map[string]any{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	task, err := getTask(w, r)
	if err != nil {
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(task.ID)
		if err != nil {
			writeJson(w, map[string]any{"error": err.Error()})
			return
		}
		writeJson(w, map[string]any{})
		return // ← вот этого return не было
	}
	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	task.Date = next
	err = db.UpdateTask(task)
	if err != nil {
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, map[string]any{})
}

func getTask(w http.ResponseWriter, r *http.Request) (*db.Task, error) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]any{"error": "Не указан идентификатор"})
		log.Println(map[string]any{"error": "Не указан идентификатор"})
		return nil, errors.New("errGettask")
	}
	task, err := db.GetTask(id)
	if err != nil {
		log.Println("error: Не указан идентификатор")
		writeJson(w, map[string]any{"error": "Задача не найдена"})
		return nil, errors.New("errGettask")
	}
	return task, nil
}
