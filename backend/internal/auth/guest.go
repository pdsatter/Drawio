package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"sync"
	"time"

	"drawio/backend/internal/lobby"
)

const (
	cookieName = "drawio_session"
	sessionTTL = 7 * 24 * time.Hour
)

type session struct {
	Player    lobby.Player
	ExpiresAt time.Time
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]session),
	}
}

func randomID() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (s *SessionStore) CreateGuest(w http.ResponseWriter) (lobby.Player, error) {
	guestID, err := randomID()
	if err != nil {
		return lobby.Player{}, err
	}

	sessionID, err := randomID()
	if err != nil {
		return lobby.Player{}, err
	}

	player := lobby.Player{
		ID:      nil,
		GuestID: guestID,
		Name:    "Guest-" + guestID[:6],
		Icon:    "default",
	}

	s.mu.Lock()
	s.sessions[sessionID] = session{
		Player:    player,
		ExpiresAt: time.Now().Add(sessionTTL),
	}
	s.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: Set to true when using HTTPS in production.
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})

	return player, nil
}

func (s *SessionStore) ResolvePlayer(r *http.Request) (lobby.Player, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return lobby.Player{}, err
	}

	s.mu.RLock()
	currentSession, exists := s.sessions[cookie.Value]
	s.mu.RUnlock()

	if !exists || time.Now().After(currentSession.ExpiresAt) {
		return lobby.Player{}, errors.New("session expired or invalid")
	}

	return currentSession.Player, nil
}
