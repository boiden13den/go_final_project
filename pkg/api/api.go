package api

import (
	"net/http"
)

var appPassword string

// Init initializes the API routes
func Init(mux *http.ServeMux, password string) {
	appPassword = password
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", auth(taskHandler))
	mux.HandleFunc("/api/task/done", auth(taskDoneHandler))
	mux.HandleFunc("/api/tasks", auth(tasksHandler))
	mux.HandleFunc("/api/signin", loginHandler)
}
