package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"fmt"
	"math"
)

// TowerButton represents a selectable tower in the UI
type TowerButton struct {
	tower      TowerType
	x, y       int
	width      int
	height     int
	selected   bool
	sprite     *ebiten.Image
	name       string
	cost       int
}

// NewTowerButton creates a new tower selection button
func NewTowerButton(tower TowerType, x, y int) *TowerButton {
	var name string
	switch tower {
	case DartTower:
		name = "Dart"
	case BulletTower:
		name = "Bullet"
	case LightningTower:
		name = "Lightning"
	case FlameTower:
		name = "Flame"
	case FreezeTower:
		name = "Freeze"
	case ForkTower:
		name = "Fork"
	}

	// Calculate cost based on tower type
	var cost int
	switch tower {
	case DartTower:
		cost = 10  // Basic tower
	case BulletTower:
		cost = 25  // Better range and speed
	case LightningTower:
		cost = 40  // High damage
	case FlameTower:
		cost = 60  // Area damage
	case FreezeTower:
		cost = 75  // Most expensive - powerful effect
	case ForkTower:
		cost = 150 // Electric fork tower with large range
	}

	return &TowerButton{
		tower:    tower,
		x:        x,
		y:        y-5,     // Move up slightly to make room for extended box
		width:    80,      // Width of tower button
		height:   65,      // Extended height to cover tower base
		sprite:   getTowerSprite(tower),
		name:     name,
		cost:     cost,
	}
}

// Draw draws the tower button
func (tb *TowerButton) Draw(screen *ebiten.Image, canAfford bool) {
	// Draw button background and selection highlight
	if tb.selected && canAfford {
		// Draw bright outline for selected tower
		vector.DrawFilledRect(screen, 
			float32(tb.x-2), float32(tb.y-2), 
			float32(tb.width+4), float32(tb.height+4), 
			color.NRGBA{0, 180, 0, 200}, true)  // Softer green outline
		// Inner background
		vector.DrawFilledRect(screen, 
			float32(tb.x), float32(tb.y), 
			float32(tb.width), float32(tb.height), 
			color.NRGBA{30, 60, 30, 255}, true)  // Darker inside
	} else {
		// Just dark background for non-selected or can't afford
		vector.DrawFilledRect(screen, 
			float32(tb.x), float32(tb.y), 
			float32(tb.width), float32(tb.height), 
			color.NRGBA{20, 20, 20, 255}, true)
	}

	// Draw tower sprite
	if tb.sprite != nil {
		op := &ebiten.DrawImageOptions{}
		
		// Scale sprite to fit button while maintaining aspect ratio
		w := float64(tb.sprite.Bounds().Dx())
		h := float64(tb.sprite.Bounds().Dy())
		scale := float64(tb.width-30) / math.Max(w, h)  // Smaller scale to ensure fit
		op.GeoM.Scale(scale, scale)
		
		// Center in button (kept within button bounds)
		op.GeoM.Translate(
			float64(tb.x) + float64(tb.width)/2 - (w*scale)/2,
			float64(tb.y) + 8,  // Move sprite up a bit more to create better balance with cost
		)
		
		// Apply effects based on affordability first, then selection
		baseScale := 1.0
		if !canAfford {
			// Gray out if can't afford
			op.ColorM.Scale(0.3, 0.3, 0.3, 1.0)
		} else {
			// Full brightness for affordable, extra bright for selected
			if tb.selected {
				baseScale = 1.5
			}
			op.ColorM.Scale(baseScale, baseScale, baseScale, 1.0)
		}
		
		screen.DrawImage(tb.sprite, op)
	}

	// Draw cost in the middle of the button
	costText := fmt.Sprintf("%d", tb.cost)
	
	var textColor color.Color = color.White
	if !canAfford {
		textColor = color.RGBA{150, 150, 150, 255}  // Dim gray when can't afford
	}

	// Draw cost near the bottom of the button
	costWidth := MeasureTextWidth(costText, false) // Using regular font for better visibility
	centerX := tb.x + (tb.width-costWidth)/2
	centerY := tb.y + tb.height - 12  // Slightly higher to balance with raised sprite
	DrawText(screen, costText, centerX, centerY, textColor)
}

// Contains checks if a point is inside the button
func (tb *TowerButton) Contains(x, y int) bool {
	return x >= tb.x && x < tb.x+tb.width &&
		   y >= tb.y && y < tb.y+tb.height
}