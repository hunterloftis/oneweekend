package tex

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Color represents an RGB color value
type Color struct {
	geom.Vec
}

// NewColor creates a Color from 3 float values
func NewColor(e0, e1, e2 float64) (c Color) {
	c.E[0] = e0
	c.E[1] = e1
	c.E[2] = e2
	return
}

// R returns the first element (Red)
func (c Color) R() float64 {
	return c.E[0]
}

// G returns the second element (Green)
func (c Color) G() float64 {
	return c.E[1]
}

// B returns the third element (Blue)
func (c Color) B() float64 {
	return c.E[2]
}

// Plus returns the sum of two colors
func (c Color) Plus(c2 Color) Color {
	return Color{Vec: c.Vec.Plus(c2.Vec)}
}

// Times returns the product of two colors
func (c Color) Times(c2 Color) Color {
	return Color{Vec: c.Vec.Times(c2.Vec)}
}

// Scaled returns the color scaled
func (c Color) Scaled(n float64) Color {
	return Color{Vec: c.Vec.Scaled(n)}
}

// Gamma raises each of R, G, and B to 1/n
func (c Color) Gamma(n float64) Color {
	ni := 1 / n
	return NewColor(
		math.Pow(c.R(), ni),
		math.Pow(c.G(), ni),
		math.Pow(c.B(), ni),
	)
}
