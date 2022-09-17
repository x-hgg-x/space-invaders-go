package loader

import (
	"github.com/x-hgg-x/goecsengine/loader"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// LoadSounds loads music and sfx
func LoadSounds(world w.World, sound bool) {
	world.Resources.AudioContext = loader.InitAudio(44100)
	audioPlayers := make(map[string]*audio.Player)
	audioPlayers["music"] = loader.LoadAudio(world.Resources.AudioContext, "assets/lfs/audio/Wave After Wave!.ogg")
	audioPlayers["shoot"] = loader.LoadAudio(world.Resources.AudioContext, "assets/lfs/audio/shoot.wav")
	audioPlayers["killed"] = loader.LoadAudio(world.Resources.AudioContext, "assets/lfs/audio/killed.wav")
	audioPlayers["explosion"] = loader.LoadAudio(world.Resources.AudioContext, "assets/lfs/audio/explosion.wav")
	world.Resources.AudioPlayers = &audioPlayers

	audioPlayers["music"].Play()
	if sound {
		EnableSound(world)
	} else {
		DisableSound(world)
	}
}

// EnableSound enables sound
func EnableSound(world w.World) {
	audioPlayers := *world.Resources.AudioPlayers
	audioPlayers["music"].SetVolume(1)
	audioPlayers["shoot"].SetVolume(0.15)
	audioPlayers["killed"].SetVolume(0.15)
	audioPlayers["explosion"].SetVolume(0.15)
}

// DisableSound disables sound
func DisableSound(world w.World) {
	audioPlayers := *world.Resources.AudioPlayers
	for key := range audioPlayers {
		audioPlayers[key].SetVolume(0)
	}
}
