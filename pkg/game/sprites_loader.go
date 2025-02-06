package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

// getTowerSprite returns the appropriate sprite for the tower type
func getTowerSprite(towerType TowerType) *ebiten.Image {
	var art string
	var primaryColor color.Color
	
	switch towerType {
	case DartTower:
		art = dartTowerPixelArt
		primaryColor = color.RGBA{220, 180, 100, 255}    // Richer bronze with better contrast
		return createSpriteFromArt(art, primaryColor, color.RGBA{180, 140, 60, 255})  // Darker bronze detail

	case BulletTower:
		art = bulletTowerPixelArt
		primaryColor = color.RGBA{255, 215, 100, 255}   // More vivid metallic gold
		return createSpriteFromArt(art, primaryColor, color.RGBA{215, 175, 60, 255})  // Darker gold detail

	case LightningTower:
		art = lightningTowerPixelArt
		primaryColor = color.RGBA{80, 220, 255, 255}     // Intense electric blue
		return createSpriteFromArt(art, primaryColor, color.RGBA{40, 180, 255, 255})  // Deeper electric detail

	case FlameTower:
		art = flameTowerPixelArt
		primaryColor = color.RGBA{255, 120, 50, 255}     // More vibrant orange-red flame
		return createSpriteFromArt(art, primaryColor, color.RGBA{255, 80, 30, 255})   // Deeper flame detail

	case FreezeTower:
		art = freezeTowerPixelArt
		primaryColor = color.RGBA{160, 240, 255, 255}   // Brighter ice blue
		return createSpriteFromArt(art, primaryColor, color.RGBA{100, 180, 255, 255}) // Deeper ice detail

	case ForkTower:
		art = forkTowerPixelArt
		primaryColor = color.RGBA{40, 255, 220, 255}     // Vibrant electric turquoise
		return createSpriteFromArt(art, primaryColor, color.RGBA{20, 215, 180, 255})  // Deeper turquoise detail
	}
	
	return nil
}