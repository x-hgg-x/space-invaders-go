package states

import (
	"fmt"
	"math/rand"
	"time"

	gloader "github.com/x-hgg-x/space-invaders-go/lib/loader"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"
	g "github.com/x-hgg-x/space-invaders-go/lib/systems"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// GameplayState is the main game state
type GameplayState struct {
	game              *resources.Game
	runningAnimations []*ec.AnimationControl
}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
	// Init rand seed
	rand.Seed(time.Now().UnixNano())

	// Load game and ui entities
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Game.Background)
	loader.AddEntities(world, prefabs.Game.Alien)
	loader.AddEntities(world, prefabs.Game.Player)
	loader.AddEntities(world, prefabs.Game.PlayerLine)
	scoreEntity := loader.AddEntities(world, prefabs.Game.Score)
	lifeEntity := loader.AddEntities(world, prefabs.Game.Life)
	difficultyEntity := loader.AddEntities(world, prefabs.Game.Difficulty)

	// Load bunkers
	gloader.LoadBunkers(world)

	// Set game
	world.Resources.Game = st.game

	// Set score text
	for iEntity := range scoreEntity {
		world.Components.Engine.Text.Get(scoreEntity[iEntity]).(*ec.Text).Text = fmt.Sprintf("SCORE: %d", st.game.Score)
	}

	// Set life text
	for iEntity := range lifeEntity {
		world.Components.Engine.Text.Get(lifeEntity[iEntity]).(*ec.Text).Text = fmt.Sprintf("LIVES: %d", st.game.Lives)
	}

	// Set difficulty text
	var difficulty string
	switch st.game.Difficulty {
	case resources.DifficultyEasy:
		difficulty = "EASY"
	case resources.DifficultyNormal:
		difficulty = "NORMAL"
	case resources.DifficultyHard:
		difficulty = "HARD"
	default:
		utils.LogError(fmt.Errorf("unknown difficulty: %v", st.game.Difficulty))
	}
	for iEntity := range difficultyEntity {
		world.Components.Engine.Text.Get(difficultyEntity[iEntity]).(*ec.Text).Text = difficulty
	}

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// OnPause method
func (st *GameplayState) OnPause(world w.World) {
	// Pause running animations
	st.runningAnimations = []*ec.AnimationControl{}
	world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
		animationControl := world.Components.Engine.AnimationControl.Get(entity).(*ec.AnimationControl)
		if animationControl.GetState().Type == ec.ControlStateRunning {
			animationControl.Command.Type = ec.AnimationCommandPause
			st.runningAnimations = append(st.runningAnimations, animationControl)
		}
	}))

	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// OnResume method
func (st *GameplayState) OnResume(world w.World) {
	// Resume running animations
	for _, animationControl := range st.runningAnimations {
		animationControl.Command.Type = ec.AnimationCommandStart
	}
	st.runningAnimations = []*ec.AnimationControl{}

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// OnStop method
func (st *GameplayState) OnStop(world w.World) {
	world.Resources.Game = nil
	world.Manager.DeleteAllEntities()

	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// Update method
func (st *GameplayState) Update(world w.World, screen *ebiten.Image) states.Transition {
	g.MovePlayerSystem(world)
	g.SpawnAlienMasterSystem(world)
	g.MoveAlienMasterSystem(world)
	g.MoveAlienSystem(world)
	g.ShootPlayerBulletSystem(world)
	g.ShootEnemyBulletSystem(world)
	g.MoveBulletSystem(world)
	g.CollisionSystem(world)
	g.LifeSystem(world)
	g.ScoreSystem(world)
	g.DeleteSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&PauseMenuState{}}}
	}

	gameResources := world.Resources.Game.(*resources.Game)
	switch gameResources.StateEvent {
	case resources.StateEventDeath:
		gameResources.StateEvent = resources.StateEventNone
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DeathState{}}}
	case resources.StateEventGameOver:
		gameResources.StateEvent = resources.StateEventNone
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&GameOverState{difficulty: gameResources.Difficulty}}}
	case resources.StateEventLevelComplete:
		gameResources.StateEvent = resources.StateEventNone
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&LevelCompleteState{game: gameResources}}}
	}

	return states.Transition{}
}
