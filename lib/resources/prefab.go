package resources

import "github.com/x-hgg-x/goecsengine/loader"

// MenuPrefabs contains menu prefabs
type MenuPrefabs struct {
	MuteMenu          loader.EntityComponentList
	MainMenu          loader.EntityComponentList
	DifficultyMenu    loader.EntityComponentList
	PauseMenu         loader.EntityComponentList
	GameOverMenu      loader.EntityComponentList
	LevelCompleteMenu loader.EntityComponentList
	HighscoresMenu    loader.EntityComponentList
}

// GamePrefabs contains game prefabs
type GamePrefabs struct {
	Background   loader.EntityComponentList
	Alien        loader.EntityComponentList
	Player       loader.EntityComponentList
	PlayerLine   loader.EntityComponentList
	Bunker       loader.EntityComponentList
	AlienMaster  loader.EntityComponentList
	PlayerBullet loader.EntityComponentList
	EnemyBullet  loader.EntityComponentList
	Score        loader.EntityComponentList
	Life         loader.EntityComponentList
	Difficulty   loader.EntityComponentList
}

// Prefabs contains menu and game prefabs
type Prefabs struct {
	Menu MenuPrefabs
	Game GamePrefabs
}
