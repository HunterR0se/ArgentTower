package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "image"
    "image/color"
    "argent/pkg/game"
    "log"
)

func init() {
    // Set window properties before anything else
    ebiten.SetWindowSize(1024, 732)
    // Using Unicode characters to try to force proper title display
    ebiten.SetWindowTitle("\u0041\u0072\u0067\u0065\u006E\u0074\u0020\u0054\u006F\u0077\u0065\u0072")
}

func createIcon() image.Image {
    icon := image.NewRGBA(image.Rect(0, 0, 32, 32))
    
    // Draw a styled tower shape
    for y := 0; y < 32; y++ {
        for x := 0; x < 32; x++ {
            // Core tower body
            if x >= 12 && x < 20 {
                if y >= 8 && y < 28 {
                    icon.Set(x, y, color.RGBA{220, 220, 220, 255})
                }
            }
            // Wider base
            if y >= 24 && y < 32 {
                if x >= 8 && x < 24 {
                    icon.Set(x, y, color.RGBA{200, 200, 200, 255})
                }
            }
            // Top spire
            if y >= 2 && y < 8 {
                if x >= 14 && x < 18 {
                    icon.Set(x, y, color.RGBA{240, 240, 240, 255})
                }
            }
        }
    }
    return icon
}

func main() {
    // Set window icon
    ebiten.SetWindowIcon([]image.Image{createIcon()})
    
    // Create and run game
    g := game.NewGame()
    if err := ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}