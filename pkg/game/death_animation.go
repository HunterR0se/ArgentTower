package game

import (
    "github.com/hajimehoshi/ebiten/v2"
    "image/color"
    "math"
    "strings"
)

type DeathAnimation struct {
    x, y          float64
    size          float64
    deathSprite   *ebiten.Image
    frameLife     int
    rotation      float64   // Track rotation angle
    rotationSpeed float64   // Rotation speed
}

func NewDeathAnimation(x, y float64, size float64, enemySprite *ebiten.Image) *DeathAnimation {
    // Create a slightly lighter version of the sprite
    deathColor := color.RGBA{200, 200, 200, 255}  // Light gray instead of white
    
    // Get the sprite's art based on its size
    var spriteArt string
    w := enemySprite.Bounds().Dx()
    
    // Match sprite dimensions to art
    spiderWidth := len(strings.Split(spiderPixelArt, "\n")[1])
    snakeWidth := len(strings.Split(snakePixelArt, "\n")[1])
    hawkWidth := len(strings.Split(hawkPixelArt, "\n")[1])
    
    if w == spiderWidth {
        spriteArt = spiderPixelArt
    } else if w == snakeWidth {
        spriteArt = snakePixelArt
    } else if w == hawkWidth {
        spriteArt = hawkPixelArt
    } else {
        spriteArt = ghoulPixelArt
    }
    
    // Create sprite
    deathSprite := createSpriteFromArt(spriteArt, deathColor, color.RGBA{0, 0, 0, 0})

    return &DeathAnimation{
        x:             x,
        y:             y,
        size:          size,
        deathSprite:   deathSprite,
        frameLife:     45,  // Longer fade out time
        rotation:      0,
        rotationSpeed: math.Pi / 45, // Much slower rotation
    }
}

func (da *DeathAnimation) Update() bool {
    da.frameLife--
    
    // Keep rotation speed constant
    da.rotation += da.rotationSpeed
    
    return da.frameLife > 0
}

func (da *DeathAnimation) Draw(screen *ebiten.Image) {
    if da.deathSprite == nil {
        return
    }

    lifePercent := float64(da.frameLife) / 45.0  // Using the new 45 frame lifetime
    shrinkScale := lifePercent  // Shrink to nothing as we fade

    w := float64(da.deathSprite.Bounds().Dx())
    h := float64(da.deathSprite.Bounds().Dy())
    baseScale := (da.size * shrinkScale) / math.Max(w, h)  // Shrink as we die

    // Draw ghost trails with scaling consistent with main sprite
    numGhosts := 3
    for i := 0; i < numGhosts; i++ {
        ghostOp := &ebiten.DrawImageOptions{}
        
        // Scale, shrinking as we die
        ghostOp.GeoM.Scale(baseScale, baseScale)
        
        // Calculate ghost offset based on rotation
        ghostRotation := da.rotation - float64(i)*0.2
        
        // Center sprite for rotation
        ghostOp.GeoM.Translate(-w/2, -h/2)
        ghostOp.GeoM.Rotate(ghostRotation)
        ghostOp.GeoM.Translate(da.x, da.y)
        
        // Calculate ghost opacity
        ghostAlpha := float32(lifePercent * 0.15 * float64(numGhosts-i)/float64(numGhosts))
        
        // Make ghosts fade out
        ghostOp.ColorScale.Scale(1.0, 1.0, 1.0, ghostAlpha)
        
        screen.DrawImage(da.deathSprite, ghostOp)
    }

    // Draw main sprite
    op := &ebiten.DrawImageOptions{}
    
    // Scale, shrinking as we die
    op.GeoM.Scale(baseScale, baseScale)
    
    // Center sprite for rotation
    op.GeoM.Translate(-w/2, -h/2)
    op.GeoM.Rotate(da.rotation)
    op.GeoM.Translate(da.x, da.y)
    
    // Simple fade out
    op.ColorScale.Scale(1.0, 1.0, 1.0, float32(lifePercent))
    
    screen.DrawImage(da.deathSprite, op)
}