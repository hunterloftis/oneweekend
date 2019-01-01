package trace

import (
	"fmt"
	"io"
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

const bias = 0.001

// Hitter represents something that can be Hit by a Ray.
type Hitter interface {
	Hit(r geom.Ray, tMin, tMax float64) (t float64, p geom.Vec, n geom.Unit)
}

// Frame gathers the results of ray traces on a W x H grid.
type Frame struct {
	W, H int
}

// NewFrame creates a new frame with specific dimensions
func NewFrame(width, height int) Frame {
	return Frame{W: width, H: height}
}

// WritePPM traces each pixel in the frame and writes the results to w in PPM format
func (f Frame) WritePPM(w io.Writer, h Hitter, samples int) error {
	if _, err := fmt.Fprintln(w, "P3"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, f.W, f.H); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "255"); err != nil {
		return err
	}

	cam := Camera{}

	for y := f.H - 1; y >= 0; y-- {
		for x := 0; x < f.W; x++ {
			c := NewColor(0, 0, 0)
			for s := 0; s < samples; s++ {
				u := (float64(x) + rand.Float64()) / float64(f.W)
				v := (float64(y) + rand.Float64()) / float64(f.H)
				r := cam.Ray(u, v)
				c = c.Plus(color(r, h))
			}
			c = c.Scaled(1 / float64(samples)).Gamma(2)
			ir := int(255.99 * c.R())
			ig := int(255.99 * c.G())
			ib := int(255.99 * c.B())
			if _, err := fmt.Fprintln(w, ir, ig, ib); err != nil {
				return err
			}
		}
	}
	return nil
}

func color(r geom.Ray, h Hitter) Color {
	if t, p, n := h.Hit(r, bias, math.MaxFloat64); t > 0 {
		target := p.Plus(n.Vec).Plus(geom.RandVecInSphere())
		r2 := geom.NewRay(p, target.Minus(p).ToUnit())
		return color(r2, h).Scaled(0.5)
	}
	t := 0.5 * (r.Dir.Y() + 1.0)
	white := NewColor(1, 1, 1).Scaled(1 - t)
	blue := NewColor(0.5, 0.7, 1).Scaled(t)
	return white.Plus(blue)
}
