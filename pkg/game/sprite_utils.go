package game

import (
    "github.com/hajimehoshi/ebiten/v2"
    "image/color"
    "strings"
)

// createSpriteFromArt converts pixel art string into an image
func createSpriteFromArt(art string, primaryColor, secondaryColor color.Color) *ebiten.Image {
    lines := strings.Split(strings.TrimSpace(art), "\n")
    height := len(lines)
    width := len(lines[0])

    img := ebiten.NewImage(width, height)

    for y, line := range lines {
        for x, char := range line {
            switch char {
            case '#':
                img.Set(x, y, primaryColor)
            case '.':
                // Make background transparent
                img.Set(x, y, color.RGBA{0, 0, 0, 0})
            }
        }
    }

    return img
}