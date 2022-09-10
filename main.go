package main

import (
	_ "image/png"

	gc "github.com/x-hgg-x/space-invaders-go/lib/components"
	gloader "github.com/x-hgg-x/space-invaders-go/lib/loader"
	gr "github.com/x-hgg-x/space-invaders-go/lib/resources"
	gs "github.com/x-hgg-x/space-invaders-go/lib/states"

	"github.com/x-hgg-x/goecsengine/loader"
	er "github.com/x-hgg-x/goecsengine/resources"
	es "github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowWidth  = 1000
	windowHeight = 800
)

type mainGame struct {
	world        w.World
	stateMachine es.StateMachine
}

func (game *mainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	ebiten.SetWindowSize(outsideWidth, outsideHeight)
	return windowWidth, windowHeight
}

func (game *mainGame) Update() error {
	game.stateMachine.Update(game.world)
	return nil
}

func (game *mainGame) Draw(screen *ebiten.Image) {
	game.stateMachine.Draw(game.world, screen)
}

func main() {
	world := w.InitWorld(&gc.Components{})

	// Init screen dimensions
	world.Resources.ScreenDimensions = &er.ScreenDimensions{Width: windowWidth, Height: windowHeight}

	// Load controls
	axes := []string{gr.PlayerAxis}
	actions := []string{gr.ShootAction, gr.EnableDisableSoundAction}
	controls, inputHandler := loader.LoadControls("config/controls.toml", axes, actions)
	world.Resources.Controls = &controls
	world.Resources.InputHandler = &inputHandler

	// Load sprite sheets
	spriteSheets := loader.LoadSpriteSheets("assets/metadata/spritesheets/spritesheets.toml")
	world.Resources.SpriteSheets = &spriteSheets

	// Load fonts
	fonts := loader.LoadFonts("assets/metadata/fonts/fonts.toml")
	world.Resources.Fonts = &fonts

	// Load prefabs
	world.Resources.Prefabs = &gr.Prefabs{
		Menu: gr.MenuPrefabs{
			MuteMenu:          gloader.PreloadEntities("assets/metadata/entities/ui/mute_menu.toml", world),
			MainMenu:          gloader.PreloadEntities("assets/metadata/entities/ui/main_menu.toml", world),
			DifficultyMenu:    gloader.PreloadEntities("assets/metadata/entities/ui/difficulty_menu.toml", world),
			PauseMenu:         gloader.PreloadEntities("assets/metadata/entities/ui/pause_menu.toml", world),
			GameOverMenu:      gloader.PreloadEntities("assets/metadata/entities/ui/game_over_menu.toml", world),
			LevelCompleteMenu: gloader.PreloadEntities("assets/metadata/entities/ui/level_complete_menu.toml", world),
			HighscoresMenu:    gloader.PreloadEntities("assets/metadata/entities/ui/highscores_menu.toml", world),
		},
		Game: gr.GamePrefabs{
			Background:   gloader.PreloadEntities("assets/metadata/entities/background.toml", world),
			Alien:        gloader.PreloadEntities("assets/metadata/entities/alien.toml", world),
			Player:       gloader.PreloadEntities("assets/metadata/entities/player.toml", world),
			PlayerLine:   gloader.PreloadEntities("assets/metadata/entities/player_line.toml", world),
			Bunker:       gloader.PreloadEntities("assets/metadata/entities/bunker.toml", world),
			AlienMaster:  gloader.PreloadEntities("assets/metadata/entities/alien_master.toml", world),
			PlayerBullet: gloader.PreloadEntities("assets/metadata/entities/player_bullet.toml", world),
			EnemyBullet:  gloader.PreloadEntities("assets/metadata/entities/enemy_bullet.toml", world),
			Score:        gloader.PreloadEntities("assets/metadata/entities/ui/score.toml", world),
			Life:         gloader.PreloadEntities("assets/metadata/entities/ui/life.toml", world),
			Difficulty:   gloader.PreloadEntities("assets/metadata/entities/ui/difficulty.toml", world),
		},
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Space Invaders")

	utils.LogError(ebiten.RunGame(&mainGame{world, es.Init(&gs.MuteMenuState{}, world)}))
}
