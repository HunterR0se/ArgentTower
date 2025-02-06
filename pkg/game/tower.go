package game

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TowerType represents different types of towers
type TowerType int

const (
	DartTower TowerType = iota
	BulletTower
	LightningTower
	FlameTower
	FreezeTower
	ForkTower
)

// Tower represents a defensive tower
type Tower struct {
	position        Point
	towerType       TowerType
	damage          float64
	attackRange     float64
	fireRate        float64
	cost            int
	lastShot        float64
	level           int
	sprite          *ebiten.Image
	canFireDiagonal bool
	health          float64
	maxHealth       float64
	damageState     int
	underAttack     int
}

var forkTowerCount = 0

// Draw draws the tower
func (t *Tower) Draw(screen *ebiten.Image, selected bool) {
	gameMap := GetGameMap()
	if gameMap == nil {
		return
	}

	// Calculate tower position
	cellSize := float64(gameMap.cellSize)
	x := float64(t.position.X) * cellSize
	y := float64(t.position.Y) * cellSize + float64(gameMap.uiHeight)

	// Draw range circle if selected
	if selected {
		centerX := x + cellSize/2
		centerY := y + cellSize/2
		vector.StrokeCircle(screen,
			float32(centerX),
			float32(centerY),
			float32(t.attackRange),
			1.5,
			color.RGBA{160, 160, 160, 100},
			true)
	}

	if t.sprite == nil {
		t.sprite = getTowerSprite(t.towerType)
	}

	if t.sprite != nil {
		op := &ebiten.DrawImageOptions{}

		// Scale sprite to fit cell (slightly smaller)
		spriteW := float64(t.sprite.Bounds().Dx())
		spriteH := float64(t.sprite.Bounds().Dy())
		scale := (cellSize * 0.9) / math.Max(spriteW, spriteH)
		op.GeoM.Scale(scale, scale)

		// Center in cell
		scaledW := spriteW * scale
		scaledH := spriteH * scale
		op.GeoM.Translate(
			x + (cellSize-scaledW)/2,
			y + (cellSize-scaledH)/2,
		)

		// Apply red flash if under attack
		if t.underAttack > 0 {
			// Get the original color from the tower sprite
			origR, origG, origB, _ := t.GetTowerColor()
			// Calculate brightness using perceived luminance
			brightness := (float64(origR)*0.299 + float64(origG)*0.587 + float64(origB)*0.114) / 255.0
			// Keep same brightness but shift to red
			redValue := uint8(brightness * 255)
			op.ColorScale.Scale(float32(redValue)/255.0, 0, 0, 1.0)
			t.underAttack--
		}

		screen.DrawImage(t.sprite, op)

		// Draw damage overlay
		t.DrawDamageOverlay(screen, x, y, cellSize)

		// Draw small skull if under attack
		if t.underAttack > 0 {
			skullSprite := createSpriteFromArt(SharedSkull, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 0, 0})
			if skullSprite != nil {
				skullW := float64(skullSprite.Bounds().Dx())
				skullH := float64(skullSprite.Bounds().Dy())
				skullScale := (cellSize * 0.3) / skullW  // Slightly larger

				skullOp := &ebiten.DrawImageOptions{}
				skullOp.GeoM.Scale(skullScale, skullScale)
				skullOp.GeoM.Translate(
					x + cellSize - (skullW * skullScale),
					y + cellSize - (skullH * skullScale),
				)

				screen.DrawImage(skullSprite, skullOp)
			}
		}
	}
}

// Rest of the tower.go code remains unchanged
func (t *Tower) Update(enemies []*Enemy) []*Projectile {
	var closestEnemy *Enemy
	closestDist := t.attackRange

	gameMap := GetGameMap()
	if gameMap == nil {
		return nil
	}

	// Calculate tower center position
	towerX := float64(t.position.X*gameMap.cellSize + gameMap.cellSize/2)
	towerY := float64(t.position.Y*gameMap.cellSize + gameMap.cellSize/2 + gameMap.uiHeight)

	for _, enemy := range enemies {
		if enemy == nil {
			continue
		}

		dx := enemy.x - towerX
		dy := enemy.y - towerY
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist >= closestDist {
			continue
		}

		if !t.canFireDiagonal {
			angle := math.Atan2(dy, dx) * 180 / math.Pi
			if angle < 0 {
				angle += 360
			}

			isCardinal := false
			for _, cardinal := range []float64{0, 90, 180, 270} {
				angleDiff := math.Abs(angle - cardinal)
				if angleDiff <= 22.5 || angleDiff >= 337.5 {
					isCardinal = true
					break
				}
			}

			if !isCardinal {
				continue
			}
		}

		closestDist = dist
		closestEnemy = enemy
	}

	if closestEnemy != nil && t.canShoot() {
		var projType ProjectileType
		switch t.towerType {
		case DartTower:
			projType = DartProjectile
		case BulletTower:
			projType = BulletProjectile
		case LightningTower:
			projType = LightningProjectile
		case FlameTower:
			projType = FlameProjectile
		case FreezeTower:
			projType = FreezeProjectile
		}

		proj := NewProjectile(
			towerX,
			towerY,
			closestEnemy.x,
			closestEnemy.y,
			projType,
			t.damage,
		)

		t.lastShot = float64(time.Now().UnixNano()) / 1e9
		PlayTowerShootSound()
		return []*Projectile{proj}
	}

	return nil
}

func (t *Tower) canShoot() bool {
	currentTime := float64(time.Now().UnixNano()) / 1e9
	return currentTime - t.lastShot >= 1.0/t.fireRate
}

func (t *Tower) GetPosition() Point {
	return t.position
}

func NewTower(towerType TowerType, x, y int) *Tower {
	if towerType == ForkTower && forkTowerCount >= 10 {
		return nil
	}

	tower := &Tower{
		position:  Point{x, y},
		towerType: towerType,
		level:     1,
		health:    100,
		maxHealth: 100,
	}

	gameMap := GetGameMap()
	if gameMap == nil {
		return nil
	}
	cellSize := gameMap.cellSize

	switch towerType {
	case DartTower:
		tower.level = 1
		tower.damage = 1.0
		tower.attackRange = 2 * float64(cellSize)
		tower.fireRate = 1.0
		tower.sprite = getTowerSprite(DartTower)
		tower.canFireDiagonal = true
	case BulletTower:
		tower.level = 2
		tower.damage = 1.0 * 1.1
		tower.attackRange = 3 * float64(cellSize)
		tower.fireRate = 1.2
		tower.sprite = getTowerSprite(BulletTower)
		tower.canFireDiagonal = true
	case LightningTower:
		tower.level = 3
		tower.damage = 1.0 * 1.2
		tower.attackRange = 3 * float64(cellSize)
		tower.fireRate = 0.8
		tower.sprite = getTowerSprite(LightningTower)
		tower.canFireDiagonal = true
	case FlameTower:
		tower.level = 4
		tower.damage = 1.0 * 1.3
		tower.attackRange = 2 * float64(cellSize)
		tower.fireRate = 3.0
		tower.sprite = getTowerSprite(FlameTower)
		tower.canFireDiagonal = false
	case FreezeTower:
		tower.level = 5
		tower.damage = 0
		tower.attackRange = 2 * float64(cellSize)
		tower.fireRate = 1.0
		tower.sprite = getTowerSprite(FreezeTower)
		tower.canFireDiagonal = true
		tower.cost = 100
	case ForkTower:
		tower.level = 6
		tower.damage = 2.0
		tower.attackRange = 4 * float64(cellSize)
		tower.fireRate = 1.5
		tower.sprite = getTowerSprite(ForkTower)
		tower.canFireDiagonal = true
		tower.cost = 150
	}

	if towerType == ForkTower {
		forkTowerCount++
	}

	return tower
}

func (t *Tower) TakeDamage(damage float64) bool {
	t.health = math.Max(0, t.health - damage)  // Prevent negative health
	t.underAttack = 30  // Flash duration (0.5 seconds at 60fps)

	// Return true if tower is destroyed
	return t.health <= 0
}

func (t *Tower) DrawDamageOverlay(screen *ebiten.Image, x, y, size float64) {
	// Only draw damage effects if tower is damaged
	if t.health < t.maxHealth {
		// Calculate how damaged the tower is
		damagePercent := 1.0 - (t.health / t.maxHealth)

		// Use consistent random seeding based on tower position
		seed := int64(t.position.X*1000 + t.position.Y)
		r := rand.New(rand.NewSource(seed))

		// Create a grid of cells that will be potentially removed
		gridSize := 5 // 5x5 grid = 25 cells
		cellSize := size / float64(gridSize)

		// Calculate how many cells should be removed based on damage
		totalCells := gridSize * gridSize
		cellsToRemove := int(float64(totalCells) * damagePercent)

		// Create a list of all possible cell positions
		type Cell struct{ x, y int }
		cells := make([]Cell, 0, totalCells)
		for i := 0; i < gridSize; i++ {
			for j := 0; j < gridSize; j++ {
				cells = append(cells, Cell{i, j})
			}
		}

		// Shuffle the cells
		for i := len(cells) - 1; i > 0; i-- {
			j := r.Intn(i + 1)
			cells[i], cells[j] = cells[j], cells[i]
		}

		// Remove the first cellsToRemove cells
		for i := 0; i < cellsToRemove; i++ {
			cell := cells[i]
			cellX := x + float64(cell.x)*cellSize
			cellY := y + float64(cell.y)*cellSize

			// Draw a black rectangle to remove this cell
			vector.DrawFilledRect(screen,
				float32(cellX),
				float32(cellY),
				float32(cellSize),
				float32(cellSize),
				color.RGBA{0, 0, 0, 255},
				false)

			// Add cracks around the removed cell
			crackCount := 2 + r.Intn(3)
			for c := 0; c < crackCount; c++ {
				startX := cellX + r.Float64()*cellSize
				startY := cellY + r.Float64()*cellSize
				angle := r.Float64() * 2 * math.Pi
				length := cellSize * (0.5 + r.Float64()*0.5)
				endX := startX + math.Cos(angle)*length
				endY := startY + math.Sin(angle)*length

				vector.StrokeLine(screen,
					float32(startX),
					float32(startY),
					float32(endX),
					float32(endY),
					1.0,
					color.RGBA{0, 0, 0, 255},
					false)
			}
		}

		// Add some small debris particles near removed cells
		debrisCount := int(damagePercent * 10)
		for i := 0; i < debrisCount; i++ {
			if len(cells) > 0 && i < cellsToRemove {
				cell := cells[i]
				debrisX := x + (float64(cell.x) + r.Float64())*cellSize
				debrisY := y + (float64(cell.y) + r.Float64())*cellSize

				vector.DrawFilledCircle(screen,
					float32(debrisX),
					float32(debrisY),
					float32(2 + r.Float32()*2),
					color.RGBA{0, 0, 0, 255},
					false)
			}
		}
	}
}

func (t *Tower) GetTowerColor() (uint8, uint8, uint8, uint8) {
	switch t.towerType {
	case DartTower:
		return 200, 150, 80, 255
	case BulletTower:
		return 240, 190, 90, 255
	case LightningTower:
		return 50, 200, 255, 255
	case FlameTower:
		return 255, 100, 50, 255
	case FreezeTower:
		return 140, 220, 255, 255
	case ForkTower:
		return 20, 255, 200, 255
	default:
		return 200, 200, 200, 255
	}
}

func drawCracks(screen *ebiten.Image, x, y, size float64, crackColor color.Color, count int) {
	for i := 0; i < count; i++ {
		vector.StrokeLine(screen,
			float32(x + (size * float64(i) / float64(count))),
			float32(y + (size * float64(i) / float64(count))),
			float32(x + 10 + (size * float64(i) / float64(count))),
			float32(y + 10 + (size * float64(i) / float64(count))),
			1.0,
			crackColor,
			false)
	}
}
