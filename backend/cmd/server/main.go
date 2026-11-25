package main

import (
    "log"
    "net/http"
    "drawio/backend/internal/api"
)

func main() {
    router := api.NewRouter()
    lobbyManager := lobby.NewManager()

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
