package systems

import (
	"math"

	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// MovePlayerSystem moves player controllable sprite
func MovePlayerSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(gameComponents.Player, gameComponents.Controllable, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerControllable := gameComponents.Controllable.Get(entity).(*gc.Controllable)
		playerTransform := world.Components.Engine.Transform.Get(entity).(*ec.Transform)

		screenWidth := float64(world.Resources.ScreenDimensions.Width)
		playerX := playerTransform.Translation.X
		axisValue := world.Resources.InputHandler.Axes[resources.PlayerAxis]

		if world.Resources.Controls.Axes[resources.PlayerAxis].Type == "MouseAxis" {
			playerX = axisValue * screenWidth
		} else {
			playerX += axisValue * screenWidth / ebiten.DefaultTPS / 4
		}

		minValue := playerControllable.Width / 2
		maxValue := float64(world.Resources.ScreenDimensions.Width) - playerControllable.Width/2
		playerTransform.Translation.X = math.Min(math.Max(playerX, minValue), maxValue)
	}))
}
