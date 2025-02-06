package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// drawTowerRangePreview draws the range circle for the tower that would be placed
func (g *Game) drawTowerRangePreview(screen *ebiten.Image) {
	// Only show preview if mouse is in game area
	if g.mouseY <= g.gameMap.uiHeight {
		return
	}

	// Get grid position
	gridX, gridY := g.gameMap.GetGridPosition(float64(g.mouseX), float64(g.mouseY))

	// Only show if position is valid for tower placement
	if !g.gameMap.CanPlaceTower(gridX, gridY) {
		return
	}

	// Calculate range based on selected tower type
	var attackRange float64
	cellSize := float64(g.gameMap.cellSize)

	switch g.selectedTower {
	case DartTower:
		attackRange = 2 * cellSize
	case BulletTower:
		attackRange = 3 * cellSize
	case LightningTower:
		attackRange = 3 * cellSize
	case FlameTower:
		attackRange = 2 * cellSize
	case FreezeTower:
		attackRange = 2 * cellSize
	case ForkTower:
		attackRange = 4 * cellSize
	}

	// Draw range circle
	centerX := float32(float64(gridX*g.gameMap.cellSize) + float64(g.gameMap.cellSize)/2 + g.gameMap.gridOffsetX)
	centerY := float32(float64(gridY*g.gameMap.cellSize) + float64(g.gameMap.cellSize)/2 + float64(g.gameMap.uiHeight))

	vector.StrokeCircle(screen, centerX, centerY, float32(attackRange), 1.5,
		color.RGBA{160, 160, 160, 100}, true)  // More visible gray with higher opacity
}