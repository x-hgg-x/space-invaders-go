package resources

// LifeEvent is triggered when the player lose a life
type LifeEvent struct{}

// ScoreEvent is triggered when the score changes
type ScoreEvent struct {
	Score int
}

// Events contains game events for communication between game systems
type Events struct {
	LifeEvents  []LifeEvent
	ScoreEvents []ScoreEvent
}

// StateEvent is an event for game progression
type StateEvent int

// List of game progression events
const (
	StateEventNone StateEvent = iota
	StateEventDeath
	StateEventGameOver
	StateEventLevelComplete
)

// Game contains game resources
type Game struct {
	Events     Events
	StateEvent StateEvent
	Lives      int
	Score      int
}

// NewGame creates a new game
func NewGame() *Game {
	return &Game{Lives: 3}
}
