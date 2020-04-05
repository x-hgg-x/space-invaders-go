package systems

import (
	"math/rand"

	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	"github.com/x-hgg-x/space-invaders-go/lib/loader"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/math"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

var shootEnemyBulletFrame = 0

// ShootEnemyBulletSystem shoots enemy bullet
func ShootEnemyBulletSystem(world w.World) {
	shootEnemyBulletFrame--

	gameComponents := world.Components.Game.(*gc.Components)

	alienSet := world.Manager.Join(gameComponents.Alien, gameComponents.AlienMaster.Not())
	if alienSet.Empty() {
		return
	}

	if shootEnemyBulletFrame <= 0 {
		shootEnemyBulletFrame = ebiten.DefaultTPS * 2

		// Select random alien
		alienEntities := []ecs.Entity{}
		alienSet.Visit(ecs.Visit(func(entity ecs.Entity) {
			alienEntities = append(alienEntities, entity)
		}))
		alienEntity := alienEntities[rand.Intn(len(alienEntities))]
		alienHeight := gameComponents.Alien.Get(alienEntity).(*gc.Alien).Height
		alienTranslation := world.Components.Engine.Transform.Get(alienEntity).(*ec.Transform).Translation

		enemyBulletEntity := loader.LoadEntities("assets/metadata/entities/enemy_bullet.toml", world)
		for iEntity := range enemyBulletEntity {
			enemyBulletHeight := gameComponents.Bullet.Get(enemyBulletEntity[iEntity]).(*gc.Bullet).Height
			enemyBulletTransform := world.Components.Engine.Transform.Get(enemyBulletEntity[iEntity]).(*ec.Transform)
			enemyBulletTransform.Translation = math.Vector2{
				X: alienTranslation.X,
				Y: alienTranslation.Y - alienHeight/2 - enemyBulletHeight/2,
			}
		}
	}
}
