package systems

import (
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// MoveBulletSystem moves bullet
func MoveBulletSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(gameComponents.Bullet, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		bulletVelocity := gameComponents.Bullet.Get(entity).(*gc.Bullet).Velocity
		bulletTransform := world.Components.Engine.Transform.Get(entity).(*ec.Transform)
		bulletTransform.Translation.Y += bulletVelocity / ebiten.DefaultTPS
	}))
}
