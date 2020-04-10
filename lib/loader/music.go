package loader

import (
	"github.com/x-hgg-x/goecsengine/loader"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/audio"
)

// LoadSounds loads music and sfx
func LoadSounds(world w.World) {
	world.Resources.AudioContext = loader.InitAudio(44100)
	audioPlayers := make(map[string]*audio.Player)
	audioPlayers["music"] = loader.LoadAudio(world.Resources.AudioContext, "assets/audio/Wave After Wave!.ogg")
	audioPlayers["music"].Play()
	audioPlayers["shoot"] = loader.LoadAudio(world.Resources.AudioContext, "assets/audio/shoot.wav")
	audioPlayers["shoot"].SetVolume(0.15)
	audioPlayers["killed"] = loader.LoadAudio(world.Resources.AudioContext, "assets/audio/killed.wav")
	audioPlayers["killed"].SetVolume(0.15)
	audioPlayers["explosion"] = loader.LoadAudio(world.Resources.AudioContext, "assets/audio/explosion.wav")
	audioPlayers["explosion"].SetVolume(0.15)
	world.Resources.AudioPlayers = &audioPlayers
}
