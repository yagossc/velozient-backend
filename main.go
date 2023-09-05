package main

import (
	"fmt"
	"log"
	"velozient-backend/api"
	"velozient-backend/db"
)

func main() {
	fmt.Println("Hello, Velozient Backend!")

	// Create and initialize the database
	db := db.NewMemoryDB()
	db.PopulateDB()

	// Create, setup and run the server
	server := api.NewServer("8080", db)
	server.RegisterRoutes()
	log.Fatal(server.Run())
}
