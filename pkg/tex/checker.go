package tex

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Checker represents an alternating checkered pattern of two sub-textures
type Checker struct {
	Size      float64
	Odd, Even Mapper
}

// NewChecker returns a new checkered texture rendering sub-textures t0 and t1 in size squares.
func NewChecker(size float64, t0, t1 Mapper) Checker {
	return Checker{
		Size: size,
		Odd:  t0,
		Even: t1,
	}
}

// Map maps a u, v coordinate in 3d space p to a Color
func (c Checker) Map(u, v float64, p geom.Vec) Color {
	sines := math.Sin(c.Size*p.X()) * math.Sin(c.Size*p.Y()) * math.Sin(c.Size*p.Z())
	if sines < 0 {
		return c.Odd.Map(u, v, p)
	}
	return c.Even.Map(u, v, p)
}
