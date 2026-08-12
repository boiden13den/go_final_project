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
	api.Init()
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
