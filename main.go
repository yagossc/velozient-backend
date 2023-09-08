package main

import (
	"log"
	"velozient-backend/api"
	"velozient-backend/db"
)

func main() {
	// Create and initialize the database
	database := db.NewMemoryDB()
	database.PopulateDB(db.InitialLoad)

	// Create, setup and run the server
	server := api.NewServer("8080", database)

	// Setup routes and middlewares
	server.Setup()
	log.Fatal(server.Run())
}
