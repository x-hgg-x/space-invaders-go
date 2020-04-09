package loader

import (
	"fmt"
	"image/color"

	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	"github.com/x-hgg-x/space-invaders-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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
	_, err := toml.DecodeFile(entityMetadataPath, &entityGameMetadata)
	utils.LogError(err)

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

// LoadBunkers creates pixel bunker entities for each bunker
func LoadBunkers(world w.World) []ecs.Entity {
	gameComponents := world.Components.Game.(*gc.Components)

	// Get bunker image path
	type spriteSheetMetadata struct {
		SpriteSheets struct {
			Bunker struct {
				TextureImageName string `toml:"texture_image"`
			}
		} `toml:"sprite_sheet"`
	}

	var metadata spriteSheetMetadata
	_, err := toml.DecodeFile("assets/metadata/spritesheets/spritesheets.toml", &metadata)
	utils.LogError(err)

	// Load bunker image
	bunkerImagePath := metadata.SpriteSheets.Bunker.TextureImageName
	_, bunkerImage, err := ebitenutil.NewImageFromFile(bunkerImagePath, ebiten.FilterNearest)
	utils.LogError(err)

	// Load bunker entities
	bunkerEntities := loader.AddEntities(world, world.Resources.Prefabs.(*resources.Prefabs).Game.Bunker)
	if len(bunkerEntities) == 0 {
		return []ecs.Entity{}
	}

	// Create pixel image
	pixelSize := gameComponents.Bunker.Get(bunkerEntities[0]).(*gc.Bunker).PixelSize
	for _, bunkerEntity := range bunkerEntities {
		if pixelSize != gameComponents.Bunker.Get(bunkerEntity).(*gc.Bunker).PixelSize {
			utils.LogError(fmt.Errorf("pixel size must be the same for all bunkers"))
		}
	}
	pixelImage, err := ebiten.NewImage(pixelSize, pixelSize, ebiten.FilterNearest)
	utils.LogError(err)
	pixelImage.Fill(color.RGBA{0, 255, 0, 255})

	// Create new bunker entities for each set of bunker pixels
	newBunkerEntities := []ecs.Entity{}
	for _, bunkerEntity := range bunkerEntities {
		bunkerSprite := world.Components.Engine.SpriteRender.Get(bunkerEntity).(*ec.SpriteRender)
		bunkerTransform := world.Components.Engine.Transform.Get(bunkerEntity).(*ec.Transform)

		bunkerSpriteWidth := float64(bunkerSprite.SpriteSheet.Sprites[bunkerSprite.SpriteNumber].Width)
		bunkerSpriteHeight := float64(bunkerSprite.SpriteSheet.Sprites[bunkerSprite.SpriteNumber].Height)

		bounds := bunkerImage.Bounds()
		for x := bounds.Min.X; x < bounds.Max.X; x += pixelSize {
			for y := bounds.Min.Y; y < bounds.Max.Y; y += pixelSize {
				if _, _, _, alpha := bunkerImage.At(x, y).RGBA(); alpha > 0 {
					newBunkerEntities = append(newBunkerEntities, world.Manager.NewEntity().
						AddComponent(world.Components.Engine.SpriteRender, &ec.SpriteRender{
							SpriteSheet: &ec.SpriteSheet{
								Texture: ec.Texture{Image: pixelImage},
								Sprites: []ec.Sprite{{X: 0, Y: 0, Width: pixelSize, Height: pixelSize}},
							},
							SpriteNumber: 0,
						}).
						AddComponent(world.Components.Engine.Transform, &ec.Transform{
							Depth: bunkerTransform.Depth,
							Translation: math.Vector2{
								X: bunkerTransform.Translation.X - bunkerSpriteWidth/2 + float64(x) + float64(pixelSize)/2,
								Y: bunkerTransform.Translation.Y + bunkerSpriteHeight/2 - float64(y) - float64(pixelSize)/2,
							},
						}).
						AddComponent(gameComponents.Bunker, &gc.Bunker{PixelSize: pixelSize}))
				}
			}
		}
		// Delete old bunker entity
		world.Manager.DeleteEntity(bunkerEntity)
	}
	return newBunkerEntities
}
