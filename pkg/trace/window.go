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

// Hitter represents something that can be Hit by a Ray.
type Hitter interface {
	Hit(r Ray, tMin, tMax float64) (t float64, s Bouncer)
}

// Bouncer represents something that can bounce attenuated light rays.
// Bounces are not necessarily successful.
type Bouncer interface {
	Bounce(in Ray, dist float64) (norm geom.Unit, uv, p geom.Vec, m Material)
}

// Material represents a material that scatters light.
type Material interface {
	Scatter(in, norm geom.Unit, uv, p geom.Vec) (out geom.Unit, attenuate tex.Color, ok bool)
	Emit(uv, p geom.Vec) tex.Color
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

	from := geom.NewVec(278, 278, -800)
	at := geom.NewVec(278, 278, 0)
	focus := 10.0
	cam := NewCamera(from, at, geom.NewUnit(0, 1, 0), 40, float64(wi.W)/float64(wi.H), 0, focus, 0, 1)

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

func color(r Ray, h Hitter, depth int) tex.Color {
	if depth > 9 {
		return tex.NewColor(0, 0, 0)
	}
	if d, b := h.Hit(r, bias, math.MaxFloat64); d > 0 {
		norm, uv, p, mat := b.Bounce(r, d)
		emit := mat.Emit(uv, p)
		out, attenuate, ok := mat.Scatter(r.Dir, norm, uv, p)
		if !ok {
			return emit
		}
		indirect := color(NewRay(p, out, r.t), h, depth+1).Times(attenuate)
		return emit.Plus(indirect)
	}
	return tex.NewColor(0, 0, 0)
}
