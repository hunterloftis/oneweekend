package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

type Solid struct {
	C Color
}

func NewSolid(c Color) Solid {
	return Solid{C: c}
}

func (s Solid) At(u, v float64, p geom.Vec) Color {
	return s.C
}
