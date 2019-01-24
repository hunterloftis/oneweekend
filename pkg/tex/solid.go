package tex

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Solid represents a single, solid colored texture
type Solid struct {
	C Color
}

// NewSolid returns a new solid texture
func NewSolid(c Color) Solid {
	return Solid{C: c}
}

// Map maps a u, v coordinate in 3d space p to a Color
func (s Solid) Map(uv, p geom.Vec) Color {
	return s.C
}
