package trace

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

var (
	black = Color{0, 0, 0}
	white = Color{1, 1, 1}
)

// Color is an RGB color value.
type Color geom.Vec

// R returns the first element (Red).
func (c Color) R() float64 {
	return c[0]
}

// G returns the second element (Green).
func (c Color) G() float64 {
	return c[1]
}

// B returns the third element (Blue).
func (c Color) B() float64 {
	return c[2]
}

// Plus returns the sum of the color and c2.
func (c Color) Plus(c2 Color) Color {
	return Color(geom.Vec(c).Plus(geom.Vec(c2)))
}

// Times returns the product of the color and c2.
func (c Color) Times(c2 Color) Color {
	return Color(geom.Vec(c).Times(geom.Vec(c2)))
}

// Scaled returns the color scaled by n.
func (c Color) Scaled(n float64) Color {
	return Color(geom.Vec(c).Scaled(n))
}

// Gamma raises each of R, G, and B to 1/n.
func (c Color) Gamma(n float64) Color {
	ni := 1 / n
	return Color{
		math.Pow(c.R(), ni),
		math.Pow(c.G(), ni),
		math.Pow(c.B(), ni),
	}
}

// RGBInt returns the red, green, and blue components as integers in a 0-255 range.
func (c Color) RGBInt() (r, g, b int) {
	return int(math.Min(255, 255*c[0])), int(math.Min(255, 255*c[1])), int(math.Min(255, 255*c[2]))
}
