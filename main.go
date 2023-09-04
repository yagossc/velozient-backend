package main

import (
	"fmt"
	"log"
	"velozient-backend/api"
)

func main() {
	fmt.Println("Hello, Velozient Backend!")
	server := api.NewServer("8080")
	server.RegisterRoutes()
	log.Fatal(server.Run())
}
