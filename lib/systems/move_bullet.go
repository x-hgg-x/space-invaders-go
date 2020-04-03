package systems

import (
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// MoveBulletSystem moves bullet
func MoveBulletSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(gameComponents.PlayerBullet, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerBulletVelocity := gameComponents.PlayerBullet.Get(entity).(*gc.PlayerBullet).Velocity
		playerBulletTransform := world.Components.Engine.Transform.Get(entity).(*ec.Transform)
		playerBulletTransform.Translation.Y += playerBulletVelocity / ebiten.DefaultTPS
	}))
}
