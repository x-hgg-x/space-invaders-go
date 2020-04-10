package systems

import (
	"fmt"

	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"
)

// CollisionSystem manages collisions
func CollisionSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)
	gameEvents := &gameResources.Events

	screenHeight := float64(world.Resources.ScreenDimensions.Height)

	// Player bullet explosion at the top of the screen
	world.Manager.Join(gameComponents.Player, gameComponents.Bullet, world.Components.Engine.SpriteRender, world.Components.Engine.Transform).Visit(ecs.Visit(func(playerBulletEntity ecs.Entity) {
		playerBullet := gameComponents.Bullet.Get(playerBulletEntity).(*gc.Bullet)
		playerBulletSprite := world.Components.Engine.SpriteRender.Get(playerBulletEntity).(*ec.SpriteRender)
		playerBulletTranslation := &world.Components.Engine.Transform.Get(playerBulletEntity).(*ec.Transform).Translation

		if playerBulletTranslation.Y >= screenHeight-playerBullet.Height/2 {
			animation := playerBulletSprite.SpriteSheet.Animations[resources.PlayerBulletExplosionAnimation]
			firstSprite := playerBulletSprite.SpriteSheet.Sprites[animation.SpriteNumber[0]]

			playerBulletTranslation.Y = screenHeight - float64(firstSprite.Height)/2

			playerBulletEntity.
				RemoveComponent(gameComponents.Bullet).
				AddComponent(gameComponents.Deleted, &gc.Deleted{}).
				AddComponent(world.Components.Engine.AnimationControl, &ec.AnimationControl{
					Animation:      animation,
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

	// Collision between player bullets and aliens
	world.Manager.Join(gameComponents.Player, gameComponents.Bullet, world.Components.Engine.Transform).Visit(ecs.Visit(func(playerBulletEntity ecs.Entity) {
		playerBullet := gameComponents.Bullet.Get(playerBulletEntity).(*gc.Bullet)
		playerBulletTranslation := world.Components.Engine.Transform.Get(playerBulletEntity).(*ec.Transform).Translation

		world.Manager.Join(gameComponents.Alien, world.Components.Engine.AnimationControl).Visit(
			func(index int) (skip bool) {
				alienEntity := ecs.Entity(index)
				alien := gameComponents.Alien.Get(alienEntity).(*gc.Alien)
				alienSprite := world.Components.Engine.SpriteRender.Get(alienEntity).(*ec.SpriteRender)
				alienTranslation := world.Components.Engine.Transform.Get(alienEntity).(*ec.Transform).Translation
				alienAnimationControl := world.Components.Engine.AnimationControl.Get(alienEntity).(*ec.AnimationControl)

				if !rectangleCollision(alienTranslation.X, alienTranslation.Y, alien.Width, alien.Height, playerBulletTranslation.X, playerBulletTranslation.Y, playerBullet.Width, playerBullet.Height) {
					// Check next alien
					return false
				}

				// Only one alien is killed for each bullet
				world.Manager.DeleteEntity(playerBulletEntity)

				var newAlienAnimation *ec.Animation
				for key := range alienSprite.SpriteSheet.Animations {
					if alienSprite.SpriteSheet.Animations[key] == alienAnimationControl.Animation {
						switch key {
						case resources.AlienLoop1Animation:
							newAlienAnimation = alienSprite.SpriteSheet.Animations[resources.AlienDeath1Animation]
							gameEvents.ScoreEvents = append(gameEvents.ScoreEvents, resources.ScoreEvent{Score: 100})
						case resources.AlienLoop2Animation:
							newAlienAnimation = alienSprite.SpriteSheet.Animations[resources.AlienDeath2Animation]
							gameEvents.ScoreEvents = append(gameEvents.ScoreEvents, resources.ScoreEvent{Score: 200})
						case resources.AlienLoop3Animation:
							newAlienAnimation = alienSprite.SpriteSheet.Animations[resources.AlienDeath3Animation]
							gameEvents.ScoreEvents = append(gameEvents.ScoreEvents, resources.ScoreEvent{Score: 300})
						case resources.AlienMasterLoopAnimation:
							newAlienAnimation = alienSprite.SpriteSheet.Animations[resources.AlienMasterDeathAnimation]
							gameEvents.ScoreEvents = append(gameEvents.ScoreEvents, resources.ScoreEvent{Score: 1000})
						default:
							utils.LogError(fmt.Errorf("unknown animation name: '%s'", key))
						}
						break
					}
				}
				if newAlienAnimation == nil {
					utils.LogError(fmt.Errorf("unable to find animation"))
				}

				*alienAnimationControl = ec.AnimationControl{
					Animation:      newAlienAnimation,
					Command:        ec.AnimationCommand{Type: ec.AnimationCommandStart},
					RateMultiplier: 1,
				}
				alienEntity.RemoveComponent(gameComponents.Alien).AddComponent(gameComponents.Deleted, &gc.Deleted{})

				if world.Resources.AudioContext != nil {
					(*world.Resources.AudioPlayers)["killed"].Rewind()
					(*world.Resources.AudioPlayers)["killed"].Play()
				}

				// Skip other aliens
				return true
			})
	}))

	// Collision between player bullets and enemy bullets
	world.Manager.Join(gameComponents.Player, gameComponents.Bullet, world.Components.Engine.Transform).Visit(ecs.Visit(func(playerBulletEntity ecs.Entity) {
		playerBullet := gameComponents.Bullet.Get(playerBulletEntity).(*gc.Bullet)
		playerBulletTranslation := world.Components.Engine.Transform.Get(playerBulletEntity).(*ec.Transform).Translation

		world.Manager.Join(gameComponents.Enemy, gameComponents.Bullet, world.Components.Engine.Transform).Visit(ecs.Visit(func(enemyBulletEntity ecs.Entity) {
			enemyBullet := gameComponents.Bullet.Get(enemyBulletEntity).(*gc.Bullet)
			enemyBulletTranslation := world.Components.Engine.Transform.Get(enemyBulletEntity).(*ec.Transform).Translation

			if rectangleCollision(enemyBulletTranslation.X, enemyBulletTranslation.Y, enemyBullet.Width, enemyBullet.Height, playerBulletTranslation.X, playerBulletTranslation.Y, playerBullet.Width, playerBullet.Height) {
				world.Manager.DeleteEntity(playerBulletEntity)
				world.Manager.DeleteEntity(enemyBulletEntity)
			}
		}))
	}))

	// Collision between bullets and bunkers
	world.Manager.Join(gameComponents.Bullet, world.Components.Engine.Transform).Visit(ecs.Visit(func(bulletEntity ecs.Entity) {
		bullet := gameComponents.Bullet.Get(bulletEntity).(*gc.Bullet)
		bulletTranslation := world.Components.Engine.Transform.Get(bulletEntity).(*ec.Transform).Translation

		world.Manager.Join(gameComponents.Bunker, world.Components.Engine.Transform).Visit(ecs.Visit(func(bunkerEntity ecs.Entity) {
			bunkerPixelSize := float64(gameComponents.Bunker.Get(bunkerEntity).(*gc.Bunker).PixelSize)
			bunkerTranslation := world.Components.Engine.Transform.Get(bunkerEntity).(*ec.Transform).Translation

			if rectangleCollision(bunkerTranslation.X, bunkerTranslation.Y, bunkerPixelSize, bunkerPixelSize, bulletTranslation.X, bulletTranslation.Y, bullet.Width, bullet.Height) {
				world.Manager.DeleteEntity(bunkerEntity)
				bullet.Health -= bunkerPixelSize * bunkerPixelSize
			}
		}))

		if bullet.Health <= 0 {
			world.Manager.DeleteEntity(bulletEntity)
		}
	}))

	// Collision between player and enemy bullets
	world.Manager.Join(gameComponents.Player, gameComponents.Controllable, world.Components.Engine.SpriteRender, world.Components.Engine.Transform).Visit(ecs.Visit(func(playerEntity ecs.Entity) {
		playerControllable := gameComponents.Controllable.Get(playerEntity).(*gc.Controllable)
		playerTranslation := world.Components.Engine.Transform.Get(playerEntity).(*ec.Transform).Translation

		world.Manager.Join(gameComponents.Enemy, gameComponents.Bullet, world.Components.Engine.Transform).Visit(
			func(index int) (skip bool) {
				enemyBulletEntity := ecs.Entity(index)
				enemyBullet := gameComponents.Bullet.Get(enemyBulletEntity).(*gc.Bullet)
				enemyBulletTranslation := world.Components.Engine.Transform.Get(enemyBulletEntity).(*ec.Transform).Translation

				if !rectangleCollision(playerTranslation.X, playerTranslation.Y, playerControllable.Width, playerControllable.Height, enemyBulletTranslation.X, enemyBulletTranslation.Y, enemyBullet.Width, enemyBullet.Height) {
					return false
				}

				world.Manager.DeleteEntity(enemyBulletEntity)
				gameEvents.LifeEvents = append(gameEvents.LifeEvents, resources.LifeEvent{})
				gameEvents.ScoreEvents = append(gameEvents.ScoreEvents, resources.ScoreEvent{Score: -1000})

				if world.Resources.AudioContext != nil {
					(*world.Resources.AudioPlayers)["explosion"].Rewind()
					(*world.Resources.AudioPlayers)["explosion"].Play()
				}
				return true
			})
	}))

	// Finish level if no alien are left
	if world.Manager.Join(gameComponents.Alien, gameComponents.AlienMaster.Not()).Empty() {
		gameResources.StateEvent = resources.StateEventLevelComplete
	}
}

func rectangleCollision(r1X, r1Y, r1Width, r1Height, r2X, r2Y, r2Width, r2Height float64) bool {
	return r1X-r1Width/2-r2Width/2 <= r2X && r2X <= r1X+r1Width/2+r2Width/2 && r1Y-r1Height/2-r2Height/2 <= r2Y && r2Y <= r1Y+r1Height/2+r2Height/2
}
