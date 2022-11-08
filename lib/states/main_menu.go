package states

import (
	"fmt"

	gloader "github.com/x-hgg-x/space-invaders-go/lib/loader"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"
	g "github.com/x-hgg-x/space-invaders-go/lib/systems"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// MainMenuState is the main menu state
type MainMenuState struct {
	selection int
	sound     bool
}

//
// Menu interface
//

func (st *MainMenuState) getSelection() int {
	return st.selection
}

func (st *MainMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *MainMenuState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// New game
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&DifficultyMenuState{}}}
	case 1:
		// Highscores
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HighscoresState{
			exitTransition: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}},
		}}}
	case 2:
		// Exit
		return states.Transition{Type: states.TransQuit}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *MainMenuState) getMenuIDs() []string {
	return []string{"new_game", "highscores", "exit"}
}

func (st *MainMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_new_game", "cursor_highscores", "cursor_exit"}
}

//
// State interface
//

// OnPause method
func (st *MainMenuState) OnPause(world w.World) {}

// OnResume method
func (st *MainMenuState) OnResume(world w.World) {}

// OnStart method
func (st *MainMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Game.Background)
	loader.AddEntities(world, prefabs.Menu.MainMenu)

	// Load music and sfx (at game start only)
	if world.Resources.AudioContext == nil {
		gloader.LoadSounds(world, st.sound)
	}
}

// OnStop method
func (st *MainMenuState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

// Update method
func (st *MainMenuState) Update(world w.World) states.Transition {
	g.SoundSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
	}
	return updateMenu(st, world)
}
