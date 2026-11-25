package lobby

type Difficulty string
const (
    Easy   Difficulty = "easy"
    Medium Difficulty = "medium"
    Hard   Difficulty = "hard"
)

// Rewards are controlled by the backend.
// Never trust reward values supplied by the frontend.
var Rewards = map[Difficulty]int{
    Easy:   800,
    Medium: 900,
    Hard:   1000,
}

type Category struct {
    Name  string `json:"name"`
    Words []Word `json:"words"`
}

type Word struct {
    Name       string     `json:"name"`
    Difficulty Difficulty `json:"difficulty"`
}

// Player represents either a guest or authenticated player.
// Guests have a nil ID and a non-empty GuestID.
// Logged-in players have a non-empty ID and an empty GuestID.
type Player struct {
    ID      *string `json:"id"`
    GuestID string  `json:"guestId,omitempty"`
    Name    string  `json:"name"`
    Icon    string  `json:"icon"`
}

type Message struct {
    Player  Player `json:"player"`
    Message string `json:"message"`
}

type Lobby struct {
    Host     Player    `json:"host"`
    Code     string    `json:"code"`
    Players  []Player  `json:"players"`
    Messages []Message `json:"messages"`
}
