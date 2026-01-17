package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ccrsxx/rest-api-go/src/internal/config"
)

type HelloResponse struct {
	Data HelloData `json:"data"`
}

type HelloData struct {
	Message string `json:"message"`
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a non-domain request")

	w.Header().Set("Content-Type", "application/json")

	resp := HelloResponse{
		Data: HelloData{
			Message: "Hello, World!",
		},
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	config.LoadEnv()

	router := http.NewServeMux()

	router.HandleFunc("/", handleHello)

	serverPort := ":" + config.Env().Port

	server := http.Server{
		Addr:    serverPort,
		Handler: router,
	}

	log.Printf("Starting server on port %v\n", serverPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
