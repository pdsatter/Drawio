package handlers

import (
	"errors"
	"net/http"

	"drawio/internal/lobby"
)

type PlayerResolver func(*http.Request) (lobby.Player, error)

type LobbyHandler struct {
	Manager       *lobby.Manager
	ResolvePlayer PlayerResolver
}

func NewLobbyHandler(
	manager *lobby.Manager,
	resolvePlayer PlayerResolver,
) *LobbyHandler {
	return &LobbyHandler{
		Manager:       manager,
		ResolvePlayer: resolvePlayer,
	}
}

// POST /lobbies
func (h *LobbyHandler) CreateLobby(w http.ResponseWriter, r *http.Request) {
	player, err := h.ResolvePlayer(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	newLobby, err := h.Manager.CreateLobby(player)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, newLobby)
}

// POST /lobbies/{code}/join
func (h *LobbyHandler) JoinLobby(w http.ResponseWriter, r *http.Request) {
	player, err := h.ResolvePlayer(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	code := r.PathValue("code")

	updatedLobby, err := h.Manager.JoinLobby(code, player)
	if err != nil {
		if errors.Is(err, lobby.ErrLobbyNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updatedLobby)
}

// POST /lobbies/{code}/leave
func (h *LobbyHandler) LeaveLobby(w http.ResponseWriter, r *http.Request) {
	player, err := h.ResolvePlayer(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	code := r.PathValue("code")

	err = h.Manager.LeaveLobby(code, player)
	if err != nil {
		if errors.Is(err, lobby.ErrLobbyNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully left lobby",
	})
}

// GET /lobbies/{code}
func (h *LobbyHandler) GetLobby(w http.ResponseWriter, r *http.Request) {
	player, err := h.ResolvePlayer(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	code := r.PathValue("code")

	currentLobby, err := h.Manager.GetLobby(code)
	if err != nil {
		if errors.Is(err, lobby.ErrLobbyNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, "Could not retrieve lobby")
		return
	}

	// Only players in the lobby can view it.
	isMember := false

	for _, member := range currentLobby.Players {
		if member.ID != nil && player.ID != nil && *member.ID == *player.ID {
			isMember = true
			break
		}

		if member.GuestID != "" && member.GuestID == player.GuestID {
			isMember = true
			break
		}
	}

	if !isMember {
		writeError(w, http.StatusForbidden, "You are not a member of this lobby")
		return
	}

	writeJSON(w, http.StatusOK, currentLobby)
}