package main

import (
    "image"
    "image/color"
    "image/png"
    "os"
)

func main() {
    // Create a 32x32 RGBA image
    icon := image.NewRGBA(image.Rect(0, 0, 32, 32))

    // Draw a simple tower shape
    for y := 0; y < 32; y++ {
        for x := 0; x < 32; x++ {
            // Base shape
            if x >= 8 && x < 24 { // Tower body
                if y >= 12 && y < 30 { // Main tower body
                    icon.Set(x, y, color.RGBA{200, 200, 200, 255})
                } else if y >= 4 && y < 12 { // Tower top
                    if x >= 12 && x < 20 { // Narrower at top
                        icon.Set(x, y, color.RGBA{220, 220, 220, 255})
                    }
                }
            }
            // Base
            if y >= 28 && y < 32 && x >= 6 && x < 26 {
                icon.Set(x, y, color.RGBA{180, 180, 180, 255})
            }
        }
    }

    // Save to file
    f, _ := os.Create("pkg/game/assets/icon.png")
    png.Encode(f, icon)
    f.Close()
}