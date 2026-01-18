package main

import (
	"log"

	"github.com/ccrsxx/api-go/src/internal/server"
)

func main() {
	server := server.NewServer()

	log.Printf("Server starting on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("cannot start server: %s", err)
	}
}
