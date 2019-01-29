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

// Hit records the details of a Ray->Surface intersection.
type Hit struct {
	Dist float64
	Norm geom.Unit
	UV   geom.Vec
	Pt   geom.Vec
	Mat  Material
}

// Hitter represents something that can be Hit by a Ray.
type Hitter interface {
	Hit(r Ray, dMin, dMax float64) *Hit
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
			ir := int(math.Min(255, 255.99*c.R()))
			ig := int(math.Min(255, 255.99*c.G()))
			ib := int(math.Min(255, 255.99*c.B()))
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
	if hit := h.Hit(r, bias, math.MaxFloat64); hit != nil {
		emit := hit.Mat.Emit(hit.UV, hit.Pt)
		out, attenuate, ok := hit.Mat.Scatter(r.Dir, hit.Norm, hit.UV, hit.Pt)
		if !ok {
			return emit
		}
		indirect := color(NewRay(hit.Pt, out, r.t), h, depth+1).Times(attenuate)
		return emit.Plus(indirect)
	}
	return tex.NewColor(0, 0, 0)
}
