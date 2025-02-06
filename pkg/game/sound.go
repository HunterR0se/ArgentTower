package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 44100
)

var (
	audioContext *audio.Context
)

// Initialize the audio context
func init() {
	audioContext = audio.NewContext(sampleRate)
}

// Simple square wave generator for retro sound
func generateSquareWave(freq float64, duration float64) []byte {
	samples := int(sampleRate * duration)
	b := make([]byte, samples)
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		if math.Mod(t*freq, 1.0) < 0.5 {
			b[i] = 128 + 32 // Positive phase
		} else {
			b[i] = 128 - 32 // Negative phase
		}
	}
	return b
}

// PlayShootSound plays a simple retro shoot sound (keeping for compatibility)
func PlayShootSound() {
	PlayTowerShootSound()
}

// PlayTowerShootSound plays a simple retro tower shoot sound
func PlayTowerShootSound() {
	// Reduced from 600 to 220 Hz - a more pleasant "pew" sound
	data := generateSquareWave(220, 0.08)
	player := audioContext.NewPlayerFromBytes(data)
	player.Play()
}

// PlayAttackSound plays a simple retro attack sound
func PlayAttackSound() {
	// Reduced from 300 to 110 Hz - a deeper attack sound
	data := generateSquareWave(110, 0.15)
	player := audioContext.NewPlayerFromBytes(data)
	player.Play()
}

// PlayEnemyDeathSound plays a simple retro death sound
func PlayEnemyDeathSound() {
	// Reduced from 200 to 80 Hz - a deeper explosion sound
	data := generateSquareWave(80, 0.2)
	player := audioContext.NewPlayerFromBytes(data)
	player.Play()
}
