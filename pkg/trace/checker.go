package trace

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

type Checker struct {
	Size      float64
	Odd, Even Texture
}

func NewChecker(size float64, t0, t1 Texture) Checker {
	return Checker{
		Size: size,
		Odd:  t0,
		Even: t1,
	}
}

func (c Checker) At(u, v float64, p geom.Vec) Color {
	sines := math.Sin(c.Size*p.X()) * math.Sin(c.Size*p.Y()) * math.Sin(c.Size*p.Z())
	if sines < 0 {
		return c.Odd.At(u, v, p)
	}
	return c.Even.At(u, v, p)
}
