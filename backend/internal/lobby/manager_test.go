package lobby

import (
    "strings"
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

func TestJoinLobby(t *testing.T) {
    manager := NewManager()

    host := guest("guest-1", "Host")
    player := guest("guest-2", "Player 2")

    created, err := manager.CreateLobby(host)
    if err != nil {
        t.Fatal(err)
    }

    joined, err := manager.JoinLobby(created.Code, player)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(joined.Players) != 2 {
        t.Fatalf("expected 2 players, got %d", len(joined.Players))
    }
}

func TestJoinLobbyDuplicatePlayer(t *testing.T) {
    manager := NewManager()

    host := guest("guest-1", "Host")

    created, err := manager.CreateLobby(host)
    if err != nil {
        t.Fatal(err)
    }

    joined, err := manager.JoinLobby(created.Code, host)
    if err != nil {
        t.Fatal(err)
    }

    if len(joined.Players) != 1 {
        t.Fatalf("expected 1 player, got %d", len(joined.Players))
    }
}

func TestJoinLobbyCaseInsensitiveCode(t *testing.T) {
    manager := NewManager()

    created, err := manager.CreateLobby(guest("guest-1", "Host"))
    if err != nil {
        t.Fatal(err)
    }

    _, err = manager.JoinLobby(strings.ToLower(created.Code),
        guest("guest-2", "Player 2"))

    if err != nil {
        t.Fatalf("expected case-insensitive code, got %v", err)
    }
}

func TestJoinNonexistentLobby(t *testing.T) {
    manager := NewManager()

    _, err := manager.JoinLobby("ABC123", guest("guest-1", "Player"))

    if err != ErrLobbyNotFound {
        t.Fatalf("expected ErrLobbyNotFound, got %v", err)
    }
}

func TestLeaveLobbyTransfersHost(t *testing.T) {
    manager := NewManager()

    host := guest("guest-1", "Host")
    player2 := guest("guest-2", "Player 2")

    created, err := manager.CreateLobby(host)
    if err != nil {
        t.Fatal(err)
    }

    _, err = manager.JoinLobby(created.Code, player2)
    if err != nil {
        t.Fatal(err)
    }

    err = manager.LeaveLobby(created.Code, host)
    if err != nil {
        t.Fatal(err)
    }

    updated, err := manager.GetLobby(created.Code)
    if err != nil {
        t.Fatal(err)
    }

    if updated.Host.GuestID != player2.GuestID {
        t.Fatal("expected player 2 to become host")
    }

    if len(updated.Players) != 1 {
        t.Fatalf("expected 1 player remaining, got %d", len(updated.Players))
    }
}

func TestLeaveLobbyDeletesWhenEmpty(t *testing.T) {
    manager := NewManager()

    host := guest("guest-1", "Host")

    created, err := manager.CreateLobby(host)
    if err != nil {
        t.Fatal(err)
    }

    err = manager.LeaveLobby(created.Code, host)
    if err != nil {
        t.Fatal(err)
    }

    _, err = manager.GetLobby(created.Code)

    if err != ErrLobbyNotFound {
        t.Fatalf("expected ErrLobbyNotFound, got %v", err)
    }
}

func TestLeaveLobbyInvalidPlayer(t *testing.T) {
    manager := NewManager()

    created, err := manager.CreateLobby(guest("guest-1", "Host"))
    if err != nil {
        t.Fatal(err)
    }

    err = manager.LeaveLobby(created.Code, guest("guest-2", "Unknown"))

    if err != ErrInvalidPlayer {
        t.Fatalf("expected ErrInvalidPlayer, got %v", err)
    }
}

func TestAddMessage(t *testing.T) {
    manager := NewManager()

    host := guest("guest-1", "Host")

    created, err := manager.CreateLobby(host)
    if err != nil {
        t.Fatal(err)
    }

    updated, err := manager.AddMessage(created.Code, host, "Hello!")

    if err != nil {
        t.Fatal(err)
    }

    if len(updated.Messages) != 1 {
        t.Fatalf("expected 1 message, got %d", len(updated.Messages))
    }

    if updated.Messages[0].Message != "Hello!" {
        t.Fatal("message content does not match")
    }
}

func TestAddMessageRejectsNonMember(t *testing.T) {
    manager := NewManager()

    created, err := manager.CreateLobby(guest("guest-1", "Host"))
    if err != nil {
        t.Fatal(err)
    }

    _, err = manager.AddMessage(
        created.Code,
        guest("guest-2", "Unknown"),
        "Hello!",
    )

    if err != ErrInvalidPlayer {
        t.Fatalf("expected ErrInvalidPlayer, got %v", err)
    }
}

func TestAddMessageRejectsEmptyMessage(t *testing.T) {
    manager := NewManager()

    host := guest("guest-1", "Host")

    created, err := manager.CreateLobby(host)
    if err != nil {
        t.Fatal(err)
    }

    _, err = manager.AddMessage(created.Code, host, "   ")

    if err == nil {
        t.Fatal("expected empty message to be rejected")
    }
}

func TestMessageHistoryLimit(t *testing.T) {
    manager := NewManager()
    host := guest("guest-1", "Host")

    created, err := manager.CreateLobby(host)
    if err != nil {
        t.Fatal(err)
    }

    for i := 0; i < MaxMessages+10; i++ {
        _, err := manager.AddMessage(created.Code, host, "Hello")
        if err != nil {
            t.Fatal(err)
        }
    }

    updated, err := manager.GetLobby(created.Code)
    if err != nil {
        t.Fatal(err)
    }

    if len(updated.Messages) != MaxMessages {
        t.Fatalf("expected %d messages, got %d",
            MaxMessages, len(updated.Messages))
    }
}
