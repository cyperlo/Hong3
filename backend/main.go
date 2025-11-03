package main

import (
	"log"

	"github.com/chenhailong/hong3/api"
)

func main() {
	server := api.NewServer()
	log.Println("Starting Hong3 server on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}