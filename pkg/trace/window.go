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
	Hit(r geom.Ray, tMin, tMax float64) (t float64, s Bouncer)
}

// Bouncer represents something that can return Bounce normals and materials
type Bouncer interface {
	Bounce(p geom.Vec) (n geom.Unit, m Material)
}

// Material represents a material that scatters light.
type Material interface {
	Scatter(in geom.Unit, n geom.Unit) (out geom.Unit, attenuation Color, ok bool)
}

// Window gathers the results of ray traces on a W x H grid.
type Window struct {
	W, H int
}

// NewWindow creates a new Window with specific dimensions
func NewWindow(width, height int) Window {
	return Window{W: width, H: height}
}

// WritePPM traces each pixel in the Window and writes the results to w in PPM format
func (wi Window) WritePPM(w io.Writer, h Hitter, samples int) error {
	if _, err := fmt.Fprintln(w, "P3"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, wi.W, wi.H); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "255"); err != nil {
		return err
	}

	from := geom.NewVec(13, 2, 3)
	at := geom.NewVec(0, 0, 0)
	focus := 10.0
	cam := NewCamera(from, at, geom.NewUnit(0, 1, 0), 20, float64(wi.W)/float64(wi.H), 0.1, focus)

	for y := wi.H - 1; y >= 0; y-- {
		for x := 0; x < wi.W; x++ {
			c := NewColor(0, 0, 0)
			for s := 0; s < samples; s++ {
				u := (float64(x) + rand.Float64()) / float64(wi.W)
				v := (float64(y) + rand.Float64()) / float64(wi.H)
				r := cam.Ray(u, v)
				c = c.Plus(color(r, h, 0))
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

func color(r geom.Ray, h Hitter, depth int) Color {
	if depth > 9 {
		return NewColor(0, 0, 0)
	}
	if t, b := h.Hit(r, bias, math.MaxFloat64); t > 0 {
		p := r.At(t)
		n, m := b.Bounce(p)
		scattered, attenuation, ok := m.Scatter(r.Dir, n)
		if !ok {
			return NewColor(0, 0, 0)
		}
		r2 := geom.NewRay(p, scattered)
		return color(r2, h, depth+1).Times(attenuation)
	}
	t := 0.5 * (r.Dir.Y() + 1.0)
	white := NewColor(1, 1, 1).Scaled(1 - t)
	blue := NewColor(0.5, 0.7, 1).Scaled(t)
	return white.Plus(blue)
}
