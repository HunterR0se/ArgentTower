package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"bytes"
	"image"
	"math"
)

// GameMap represents the game map
type GameMap struct {
	width, height    int
	cellSize        int
	terrain         [][]TerrainType
	towers          []*Tower
	enemies         []*Enemy  // Track current enemies for attack coordination
	gridLineColor   color.Color
	borderColor     color.Color
	backgroundColor color.Color
	entrance        Point
	entranceStart   int
	entranceEnd     int
	exit           Point
	exitStart      int
	exitEnd        int
	uiHeight       int    // Height of UI area at top
	gridOffsetX    float64 // Horizontal offset to center grid
	backgroundImg   *ebiten.Image
}

// NewGameMap creates a new game map
func NewGameMap() *GameMap {
	// Load embedded background image
	img, _, err := image.Decode(bytes.NewReader(embeddedBackground))
	if err != nil {
		panic(err)
	}
	ebitenImg := ebiten.NewImageFromImage(img)
	m := &GameMap{
		width:           18,             // 18 cells wide
		height:          12,             // 12 cells tall
		cellSize:        56,             // 56px square cells (18 * 56 = 1008px width)
		gridLineColor:   color.RGBA{35, 35, 35, 255},
		borderColor:     color.RGBA{60, 60, 60, 255},
		backgroundColor: color.RGBA{20, 20, 20, 255},
		towers:          make([]*Tower, 0),
		enemies:         make([]*Enemy, 0),  // Initialize enemies slice
		uiHeight:        60,             // 60px UI height
		gridOffsetX:     8,              // Center the grid (1024 - 18*56 = 16px/2 = 8)
		backgroundImg:   ebitenImg,
	}
	
	// Initialize terrain - all cells are empty by default
	m.terrain = make([][]TerrainType, m.height)
	for i := range m.terrain {
		m.terrain[i] = make([]TerrainType, m.width)
		for j := range m.terrain[i] {
			m.terrain[i][j] = Empty
		}
	}

	// Calculate entrance and exit areas (1 cell high, centered)
	m.entranceStart = m.height / 2
	m.entranceEnd = m.entranceStart
	m.entrance = Point{0, m.height / 2}
	
	m.exitStart = m.entranceStart
	m.exitEnd = m.entranceEnd
	m.exit = Point{m.width - 1, m.height / 2}
	
	return m
}

// Draw draws the map
func (m *GameMap) Draw(screen *ebiten.Image) {
	// Draw background image at 20% visibility
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.1) // Set alpha to 10%
	screen.DrawImage(m.backgroundImg, op)
	
	// Draw UI background slightly darker
	vector.DrawFilledRect(
		screen,
		0,
		0,
		float32(m.width*m.cellSize),
		float32(m.uiHeight),
		color.RGBA{15, 15, 15, 255},
		false)
	
	// Draw grid lines
	for i := 0; i <= m.width; i++ {
		x := m.gridOffsetX + float64(i * m.cellSize)
		if i > 0 && i < m.width {
			vector.StrokeLine(
				screen,
				float32(x),
				float32(m.uiHeight),
				float32(x),
				float32(m.uiHeight+m.height*m.cellSize),
				1,
				m.gridLineColor,
				false)
		}
	}
	
	for i := 0; i <= m.height; i++ {
		y := float64(i*m.cellSize + m.uiHeight)
		if i > 0 && i < m.height {
			vector.StrokeLine(
				screen,
				float32(m.gridOffsetX),
				float32(y),
				float32(m.gridOffsetX + float64(m.width*m.cellSize)),
				float32(y),
				1,
				m.gridLineColor,
				false)
		}
	}

	// Draw borders with gaps for entrance and exit
	borderThickness := float32(2)

	// Top border (below UI)
	vector.DrawFilledRect(
		screen,
		float32(m.gridOffsetX),
		float32(m.uiHeight),
		float32(m.width*m.cellSize),
		borderThickness,
		m.borderColor,
		false)
	
	// Bottom border
	vector.DrawFilledRect(
		screen,
		float32(m.gridOffsetX),
		float32(m.uiHeight+m.height*m.cellSize)-borderThickness,
		float32(m.width*m.cellSize),
		borderThickness,
		m.borderColor,
		false)
	
	// Left border (with gap for entrance)
	vector.DrawFilledRect(
		screen,
		float32(m.gridOffsetX),
		float32(m.uiHeight),
		borderThickness,
		float32(m.entranceStart*m.cellSize),
		m.borderColor,
		false)
	vector.DrawFilledRect(
		screen,
		float32(m.gridOffsetX),
		float32(m.uiHeight+(m.entranceEnd+1)*m.cellSize),
		borderThickness,
		float32((m.height-m.entranceEnd-1)*m.cellSize),
		m.borderColor,
		false)
	
	// Right border (with gap for exit)
	vector.DrawFilledRect(
		screen,
		float32(m.gridOffsetX+float64(m.width*m.cellSize)-float64(borderThickness)),
		float32(m.uiHeight),
		borderThickness,
		float32(m.exitStart*m.cellSize),
		m.borderColor,
		false)
	vector.DrawFilledRect(
		screen,
		float32(m.gridOffsetX+float64(m.width*m.cellSize)-float64(borderThickness)),
		float32(m.uiHeight+(m.exitEnd+1)*m.cellSize),
		borderThickness,
		float32((m.height-m.exitEnd-1)*m.cellSize),
		m.borderColor,
		false)

	// Draw towers
	for _, tower := range m.towers {
		if tower != nil {
			// Pass true if this is the selected tower
			isSelected := false
			mouseX, mouseY := ebiten.CursorPosition()
			gridX, gridY := m.GetGridPosition(float64(mouseX), float64(mouseY))
			if tower.position.X == gridX && tower.position.Y == gridY {
				isSelected = true
			}
			tower.Draw(screen, isSelected)
		}
	}
}

// GetGridPosition converts screen coordinates to grid coordinates
func (m *GameMap) GetGridPosition(screenX, screenY float64) (int, int) {
	// Adjust for UI height and grid offset
	screenY -= float64(m.uiHeight)
	screenX -= m.gridOffsetX
	
	gridX := int(screenX) / m.cellSize
	gridY := int(screenY) / m.cellSize
	
	// Ensure coordinates are within bounds
	if gridX < 0 {
		gridX = 0
	}
	if gridX >= m.width {
		gridX = m.width - 1
	}
	if gridY < 0 {
		gridY = 0
	}
	if gridY >= m.height {
		gridY = m.height - 1
	}
	
	return gridX, gridY
}

// IsBlocked checks if a position is blocked by a tower or out of bounds
func (m *GameMap) IsBlocked(x, y int) bool {
	// Check bounds
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return true
	}
	
	// Always block cells with towers
	if m.terrain[y][x] == TowerPlacement {
		return true
	}
	
	return false
}

// CanPlaceTower checks if a position is suitable for tower placement
func (m *GameMap) CanPlaceTower(x, y int) bool {
	// Can't place in entrance or exit areas
	if x == 0 && (y >= m.entranceStart && y <= m.entranceEnd) {
		return false
	}
	if x == m.width-1 && (y >= m.exitStart && y <= m.exitEnd) {
		return false
	}
	
	// Check if position is within bounds and empty
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return false
	}
	return m.terrain[y][x] == Empty
}

// PlaceTower attempts to place a tower at the specified position
func (m *GameMap) PlaceTower(x, y int) bool {
	if !m.CanPlaceTower(x, y) {
		return false
	}
	
	// Check if placing the tower would block all paths
	if !m.checkPathExists(x, y) {
		return false
	}
	
	m.terrain[y][x] = TowerPlacement
	m.towers = append(m.towers, NewTower(currentGame.selectedTower, x, y))
	return true
}

// GetTowerAt returns the tower at the specified position or nil if there isn't one
func (m *GameMap) GetTowerAt(x, y int) *Tower {
	for _, tower := range m.towers {
		if tower.position.X == x && tower.position.Y == y {
			return tower
		}
	}
	return nil
}

// RemoveTower removes a tower from the specified position
func (m *GameMap) RemoveTower(x, y int) {
	if x >= 0 && x < m.width && y >= 0 && y < m.height {
		m.terrain[y][x] = Empty
		// Remove tower from towers slice
		for i, tower := range m.towers {
			pos := tower.GetPosition()
			if pos.X == x && pos.Y == y {
				m.towers = append(m.towers[:i], m.towers[i+1:]...)
				break
			}
		}
	}
}

// GetEntranceArea returns the entrance area coordinates
func (m *GameMap) GetEntranceArea() (int, int, int) {
	return m.entranceStart, m.entranceEnd, 0
}

// GetExitArea returns the exit area coordinates
func (m *GameMap) GetExitArea() (int, int, int) {
	return m.exitStart, m.exitEnd, m.width - 1
}

// GetTowersInRange returns all towers within range of a point
func (m *GameMap) GetTowersInRange(x, y, range_ float64) []*Tower {
	var nearbyTowers []*Tower
	for _, tower := range m.towers {
		if tower == nil {
			continue
		}
		
		towerX := float64(tower.position.X*m.cellSize + m.cellSize/2)
		towerY := float64(tower.position.Y*m.cellSize + m.cellSize/2 + m.uiHeight)
		
		dx := towerX - x
		dy := towerY - y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist <= range_ {
			nearbyTowers = append(nearbyTowers, tower)
		}
	}
	return nearbyTowers
}

// checkPathExists uses breadth-first search to verify if a path exists
func (m *GameMap) checkPathExists(testX, testY int) bool {
	// Temporarily place tower for testing
	originalTerrain := m.terrain[testY][testX]
	m.terrain[testY][testX] = TowerPlacement
	
	// Create visited array
	visited := make([][]bool, m.height)
	for i := range visited {
		visited[i] = make([]bool, m.width)
	}
	
	// Create queue for BFS
	type QueueItem struct {
		x, y int
	}
	queue := []QueueItem{}
	
	// Add all entrance points to queue
	for y := m.entranceStart; y <= m.entranceEnd; y++ {
		queue = append(queue, QueueItem{0, y})
		visited[y][0] = true
	}
	
	// BFS
	pathExists := false
	for len(queue) > 0 {
		// Pop front of queue
		current := queue[0]
		queue = queue[1:]
		
		// Check if we've reached an exit point
		if current.x == m.width-1 && current.y >= m.exitStart && current.y <= m.exitEnd {
			pathExists = true
			break
		}
		
		// Try all four directions
		directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for _, dir := range directions {
			newX := current.x + dir[0]
			newY := current.y + dir[1]
			
			// Check bounds and if not visited and not blocked
			if newX >= 0 && newX < m.width && newY >= 0 && newY < m.height && 
			   !visited[newY][newX] && m.terrain[newY][newX] != TowerPlacement {
				queue = append(queue, QueueItem{newX, newY})
				visited[newY][newX] = true
			}
		}
	}
	
	// Restore original terrain
	m.terrain[testY][testX] = originalTerrain
	
	return pathExists
}