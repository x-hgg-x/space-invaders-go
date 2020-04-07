package systems

import (
	"math/rand"

	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	"github.com/x-hgg-x/space-invaders-go/lib/loader"

	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

var spawnAlienMasterFrame = int(ebiten.DefaultTPS * 20 * rand.Float32())

// SpawnAlienMasterSystem spawns alien master
func SpawnAlienMasterSystem(world w.World) {
	spawnAlienMasterFrame--

	gameComponents := world.Components.Game.(*gc.Components)

	if !world.Manager.Join(gameComponents.AlienMaster).Empty() {
		return
	}

	if spawnAlienMasterFrame <= 0 {
		spawnAlienMasterFrame = int(ebiten.DefaultTPS * 20 * rand.Float32())

		alienMasterEntity := loader.LoadEntities("assets/metadata/entities/alien_master.toml", world)
		for iEntity := range alienMasterEntity {
			alien := gameComponents.Alien.Get(alienMasterEntity[iEntity]).(*gc.Alien)
			alienMaster := gameComponents.AlienMaster.Get(alienMasterEntity[iEntity]).(*gc.AlienMaster)
			alienMasterTransform := world.Components.Engine.Transform.Get(alienMasterEntity[iEntity]).(*ec.Transform)

			if rand.Intn(2) == 0 {
				alienMaster.Direction = 1
				alienMasterTransform.Translation.X = -alien.Width / 2
			} else {
				alienMaster.Direction = -1
				alienMasterTransform.Translation.X = float64(world.Resources.ScreenDimensions.Width) + alien.Width/2
			}
		}
	}
}
