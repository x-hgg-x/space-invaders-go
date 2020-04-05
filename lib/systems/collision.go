package systems

import (
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"
)

// CollisionSystem manages collisions
func CollisionSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	screenHeight := float64(world.Resources.ScreenDimensions.Height)

	// Player bullet explosion at the top of the screen
	world.Manager.Join(gameComponents.Player, gameComponents.Bullet, world.Components.Engine.SpriteRender, world.Components.Engine.Transform).Visit(ecs.Visit(func(playerBulletEntity ecs.Entity) {
		playerBullet := gameComponents.Bullet.Get(playerBulletEntity).(*gc.Bullet)
		playerBulletSprite := world.Components.Engine.SpriteRender.Get(playerBulletEntity).(*ec.SpriteRender)
		playerBulletTranslation := &world.Components.Engine.Transform.Get(playerBulletEntity).(*ec.Transform).Translation

		if playerBulletTranslation.Y >= screenHeight-playerBullet.Height/2 {
			animation := playerBulletSprite.SpriteSheet.Animations[resources.PlayerBulletExplosion]
			firstSprite := playerBulletSprite.SpriteSheet.Sprites[animation.SpriteNumber[0]]

			playerBulletTranslation.Y = screenHeight - float64(firstSprite.Height)/2

			playerBulletEntity.
				RemoveComponent(gameComponents.Bullet).
				AddComponent(gameComponents.Deleted, &gc.Deleted{}).
				AddComponent(world.Components.Engine.AnimationControl, &ec.AnimationControl{
					Animation:      &animation,
					Command:        ec.AnimationCommand{Type: ec.AnimationCommandStart},
					RateMultiplier: 1,
				})
		}
	}))

	// Remove enemy bullet at the bottom of the screen
	world.Manager.Join(gameComponents.Enemy, gameComponents.Bullet, world.Components.Engine.Transform).Visit(ecs.Visit(func(enemyBulletEntity ecs.Entity) {
		enemyBullet := gameComponents.Bullet.Get(enemyBulletEntity).(*gc.Bullet)
		enemyBulletTranslation := &world.Components.Engine.Transform.Get(enemyBulletEntity).(*ec.Transform).Translation

		if enemyBulletTranslation.Y <= -enemyBullet.Height/2 {
			world.Manager.DeleteEntity(enemyBulletEntity)
		}
	}))
}
