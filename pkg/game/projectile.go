package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
)

// ProjectileType represents different types of projectiles
type ProjectileType int

const (
	DartProjectile ProjectileType = iota
	BulletProjectile
	LightningProjectile
	FlameProjectile
	FreezeProjectile
)

// Projectile represents a projectile shot from a tower
type Projectile struct {
	x, y        float64   // Current position
	targetX     float64   // Target X position
	targetY     float64   // Target Y position
	speed       float64   // Movement speed
	damage      float64   // Damage amount
	projType    ProjectileType
	size        float64   // Size for drawing
	color       color.Color
}

// NewProjectile creates a new projectile
func NewProjectile(startX, startY, targetX, targetY float64, projType ProjectileType, damage float64) *Projectile {
	// Base projectile setup
	proj := &Projectile{
		x:        startX,
		y:        startY,
		targetX:  targetX,
		targetY:  targetY,
		speed:    5.0,    // Base speed for all projectiles
		damage:   damage,
		projType: projType,
		size:     8.0,    // Base size for projectiles
	}

	// Set specific colors based on tower type
	switch projType {
	case DartProjectile:
		// Brighter version of tower bronze
		proj.color = color.RGBA{255, 210, 120, 255}  // Brighter bronze
		proj.size = 12.0  // Slightly larger for visibility
		proj.speed = 6.0  // Faster than base
	case BulletProjectile:
		// Brighter version of tower gold
		proj.color = color.RGBA{255, 235, 120, 255}  // Bright metallic gold
		proj.size = 6.0   // Small but fast
		proj.speed = 8.0  // Fastest projectile
	case LightningProjectile:
		// Brighter version of tower blue
		proj.color = color.RGBA{120, 240, 255, 255}  // Intense bright electric blue
		proj.size = 14.0  // Larger for lightning effect
		proj.speed = 7.0  // Fast
	case FlameProjectile:
		// Brighter version of tower orange-red
		proj.color = color.RGBA{255, 140, 60, 255}  // Vivid orange-red
		proj.size = 7.0   // Smaller base size for flame effect
		proj.speed = 4.0  // Slower but area effect
	case FreezeProjectile:
		// Brighter version of tower ice blue
		proj.color = color.RGBA{200, 250, 255, 255}  // Brilliant ice blue
		proj.size = 11.0   // Larger for snowflake
		proj.speed = 5.0  // Medium speed
	}

	return proj
}

// Update moves the projectile and returns true if it reached its target
func (p *Projectile) Update() bool {
	// Calculate direction to target
	dx := p.targetX - p.x
	dy := p.targetY - p.y
	dist := math.Sqrt(dx*dx + dy*dy)
	
	// If we're very close to target, we've hit
	if dist < p.speed {
		return true
	}
	
	// Move towards target
	p.x += (dx / dist) * p.speed
	p.y += (dy / dist) * p.speed
	
	return false
}

// Draw draws the projectile
func (p *Projectile) Draw(screen *ebiten.Image) {
	// Calculate angle to target for rotation
	dx := p.targetX - p.x
	dy := p.targetY - p.y
	angle := math.Atan2(dy, dx)

	switch p.projType {
	case DartProjectile:
		// Draw dart as an arrow with tail
		length := p.size * 1.5
		headSize := p.size * 0.5
		
		// Main body of dart
		vector.StrokeLine(screen, 
			float32(p.x), float32(p.y),
			float32(p.x + math.Cos(angle)*length),
			float32(p.y + math.Sin(angle)*length),
			2, // Line width
			p.color,
			true)
		
		// Arrow head
		headAngle1 := angle + math.Pi*0.8 // 144 degrees
		headAngle2 := angle - math.Pi*0.8
		vector.StrokeLine(screen,
			float32(p.x + math.Cos(angle)*length),
			float32(p.y + math.Sin(angle)*length),
			float32(p.x + math.Cos(angle)*length + math.Cos(headAngle1)*headSize),
			float32(p.y + math.Sin(angle)*length + math.Sin(headAngle1)*headSize),
			2,
			p.color,
			true)
		vector.StrokeLine(screen,
			float32(p.x + math.Cos(angle)*length),
			float32(p.y + math.Sin(angle)*length),
			float32(p.x + math.Cos(angle)*length + math.Cos(headAngle2)*headSize),
			float32(p.y + math.Sin(angle)*length + math.Sin(headAngle2)*headSize),
			2,
			p.color,
			true)

	case BulletProjectile:
		// Draw bullet as a fast-moving metallic projectile with trail
		// Main bullet
		vector.DrawFilledCircle(screen,
			float32(p.x), float32(p.y),
			float32(p.size/2),
			p.color,
			true)
		// Trail effect
		trailLength := p.size * 2
		trailAngle := angle + math.Pi // Opposite direction
		vector.StrokeLine(screen,
			float32(p.x), float32(p.y),
			float32(p.x + math.Cos(trailAngle)*trailLength),
			float32(p.y + math.Sin(trailAngle)*trailLength),
			2,
			color.RGBA{255, 220, 120, 128}, // Transparent trail
			true)

	case LightningProjectile:
		// Draw lightning as multiple connected zigzag segments with more dramatic effect
		segmentLength := p.size * 0.8
		numSegments := 5  // More segments
		lastX, lastY := p.x, p.y
		
		// Draw outer glow first
		for i := 0; i < numSegments; i++ {
			offset := p.size * 0.4 * float64(1 - 2*(i%2)) // Larger zigzag
			nextX := p.x + math.Cos(angle)*segmentLength*float64(i+1)
			nextY := p.y + math.Sin(angle)*segmentLength*float64(i+1)
			
			// Add randomized offset
			nextX += math.Cos(angle+math.Pi/2) * offset * (0.8 + rand.Float64()*0.4)
			nextY += math.Sin(angle+math.Pi/2) * offset * (0.8 + rand.Float64()*0.4)
			
			// Draw wider outer glow
			vector.StrokeLine(screen,
				float32(lastX), float32(lastY),
				float32(nextX), float32(nextY),
				6, // Much wider line for glow
				color.RGBA{180, 240, 255, 64}, // Very transparent outer glow
				true)
			
			lastX, lastY = nextX, nextY
		}
		
		// Reset for core lightning
		lastX, lastY = p.x, p.y
		
		// Draw main lightning effect
		for i := 0; i < numSegments; i++ {
			offset := p.size * 0.4 * float64(1 - 2*(i%2))
			nextX := p.x + math.Cos(angle)*segmentLength*float64(i+1)
			nextY := p.y + math.Sin(angle)*segmentLength*float64(i+1)
			
			nextX += math.Cos(angle+math.Pi/2) * offset * (0.8 + rand.Float64()*0.4)
			nextY += math.Sin(angle+math.Pi/2) * offset * (0.8 + rand.Float64()*0.4)
			
			// Draw middle layer
			vector.StrokeLine(screen,
				float32(lastX), float32(lastY),
				float32(nextX), float32(nextY),
				3, // Medium width
				color.RGBA{180, 240, 255, 180}, // Semi-transparent middle layer
				true)
			
			// Draw core
			vector.StrokeLine(screen,
				float32(lastX), float32(lastY),
				float32(nextX), float32(nextY),
				1.5, // Thin core
				p.color, // Bright core
				true)
			
			lastX, lastY = nextX, nextY
			
			// Add occasional branch lightning
			if rand.Float64() < 0.3 { // 30% chance per segment
				branchAngle := angle + (rand.Float64()-0.5)*1.0 // Random branch direction
				branchLength := segmentLength * 0.5 // Half length of main bolt
				
				// Draw branch with same layered effect
				vector.StrokeLine(screen,
					float32(nextX), float32(nextY),
					float32(nextX + math.Cos(branchAngle)*branchLength),
					float32(nextY + math.Sin(branchAngle)*branchLength),
					4, // Glow
					color.RGBA{180, 240, 255, 64},
					true)
				
				vector.StrokeLine(screen,
					float32(nextX), float32(nextY),
					float32(nextX + math.Cos(branchAngle)*branchLength),
					float32(nextY + math.Sin(branchAngle)*branchLength),
					2, // Middle
					color.RGBA{180, 240, 255, 180},
					true)
				
				vector.StrokeLine(screen,
					float32(nextX), float32(nextY),
					float32(nextX + math.Cos(branchAngle)*branchLength),
					float32(nextY + math.Sin(branchAngle)*branchLength),
					1, // Core
					p.color,
					true)
			}
		}

	case FlameProjectile:
		// Draw flame as multiple particles in a flame shape with more dramatic effect
		numParticles := 5  // Fewer particles for more compact flame
		baseSize := p.size * 0.5  // Smaller base particle size
		
		for i := 0; i < numParticles; i++ {
			// Calculate particle position in flame pattern
			spread := 0.3 // Narrower flame spread
			particleDistance := float64(i) * baseSize * 0.5  // More compact spacing
			particleAngle := angle + (rand.Float64()-0.5)*spread
			
			// Vary particle size with more dramatic falloff
			particleSize := baseSize * (1.0 - float64(i)/float64(numParticles))
			
			// Create more dramatic gradient effect
			particleColor := color.RGBA{
				255,                                                    // Red always max
				uint8(140 - float64(i)/float64(numParticles)*120),    // More dramatic green falloff
				uint8(60 - float64(i)/float64(numParticles)*60),     // More dramatic blue falloff
				255,
			}
			
			// Draw glow effect first
			vector.DrawFilledCircle(screen,
				float32(p.x + math.Cos(particleAngle)*particleDistance),
				float32(p.y + math.Sin(particleAngle)*particleDistance),
				float32(particleSize*1.2),  // Smaller glow
				color.RGBA{255, 100, 30, 64}, // Outer glow
				true)
			
			// Draw main particle
			vector.DrawFilledCircle(screen,
				float32(p.x + math.Cos(particleAngle)*particleDistance),
				float32(p.y + math.Sin(particleAngle)*particleDistance),
				float32(particleSize),
				particleColor,
				true)
		}

	case FreezeProjectile:
		// Draw freeze as a snowflake pattern
		// Draw main crystal shape
		for i := 0; i < 6; i++ {
			rayAngle := angle + float64(i) * math.Pi / 3
			// Main ray
			vector.StrokeLine(screen,
				float32(p.x), float32(p.y),
				float32(p.x + math.Cos(rayAngle)*p.size),
				float32(p.y + math.Sin(rayAngle)*p.size),
				2,
				p.color,
				true)
			
			// Small branches on each ray
			branchSize := p.size * 0.4
			branchAngle1 := rayAngle + math.Pi/4
			branchAngle2 := rayAngle - math.Pi/4
			branchX := p.x + math.Cos(rayAngle)*p.size*0.6
			branchY := p.y + math.Sin(rayAngle)*p.size*0.6
			
			vector.StrokeLine(screen,
				float32(branchX), float32(branchY),
				float32(branchX + math.Cos(branchAngle1)*branchSize),
				float32(branchY + math.Sin(branchAngle1)*branchSize),
				1,
				p.color,
				true)
			
			vector.StrokeLine(screen,
				float32(branchX), float32(branchY),
				float32(branchX + math.Cos(branchAngle2)*branchSize),
				float32(branchY + math.Sin(branchAngle2)*branchSize),
				1,
				p.color,
				true)
		}
		
		// Draw center crystal
		vector.DrawFilledCircle(screen,
			float32(p.x), float32(p.y),
			float32(p.size/4),
			color.RGBA{80, 160, 255, 255}, // Deeper blue for center
			true)
	}
}

// GetPosition returns the current position of the projectile
func (p *Projectile) GetPosition() (float64, float64) {
	return p.x, p.y
}

// GetDamage returns the damage amount of the projectile
func (p *Projectile) GetDamage() float64 {
	return p.damage
}

// GetProjectileType returns the type of the projectile
func (p *Projectile) GetProjectileType() ProjectileType {
	return p.projType
}