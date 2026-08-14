package api

import (
	"net/http"
)

// Init initializes the API routes
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", auth(taskHandler))
	mux.HandleFunc("/api/task/done", auth(taskDoneHandler))
	mux.HandleFunc("/api/tasks", auth(tasksHandler))
	mux.HandleFunc("/api/signin", loginHandler)
}
