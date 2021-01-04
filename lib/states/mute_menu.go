package states

import (
	"fmt"

	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// MuteMenuState is the mute menu state
type MuteMenuState struct {
	muteMenu  []ecs.Entity
	selection int
}

//
// Menu interface
//

func (st *MuteMenuState) getSelection() int {
	return st.selection
}

func (st *MuteMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *MuteMenuState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Play with sound
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{sound: true}}}
	case 1:
		// Play muted
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{sound: false}}}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *MuteMenuState) getMenuIDs() []string {
	return []string{"play_sound", "play_muted"}
}

func (st *MuteMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_play_sound", "cursor_play_muted"}
}

//
// State interface
//

// OnPause method
func (st *MuteMenuState) OnPause(world w.World) {}

// OnResume method
func (st *MuteMenuState) OnResume(world w.World) {}

// OnStart method
func (st *MuteMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.muteMenu = append(st.muteMenu, loader.AddEntities(world, prefabs.Menu.MuteMenu)...)
}

// OnStop method
func (st *MuteMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.muteMenu...)
}

// Update method
func (st *MuteMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
	}
	return updateMenu(st, world)
}
