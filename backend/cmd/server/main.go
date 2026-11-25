package main

import (
    "log"
    "net/http"

    "drawio/backend/internal/api"
    "drawio/backend/internal/lobby"
    "drawio/backend/internal/middleware"
)

func main() {
    lobbyManager := lobby.NewManager()
    sessions := auth.NewSessionStore()

    router := api.NewRouter(lobbyManager, resolvePlayer)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", middleware.CORS(router)))
}
