package systems

import (
	"math"

	"github.com/hajimehoshi/ebiten"
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"
)

var moveAlienDirectionX = 1.0

const (
	moveAlienVelocity = 1000.0 / ebiten.DefaultTPS
	moveAlienDiffY    = -24.0
)

// MoveAlienSystem moves alien
func MoveAlienSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	alienSet := world.Manager.Join(gameComponents.Alien, gameComponents.AlienMaster.Not(), world.Components.Engine.Transform)
	if alienSet.Empty() {
		return
	}

	velocityRatio := float64(alienSet.Size())

	minAlienPosX := float64(world.Resources.ScreenDimensions.Width)
	maxAlienPosX := 0.0

	alienSet.Visit(ecs.Visit(func(entity ecs.Entity) {
		alienWidth := gameComponents.Alien.Get(entity).(*gc.Alien).Width
		alienX := world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation.X
		minAlienPosX = math.Min(minAlienPosX, alienX-alienWidth/2)
		maxAlienPosX = math.Max(maxAlienPosX, alienX+alienWidth/2)
	}))

	var movementX, movementY float64
	if moveAlienDirectionX > 0 && maxAlienPosX < float64(world.Resources.ScreenDimensions.Width) {
		movementX = math.Min(moveAlienDirectionX*moveAlienVelocity/velocityRatio, float64(world.Resources.ScreenDimensions.Width)-maxAlienPosX)
	} else if moveAlienDirectionX < 0 && minAlienPosX > 0 {
		movementX = math.Max(moveAlienDirectionX*moveAlienVelocity/velocityRatio, -minAlienPosX)
	} else if moveAlienDirectionX > 0 && maxAlienPosX >= float64(world.Resources.ScreenDimensions.Width) || moveAlienDirectionX < 0 && minAlienPosX <= 0 {
		moveAlienDirectionX *= -1
		movementY = moveAlienDiffY
	}

	alienSet.Visit(ecs.Visit(func(entity ecs.Entity) {
		alienTranslation := &world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation
		alienTranslation.X += movementX
		alienTranslation.Y += movementY
	}))
}
