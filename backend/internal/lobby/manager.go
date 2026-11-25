package lobby

type Manager struct {
    Lobbies map[string]*Lobby
}

func NewManager() *Manager {
    return &Manager{
        Lobbies: make(map[string]*Lobby),
    }
}

func (m *Manager) CreateLobby() (*Lobby, error) {
    // Generate lobby code.
    // Create lobby.
    // Register host.
    // Return lobby.

    return nil, nil
}
