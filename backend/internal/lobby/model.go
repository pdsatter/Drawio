
package lobby

import "time"

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

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

type LobbyState string

const (
	StateWaiting  LobbyState = "waiting"
	StateDrawing  LobbyState = "drawing"
	StateSwapping LobbyState = "swapping"
)

const (
	DrawingDuration = 60 * time.Second
	SwappingDuration = 10 * time.Second
)

type Lobby struct {
	Host     Player    `json:"host"`
	Code     string    `json:"code"`
	Players  []Player  `json:"players"`
	Messages []Message `json:"messages"`

	// Identity keys: user:<id> or guest:<guestId>.
	Drawer *string `json:"drawer,omitempty"`

	State LobbyState `json:"state"`

	// Unix timestamps in milliseconds for frontend countdowns.
	TurnEndsAt int64 `json:"turnEndsAt"`
	SwapEndsAt int64 `json:"swapEndsAt"`

	// Calculated by the backend whenever the lobby is returned.
	TimeLeft int `json:"timeLeft"`

	TimerReset int `json:"timerReset"`
	
	// Internal deadlines. Not exposed through JSON.
	TurnDeadline time.Time `json:"-"`
	SwapDeadline time.Time `json:"-"`
}
