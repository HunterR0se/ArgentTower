package game

import (
	"image/color"
)

// NewEnemyWithColor creates a new enemy that uses the color scheme of the specified wave type
func NewEnemyWithColor(startY int, cellSize int, uiHeight int, enemyType EnemyType, level int, colorSource EnemyType) *Enemy {
	// Get the base color scheme from the wave type
	var primaryColor, secondaryColor color.Color
	switch colorSource {
	case SpiderEnemy:
		primaryColor = color.RGBA{220, 30, 70, 255}   // Rich venomous red
		secondaryColor = color.RGBA{120, 10, 30, 255} // Deep blood accent
	case SnakeEnemy:
		primaryColor = color.RGBA{50, 200, 50, 255}   // Toxic bright green
		secondaryColor = color.RGBA{30, 120, 30, 255} // Forest green detail
	case HawkEnemy:
		primaryColor = color.RGBA{230, 140, 30, 255}  // Majestic golden
		secondaryColor = color.RGBA{160, 80, 10, 255} // Rich brown detail
	case GhoulEnemy:
		primaryColor = color.RGBA{200, 210, 255, 255}   // Ethereal blue-white
		secondaryColor = color.RGBA{100, 110, 160, 255} // Mystic blue detail
	}

	// Start at actual entrance
	startX := float64(0)
	startYPos := float64(startY*cellSize) + float64(cellSize)/2 + float64(uiHeight)

	// Enemy size should be 80% of cell size
	enemySize := float64(cellSize) * 0.8
	startingHealth := float64(10 * level)

	// Create base enemy
	enemy := &Enemy{
		x:           startX,
		y:           startYPos,
		targetX:     startX,
		targetY:     startYPos,
		health:      startingHealth,
		maxHealth:   startingHealth,
		level:       level,
		size:        enemySize,
		enemyType:   enemyType,
		pathIndex:   0,
		pathInvalid: true,
	}

	// Set type-specific properties but use wave colors
	var spriteArt string
	switch enemyType {
	case SpiderEnemy:
		enemy.speed = 0.9
		enemy.canFly = false
		enemy.canAttack = false
		spriteArt = spiderPixelArt
	case SnakeEnemy:
		enemy.speed = 1.1
		enemy.canFly = false
		enemy.canAttack = true
		enemy.attackDamage = 0.5
		enemy.attackRange = float64(cellSize) * 1.0
		enemy.attackRate = 60
		enemy.attackChance = 0.3
		spriteArt = snakePixelArt
	case HawkEnemy:
		enemy.speed = 1.4
		enemy.canFly = false
		enemy.canAttack = false
		spriteArt = hawkPixelArt
	case GhoulEnemy:
		enemy.speed = 0.8
		enemy.canFly = false
		enemy.canAttack = true
		enemy.attackDamage = 0.8
		enemy.attackRange = float64(cellSize) * 1.5
		enemy.attackRate = 60
		enemy.attackChance = 0.25
		spriteArt = ghoulPixelArt
	}

	// Create sprite with wave colors
	enemy.sprite = createSpriteFromArt(spriteArt, primaryColor, secondaryColor)
	enemy.color = primaryColor

	return enemy
}
