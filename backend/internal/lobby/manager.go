package lobby

import (
    "crypto/rand"
    "errors"
    "strings"
    "sync"
)

var (
    ErrLobbyNotFound = errors.New("lobby not found")
    ErrPlayerExists  = errors.New("player already in lobby")
    ErrInvalidPlayer = errors.New("invalid player")
)

const (
    LobbyCodeLength = 6
    MaxMessages     = 100
)

const codeCharacters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type Manager struct {
    mu      sync.RWMutex
    lobbies map[string]*Lobby
}

func NewManager() *Manager {
    return &Manager{
        lobbies: make(map[string]*Lobby),
    }
}

// CreateLobby creates a lobby and automatically adds the host.
func (m *Manager) CreateLobby(host Player) (*Lobby, error) {
    if err := validatePlayer(host); err != nil {
        return nil, err
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    code, err := m.generateCode()
    if err != nil {
        return nil, err
    }

    newLobby := &Lobby{
        Host:     host,
        Code:     code,
        Players:  []Player{host},
        Messages: []Message{},
    }

    m.lobbies[code] = newLobby

    return cloneLobby(newLobby), nil
}

// JoinLobby adds a player to an existing lobby.
func (m *Manager) JoinLobby(code string, player Player) (*Lobby, error) {
    if err := validatePlayer(player); err != nil {
        return nil, err
    }

    code = normalizeCode(code)

    m.mu.Lock()
    defer m.mu.Unlock()

    currentLobby, exists := m.lobbies[code]
    if !exists {
        return nil, ErrLobbyNotFound
    }

    // Prevent the same player from joining twice.
    for i, existing := range currentLobby.Players {
        if playerKey(existing) == playerKey(player) {
            // Refresh their displayed information.
            currentLobby.Players[i] = player

            if playerKey(currentLobby.Host) == playerKey(player) {
                currentLobby.Host = player
            }

            return cloneLobby(currentLobby), nil
        }
    }

    currentLobby.Players = append(currentLobby.Players, player)

    return cloneLobby(currentLobby), nil
}

// LeaveLobby removes a player from a lobby.
func (m *Manager) LeaveLobby(code string, player Player) error {
    code = normalizeCode(code)

    m.mu.Lock()
    defer m.mu.Unlock()

    currentLobby, exists := m.lobbies[code]
    if !exists {
        return ErrLobbyNotFound
    }

    playerID := playerKey(player)

    playerIndex := -1

    for i, existing := range currentLobby.Players {
        if playerKey(existing) == playerID {
            playerIndex = i
            break
        }
    }

    if playerIndex == -1 {
        return ErrInvalidPlayer
    }

    // Remove the player.
    currentLobby.Players = append(
        currentLobby.Players[:playerIndex],
        currentLobby.Players[playerIndex+1:]...,
    )

    // Delete the lobby if nobody remains.
    if len(currentLobby.Players) == 0 {
        delete(m.lobbies, code)
        return nil
    }

    // Transfer host ownership if necessary.
    if playerKey(currentLobby.Host) == playerID {
        currentLobby.Host = currentLobby.Players[0]
    }

    return nil
}

// GetLobby retrieves a lobby by its code.
func (m *Manager) GetLobby(code string) (*Lobby, error) {
    code = normalizeCode(code)

    m.mu.RLock()
    defer m.mu.RUnlock()

    currentLobby, exists := m.lobbies[code]
    if !exists {
        return nil, ErrLobbyNotFound
    }

    return cloneLobby(currentLobby), nil
}

// AddMessage adds a chat message to a lobby.
func (m *Manager) AddMessage(code string, player Player, message string) (*Lobby, error) {
    code = normalizeCode(code)
    message = strings.TrimSpace(message)

    if message == "" {
        return nil, errors.New("message cannot be empty")
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    currentLobby, exists := m.lobbies[code]
    if !exists {
        return nil, ErrLobbyNotFound
    }

    // Only players in the lobby can send messages.
    found := false

    for _, existing := range currentLobby.Players {
        if playerKey(existing) == playerKey(player) {
            found = true
            break
        }
    }

    if !found {
        return nil, ErrInvalidPlayer
    }

    currentLobby.Messages = append(currentLobby.Messages, Message{
        Player:  player,
        Message: message,
    })

    // Keep chat history bounded.
    if len(currentLobby.Messages) > MaxMessages {
        currentLobby.Messages = currentLobby.Messages[1:]
    }

    return cloneLobby(currentLobby), nil
}

// generateCode creates a unique lobby code.
// Must be called while holding the manager lock.
func (m *Manager) generateCode() (string, error) {
    for {
        randomBytes := make([]byte, LobbyCodeLength)

        if _, err := rand.Read(randomBytes); err != nil {
            return "", err
        }

        var code strings.Builder

        for _, b := range randomBytes {
            code.WriteByte(codeCharacters[int(b)%len(codeCharacters)])
        }

        result := code.String()

        if _, exists := m.lobbies[result]; !exists {
            return result, nil
        }
    }
}

// normalizeCode ensures lobby codes are case-insensitive.
func normalizeCode(code string) string {
    return strings.ToUpper(strings.TrimSpace(code))
}

// validatePlayer ensures the player has a valid identity.
func validatePlayer(player Player) error {
    hasAccountID := player.ID != nil && *player.ID != ""
    hasGuestID := player.GuestID != ""

    if hasAccountID == hasGuestID {
        return ErrInvalidPlayer
    }

    if strings.TrimSpace(player.Name) == "" {
        return ErrInvalidPlayer
    }

    return nil
}

// playerKey returns a consistent identity for comparing players.
func playerKey(player Player) string {
    if player.ID != nil && *player.ID != "" {
        return "user:" + *player.ID
    }

    return "guest:" + player.GuestID
}

// cloneLobby returns a copy so callers cannot modify internal state.
func cloneLobby(source *Lobby) *Lobby {
    result := *source

    result.Host = clonePlayer(source.Host)

    result.Players = make([]Player, len(source.Players))
    for i, player := range source.Players {
        result.Players[i] = clonePlayer(player)
    }

    result.Messages = make([]Message, len(source.Messages))
    for i, message := range source.Messages {
        result.Messages[i] = Message{
            Player:  clonePlayer(message.Player),
            Message: message.Message,
        }
    }

    return &result
}

func clonePlayer(source Player) Player {
    result := source

    if source.ID != nil {
        id := *source.ID
        result.ID = &id
    }

    return result
}
