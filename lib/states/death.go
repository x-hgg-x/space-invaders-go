package states

import (
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	"github.com/x-hgg-x/space-invaders-go/lib/loader"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// DeathState is the player death state
type DeathState struct {
	playerEntity    ecs.Entity
	playerAnimation *ec.AnimationControl
}

// OnStart method
func (st *DeathState) OnStart(world w.World) {
	// Restart player death animation
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Player, gameComponents.Controllable, world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(playerEntity ecs.Entity) {
		st.playerEntity = playerEntity
		st.playerAnimation = world.Components.Engine.AnimationControl.Get(playerEntity).(*ec.AnimationControl)
		st.playerAnimation.Command.Type = ec.AnimationCommandStart
	}))

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// OnPause method
func (st *DeathState) OnPause(world w.World) {
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// OnResume method
func (st *DeathState) OnResume(world w.World) {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// OnStop method
func (st *DeathState) OnStop(world w.World) {
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// Update method
func (st *DeathState) Update(world w.World, screen *ebiten.Image) states.Transition {
	if st.playerAnimation.GetState().Type == ec.ControlStateDone {
		world.Manager.DeleteEntity(st.playerEntity)
		resurrectPlayer(world)
		return states.Transition{Type: states.TransPop}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&PauseMenuState{}}}
	}

	return states.Transition{}
}

func resurrectPlayer(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	// Reset position of remaining aliens
	world.Manager.Join(gameComponents.Alien, gameComponents.AlienMaster.Not(), world.Components.Engine.Transform).Visit(ecs.Visit(func(alienEntity ecs.Entity) {
		alien := gameComponents.Alien.Get(alienEntity).(*gc.Alien)
		alienTranslation := &world.Components.Engine.Transform.Get(alienEntity).(*ec.Transform).Translation

		alienTranslation.X -= alien.Translation.X
		alienTranslation.Y -= alien.Translation.Y
		alien.Translation = math.Vector2{}
	}))

	// Clear enemy bullets
	world.Manager.Join(gameComponents.Enemy, gameComponents.Bullet).Visit(ecs.Visit(func(enemyBulletEntity ecs.Entity) {
		world.Manager.DeleteEntity(enemyBulletEntity)
	}))

	// Clear alien master
	world.Manager.Join(gameComponents.AlienMaster).Visit(ecs.Visit(func(enemyBulletEntity ecs.Entity) {
		world.Manager.DeleteEntity(enemyBulletEntity)
	}))

	// Resurrect player
	loader.LoadEntities("assets/metadata/entities/player.toml", world)
}
