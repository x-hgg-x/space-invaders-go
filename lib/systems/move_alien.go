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

var moveAlienDirectionX = 1.0

const (
	moveAlienVelocity = 1000.0 / ebiten.DefaultTPS
	moveAlienDiffY    = -24.0
)

// MoveAlienSystem moves alien
func MoveAlienSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameEvents := &world.Resources.Game.(*resources.Game).Events

	alienSet := world.Manager.Join(gameComponents.Alien, gameComponents.AlienMaster.Not(), world.Components.Engine.Transform)
	if alienSet.Empty() {
		return
	}

	velocityRatio := float64(alienSet.Size())

	minAlienPosX := float64(world.Resources.ScreenDimensions.Width)
	maxAlienPosX := 0.0
	minAlienPosY := float64(world.Resources.ScreenDimensions.Height)

	alienSet.Visit(ecs.Visit(func(entity ecs.Entity) {
		alien := gameComponents.Alien.Get(entity).(*gc.Alien)
		alienTranslation := world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation
		minAlienPosX = math.Min(minAlienPosX, alienTranslation.X-alien.Width/2)
		maxAlienPosX = math.Max(maxAlienPosX, alienTranslation.X+alien.Width/2)
		minAlienPosY = math.Min(minAlienPosY, alienTranslation.Y-alien.Height/2)
	}))

	var movementX, movementY float64
	if moveAlienDirectionX > 0 && maxAlienPosX < float64(world.Resources.ScreenDimensions.Width) {
		movementX = math.Min(moveAlienDirectionX*moveAlienVelocity/velocityRatio, float64(world.Resources.ScreenDimensions.Width)-maxAlienPosX)
	} else if moveAlienDirectionX < 0 && minAlienPosX > 0 {
		movementX = math.Max(moveAlienDirectionX*moveAlienVelocity/velocityRatio, -minAlienPosX)
	} else if moveAlienDirectionX > 0 && maxAlienPosX >= float64(world.Resources.ScreenDimensions.Width) || moveAlienDirectionX < 0 && minAlienPosX <= 0 {
		moveAlienDirectionX *= -1
		movementY = moveAlienDiffY

		// Lose a life when aliens reach the player line
		if playerLineEntity := ecs.GetFirst(world.Manager.Join(gameComponents.PlayerLine, world.Components.Engine.Transform)); playerLineEntity != nil {
			playerLineY := world.Components.Engine.Transform.Get(*playerLineEntity).(*ec.Transform).Translation.Y

			if minAlienPosY+moveAlienDiffY < playerLineY {
				gameEvents.LifeEvents = append(gameEvents.LifeEvents, resources.LifeEvent{})
			}
		}
	}

	alienSet.Visit(ecs.Visit(func(entity ecs.Entity) {
		alien := gameComponents.Alien.Get(entity).(*gc.Alien)
		alienTranslation := &world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation
		alienTranslation.X += movementX
		alienTranslation.Y += movementY
		alien.Translation.X += movementX
		alien.Translation.Y += movementY
	}))
}
