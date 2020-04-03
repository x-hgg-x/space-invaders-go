package states

import (
	"github.com/x-hgg-x/space-invaders-go/lib/loader"
	g "github.com/x-hgg-x/space-invaders-go/lib/systems"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// GameplayState is the main game state
type GameplayState struct {
	runningAnimations []*ec.AnimationControl
}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
	// Load game and ui entities
	loader.LoadEntities("assets/metadata/entities/background.toml", world)
	loader.LoadEntities("assets/metadata/entities/level.toml", world)
	loader.LoadEntities("assets/metadata/entities/bunker.toml", world)
	loader.LoadEntities("assets/metadata/entities/ui/score.toml", world)
	loader.LoadEntities("assets/metadata/entities/ui/life.toml", world)

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
	world.Manager.DeleteAllEntities()

	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// Update method
func (st *GameplayState) Update(world w.World, screen *ebiten.Image) states.Transition {
	g.MovePlayerSystem(world)
	g.ShootPlayerBulletSystem(world)
	g.MoveBulletSystem(world)
	g.CollisionSystem(world)
	g.DeleteSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&PauseMenuState{}}}
	}
	return states.Transition{}
}
