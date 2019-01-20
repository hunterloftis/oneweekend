package tex

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

type Checker struct {
	Size      float64
	Odd, Even Mapper
}

func NewChecker(size float64, t0, t1 Mapper) Checker {
	return Checker{
		Size: size,
		Odd:  t0,
		Even: t1,
	}
}

func (c Checker) Map(u, v float64, p geom.Vec) Color {
	sines := math.Sin(c.Size*p.X()) * math.Sin(c.Size*p.Y()) * math.Sin(c.Size*p.Z())
	if sines < 0 {
		return c.Odd.Map(u, v, p)
	}
	return c.Even.Map(u, v, p)
}
