package api

import (
    "net/http"

    "drawio/backend/internal/api/handlers"
    "drawio/backend/internal/lobby"
)

func NewRouter(
    lobbyManager *lobby.Manager,
    resolvePlayer handlers.PlayerResolver,
) http.Handler {
    mux := http.NewServeMux()

    lobbyHandler := handlers.NewLobbyHandler(
        lobbyManager,
        resolvePlayer,
    )

    mux.HandleFunc("POST /lobbies", lobbyHandler.CreateLobby)
    mux.HandleFunc("POST /lobbies/{code}/join", lobbyHandler.JoinLobby)
    mux.HandleFunc("POST /lobbies/{code}/leave", lobbyHandler.LeaveLobby)
    mux.HandleFunc("GET /lobbies/{code}", lobbyHandler.GetLobby)

    return mux
}
