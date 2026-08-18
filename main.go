package main

import (
	"log"
	"main/pkg/db"
	"main/pkg/server"
	"os"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	_, err := os.Stat(dbFile)
	if err := db.InitDB(dbFile, err != nil); err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	if err := server.App(); err != nil {
		log.Fatal(err)
	}
}
