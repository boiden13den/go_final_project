package server

import (
	"log"
	"net/http"
	"os"
)

// Starts the server
func App() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
