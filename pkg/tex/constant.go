package tex

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Constant represents a single, Constant colored texture
type Constant struct {
	C Color
}

// NewConstant returns a new Constant texture
func NewConstant(c Color) Constant {
	return Constant{C: c}
}

// Map maps a u, v coordinate in 3d space p to a Color
func (c Constant) Map(uv, p geom.Vec) Color {
	return c.C
}
