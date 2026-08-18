package api

import (
	"log"
	"main/pkg/db"
	"net/http"
)

const limit int = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(limit, search) // в параметре максимальное количество записей
	if err != nil {
		log.Print(err)
		writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	if len(tasks) == 0 {
		writeJson(w, http.StatusOK, TasksResp{Tasks: []*db.Task{}})
		return
	}
	writeJson(w, http.StatusOK, TasksResp{Tasks: tasks})
}
