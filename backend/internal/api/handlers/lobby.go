
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"drawio/internal/lobby"
)

type PlayerResolver func(*http.Request) (lobby.Player, error)

type LobbyHandler struct {
	Manager       *lobby.Manager
	ResolvePlayer PlayerResolver
}

type messageRequest struct {
	Message string `json:"message"`
}

type messageResponse struct {
	Lobby   *lobby.Lobby `json:"lobby"`
	Correct bool         `json:"correct"`
	Points  int          `json:"points"`
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

func (h *LobbyHandler) resolveRequestPlayer(
	w http.ResponseWriter,
	r *http.Request,
) (lobby.Player, bool) {
	player, err := h.ResolvePlayer(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return lobby.Player{}, false
	}

	return player, true
}

// playerKey identifies a player by account ID or guest ID.
// Names are never used as identifiers.
func playerKey(player lobby.Player) string {
	if player.ID != nil && strings.TrimSpace(*player.ID) != "" {
		return "user:" + strings.TrimSpace(*player.ID)
	}

	return "guest:" + strings.TrimSpace(player.GuestID)
}

func isLobbyMember(currentLobby *lobby.Lobby, player lobby.Player) bool {
	key := playerKey(player)

	for _, member := range currentLobby.Players {
		if playerKey(member) == key {
			return true
		}
	}

	return false
}

// requireLobbyMember retrieves a lobby and verifies membership.
func (h *LobbyHandler) requireLobbyMember(
	w http.ResponseWriter,
	r *http.Request,
	player lobby.Player,
) (*lobby.Lobby, bool) {
	currentLobby, err := h.Manager.GetLobby(r.PathValue("code"))

	if err != nil {
		if errors.Is(err, lobby.ErrLobbyNotFound) {
			writeError(w, http.StatusNotFound, "Lobby not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Could not retrieve lobby")
		}

		return nil, false
	}

	if !isLobbyMember(currentLobby, player) {
		writeError(w, http.StatusForbidden, "You are not a member of this lobby")
		return nil, false
	}

	return currentLobby, true
}

// POST /lobbies
func (h *LobbyHandler) CreateLobby(w http.ResponseWriter, r *http.Request) {
	player, ok := h.resolveRequestPlayer(w, r)
	if !ok {
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
	player, ok := h.resolveRequestPlayer(w, r)
	if !ok {
		return
	}

	updatedLobby, err := h.Manager.JoinLobby(
		r.PathValue("code"),
		player,
	)

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
	player, ok := h.resolveRequestPlayer(w, r)
	if !ok {
		return
	}

	err := h.Manager.LeaveLobby(
		r.PathValue("code"),
		player,
	)

	if err != nil {
		if errors.Is(err, lobby.ErrLobbyNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, lobby.ErrInvalidPlayer) {
			writeError(w, http.StatusForbidden, err.Error())
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
	player, ok := h.resolveRequestPlayer(w, r)
	if !ok {
		return
	}

	currentLobby, ok := h.requireLobbyMember(w, r, player)
	if !ok {
		return
	}

	// The secret animal is never part of this response.
	writeJSON(w, http.StatusOK, currentLobby)
}

// POST /lobbies/{code}/start
func (h *LobbyHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	player, ok := h.resolveRequestPlayer(w, r)
	if !ok {
		return
	}

	currentLobby, ok := h.requireLobbyMember(w, r, player)
	if !ok {
		return
	}

	if playerKey(currentLobby.Host) != playerKey(player) {
		writeError(w, http.StatusForbidden, "Only the host can start the game")
		return
	}

	startedLobby, err := h.Manager.StartGame(
		r.PathValue("code"),
		player,
	)

	if err != nil {
		switch {
		case errors.Is(err, lobby.ErrLobbyNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, lobby.ErrInvalidPlayer):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, startedLobby)
}

// GET /lobbies/{code}/animal
// Only the current drawer can retrieve the secret animal.
func (h *LobbyHandler) GetAnimal(w http.ResponseWriter, r *http.Request) {
	player, ok := h.resolveRequestPlayer(w, r)
	if !ok {
		return
	}

	_, ok = h.requireLobbyMember(w, r, player)
	if !ok {
		return
	}

	animal, err := h.Manager.GetAnimal(
		r.PathValue("code"),
		player,
	)

	if err != nil {
		switch {
		case errors.Is(err, lobby.ErrLobbyNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, lobby.ErrInvalidPlayer):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, lobby.ErrNotDrawer):
			writeError(w, http.StatusForbidden, "Only the current drawer can access the animal")
		default:
			writeError(w, http.StatusConflict, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"animal": animal,
	})
}

// POST /lobbies/{code}/messages
// A message is both a chat message and a possible guess.
func (h *LobbyHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	player, ok := h.resolveRequestPlayer(w, r)
	if !ok {
		return
	}

	_, ok = h.requireLobbyMember(w, r, player)
	if !ok {
		return
	}

	var request messageRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	request.Message = strings.TrimSpace(request.Message)

	if request.Message == "" {
		writeError(w, http.StatusBadRequest, "Message cannot be empty")
		return
	}

	updatedLobby, correct, points, err := h.Manager.SendMessage(
		r.PathValue("code"),
		player,
		request.Message,
	)

	if err != nil {
		switch {
		case errors.Is(err, lobby.ErrLobbyNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, lobby.ErrInvalidPlayer):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, lobby.ErrCannotChat):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, messageResponse{
		Lobby:   updatedLobby,
		Correct: correct,
		Points:  points,
	})
}
