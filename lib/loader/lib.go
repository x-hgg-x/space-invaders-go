package loader

import (
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/x-hgg-x/go-toml"
)

type gameComponentList struct {
	Player       *gc.Player
	Enemy        *gc.Enemy
	Controllable *gc.Controllable
	Alien        *gc.Alien
	AlienMaster  *gc.AlienMaster
	Bunker       *gc.Bunker
	Bullet       *gc.Bullet
	PlayerLine   *gc.PlayerLine
	Deleted      *gc.Deleted
}

type entity struct {
	Components gameComponentList
}

type entityGameMetadata struct {
	Entities []entity `toml:"entity"`
}

func loadGameComponents(entityMetadataPath string, world w.World) []interface{} {
	var entityGameMetadata entityGameMetadata
	tree, err := toml.LoadFile(entityMetadataPath)
	utils.LogError(err)
	utils.LogError(tree.Unmarshal(&entityGameMetadata))

	gameComponentList := make([]interface{}, len(entityGameMetadata.Entities))
	for iEntity, entity := range entityGameMetadata.Entities {
		gameComponentList[iEntity] = entity.Components
	}
	return gameComponentList
}

// PreloadEntities preloads entities with components
func PreloadEntities(entityMetadataPath string, world w.World) loader.EntityComponentList {
	return loader.EntityComponentList{
		Engine: loader.LoadEngineComponents(entityMetadataPath, world),
		Game:   loadGameComponents(entityMetadataPath, world),
	}
}
