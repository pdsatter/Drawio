package lobby

import (
    "strings"
    "sync"
    "testing"
)

func guest(id, name string) Player {
    return Player{
        GuestID: id,
        Name:    name,
        Icon:    "🐱",
    }
}

func TestCreateLobby(t *testing.T) {
    manager := NewManager()
    host := guest("guest-1", "Player 1")

    created, err := manager.CreateLobby(host)

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(created.Code) != LobbyCodeLength {
        t.Fatalf("expected code length %d, got %d",
            LobbyCodeLength, len(created.Code))
    }

    if created.Host.GuestID != host.GuestID {
        t.Fatal("expected host to be the creator")
    }

    if len(created.Players) != 1 {
        t.Fatalf("expected 1 player, got %d", len(created.Players))
    }

    if created.Players[0].GuestID != host.GuestID {
        t.Fatal("creator should be in players list")
    }
}

func TestCreateLobbyInvalidPlayer(t *testing.T) {
    manager := NewManager()

    _, err := manager.CreateLobby(Player{
        Name: "Player 1",
    })

    if err != ErrInvalidPlayer {
        t.Fatalf("expected ErrInvalidPlayer, got %v", err)
    }
}
