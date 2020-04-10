package states

import (
	"fmt"
	"image/color"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/x-hgg-x/space-invaders-go/lib/math"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const highscoreNum = 9

type highscore struct {
	difficulty resources.Difficulty
	score      int
	author     string
	position   int
	highscore  *resources.Score
}

// HighscoresState is the highscores state
type HighscoresState struct {
	highscoresMenu      []ecs.Entity
	difficulties        []resources.Difficulty
	difficultySelection int
	highscores          resources.Highscores
	newScore            *highscore
	exitTransition      states.Transition
}

//
// State interface
//

// OnPause method
func (st *HighscoresState) OnPause(world w.World) {}

// OnResume method
func (st *HighscoresState) OnResume(world w.World) {}

// OnStart method
func (st *HighscoresState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.highscoresMenu = append(st.highscoresMenu, loader.AddEntities(world, prefabs.Menu.HighscoresMenu)...)
	st.difficulties = []resources.Difficulty{resources.DifficultyEasy, resources.DifficultyNormal, resources.DifficultyHard}

	// Load highscores
	toml.DecodeFile("config/highscores.toml", &st.highscores)

	if st.newScore != nil {
		st.difficultySelection = find(st.difficulties, st.newScore.difficulty)
	} else {
		st.difficultySelection = 1
	}

	// Display highscores and check if a new highscore has been made
	if newHighscore := st.displayHighScores(world); !newHighscore {
		st.newScore = nil
	}

	// Hide game ui
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "game_score" || text.ID == "game_life" || text.ID == "game_difficulty" {
			text.Color.A = 0
		}
	}))
}

// OnStop method
func (st *HighscoresState) OnStop(world w.World) {
	// Show game ui
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "game_score" || text.ID == "game_life" || text.ID == "game_difficulty" {
			text.Color.A = 255
		}
	}))

	world.Manager.DeleteEntities(st.highscoresMenu...)
}

// Update method
func (st *HighscoresState) Update(world w.World, screen *ebiten.Image) states.Transition {
	if st.newScore != nil {
		// Set highscore author
		// Get user input
		re := regexp.MustCompile("[[:^alnum:]]")
		st.newScore.author += strings.ToUpper(re.ReplaceAllLiteralString(string(ebiten.InputChars()), ""))
		if len(st.newScore.author) > 6 {
			st.newScore.author = st.newScore.author[:6]
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(st.newScore.author) > 0 {
			st.newScore.author = st.newScore.author[:len(st.newScore.author)-1]
		}

		// Set new score text
		world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
			text := world.Components.Engine.Text.Get(entity).(*ec.Text)
			if text.ID == fmt.Sprintf("score%d", st.newScore.position+1) {
				text.Text = fmt.Sprintf("%d. %-6s %5d", st.newScore.position+1, st.newScore.author, st.newScore.score)
			}
		}))

		// Validate score
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			st.newScore.highscore.Author = st.newScore.author
			st.newScore = nil

			// Save highscores
			f, err := os.Create("config/highscores.toml")
			utils.LogError(err)
			e := toml.NewEncoder(f)
			e.Indent = "	"
			err = e.Encode(st.highscores)
			utils.LogError(err)

			st.displayHighScores(world)
		}
	} else {
		// View all scores by looping difficulties
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			st.difficultySelection = math.Mod(st.difficultySelection-1, len(st.difficulties))
			st.displayHighScores(world)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			st.difficultySelection = math.Mod(st.difficultySelection+1, len(st.difficulties))
			st.displayHighScores(world)
		}

		// Exit
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return st.exitTransition
		}
	}
	return states.Transition{}
}

func (st *HighscoresState) displayHighScores(world w.World) bool {
	var difficultyText string
	var scores *[]resources.Score
	switch st.difficulties[st.difficultySelection] {
	case resources.DifficultyEasy:
		difficultyText = "EASY"
		scores = &st.highscores.Easy.Scores
	case resources.DifficultyNormal:
		difficultyText = "NORMAL"
		scores = &st.highscores.Normal.Scores
	case resources.DifficultyHard:
		difficultyText = "HARD"
		scores = &st.highscores.Hard.Scores
	default:
		utils.LogError(fmt.Errorf("unknown difficulty: %v", st.difficulties[st.difficultySelection]))
	}

	// Sort scores
	sort.SliceStable(*scores, func(i, j int) bool {
		return (*scores)[i].Score > (*scores)[j].Score
	})

	// Get new score position
	newHighscore := false
	if st.newScore != nil {
		position := 0
		for _, score := range *scores {
			if st.newScore.score <= score.Score {
				position++
			}
		}

		if position < highscoreNum {
			(*scores) = append((*scores), resources.Score{})
			copy((*scores)[position+1:], (*scores)[position:])
			(*scores)[position] = resources.Score{Score: st.newScore.score}
			st.newScore.position = position
			st.newScore.highscore = &(*scores)[position]
			newHighscore = true
		}
	}

	// Set score texts
	for iScore := 0; iScore < highscoreNum; iScore++ {
		world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
			text := world.Components.Engine.Text.Get(entity).(*ec.Text)
			if text.ID == fmt.Sprintf("score%d", iScore+1) {
				text.Color.A = 0
				if iScore < len((*scores)) {
					text.Text = fmt.Sprintf("%d. %-6s %5d", iScore+1, (*scores)[iScore].Author, (*scores)[iScore].Score)
					text.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
					if newHighscore && iScore == st.newScore.position {
						text.Color = color.RGBA{R: 255, A: 255}
					}
				}
			}
		}))
	}

	// Set other texts
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "score_difficulty" {
			text.Text = difficultyText
		}
		if !newHighscore && (text.ID == "arrow_left" || text.ID == "arrow_right") {
			text.Color.A = 255
		}
	}))

	return newHighscore
}

func find(slice []resources.Difficulty, x resources.Difficulty) int {
	for i, e := range slice {
		if x == e {
			return i
		}
	}
	return len(slice)
}
