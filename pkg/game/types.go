package game

// Point represents a position on the grid
type Point struct {
	X, Y int
}

// TerrainType represents different types of terrain
type TerrainType int

const (
	Empty TerrainType = iota
	TowerPlacement
)