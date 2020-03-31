package states

import (
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// GameplayState is the main game state
type GameplayState struct{}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
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
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// Update method
func (st *GameplayState) Update(world w.World, screen *ebiten.Image) states.Transition {
	return states.Transition{}
}
