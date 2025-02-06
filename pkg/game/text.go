package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
)

var (
	gameFont font.Face
	smallFont font.Face
	largeFont font.Face
)

func init() {
	// Use embedded font data
	tt, err := opentype.Parse(embeddedFont)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	// Create normal font face
	const dpi = 72
	gameFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create font face: %v", err)
	}

	// Create small font face
	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create small font face: %v", err)
	}

	// Create large font face
	largeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create large font face: %v", err)
	}
}

// DrawText draws text with the game font
func DrawText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	// Draw shadow first
	shadowColor := color.RGBA{0, 0, 0, 180}
	text.Draw(screen, str, gameFont, x+1, y+1, shadowColor)
	text.Draw(screen, str, gameFont, x+1, y, shadowColor)
	text.Draw(screen, str, gameFont, x, y+1, shadowColor)
	
	// Draw main text
	text.Draw(screen, str, gameFont, x, y, clr)
}

// DrawSmallText draws text with the small game font
func DrawSmallText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	// Draw shadow first
	shadowColor := color.RGBA{0, 0, 0, 180}
	text.Draw(screen, str, smallFont, x+1, y+1, shadowColor)
	text.Draw(screen, str, smallFont, x+1, y, shadowColor)
	text.Draw(screen, str, smallFont, x, y+1, shadowColor)
	
	// Draw main text
	text.Draw(screen, str, smallFont, x, y, clr)
}

// MeasureTextWidth returns the width of text in pixels
func MeasureTextWidth(str string, large bool) int {
	var face font.Face
	if large {
		face = largeFont
	} else {
		face = gameFont
	}
	bounds := text.BoundString(face, str)
	return bounds.Dx()
}

// DrawLargeText draws text with the large game font
func DrawLargeText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	// Draw text with large font
	text.Draw(screen, str, largeFont, x, y, clr)
}