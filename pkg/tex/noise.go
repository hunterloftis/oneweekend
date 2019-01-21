package tex

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Noise represents a Perlin noise texture
type Noise struct {
	Scale float64
}

// NewNoise returns a new noise texture with frequency scaled by scale
func NewNoise(scale float64) Noise {
	return Noise{Scale: scale}
}

// Map maps a u, v coordinate in 3d space p to a Color
func (n Noise) Map(u, v float64, p geom.Vec) Color {
	bright := 0.5 * (1 + math.Sin(n.Scale*p.Z()+10*turb(p, 7)))
	return NewColor(1, 1, 1).Scaled(bright)
}
