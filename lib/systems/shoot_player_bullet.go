package systems

import (
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	"github.com/x-hgg-x/space-invaders-go/lib/loader"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

var shootPlayerBulletFrame = 0

// ShootPlayerBulletSystem shoots player bullet
func ShootPlayerBulletSystem(world w.World) {
	shootPlayerBulletFrame--

	gameComponents := world.Components.Game.(*gc.Components)

	if world.Resources.InputHandler.Actions[resources.ShootAction] && (shootPlayerBulletFrame <= 0 || world.Manager.Join(gameComponents.Player, gameComponents.Bullet).Empty()) {
		shootPlayerBulletFrame = ebiten.DefaultTPS

		firstPlayer := ecs.GetFirst(world.Manager.Join(gameComponents.Player, gameComponents.Controllable, world.Components.Engine.Transform))
		if firstPlayer == nil {
			return
		}
		playerX := world.Components.Engine.Transform.Get(ecs.Entity(*firstPlayer)).(*ec.Transform).Translation.X

		playerBulletEntity := loader.LoadEntities("assets/metadata/entities/player_bullet.toml", world)
		for iEntity := range playerBulletEntity {
			playerBulletTransform := world.Components.Engine.Transform.Get(playerBulletEntity[iEntity]).(*ec.Transform)
			playerBulletTransform.Translation.X = playerX
		}
	}
}
