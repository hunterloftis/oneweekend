package tex

import "github.com/hunterloftis/oneweekend/pkg/geom"

type Solid struct {
	C Color
}

func NewSolid(c Color) Solid {
	return Solid{C: c}
}

func (s Solid) Map(u, v float64, p geom.Vec) Color {
	return s.C
}
