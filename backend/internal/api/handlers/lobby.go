package handlers

import (
    "encoding/json"
    "net/http"

    "your-module/internal/lobby"
)

type LobbyHandler struct {
    Manager *lobby.Manager
}

func (h *LobbyHandler) CreateLobby(w http.ResponseWriter, r *http.Request) {
    newLobby, err := h.Manager.CreateLobby()

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(newLobby)
}
