package game

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// currentGame holds the current game instance for global access
var currentGame *Game

func GetGameMap() *GameMap {
	if currentGame != nil {
		return currentGame.gameMap
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GameState represents the current state of the game
type GameState int

const (
	BuildState GameState = iota
	PlayState
	PausedState
	GameOverState
)

// Game represents the main game state
type Game struct {
	gameMap         *GameMap
	enemies         []*Enemy
	projectiles     []*Projectile     // Active projectiles
	deathAnims      []*DeathAnimation // Death animations
	towerButtons    []*TowerButton    // Tower selection buttons
	score           int
	lives           int
	confirmingReset bool
	money           int // Points available for tower placement
	gameState       GameState
	spawnTimer      int
	spawnInterval   int
	currentWave     int
	enemiesInWave   int
	enemiesSpawned  int
	waveType        EnemyType
	startButton     Button
	pauseButton     Button
	selectedTower   TowerType // Currently selected tower type
	mouseX, mouseY  int       // Current mouse position for tower preview
}

// Button represents a clickable button
type Button struct {
	x, y, width, height int
	text                string
	color               color.Color
	hovered             bool
}

// NewGame creates a new game instance
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())

	// Create start button in left section
	startBtn := Button{
		x:      20,  // Left margin
		y:      15,  // Centered in 60px UI height
		width:  160, // Even wider to fit text
		height: 30,  // Standard height
		text:   "Begin!",
		color:  color.RGBA{0, 200, 0, 255},
	}

	// Create pause button next to start button
	pauseBtn := Button{
		x:      190, // Next to start button with small gap
		y:      15,  // Aligned with start button
		width:  160, // Same width as start button
		height: 30,  // Standard height
		text:   "Pause",
		color:  color.RGBA{200, 200, 0, 255},
	}

	// Create tower buttons at bottom of screen
	towerButtons := make([]*TowerButton, 0)
	btnX := 250       // Starting X position
	btnY := 832 - 90  // 90 pixels from bottom of new height
	btnSpacing := 100 // Space between buttons

	// Create a button for each tower type
	towerTypes := []TowerType{DartTower, BulletTower, LightningTower, FlameTower, FreezeTower, ForkTower}
	for _, tType := range towerTypes {
		btn := NewTowerButton(tType, btnX, btnY)
		// Set initial selection
		if tType == DartTower {
			btn.selected = true
		}
		towerButtons = append(towerButtons, btn)
		btnX += btnSpacing
	}

	game := &Game{
		gameMap:        NewGameMap(),
		enemies:        make([]*Enemy, 0),
		projectiles:    make([]*Projectile, 0),
		deathAnims:     make([]*DeathAnimation, 0),
		towerButtons:   towerButtons,
		score:          0,
		lives:          20,
		money:          200, // Starting points - enough for any basic tower setup
		gameState:      BuildState,
		spawnTimer:     0,
		spawnInterval:  60, // Start with 1 second between spawns
		currentWave:    0,
		enemiesInWave:  10, // 10 enemies per wave
		enemiesSpawned: 0,
		waveType:       SpiderEnemy, // Start with spiders
		startButton:    startBtn,
		pauseButton:    pauseBtn,
		selectedTower:  DartTower, // Default to dart tower
	}

	currentGame = game
	return game
}

// Update updates the game state
func (g *Game) Update() error {
	// Update mouse position
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Handle mouse position for button hover
	mouseX, mouseY := ebiten.CursorPosition()
	g.startButton.hovered = g.startButton.contains(mouseX, mouseY)
	g.pauseButton.hovered = g.pauseButton.contains(mouseX, mouseY)

	// Handle tower selection clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Check tower buttons first
		for _, btn := range g.towerButtons {
			if btn.Contains(mouseX, mouseY) && g.money >= btn.cost {
				g.selectedTower = btn.tower
				// Update selected state for all buttons
				for _, otherBtn := range g.towerButtons {
					otherBtn.selected = (otherBtn.tower == btn.tower)
				}
				return nil
			}
		}

		// Handle other button clicks
		if g.startButton.contains(mouseX, mouseY) {
			if g.gameState == BuildState {
				g.gameState = PlayState
				g.startButton.text = "Reset"
				g.confirmingReset = false
			} else if g.gameState == PlayState || g.gameState == PausedState {
				if g.confirmingReset {
					// Actually reset the game
					g.gameState = BuildState
					g.enemies = make([]*Enemy, 0)
					g.projectiles = make([]*Projectile, 0)
					g.deathAnims = make([]*DeathAnimation, 0)
					g.lives = 20
					g.money = 200 // Reset to initial money amount
					g.currentWave = 0
					g.enemiesInWave = 10
					g.enemiesSpawned = 0
					// Add this check before setting wave type
					if g.waveType == BlobEnemy {
						g.waveType = SpiderEnemy // Reset to normal enemy type if resetting during boss wave
					} else {
						g.waveType = SpiderEnemy // Normal reset
					}
					g.spawnInterval = 60
					g.startButton.text = "Begin!" // Reset to initial text
					g.pauseButton.text = "Pause"  // Reset pause button text
					g.confirmingReset = false
					g.score = 0
					// Clear all towers from the map
					g.gameMap = NewGameMap()
					// Reset selected tower to default
					g.selectedTower = DartTower
					// Reset tower button selection
					for _, btn := range g.towerButtons {
						btn.selected = (btn.tower == DartTower)
					}
				} else {
					// Ask for confirmation
					g.confirmingReset = true
					g.startButton.text = "Sure?"
				}
			}
			return nil
		}

		if g.pauseButton.contains(mouseX, mouseY) {
			if g.gameState == PlayState {
				g.gameState = PausedState
				g.pauseButton.text = "Continue"
				g.confirmingReset = false // Cancel reset confirmation when pausing
			} else if g.gameState == PausedState {
				g.gameState = PlayState
				g.pauseButton.text = "Pause"
				g.confirmingReset = false // Cancel reset confirmation when unpausing
			}
			return nil
		}
	}

	// Handle tower placement
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if mouseY > g.gameMap.uiHeight { // Don't place towers in the UI area
			gridX, gridY := g.gameMap.GetGridPosition(float64(mouseX), float64(mouseY))
			if err := g.tryPlaceTower(gridX, gridY); err != nil {
				log.Printf("Tower placement failed: %v", err)
			}
		}
	}

	// Handle tower removal
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if mouseY > g.gameMap.uiHeight { // Don't remove towers in the UI area
			gridX, gridY := g.gameMap.GetGridPosition(float64(mouseX), float64(mouseY))
			// Get the tower before removing it to calculate refund
			if tower := g.gameMap.GetTowerAt(gridX, gridY); tower != nil {
				// Calculate refund based on tower type and remaining health
				var baseRefund int
				switch tower.towerType {
				case DartTower:
					baseRefund = 10
				case BulletTower:
					baseRefund = 25
				case LightningTower:
					baseRefund = 40
				case FlameTower:
					baseRefund = 60
				case FreezeTower:
					baseRefund = 75
				case ForkTower:
					baseRefund = 150
				}

				// Calculate actual refund based on remaining health percentage
				healthPercent := tower.health / tower.maxHealth
				actualRefund := int(float64(baseRefund) * healthPercent)

				g.money += actualRefund
				g.gameMap.RemoveTower(gridX, gridY)
			}
		}
	}

	// Handle game state logic
	if g.gameState == BuildState || g.gameState == PausedState {
	}

	// Update game logic only in play state
	if g.gameState == PlayState {
		// Update enemies
		remainingEnemies := make([]*Enemy, 0, len(g.enemies))
		for _, enemy := range g.enemies {
			if enemy == nil {
				continue
			}

			reached := enemy.Update(g.gameMap)
			if reached {
				g.lives--
				if g.lives <= 0 {
					g.gameState = GameOverState
				}
			} else if enemy.health <= 0 {
				// Enemy was killed, create death animation with enemy sprite
				g.deathAnims = append(g.deathAnims, NewDeathAnimation(enemy.x, enemy.y, enemy.size, enemy.sprite))

				// death sound
				PlayEnemyDeathSound()

				// Award points - extra for boss
				reward := 2 * enemy.level // Base: 2 points per level
				if enemy.enemyType == BlobEnemy {
					reward = reward * 5 // 5x reward for boss
					g.score += 1000     // Bonus score for boss kill
				}
				g.money += reward
			} else {
				remainingEnemies = append(remainingEnemies, enemy)
			}
		}
		g.enemies = remainingEnemies

		// Update death animations
		remainingAnims := make([]*DeathAnimation, 0)
		for _, anim := range g.deathAnims {
			if anim.Update() {
				remainingAnims = append(remainingAnims, anim)
			}
		}
		g.deathAnims = remainingAnims

		// Update towers and generate projectiles
		for _, tower := range g.gameMap.towers {
			newProjectiles := tower.Update(g.enemies)
			if newProjectiles != nil {
				g.projectiles = append(g.projectiles, newProjectiles...)
			}
		}

		// Update projectiles and check for hits
		remainingProjectiles := make([]*Projectile, 0, len(g.projectiles))
		for _, proj := range g.projectiles {
			hit := proj.Update()
			if !hit {
				// Keep projectile if it hasn't hit
				remainingProjectiles = append(remainingProjectiles, proj)
			} else {
				// Check for collision with enemies
				for _, enemy := range g.enemies {
					if enemy == nil {
						continue
					}

					// Simple collision check
					ex, ey := enemy.x, enemy.y
					px, py := proj.GetPosition()
					dx := ex - px
					dy := ey - py
					dist := math.Sqrt(dx*dx + dy*dy)

					if dist < enemy.size/2 { // If within enemy radius
						if proj.GetProjectileType() == FreezeProjectile {
							if enemy.frozenTimer <= 0 { // Only freeze if not already frozen
								enemy.frozenTimer = 60 // Freeze for 1 second (60 frames)
							}
						} else {
							enemy.health -= proj.GetDamage()
						}
						break // Exit after first hit
					}
				}
			}
		}
		g.projectiles = remainingProjectiles

		// Handle wave spawning
		if g.enemiesSpawned < g.enemiesInWave {
			g.spawnTimer++
			if g.spawnTimer >= g.spawnInterval {
				g.spawnTimer = 0
				// Random spawn intervals between 1.5-3 seconds (90-180 frames)
				g.spawnInterval = 90 + rand.Intn(90)

				// After level 5, sometimes spawn a different enemy type
				spawnType := g.waveType
				if g.currentWave >= 4 && g.enemiesSpawned > 0 { // Start at wave 5 (index 4)
					if rand.Float64() < 0.05 { // 5% chance for different enemy
						// Create list of enemy types excluding current wave type
						availableTypes := []EnemyType{}
						for _, t := range []EnemyType{SpiderEnemy, SnakeEnemy, HawkEnemy, GhoulEnemy} {
							if t != g.waveType {
								availableTypes = append(availableTypes, t)
							}
						}
						if len(availableTypes) > 0 {
							spawnType = availableTypes[rand.Intn(len(availableTypes))]
						}
					}
				}

				// Spawn new enemy with current wave's colors
				entranceStart, entranceEnd, _ := g.gameMap.GetEntranceArea()
				randomY := entranceStart + rand.Intn(entranceEnd-entranceStart+1)
				newEnemy := NewEnemyWithColor(randomY, g.gameMap.cellSize, g.gameMap.uiHeight,
					spawnType, g.currentWave+1, g.waveType)
				if newEnemy != nil {
					g.enemies = append(g.enemies, newEnemy)
					g.enemiesSpawned++
				}
			}

		} else if len(g.enemies) == 0 {
			// Wave completed
			g.currentWave++
			g.enemiesSpawned = 0

			if IsBossWave(g.currentWave) {
				// Boss wave - spawn a single powerful blob enemy
				g.waveType = BlobEnemy
				g.enemiesInWave = 1 // Only one boss
			} else {
				// Normal wave - cycle through enemy types
				switch g.waveType {
				case SpiderEnemy:
					g.waveType = SnakeEnemy
				case SnakeEnemy:
					g.waveType = HawkEnemy
				case HawkEnemy:
					g.waveType = GhoulEnemy
				case GhoulEnemy:
					g.waveType = SpiderEnemy
				case BlobEnemy:
					g.waveType = SpiderEnemy // Reset to normal enemy type if resetting during boss wave
				}
				// Regular wave size calculation
				g.enemiesInWave = 10 + g.currentWave*2
			}

			// Keep spawn interval adjustment
			g.spawnInterval = max(60, 120-g.currentWave*10) // Minimum 1 second between spawns
		}
	}

	return nil
}

// Draw draws the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	if screen == nil {
		return
	}

	// Draw map
	g.gameMap.Draw(screen)

	// Draw enemies if not in build state
	if g.gameState != BuildState {
		for _, enemy := range g.enemies {
			if enemy != nil {
				enemy.Draw(screen)
			}
		}

		// Draw death animations
		for _, anim := range g.deathAnims {
			if anim != nil {
				anim.Draw(screen)
			}
		}

		// Draw projectiles
		for _, proj := range g.projectiles {
			if proj != nil {
				proj.Draw(screen)
			}
		}
	}

	// LEFT SECTION (0-320px) - Buttons and game state
	// Draw start button
	buttonColor := g.startButton.color
	if g.startButton.hovered {
		buttonColor = color.RGBA{0, 255, 0, 255}
	}
	vector.DrawFilledRect(
		screen,
		float32(g.startButton.x),
		float32(g.startButton.y),
		float32(g.startButton.width),
		float32(g.startButton.height),
		buttonColor,
		true,
	)

	// Draw pause button if game has started
	if g.gameState != BuildState {
		buttonColor = g.pauseButton.color
		if g.pauseButton.hovered {
			buttonColor = color.RGBA{255, 255, 0, 255}
		}
		vector.DrawFilledRect(
			screen,
			float32(g.pauseButton.x),
			float32(g.pauseButton.y),
			float32(g.pauseButton.width),
			float32(g.pauseButton.height),
			buttonColor,
			true,
		)
	}

	// Draw button texts centered
	textWidth := MeasureTextWidth(g.startButton.text, false)

	DrawText(screen, g.startButton.text,
		g.startButton.x+(g.startButton.width-textWidth)/2,
		g.startButton.y+20, color.Black)

	if g.gameState != BuildState {
		textWidth = MeasureTextWidth(g.pauseButton.text, false)
		DrawText(screen, g.pauseButton.text,
			g.pauseButton.x+(g.pauseButton.width-textWidth)/2,
			g.pauseButton.y+20, color.Black)
	}

	// Game state text well below buttons using small font
	stateText := ""
	switch g.gameState {
	case BuildState:
		stateText = "Building Mode"
	case PlayState:
		stateText = "Wave in Progress"
	case PausedState:
		stateText = "Game Paused"
	case GameOverState:
		stateText = "Game Over"
	}
	DrawSmallText(screen, stateText, 20, 58, color.White) // Even lower and smaller font

	// RIGHT SECTION (640-1024px) - Points and Lives with right-justified numbers
	DrawText(screen, "Points:", 750, 20, color.White)
	DrawText(screen, "Lives:", 750, 40, color.White)

	// Right-justify the numbers at x=900
	moneyText := fmt.Sprintf("%d", g.money)
	livesText := fmt.Sprintf("%d", g.lives)
	moneyWidth := len(moneyText) * 12 // Approximate width for right justification
	livesWidth := len(livesText) * 12 // Approximate width for right justification

	DrawText(screen, moneyText, 900-moneyWidth, 20, color.White) // Points value
	DrawText(screen, livesText, 900-livesWidth, 40, color.White) // Lives value

	// MIDDLE SECTION (320-640px) - Wave Information
	if g.gameState != BuildState {
		waveTypeText := ""
		switch g.waveType {
		case SpiderEnemy:
			waveTypeText = "Spiders"
		case SnakeEnemy:
			waveTypeText = "Snakes"
		case HawkEnemy:
			waveTypeText = "Hawks"
		case GhoulEnemy:
			waveTypeText = "Ghouls"
		case BlobEnemy:
			waveTypeText = "BOSS"
		}
		waveText := fmt.Sprintf("Wave %d", g.currentWave+1)
		enemyInfo := fmt.Sprintf("%s: %d/%d", waveTypeText, g.enemiesSpawned, g.enemiesInWave)

		// Center wave info
		if IsBossWave(g.currentWave+1) && len(g.enemies) == 0 {
			// Draw warning text in red when next wave will be boss
			warningText := "! BOSS INCOMING !"
			warningWidth := MeasureTextWidth(warningText, false)
			DrawText(screen, warningText,
				400-warningWidth/2, 60,
				color.RGBA{255, 0, 0, 255}) // Bright red
		}

		DrawText(screen, waveText, 400, 20, color.White)  // Wave number at top
		DrawText(screen, enemyInfo, 400, 40, color.White) // Enemy info below
	}

	// Draw tower selection buttons
	for _, btn := range g.towerButtons {
		btn.Draw(screen, g.money >= btn.cost)
	}

	// Draw tower range preview during build or pause states
	if g.gameState == BuildState || g.gameState == PausedState {
		g.drawTowerRangePreview(screen)
	}

	// Show range for clicked tower
	if g.mouseY > g.gameMap.uiHeight {
		gridX, gridY := g.gameMap.GetGridPosition(float64(g.mouseX), float64(g.mouseY))
		if tower := g.gameMap.GetTowerAt(gridX, gridY); tower != nil {
			tower.Draw(screen, true) // Pass true to show range
		}
	}

	// Draw game over screen if dead
	if g.gameState == GameOverState {
		if clickHandled := drawGameOver(screen); clickHandled {
			return // Skip processing other input if game over screen handled it
		}
	}
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1024, 832 // 12 rows * 56px + 60px UI + 100px tower selection = 832px
}

// PlaceTower attempts to place a tower and returns an error if it fails
func (g *Game) tryPlaceTower(x, y int) error {
	// Calculate tower cost based on type/level
	var towerCost int
	switch g.selectedTower {
	case DartTower:
		towerCost = 10 // Basic tower
	case BulletTower:
		towerCost = 25 // Better range and speed
	case LightningTower:
		towerCost = 40 // High damage
	case FlameTower:
		towerCost = 60 // Area damage
	case FreezeTower:
		towerCost = 75 // Most expensive basic tower
	case ForkTower:
		towerCost = 150 // Electric fork tower with large range
	}

	// Check if we have enough points
	if g.money < towerCost {
		return fmt.Errorf("not enough points: need %d, have %d", towerCost, g.money)
	}

	if g.gameMap.PlaceTower(x, y) {
		// Tower was placed successfully, deduct points
		g.money -= towerCost
		// Force all enemies to recalculate their paths
		for _, enemy := range g.enemies {
			if enemy != nil {
				enemy.InvalidatePath()
			}
		}
		return nil
	}
	return fmt.Errorf("cannot place tower at position %d,%d", x, y)
}

// Button.contains checks if a point is inside the button
func (b *Button) contains(x, y int) bool {
	return x >= b.x && x <= b.x+b.width &&
		y >= b.y && y <= b.y+b.height
}
