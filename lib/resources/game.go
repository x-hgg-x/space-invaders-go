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
	StateEventLevelComplete
)

// Difficulty is a game difficulty
type Difficulty float64

// List of game difficulties
const (
	DifficultyEasy   Difficulty = 0.5
	DifficultyNormal Difficulty = 1
	DifficultyHard   Difficulty = 2
)

// Game contains game resources
type Game struct {
	Events     Events
	StateEvent StateEvent
	Difficulty Difficulty
	Lives      int
	Score      int
}

// NewGame creates a new game
func NewGame(difficulty Difficulty) *Game {
	return &Game{Difficulty: difficulty, Lives: 3}
}

// Score is a game score
type Score struct {
	Score  int
	Author string
}

// ScoreTable contains highscores for a single difficulty
type ScoreTable struct {
	Scores []Score
}

// Highscores contains all highscores
type Highscores struct {
	Easy   ScoreTable
	Normal ScoreTable
	Hard   ScoreTable
}
