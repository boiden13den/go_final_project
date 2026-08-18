package server

import (
	"log"
	"main/pkg/api"
	"net/http"
	"os"
)

// Starts the server
func App() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	password := os.Getenv("TODO_PASSWORD")
	mux := http.NewServeMux()
	api.Init(mux, password)
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(":"+port, mux)
}
