package states

import (
	"fmt"

	"github.com/x-hgg-x/space-invaders-go/lib/loader"

	ecs "github.com/x-hgg-x/goecs"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// GameOverState is the game over menu state
type GameOverState struct {
	gameOverMenu []ecs.Entity
	selection    int
}

//
// Menu interface
//

func (st *GameOverState) getSelection() int {
	return st.selection
}

func (st *GameOverState) setSelection(selection int) {
	st.selection = selection
}

func (st *GameOverState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Restart
		return states.Transition{Type: states.TransReplace, NewStates: []states.State{&GameplayState{}}}
	case 1:
		// Main Menu
		return states.Transition{Type: states.TransReplace, NewStates: []states.State{&MainMenuState{}}}
	case 2:
		// Exit
		return states.Transition{Type: states.TransQuit}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *GameOverState) getMenuIDs() []string {
	return []string{"restart", "main_menu", "exit"}
}

func (st *GameOverState) getCursorMenuIDs() []string {
	return []string{"cursor_restart", "cursor_main_menu", "cursor_exit"}
}

//
// State interface
//

// OnPause method
func (st *GameOverState) OnPause(world w.World) {}

// OnResume method
func (st *GameOverState) OnResume(world w.World) {}

// OnStart method
func (st *GameOverState) OnStart(world w.World) {
	st.gameOverMenu = loader.LoadEntities("assets/metadata/entities/ui/game_over_menu.toml", world)
}

// OnStop method
func (st *GameOverState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.gameOverMenu...)
}

// Update method
func (st *GameOverState) Update(world w.World, screen *ebiten.Image) states.Transition {
	return updateMenu(st, world)
}
