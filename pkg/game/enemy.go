package game

import (
	"image/color"
	"math"
	mathrand "math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

// EnemyType represents different types of enemies
type EnemyType int

const (
	SpiderEnemy EnemyType = iota
	SnakeEnemy
	HawkEnemy  // Can fly over obstacles
	GhoulEnemy // Will attack towers (renamed from Wolf)
	BlobEnemy  // Boss type enemy
)

// Enemy represents an enemy unit
type Enemy struct {
	x, y              float64 // Precise position for smooth movement
	targetX, targetY  float64 // Next target point
	speed             float64
	health            float64
	maxHealth         float64
	level             int     // Enemy level for health and rewards
	size              float64 // Size of the enemy for drawing
	enemyType         EnemyType
	sprite            *ebiten.Image
	color             color.Color // Fallback if sprite not loaded
	path              []Point     // Current path to follow
	pathIndex         int         // Current position in path
	pathInvalid       bool        // Flag to indicate if path needs recalculation
	canFly            bool        // For flying enemies like hawks
	frozenTimer       float64     // Time until enemy unfreezes (0 if not frozen)
	canAttack         bool        // Whether this enemy can attack towers
	attackDamage      float64     // How much damage this enemy does to towers
	attackRange       float64     // How close enemy needs to be to attack tower
	attackRate        int         // How often the enemy can attack (in frames)
	lastAttack        int         // Frames since last attack
	attackChance      float64     // Probability to choose to attack (0-1)
	targetTower       *Tower      // Current tower being targeted
	attackDuration    int         // How long to stay in one place attacking (in frames)
	currentAttackTime int         // Current time spent attacking
	eyeFlashTimer     int         // Timer for eye flash effect
	eyeFlashing       bool        // Whether eyes are currently flashing
	moveTimer         float64     // Timer for movement animation
	moveOffset        float64     // Current movement offset
}

// NewEnemy creates a new enemy at the entrance
func NewEnemy(startY int, cellSize int, uiHeight int, enemyType EnemyType, level int) *Enemy {
	// Start at actual entrance
	startX := float64(0) // Start at the edge
	startYPos := float64(startY*cellSize) + float64(cellSize)/2 + float64(uiHeight)

	// Enemy size should be 80% of cell size
	enemySize := float64(cellSize) * 0.8

	// Set initial health to 10 * Level
	startingHealth := float64(10 * level)

	// Default enemy properties
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

	// Set type-specific properties
	switch enemyType {
	case SpiderEnemy:
		enemy.speed = 0.9
		primaryColor := color.RGBA{180, 20, 60, 255} // Deep crimson
		secondaryColor := color.RGBA{40, 0, 0, 255}  // Dark accent for details
		enemy.color = primaryColor
		enemy.canFly = false
		enemy.canAttack = false
		enemy.sprite = createSpriteFromArt(spiderPixelArt, primaryColor, secondaryColor)

	case SnakeEnemy:
		enemy.speed = 1.1
		primaryColor := color.RGBA{40, 180, 40, 255}  // Vibrant green
		secondaryColor := color.RGBA{20, 60, 20, 255} // Dark green scales
		enemy.color = primaryColor
		enemy.canFly = false
		enemy.canAttack = true
		enemy.attackDamage = 0.5                    // Low damage
		enemy.attackRange = float64(cellSize) * 1.0 // 1 cell range
		enemy.attackRate = 60                       // Attack every 1 second (60 frames)
		enemy.attackDuration = 120                  // Stay in place for 2 seconds while attacking
		enemy.attackChance = 0.3                    // 30% chance to attack when in range
		enemy.sprite = createSpriteFromArt(snakePixelArt, primaryColor, secondaryColor)

	case HawkEnemy:
		enemy.speed = 1.4                             // Faster than snake (1.1) but must path around towers
		primaryColor := color.RGBA{200, 120, 20, 255} // Rich amber
		secondaryColor := color.RGBA{80, 40, 0, 255}  // Dark brown wings
		enemy.color = primaryColor
		enemy.canFly = false // Hawks must path around towers
		enemy.canAttack = false
		enemy.sprite = createSpriteFromArt(hawkPixelArt, primaryColor, secondaryColor)

	case GhoulEnemy:
		enemy.speed = 0.8
		primaryColor := color.RGBA{180, 180, 220, 255} // Pale ghostly color
		secondaryColor := color.RGBA{60, 60, 100, 255} // Dark ethereal color
		enemy.color = primaryColor
		enemy.canFly = false
		enemy.canAttack = true
		enemy.attackDamage = 0.8                    // Higher damage
		enemy.attackRange = float64(cellSize) * 1.5 // 1.5 cell range
		enemy.attackRate = 60                       // Attack every 1 second
		enemy.attackDuration = 180                  // Stay in place for 3 seconds while attacking
		enemy.attackChance = 0.25                   // 25% chance to attack when in range
		enemy.sprite = createSpriteFromArt(ghoulPixelArt, primaryColor, secondaryColor)

	case BlobEnemy:
		// Multiply base health by 5 for boss-type enemy
		enemy.health = startingHealth * 5
		enemy.maxHealth = enemy.health
		enemy.speed = 0.6 // Slower but more intimidating
		// More intense purple colors for boss
		primaryColor := color.RGBA{180, 0, 180, 255}   // Brighter purple
		secondaryColor := color.RGBA{100, 0, 100, 255} // Darker purple for depth
		enemy.color = primaryColor
		enemy.canFly = false
		enemy.canAttack = true
		enemy.attackDamage = 2.0                    // High base damage
		enemy.attackRange = float64(cellSize) * 2.0 // 2 cell range
		enemy.attackRate = 180                      // Attack every 3 seconds
		enemy.attackChance = 0.5                    // 50% chance to attack
		enemy.size = enemySize * 1.5                // 50% larger than normal enemies
		enemy.sprite = createSpriteFromArt(blobPixelArt, primaryColor, secondaryColor)

	}

	return enemy
}

// IsBossWave checks if the given wave number should spawn a boss
func IsBossWave(waveNumber int) bool {
	// Boss appears after wave 5 and then every 5 waves
	return waveNumber > 5 && waveNumber%5 == 0
}

// Update updates the enemy position and handles pathfinding
func (e *Enemy) Update(gameMap *GameMap) bool {
	// Handle frozen state
	if e.frozenTimer > 0 {
		e.frozenTimer--
		return false
	}

	// Check if we need to recalculate the path
	if e.pathInvalid || len(e.path) == 0 || e.pathIndex >= len(e.path) {
		if !e.findPath(gameMap) {
			return false
		}
	}

	// Move towards current target
	if e.pathIndex < len(e.path) {
		targetCell := e.path[e.pathIndex]
		e.targetX = float64(targetCell.X*gameMap.cellSize + gameMap.cellSize/2)
		e.targetY = float64(targetCell.Y*gameMap.cellSize + gameMap.cellSize/2 + gameMap.uiHeight)

		dx := e.targetX - e.x
		dy := e.targetY - e.y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < e.speed {
			// Reached target point
			e.x = e.targetX
			e.y = e.targetY
			e.pathIndex++
		} else {
			// Move towards target
			e.x += (dx / dist) * e.speed
			e.y += (dy / dist) * e.speed

			// Update movement animation
			e.moveTimer += 0.05 // Slower animation
			if e.moveTimer > 2*math.Pi {
				e.moveTimer -= 2 * math.Pi
			}
			// Calculate movement offset based on enemy type
			switch e.enemyType {
			case SpiderEnemy:
				// Tiny bobbing motion
				e.moveOffset = math.Sin(e.moveTimer) * 0.5
			case SnakeEnemy:
				// Subtle slithering
				e.moveOffset = math.Sin(e.moveTimer*1.5) * 0.8
			case HawkEnemy:
				// Very slight wave motion
				e.moveOffset = math.Sin(e.moveTimer*0.8) * 1.0
			case GhoulEnemy:
				// Subtle floating
				e.moveOffset = math.Sin(e.moveTimer*0.5) * 0.7
			}
		}
	}

	// Calculate grid position
	gridX := int((e.x - float64(gameMap.cellSize/2)) / float64(gameMap.cellSize))
	gridY := int((e.y - float64(gameMap.uiHeight)) / float64(gameMap.cellSize))

	// Check if we're too close to entrance to allow attacks
	if gridX < 4 { // No attacks in first 4 squares
		e.targetTower = nil // Clear any existing target
		return false
	}

	// If enemy can attack and has no target, look for towers to attack
	if e.canAttack && e.targetTower == nil {
		// Get current game enemies for coordination
		if currentGame != nil {
			gameMap.enemies = currentGame.enemies
		}

		// Random roll to decide if we look for a tower to attack
		if randFloat() < e.attackChance {
			// Get nearby towers and check if any are available for attack
			towers := gameMap.GetTowersInRange(e.x, e.y, e.attackRange)
			availableTowers := make([]*Tower, 0)

			// Filter out towers that are already being attacked
			for _, t := range towers {
				isBeingAttacked := false
				for _, other := range gameMap.enemies {
					if other != nil && other != e && other.targetTower == t && other.currentAttackTime > 0 {
						isBeingAttacked = true
						break
					}
				}
				if !isBeingAttacked {
					availableTowers = append(availableTowers, t)
				}
			}

			// Only attack if there are available towers
			if len(availableTowers) > 0 {
				// Randomly select one tower to attack
				e.targetTower = availableTowers[mathrand.Intn(len(availableTowers))]
				e.lastAttack = e.attackRate // Start with full cooldown
			}
		}
	}

	// If we have a target tower, try to attack it
	if e.targetTower != nil {
		// First verify tower still exists
		towerStillExists := false
		for _, t := range gameMap.towers {
			if t == e.targetTower {
				towerStillExists = true
				break
			}
		}
		if !towerStillExists {
			e.targetTower = nil
			e.currentAttackTime = 0
			e.pathInvalid = true
			return false
		}

		// Calculate tower's grid position
		towerGridX := e.targetTower.position.X
		if towerGridX < 4 { // Don't attack towers in first 4 squares
			e.targetTower = nil
			e.currentAttackTime = 0
			return false
		}

		// Calculate distance to tower
		towerX := float64(e.targetTower.position.X*gameMap.cellSize + gameMap.cellSize/2)
		towerY := float64(e.targetTower.position.Y*gameMap.cellSize + gameMap.cellSize/2 + gameMap.uiHeight)
		dx := towerX - e.x
		dy := towerY - e.y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist <= e.attackRange {
			// In range, handle attack sequence
			if e.currentAttackTime < e.attackDuration {
				// Still in attack animation
				e.currentAttackTime++

				// Only deal damage on specific intervals
				if e.lastAttack <= 0 {
					// Attack the tower
					if e.targetTower.TakeDamage(e.attackDamage) {
						// Tower was destroyed
						gameMap.RemoveTower(e.targetTower.position.X, e.targetTower.position.Y)
						e.targetTower = nil
						e.currentAttackTime = 0
						e.lastAttack = e.attackRate
					} else {
						// Continue attack sequence
						PlayAttackSound()
						e.lastAttack = e.attackRate
					}
				} else {
					e.lastAttack--
				}
				return false // Stay in place while attacking
			} else {
				// Attack sequence complete
				e.targetTower = nil
				e.currentAttackTime = 0
				e.lastAttack = e.attackRate
			}
		} else {
			// Tower out of range, stop attacking
			e.targetTower = nil
			e.currentAttackTime = 0
			e.lastAttack = e.attackRate
		}
	}

	// Check if we've reached the exit
	exitStart, exitEnd, exitX := gameMap.GetExitArea()

	// Get current game for enemy list reference
	if currentGame == nil {
		return false
	}
	gameMap.enemies = currentGame.enemies

	// Must be at the exit X position and within the exit Y range
	if gridX == exitX && gridY >= exitStart && gridY <= exitEnd {
		// Must be close to the center of the exit cell
		cellCenterX := float64(exitX*gameMap.cellSize + gameMap.cellSize/2)
		cellCenterY := float64(gridY*gameMap.cellSize + gameMap.cellSize/2 + gameMap.uiHeight)

		dx := e.x - cellCenterX
		dy := e.y - cellCenterY
		distanceToCenter := math.Sqrt(dx*dx + dy*dy)

		// Only count as reached if very close to center
		return distanceToCenter < 5.0
	}
	return false
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	if e.sprite != nil {
		// Draw sprite
		op := &ebiten.DrawImageOptions{}

		// Get the cell size from game map
		gameMap := GetGameMap()
		cellSize := float64(gameMap.cellSize)

		// Scale sprite to 80% of cell size for better visibility
		spriteSize := cellSize * 0.8 // Increased from 70% to 80%

		// Calculate sprite dimensions after scaling
		spriteW := float64(e.sprite.Bounds().Dx())
		spriteH := float64(e.sprite.Bounds().Dy())
		scale := spriteSize / math.Max(spriteW, spriteH)

		// Calculate position to center in cell
		scaledW := spriteW * scale
		scaledH := spriteH * scale

		// Calculate rotation angle based on movement direction
		dx := e.targetX - e.x
		dy := e.targetY - e.y
		moveAngle := math.Atan2(dy, dx)

		// Center the sprite in its position with movement animation
		op.GeoM.Scale(scale, scale)

		// Calculate offset direction perpendicular to movement
		perpX := math.Cos(moveAngle + math.Pi/2)
		perpY := math.Sin(moveAngle + math.Pi/2)

		// Apply movement offset if not frozen
		var offsetX, offsetY float64
		if e.frozenTimer <= 0 {
			offsetX = perpX * e.moveOffset
			offsetY = perpY * e.moveOffset

			// Add very subtle tilt based on movement
			if e.enemyType != GhoulEnemy { // Ghouls don't tilt
				tiltAngle := e.moveOffset * 0.02 // Much smaller tilt angle
				op.GeoM.Rotate(tiltAngle)
			}
		}

		op.GeoM.Translate(
			e.x-scaledW/2+offsetX,
			e.y-scaledH/2+offsetY,
		)

		// Update eye flash state
		e.updateEyeFlash()

		// If frozen, draw ice effect with improved visuals
		if e.frozenTimer > 0 {
			// First draw a more pronounced light blue border/glow
			borderOp := &ebiten.DrawImageOptions{}
			borderOp.GeoM.Scale(scale*1.15, scale*1.15) // Slightly larger glow
			borderOp.GeoM.Translate(
				e.x-(scaledW*1.15)/2,
				e.y-(scaledH*1.15)/2,
			)
			borderOp.ColorScale.Scale(0.7, 0.9, 1.0, 0.6) // Brighter ice blue, more visible
			screen.DrawImage(e.sprite, borderOp)

			// Then draw the main sprite with improved ice tint
			op.ColorScale.Scale(0.3, 0.5, 0.9, 1.0) // More vibrant ice blue
		}

		// Draw main sprite
		screen.DrawImage(e.sprite, op)

		// Draw eye flash effect if active
		if e.eyeFlashing {
			// Create a bright red eye glow
			glowOp := &ebiten.DrawImageOptions{}
			glowOp.GeoM.Scale(scale*1.1, scale*1.1) // Slightly larger for glow
			glowOp.GeoM.Translate(
				e.x-(scaledW*1.1)/2,
				e.y-(scaledH*1.1)/2,
			)
			glowOp.ColorScale.Scale(1.0, 0.0, 0.0, float32(e.eyeFlashTimer)/4.0) // Red glow, fading with timer
			screen.DrawImage(e.sprite, glowOp)
		}
	} else {
		// Fallback: draw colored rectangle if sprite not loaded
		vector.DrawFilledRect(screen,
			float32(e.x-e.size/2),
			float32(e.y-e.size/2),
			float32(e.size),
			float32(e.size),
			e.color,
			false)
	}

	// Draw health bar
	healthBarWidth := e.size
	healthBarHeight := 4.0
	healthBarY := e.y + e.size/2 + 2 // Position below enemy

	// Draw background (empty health bar)
	vector.DrawFilledRect(screen,
		float32(e.x-e.size/2),
		float32(healthBarY),
		float32(healthBarWidth),
		float32(healthBarHeight),
		color.RGBA{100, 0, 0, 255},
		false)

	// Draw filled portion based on current health
	healthPercent := e.health / e.maxHealth
	if healthPercent > 0 {
		vector.DrawFilledRect(screen,
			float32(e.x-e.size/2),
			float32(healthBarY),
			float32(healthBarWidth*healthPercent),
			float32(healthBarHeight),
			color.RGBA{0, 255, 0, 255},
			false)
	}

	// Draw attack indicator ONLY if actively attacking and dealing damage
	if e.targetTower != nil && e.currentAttackTime > 0 && e.lastAttack <= 0 {
		// Calculate distance to tower to verify we're in range
		gm := GetGameMap()
		if gm == nil {
			return
		}
		towerX := float64(e.targetTower.position.X*gm.cellSize + gm.cellSize/2)
		towerY := float64(e.targetTower.position.Y*gm.cellSize + gm.cellSize/2 + gm.uiHeight)
		dx := towerX - e.x
		dy := towerY - e.y
		dist := math.Sqrt(dx*dx + dy*dy)

		// Only show indicator if we're actually in range and attacking
		if dist <= e.attackRange {
			// Draw a red attack indicator that pulses with the attack cycle
			attackSize := float32(8.0 + float64(e.lastAttack%20)/20.0*4.0) // Size pulses between 8-12
			enemyX := float32(e.x)
			enemyY := float32(e.y)
			vector.DrawFilledCircle(screen,
				enemyX,
				enemyY+float32(e.size/2)+8, // Position below health bar
				attackSize/2,
				color.RGBA{255, 0, 0, 192}, // Semi-transparent red
				false)
		}
	}
}

// InvalidatePath marks the current path as invalid
func (e *Enemy) InvalidatePath() {
	e.pathInvalid = true
}

// updateEyeFlash handles the timing of eye flashing
func (e *Enemy) updateEyeFlash() {
	if e.eyeFlashing {
		// If currently flashing, decrease timer
		e.eyeFlashTimer--
		if e.eyeFlashTimer <= 0 {
			e.eyeFlashing = false
			e.eyeFlashTimer = 0
		}
	} else {
		// Not flashing, randomly start a flash
		if mathrand.Float64() < 0.01 { // 1% chance each frame
			e.eyeFlashing = true
			e.eyeFlashTimer = 4 // Flash for 4 frames
		}
	}
}

// findPath uses breadth-first search to find a path to the exit
func (e *Enemy) findPath(gameMap *GameMap) bool {
	// Calculate grid position for pathfinding
	currentX := int(e.x / float64(gameMap.cellSize))
	currentY := int((e.y - float64(gameMap.uiHeight)) / float64(gameMap.cellSize))

	// Create visited array and parent map for path reconstruction
	visited := make([][]bool, gameMap.height)
	parent := make([][][2]int, gameMap.height)
	for i := range visited {
		visited[i] = make([]bool, gameMap.width)
		parent[i] = make([][2]int, gameMap.width)
	}

	// Queue for BFS
	queue := [][2]int{{currentX, currentY}}
	visited[currentY][currentX] = true

	// Target is any point on the right edge within exit area
	exitStart, exitEnd, exitX := gameMap.GetExitArea()
	targetFound := false
	var targetX, targetY int

	// BFS
	for len(queue) > 0 && !targetFound {
		current := queue[0]
		queue = queue[1:]

		// Found a valid exit cell
		if current[0] == exitX && current[1] >= exitStart && current[1] <= exitEnd {
			targetX, targetY = current[0], current[1]
			targetFound = true
			break
		}

		// Try cardinal directions: right, down, up, left (prioritize moving right)
		directions := [][2]int{{1, 0}, {0, 1}, {0, -1}, {-1, 0}}
		for _, dir := range directions {
			nextX := current[0] + dir[0]
			nextY := current[1] + dir[1]

			// Check bounds and if cell is valid
			if nextX >= 0 && nextX < gameMap.width &&
				nextY >= 0 && nextY < gameMap.height &&
				!visited[nextY][nextX] &&
				(!gameMap.IsBlocked(nextX, nextY) || (e.canFly && nextX > currentX)) { // Flying only helps when moving forward
				queue = append(queue, [2]int{nextX, nextY})
				visited[nextY][nextX] = true
				parent[nextY][nextX] = current
			}
		}
	}

	if !targetFound {
		return false
	}

	// Reconstruct path
	path := make([]Point, 0)
	currentPos := [2]int{targetX, targetY}
	for currentPos[0] != currentX || currentPos[1] != currentY {
		path = append([]Point{{currentPos[0], currentPos[1]}}, path...)
		currentPos = parent[currentPos[1]][currentPos[0]]
	}

	e.path = path
	e.pathIndex = 0
	e.pathInvalid = false
	return true
}
