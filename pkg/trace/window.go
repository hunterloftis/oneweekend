package trace

import (
	"fmt"
	"io"
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/tex"
)

const bias = 0.001

// HitBoxer represents something that can be Hit by a Ray.
type HitBoxer interface {
	Hit(r Ray, tMin, tMax float64) (t float64, s Bouncer)
	Box(t0, t1 float64) (box *AABB)
}

// Bouncer represents something that can return Bounce rays out in a direction with attenuation
type Bouncer interface {
	Bounce(in Ray, dist float64) (out Ray, attenuation tex.Color, ok bool)
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
func (wi Window) WritePPM(w io.Writer, h HitBoxer, samples int) error {
	if _, err := fmt.Fprintln(w, "P3"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, wi.W, wi.H); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "255"); err != nil {
		return err
	}

	from := geom.NewVec(0, 0, 6)
	at := geom.NewVec(0, 0, 0)
	focus := 10.0
	cam := NewCamera(from, at, geom.NewUnit(0, 1, 0), 20, float64(wi.W)/float64(wi.H), 0, focus, 0, 1)

	for y := wi.H - 1; y >= 0; y-- {
		for x := 0; x < wi.W; x++ {
			c := tex.NewColor(0, 0, 0)
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

func color(r Ray, h HitBoxer, depth int) tex.Color {
	if depth > 9 {
		return tex.NewColor(0, 0, 0)
	}
	if d, b := h.Hit(r, bias, math.MaxFloat64); d > 0 {
		r2, attenuation, ok := b.Bounce(r, d)
		return attenuation
		if !ok {
			return tex.NewColor(0, 0, 0)
		}
		return color(r2, h, depth+1).Times(attenuation)
	}
	return tex.NewColor(0, 0, 0)
	t := 0.5 * (r.Dir.Y() + 1.0)
	white := tex.NewColor(1, 1, 1).Scaled(1 - t)
	blue := tex.NewColor(0.5, 0.7, 1).Scaled(t)
	return white.Plus(blue)
}
