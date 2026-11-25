package main

import (
    "log"
    "net/http"

    "drawio/backend/internal/api"
    "drawio/backend/internal/lobby"
)

func main() {
    lobbyManager := lobby.NewManager()

    // TODO: Replace with authentication/session resolver.
    resolvePlayer := func(r *http.Request) (lobby.Player, error) {
        // Retrieve the player from their session cookie.
        return lobby.Player{}, http.ErrNoCookie
    }

    router := api.NewRouter(lobbyManager, resolvePlayer)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
