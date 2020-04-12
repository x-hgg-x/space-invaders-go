package systems

import (
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/math"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

var shootPlayerBulletFrame = 0

// ShootPlayerBulletSystem shoots player bullet
func ShootPlayerBulletSystem(world w.World) {
	shootPlayerBulletFrame--

	gameComponents := world.Components.Game.(*gc.Components)
	audioPlayers := *world.Resources.AudioPlayers

	if world.Manager.Join(gameComponents.Player, gameComponents.Bullet).Empty() {
		shootPlayerBulletFrame = math.Min(ebiten.DefaultTPS/20, shootPlayerBulletFrame)
	}

	if world.Resources.InputHandler.Actions[resources.ShootAction] && shootPlayerBulletFrame <= 0 {
		shootPlayerBulletFrame = ebiten.DefaultTPS

		firstPlayer := ecs.GetFirst(world.Manager.Join(gameComponents.Player, gameComponents.Controllable, world.Components.Engine.Transform))
		if firstPlayer == nil {
			return
		}
		playerX := world.Components.Engine.Transform.Get(ecs.Entity(*firstPlayer)).(*ec.Transform).Translation.X

		playerBulletEntity := loader.AddEntities(world, world.Resources.Prefabs.(*resources.Prefabs).Game.PlayerBullet)
		for iEntity := range playerBulletEntity {
			playerBulletTransform := world.Components.Engine.Transform.Get(playerBulletEntity[iEntity]).(*ec.Transform)
			playerBulletTransform.Translation.X = playerX
		}

		audioPlayers["shoot"].Rewind()
		audioPlayers["shoot"].Play()
	}
}
