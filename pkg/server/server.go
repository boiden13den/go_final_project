package server

import (
	"log"
	"main/pkg/api"
	"net/http"
	"os"
)

// Starts the server
func App() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	mux := http.NewServeMux()
	api.Init(mux)
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
