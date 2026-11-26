package handlers

import (
	"net/http"

	"drawio/internal/auth"
)

type GuestHandler struct {
	Sessions *auth.SessionStore
}

func NewGuestHandler(sessions *auth.SessionStore) *GuestHandler {
	return &GuestHandler{Sessions: sessions}
}

func (h *GuestHandler) CreateGuest(w http.ResponseWriter, r *http.Request) {
	player, err := h.Sessions.CreateGuest(w)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Could not create guest session")
		return
	}

	writeJSON(w, http.StatusCreated, player)
}
