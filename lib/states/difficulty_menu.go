package states

import (
	"fmt"

	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// DifficultyMenuState is the difficulty menu state
type DifficultyMenuState struct {
	difficultyMenu []ecs.Entity
	selection      int
}

//
// Menu interface
//

func (st *DifficultyMenuState) getSelection() int {
	return st.selection
}

func (st *DifficultyMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *DifficultyMenuState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Easy
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameplayState{game: resources.NewGame(resources.DifficultyEasy)}}}
	case 1:
		// Normal
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameplayState{game: resources.NewGame(resources.DifficultyNormal)}}}
	case 2:
		// Hard
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameplayState{game: resources.NewGame(resources.DifficultyHard)}}}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *DifficultyMenuState) getMenuIDs() []string {
	return []string{"easy", "normal", "hard"}
}

func (st *DifficultyMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_easy", "cursor_normal", "cursor_hard"}
}

//
// State interface
//

// OnPause method
func (st *DifficultyMenuState) OnPause(world w.World) {}

// OnResume method
func (st *DifficultyMenuState) OnResume(world w.World) {}

// OnStart method
func (st *DifficultyMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.difficultyMenu = append(st.difficultyMenu, loader.AddEntities(world, prefabs.Game.Background)...)
	st.difficultyMenu = append(st.difficultyMenu, loader.AddEntities(world, prefabs.Menu.DifficultyMenu)...)

	// Default difficulty is normal
	st.setSelection(1)
}

// OnStop method
func (st *DifficultyMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.difficultyMenu...)
}

// Update method
func (st *DifficultyMenuState) Update(world w.World, screen *ebiten.Image) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}
	return updateMenu(st, world)
}
