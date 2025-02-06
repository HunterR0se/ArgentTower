package game

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawGameOver renders the game over screen with darkened background and skull
func drawGameOver(screen *ebiten.Image) bool {
	w := screen.Bounds().Dx()
	h := screen.Bounds().Dy()

	// Draw very dark overlay with slight transparency for depth
	vector.DrawFilledRect(screen, 0, 0, float32(w), float32(h),
		color.RGBA{0, 0, 0, 252}, false)

	// Calculate pulsing glow based on time
	currentTime := float64(time.Now().UnixNano()) / 1e9
	pulseIntensity := (math.Sin(currentTime*3.0) + 1.0) / 2.0 // Oscillates between 0 and 1
	glowIntensity := uint8(20 + pulseIntensity*35)            // Oscillates between 20 and 55

	// Draw "YOU ARE DEAD" text
	text := "YOU ARE DEAD"
	textWidth := MeasureTextWidth(text, true) // Using large font
	x := (w - textWidth) / 2
	y := h/2 - 100 // Higher up to make room for skull and button

	// Draw outer glow (varies with pulse)
	offsets := []int{-4, -3, -2, -1, 0, 1, 2, 3, 4}
	for _, dx := range offsets {
		for _, dy := range offsets {
			if dx == 0 && dy == 0 {
				continue
			}
			// Larger offsets get more transparency
			alpha := uint8(255 - (math.Abs(float64(dx))+math.Abs(float64(dy)))*20)
			DrawLargeText(screen, text, x+dx, y+dy,
				color.RGBA{glowIntensity, 0, 0, alpha})
		}
	}

	// Draw bright red core text
	DrawLargeText(screen, text, x, y, color.RGBA{255, 30 + glowIntensity, 30, 255})

	// Draw skull below text with pulsing glow
	skull := createSpriteFromArt(SharedSkull, color.RGBA{220, 220, 220, 255}, color.RGBA{0, 0, 0, 0})
	if skull != nil {
		skullW := float64(skull.Bounds().Dx())
		scale := float64(w) * 0.18 / skullW // Slightly larger skull (18% of screen width)

		// Draw outer glow effect
		glowOp := &ebiten.DrawImageOptions{}
		glowScale := scale * (1.0 + pulseIntensity*0.1) // Pulsing size
		glowOp.GeoM.Scale(glowScale, glowScale)
		glowOp.GeoM.Translate(
			float64(x)+float64(textWidth)/2-(skullW*glowScale)/2,
			float64(y)+80,
		)
		glowOp.ColorScale.Scale(0.8, 0.0, 0.0, float32(0.3+pulseIntensity*0.2)) // Pulsing opacity
		screen.DrawImage(skull, glowOp)

		// Draw second glow layer
		glow2Op := &ebiten.DrawImageOptions{}
		glow2Scale := scale * (1.0 + pulseIntensity*0.05) // Smaller pulse
		glow2Op.GeoM.Scale(glow2Scale, glow2Scale)
		glow2Op.GeoM.Translate(
			float64(x)+float64(textWidth)/2-(skullW*glow2Scale)/2,
			float64(y)+80,
		)
		glow2Op.ColorScale.Scale(0.9, 0.0, 0.0, float32(0.4+pulseIntensity*0.3))
		screen.DrawImage(skull, glow2Op)

		// Draw main skull
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(
			float64(x)+float64(textWidth)/2-(skullW*scale)/2,
			float64(y)+80,
		)
		screen.DrawImage(skull, op)
	}

	// Draw "New Game?" button at bottom with pulsing effect
	buttonWidth := 240
	buttonHeight := 50
	buttonX := (w - buttonWidth) / 2
	buttonY := h - buttonHeight - 60

	// Check if mouse is over button
	mouseX, mouseY := ebiten.CursorPosition()
	hovered := mouseX >= buttonX && mouseX < buttonX+buttonWidth &&
		mouseY >= buttonY && mouseY < buttonY+buttonHeight

	// Draw button glow
	glowColor := color.RGBA{
		uint8(180 + glowIntensity),
		uint8(glowIntensity / 2),
		uint8(glowIntensity / 2),
		128 + uint8(pulseIntensity*64)}

	if hovered {
		// Enhanced glow when hovered
		vector.DrawFilledRect(screen,
			float32(buttonX-4), float32(buttonY-4),
			float32(buttonWidth+8), float32(buttonHeight+8),
			glowColor, true)
	} else {
		// Normal glow
		vector.DrawFilledRect(screen,
			float32(buttonX-2), float32(buttonY-2),
			float32(buttonWidth+4), float32(buttonHeight+4),
			glowColor, true)
	}

	// Draw button background
	buttonBgColor := color.RGBA{60, 0, 0, 255}
	if hovered {
		buttonBgColor = color.RGBA{80, 0, 0, 255}
	}
	vector.DrawFilledRect(screen,
		float32(buttonX), float32(buttonY),
		float32(buttonWidth), float32(buttonHeight),
		buttonBgColor, true)

	// Draw button text
	buttonText := "New Game?"
	buttonTextW := MeasureTextWidth(buttonText, false)
	textColor := color.RGBA{255, uint8(200 + glowIntensity/2), uint8(200 + glowIntensity/2), 255}
	DrawText(screen, buttonText,
		buttonX+(buttonWidth-buttonTextW)/2,
		buttonY+buttonHeight/2+5,
		textColor)

	// Handle button click
	if hovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Create a completely new game instance and replace the current one
		newGame := NewGame()
		*currentGame = *newGame
		return true // indicate we handled the click
	}
	return false // no click was handled

}
