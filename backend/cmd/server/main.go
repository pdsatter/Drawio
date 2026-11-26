package main

import (
	"log"
	"net/http"

	"drawio/internal/api"
	"drawio/internal/auth"
	"drawio/internal/lobby"
	"drawio/internal/middleware"
)

func main() {
	lobbyManager := lobby.NewManager()
	sessions := auth.NewSessionStore()

	router := api.NewRouter(lobbyManager, sessions)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.CORS(router)))
}