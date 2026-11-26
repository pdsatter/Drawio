package api

import (
	"net/http"

	"drawio/internal/api/handlers"
	"drawio/internal/auth"
	"drawio/internal/lobby"
)

func NewRouter(
	lobbyManager *lobby.Manager,
	sessions *auth.SessionStore,
) http.Handler {
	mux := http.NewServeMux()

	lobbyHandler := handlers.NewLobbyHandler(
		lobbyManager,
		sessions.ResolvePlayer,
	)

	guestHandler := handlers.NewGuestHandler(sessions)

	mux.HandleFunc("POST /auth/guest", guestHandler.CreateGuest)

	mux.HandleFunc("POST /lobbies", lobbyHandler.CreateLobby)
	mux.HandleFunc("POST /lobbies/{code}/join", lobbyHandler.JoinLobby)
	mux.HandleFunc("POST /lobbies/{code}/leave", lobbyHandler.LeaveLobby)
	mux.HandleFunc("GET /lobbies/{code}", lobbyHandler.GetLobby)

	return mux
}