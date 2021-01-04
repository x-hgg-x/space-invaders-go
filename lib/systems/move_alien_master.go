package systems

import (
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
)

// MoveAlienMasterSystem moves alien master
func MoveAlienMasterSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(gameComponents.Alien, gameComponents.AlienMaster, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		alien := gameComponents.Alien.Get(entity).(*gc.Alien)
		alienMaster := gameComponents.AlienMaster.Get(entity).(*gc.AlienMaster)
		alienMasterTranslation := &world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation
		alienMasterTranslation.X += alienMaster.Direction * float64(world.Resources.ScreenDimensions.Width) / 4 / ebiten.DefaultTPS

		if alienMasterTranslation.X <= -alien.Width/2 || alienMasterTranslation.X >= float64(world.Resources.ScreenDimensions.Width)+alien.Width/2 {
			world.Manager.DeleteEntity(entity)
		}
	}))
}
