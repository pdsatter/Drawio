
package lobby

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"sync"
	"time"
	"unicode"
)

var (
	ErrLobbyNotFound = errors.New("lobby not found")
	ErrPlayerExists  = errors.New("player already in lobby")
	ErrInvalidPlayer = errors.New("invalid player")
	ErrAlreadyDrawer = errors.New("player is already the drawer")
	ErrNoNextDrawer  = errors.New("no eligible next drawer")
	ErrGameStarted   = errors.New("game has already started")
	ErrGameNotStarted = errors.New("game has not started")
	ErrNotDrawer     = errors.New("only the current drawer can access the animal")
	ErrCannotChat    = errors.New("current drawer cannot send messages")
)

const (
	LobbyCodeLength    = 6
	MaxMessages        = 100
	CorrectGuessPoints = 10
)

const codeCharacters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

var animals = []string{
	"alligator", "ant", "bear", "bee", "camel",
	"cat", "cheetah", "chicken", "cow", "crocodile",
	"deer", "dog", "dolphin", "donkey", "duck",
	"eagle", "elephant", "flamingo", "fox", "frog",
	"giraffe", "goat", "gorilla", "hamster", "hippo",
	"horse", "kangaroo", "koala", "lion", "monkey",
	"mouse", "octopus", "owl", "panda", "penguin",
	"pig", "rabbit", "raccoon", "shark", "sheep",
	"snake", "spider", "squirrel", "tiger", "turtle",
	"whale", "wolf", "zebra",
}

type Manager struct {
	mu      sync.RWMutex
	lobbies map[string]*Lobby

	// Private game data. Never include these values in Lobby JSON.
	animals map[string]string
	scores  map[string]map[string]int
}

func NewManager() *Manager {
	m := &Manager{
		lobbies: make(map[string]*Lobby),
		animals: make(map[string]string),
		scores:  make(map[string]map[string]int),
	}

	go m.runTimers()

	return m
}

// runTimers manages all lobby transitions.
func (m *Manager) runTimers() {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()

		now := time.Now()

		for _, currentLobby := range m.lobbies {
			switch currentLobby.State {
			case StateDrawing:
				if !currentLobby.TurnDeadline.IsZero() &&
					!now.Before(currentLobby.TurnDeadline) {
					m.beginNextDrawer(currentLobby, now)
				}

			case StateSwapping:
				if !currentLobby.SwapDeadline.IsZero() &&
					!now.Before(currentLobby.SwapDeadline) {
					m.startDrawing(currentLobby, now)
				}
			}
		}

		m.mu.Unlock()
	}
}

// CreateLobby creates a lobby in the waiting state.
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
		Host:       clonePlayer(host),
		Code:       code,
		Players:    []Player{clonePlayer(host)},
		Messages:   []Message{},
		State:      StateWaiting,
		TimerReset: int(DrawingDuration.Seconds()),
	}

	m.lobbies[code] = newLobby
	m.scores[code] = make(map[string]int)

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

	key := playerKey(player)

	for i, existing := range currentLobby.Players {
		if playerKey(existing) == key {
			currentLobby.Players[i] = clonePlayer(player)

			if playerKey(currentLobby.Host) == key {
				currentLobby.Host = clonePlayer(player)
			}

			return cloneLobby(currentLobby), nil
		}
	}

	currentLobby.Players = append(currentLobby.Players, clonePlayer(player))

	if m.scores[code] == nil {
		m.scores[code] = make(map[string]int)
	}

	return cloneLobby(currentLobby), nil
}

// StartGame starts the game. Only the host can start it.
func (m *Manager) StartGame(code string, player Player) (*Lobby, error) {
	code = normalizeCode(code)

	m.mu.Lock()
	defer m.mu.Unlock()

	currentLobby, exists := m.lobbies[code]
	if !exists {
		return nil, ErrLobbyNotFound
	}

	if !containsPlayer(currentLobby, playerKey(player)) {
		return nil, ErrInvalidPlayer
	}

	if playerKey(currentLobby.Host) != playerKey(player) {
		return nil, ErrInvalidPlayer
	}

	if currentLobby.State != StateWaiting {
		return nil, ErrGameStarted
	}

	if len(currentLobby.Players) == 0 {
		return nil, ErrInvalidPlayer
	}

	index, err := randomIndex(len(currentLobby.Players))
	if err != nil {
		return nil, err
	}

	currentLobby.Drawer = stringPointer(
		playerKey(currentLobby.Players[index]),
	)

	m.startDrawing(currentLobby, time.Now())

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

	key := playerKey(player)
	playerIndex := -1

	for i, existing := range currentLobby.Players {
		if playerKey(existing) == key {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		return ErrInvalidPlayer
	}

	wasDrawer := currentLobby.Drawer != nil &&
		*currentLobby.Drawer == key

	wasHost := playerKey(currentLobby.Host) == key

	currentLobby.Players = append(
		currentLobby.Players[:playerIndex],
		currentLobby.Players[playerIndex+1:]...,
	)

	if len(currentLobby.Players) == 0 {
		delete(m.lobbies, code)
		delete(m.animals, code)
		delete(m.scores, code)
		return nil
	}

	if wasHost {
		currentLobby.Host = clonePlayer(currentLobby.Players[0])
	}

	if wasDrawer {
		if currentLobby.State == StateWaiting {
			currentLobby.Drawer = nil
			return nil
		}

		// Keep the previous position in the rotation. The next player
		// now occupies this index after the departing player is removed.
		previousIndex := (playerIndex - 1 + len(currentLobby.Players)) %
			len(currentLobby.Players)

		currentLobby.Drawer = stringPointer(
			playerKey(currentLobby.Players[previousIndex]),
		)

		m.beginNextDrawer(currentLobby, time.Now())
	}

	return nil
}

// GetLobby retrieves a lobby by its code.
// The secret animal is never stored in the Lobby response.
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

// GetAnimal returns the secret animal only to the current drawer.
func (m *Manager) GetAnimal(code string, player Player) (string, error) {
	code = normalizeCode(code)

	m.mu.RLock()
	defer m.mu.RUnlock()

	currentLobby, exists := m.lobbies[code]
	if !exists {
		return "", ErrLobbyNotFound
	}

	key := playerKey(player)

	if !containsPlayer(currentLobby, key) {
		return "", ErrInvalidPlayer
	}

	if currentLobby.State != StateDrawing ||
		currentLobby.Drawer == nil ||
		*currentLobby.Drawer != key {
		return "", ErrNotDrawer
	}

	animal, exists := m.animals[code]
	if !exists {
		return "", ErrGameNotStarted
	}

	return animal, nil
}

// SetDrawer selects a player and starts the swapping state.
func (m *Manager) SetDrawer(code string, player Player) error {
	code = normalizeCode(code)

	m.mu.Lock()
	defer m.mu.Unlock()

	currentLobby, exists := m.lobbies[code]
	if !exists {
		return ErrLobbyNotFound
	}

	if currentLobby.State == StateWaiting {
		return ErrGameNotStarted
	}

	key := playerKey(player)

	if !containsPlayer(currentLobby, key) {
		return ErrInvalidPlayer
	}

	if currentLobby.Drawer != nil && *currentLobby.Drawer == key {
		return ErrAlreadyDrawer
	}

	m.beginSwapTo(currentLobby, key, time.Now())

	return nil
}

// ChangeDrawer begins the next drawer transition.
func (m *Manager) ChangeDrawer(code string) error {
	code = normalizeCode(code)

	m.mu.Lock()
	defer m.mu.Unlock()

	currentLobby, exists := m.lobbies[code]
	if !exists {
		return ErrLobbyNotFound
	}

	if currentLobby.State != StateDrawing {
		return ErrGameNotStarted
	}

	m.beginNextDrawer(currentLobby, time.Now())

	return nil
}

// beginNextDrawer selects the next player in lobby order.
func (m *Manager) beginNextDrawer(currentLobby *Lobby, now time.Time) {
	if len(currentLobby.Players) == 0 {
		return
	}

	if len(currentLobby.Players) == 1 {
		currentLobby.Drawer = stringPointer(
			playerKey(currentLobby.Players[0]),
		)
		m.startDrawing(currentLobby, now)
		return
	}

	currentIndex := -1

	if currentLobby.Drawer != nil {
		for i, player := range currentLobby.Players {
			if playerKey(player) == *currentLobby.Drawer {
				currentIndex = i
				break
			}
		}
	}

	for offset := 1; offset <= len(currentLobby.Players); offset++ {
		index := (currentIndex + offset) % len(currentLobby.Players)
		key := playerKey(currentLobby.Players[index])

		if currentLobby.Drawer == nil || key != *currentLobby.Drawer {
			m.beginSwapTo(currentLobby, key, now)
			return
		}
	}
}

// beginSwapTo selects the incoming drawer and starts the transition.
func (m *Manager) beginSwapTo(
	currentLobby *Lobby,
	nextDrawer string,
	now time.Time,
) {
	currentLobby.Drawer = stringPointer(nextDrawer)
	currentLobby.State = StateSwapping

	currentLobby.TurnDeadline = time.Time{}
	currentLobby.SwapDeadline = now.Add(SwappingDuration)

	currentLobby.TurnEndsAt = 0
	currentLobby.SwapEndsAt = currentLobby.SwapDeadline.UnixMilli()
	currentLobby.TimeLeft = int(SwappingDuration.Seconds())
}

// startDrawing begins a new drawing turn and selects a new secret animal.
func (m *Manager) startDrawing(currentLobby *Lobby, now time.Time) {
	currentLobby.State = StateDrawing

	currentLobby.TurnDeadline = now.Add(DrawingDuration)
	currentLobby.SwapDeadline = time.Time{}

	currentLobby.TurnEndsAt = currentLobby.TurnDeadline.UnixMilli()
	currentLobby.SwapEndsAt = 0
	currentLobby.TimeLeft = int(DrawingDuration.Seconds())
	currentLobby.TimerReset = int(DrawingDuration.Seconds())

	animal, err := randomAnimal()
	if err == nil {
		m.animals[currentLobby.Code] = animal
	}
}

// IsPlayerCurrentDrawer checks the current drawer's identity.
func (m *Manager) IsPlayerCurrentDrawer(
	code string,
	player Player,
) (bool, error) {
	code = normalizeCode(code)

	m.mu.RLock()
	defer m.mu.RUnlock()

	currentLobby, exists := m.lobbies[code]
	if !exists {
		return false, ErrLobbyNotFound
	}

	if currentLobby.Drawer == nil {
		return false, nil
	}

	return playerKey(player) == *currentLobby.Drawer, nil
}

// CanPlayerSendChat checks whether a player can send messages.
func (m *Manager) CanPlayerSendChat(code string, player Player) bool {
	code = normalizeCode(code)

	m.mu.RLock()
	defer m.mu.RUnlock()

	currentLobby, exists := m.lobbies[code]
	if !exists {
		return false
	}

	if !containsPlayer(currentLobby, playerKey(player)) {
		return false
	}

	if len(currentLobby.Players) <= 1 {
		return true
	}

	return currentLobby.Drawer == nil ||
		playerKey(player) != *currentLobby.Drawer
}

// SendMessage adds a message and checks whether it is a correct guess.
// Correct guesses are still added to chat, but the secret is never returned.
func (m *Manager) SendMessage(
	code string,
	player Player,
	message string,
) (*Lobby, bool, int, error) {
	code = normalizeCode(code)
	message = strings.TrimSpace(message)

	if message == "" {
		return nil, false, 0, errors.New("message cannot be empty")
	}

	if err := validatePlayer(player); err != nil {
		return nil, false, 0, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	currentLobby, exists := m.lobbies[code]
	if !exists {
		return nil, false, 0, ErrLobbyNotFound
	}

	key := playerKey(player)

	if !containsPlayer(currentLobby, key) {
		return nil, false, 0, ErrInvalidPlayer
	}

	if len(currentLobby.Players) > 1 &&
		currentLobby.Drawer != nil &&
		*currentLobby.Drawer == key {
		return nil, false, 0, ErrCannotChat
	}

	correct := false
	points := 0

	if currentLobby.State == StateDrawing {
		animal, hasAnimal := m.animals[code]

		if hasAnimal && normalizeGuess(message) == normalizeGuess(animal) {
			correct = true
			points = CorrectGuessPoints

			if m.scores[code] == nil {
				m.scores[code] = make(map[string]int)
			}

			m.scores[code][key] += points
		}
	}

	currentLobby.Messages = append(currentLobby.Messages, Message{
		Player:  clonePlayer(player),
		Message: message,
	})

	if len(currentLobby.Messages) > MaxMessages {
		currentLobby.Messages = currentLobby.Messages[1:]
	}

	return cloneLobby(currentLobby), correct, points, nil
}

// AddMessage is retained for compatibility with existing callers.
// New API code should use SendMessage so guesses are checked.
func (m *Manager) AddMessage(
	code string,
	player Player,
	message string,
) (*Lobby, error) {
	updatedLobby, _, _, err := m.SendMessage(code, player, message)
	return updatedLobby, err
}

// GetPlayerScore returns a player's server-side score.
func (m *Manager) GetPlayerScore(code string, player Player) (int, error) {
	code = normalizeCode(code)

	m.mu.RLock()
	defer m.mu.RUnlock()

	currentLobby, exists := m.lobbies[code]
	if !exists {
		return 0, ErrLobbyNotFound
	}

	if !containsPlayer(currentLobby, playerKey(player)) {
		return 0, ErrInvalidPlayer
	}

	return m.scores[code][playerKey(player)], nil
}

// generateCode creates a unique lobby code.
func (m *Manager) generateCode() (string, error) {
	for {
		var code strings.Builder

		for i := 0; i < LobbyCodeLength; i++ {
			n, err := rand.Int(
				rand.Reader,
				big.NewInt(int64(len(codeCharacters))),
			)
			if err != nil {
				return "", err
			}

			code.WriteByte(codeCharacters[n.Int64()])
		}

		result := code.String()

		if _, exists := m.lobbies[result]; !exists {
			return result, nil
		}
	}
}

func randomIndex(length int) (int, error) {
	if length <= 0 {
		return 0, ErrInvalidPlayer
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}

func randomAnimal() (string, error) {
	index, err := randomIndex(len(animals))
	if err != nil {
		return "", err
	}

	return animals[index], nil
}

func normalizeGuess(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.TrimFunc(value, unicode.IsPunct)
	return strings.Join(strings.Fields(value), " ")
}

func normalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func validatePlayer(player Player) error {
	hasAccountID := player.ID != nil &&
		strings.TrimSpace(*player.ID) != ""

	hasGuestID := strings.TrimSpace(player.GuestID) != ""

	if hasAccountID == hasGuestID {
		return ErrInvalidPlayer
	}

	if strings.TrimSpace(player.Name) == "" {
		return ErrInvalidPlayer
	}

	return nil
}

// playerKey creates a unique identity key.
func playerKey(player Player) string {
	if player.ID != nil && strings.TrimSpace(*player.ID) != "" {
		return "user:" + strings.TrimSpace(*player.ID)
	}

	return "guest:" + strings.TrimSpace(player.GuestID)
}

func stringPointer(value string) *string {
	return &value
}

func containsPlayer(currentLobby *Lobby, key string) bool {
	for _, player := range currentLobby.Players {
		if playerKey(player) == key {
			return true
		}
	}

	return false
}

// cloneLobby prevents callers from modifying manager-owned data.
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

	if source.Drawer != nil {
		result.Drawer = stringPointer(*source.Drawer)
	}

	now := time.Now()

	switch source.State {
	case StateDrawing:
		result.TimeLeft = remainingSeconds(source.TurnDeadline, now)
	case StateSwapping:
		result.TimeLeft = remainingSeconds(source.SwapDeadline, now)
	}

	return &result
}

func remainingSeconds(deadline, now time.Time) int {
	if deadline.IsZero() || !deadline.After(now) {
		return 0
	}

	remaining := deadline.Sub(now)

	return int((remaining + time.Second - 1) / time.Second)
}

func clonePlayer(source Player) Player {
	result := source

	if source.ID != nil {
		result.ID = stringPointer(*source.ID)
	}

	return result
}
