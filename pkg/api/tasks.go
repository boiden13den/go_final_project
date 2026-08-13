package api

import (
	"main/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		// здесь вызываете функцию, которая возвращает ошибку в JSON
		// её желательно было реализовать на предыдущем шаге
		writeJson(w, map[string]any{"error": err})
		return
	}
	if len(tasks) == 0 {

		writeJson(w, TasksResp{Tasks: []*db.Task{}})
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
