package states

import (
	"fmt"

	"github.com/x-hgg-x/space-invaders-go/lib/loader"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// LevelCompleteState is the level complete menu state
type LevelCompleteState struct {
	game              *resources.Game
	levelCompleteMenu []ecs.Entity
	selection         int
}

//
// Menu interface
//

func (st *LevelCompleteState) getSelection() int {
	return st.selection
}

func (st *LevelCompleteState) setSelection(selection int) {
	st.selection = selection
}

func (st *LevelCompleteState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Continue
		return states.Transition{Type: states.TransReplace, NewStates: []states.State{&GameplayState{game: st.game}}}
	case 1:
		// Main Menu
		return states.Transition{Type: states.TransReplace, NewStates: []states.State{&MainMenuState{}}}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *LevelCompleteState) getMenuIDs() []string {
	return []string{"continue", "main_menu"}
}

func (st *LevelCompleteState) getCursorMenuIDs() []string {
	return []string{"cursor_continue", "cursor_main_menu"}
}

//
// State interface
//

// OnPause method
func (st *LevelCompleteState) OnPause(world w.World) {}

// OnResume method
func (st *LevelCompleteState) OnResume(world w.World) {}

// OnStart method
func (st *LevelCompleteState) OnStart(world w.World) {
	st.levelCompleteMenu = loader.LoadEntities("assets/metadata/entities/ui/level_complete_menu.toml", world)
}

// OnStop method
func (st *LevelCompleteState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.levelCompleteMenu...)
}

// Update method
func (st *LevelCompleteState) Update(world w.World, screen *ebiten.Image) states.Transition {
	return updateMenu(st, world)
}
