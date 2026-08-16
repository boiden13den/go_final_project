package api

import (
	"net/http"
)

// Init initializes the API routes
func Init(mux *http.ServeMux) {
<<<<<<< HEAD
	mux.HandleFunc("/api/nextdate", auth(nextDayHandler))
=======
	mux.HandleFunc("/api/nextdate", nextDayHandler)
>>>>>>> 1ec9f59b5b6f58a7bd148c21a317831a3eade709
	mux.HandleFunc("/api/task", auth(taskHandler))
	mux.HandleFunc("/api/task/done", auth(taskDoneHandler))
	mux.HandleFunc("/api/tasks", auth(tasksHandler))
	mux.HandleFunc("/api/signin", loginHandler)
}
