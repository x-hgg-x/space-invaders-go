package states

import (
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// GameplayState is the main game state
type GameplayState struct{}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
	// Load game and ui entities
	loader.LoadEntities("assets/metadata/entities/background.toml", world, nil)
	loader.LoadEntities("assets/metadata/entities/level.toml", world, nil)
	loader.LoadEntities("assets/metadata/entities/bunker.toml", world, nil)
	loader.LoadEntities("assets/metadata/entities/ui/score.toml", world, nil)
	loader.LoadEntities("assets/metadata/entities/ui/life.toml", world, nil)

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// OnPause method
func (st *GameplayState) OnPause(world w.World) {
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// OnResume method
func (st *GameplayState) OnResume(world w.World) {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// OnStop method
func (st *GameplayState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()

	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// Update method
func (st *GameplayState) Update(world w.World, screen *ebiten.Image) states.Transition {
	return states.Transition{}
}
