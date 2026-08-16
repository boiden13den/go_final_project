package api

import (
	"log"
	"main/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search) // в параметре максимальное количество записей
	if err != nil {
		log.Print(err)
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
