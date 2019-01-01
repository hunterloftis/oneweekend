package trace

import (
	"fmt"
	"io"
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Surface represents something that can be Hit by a Ray.
type Surface interface {
	Hit(r geom.Ray, tMin, tMax float64) (t, p geom.Vec, n geom.Unit)
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
func (f Frame) WritePPM(w io.Writer, s Surface) error {
	if _, err := fmt.Fprintln(w, "P3"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, f.W, f.H); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "255"); err != nil {
		return err
	}

	lowerLeft := geom.NewVec(-2, -1, -1)
	horizontal := geom.NewVec(4, 0, 0)
	vertical := geom.NewVec(0, 2, 0)
	origin := geom.NewVec(0, 0, 0)

	for j := f.H - 1; j >= 0; j-- {
		for i := 0; i < f.W; i++ {
			u := float64(i) / float64(f.W)
			v := float64(j) / float64(f.H)
			r := geom.NewRay(
				origin,
				lowerLeft.Plus((horizontal.Scaled(u)).Plus(vertical.Scaled(v))).ToUnit(),
			)
			col := color(r)
			ir := int(255.99 * col.R())
			ig := int(255.99 * col.G())
			ib := int(255.99 * col.B())
			if _, err := fmt.Fprintln(w, ir, ig, ib); err != nil {
				return err
			}
		}
	}
	return nil
}

func color(r geom.Ray) Color {
	if t, ok := hitSphere(geom.NewVec(0, 0, -1), 0.5, r); ok {
		n := r.At(t).Minus(geom.NewVec(0, 0, -1)).ToUnit()
		return NewColor(n.X()+1, n.Y()+1, n.Z()+1).Scaled(0.5)
	}
	t := 0.5 * (r.Dir.Y() + 1.0)
	white := NewColor(1, 1, 1).Scaled(1 - t)
	blue := NewColor(0.5, 0.7, 1).Scaled(t)
	return white.Plus(blue)
}

func hitSphere(center geom.Vec, radius float64, r geom.Ray) (t float64, ok bool) {
	oc := r.Or.Minus(center)
	a := r.Dir.Dot(r.Dir.Vec)
	b := 2 * oc.Dot(r.Dir.Vec)
	c := oc.Dot(oc) - radius*radius
	disc := b*b - 4*a*c
	if disc < 0 {
		return 0, false
	}
	return (-b - math.Sqrt(disc)) / (2 * a), true
}
