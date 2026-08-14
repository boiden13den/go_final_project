package api

import (
	"net/http"
)

// Init initializes the API routes
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/task/done", taskDoneHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
}
