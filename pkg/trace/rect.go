package trace

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Rect represents an axis-aligned rectangle
type Rect struct {
	Min, Max geom.Vec
	Axis     int
	Mat      Material
}

// NewRect creates a new Rect with the given min and max points and a Material.
// One axis of min and max should be equal - that's the axis of the plane of the rectangle.
// Rects default to XY rects (with an axis of Z).
func NewRect(min, max geom.Vec, m Material) *Rect {
	r := Rect{
		Min: min,
		Max: max,
		Mat: m,
	}
	r.Axis = 2
	if min.X() == max.X() {
		r.Axis = 0
	} else if min.Y() == max.Y() {
		r.Axis = 1
	}
	return &r
}

// Hit finds the distnace to the first intersection (if any) between Ray in and the Rect.
// If no intersection is found, d = 0.
func (r *Rect) Hit(in Ray, dMin, dMax float64) (d float64, b Bouncer) {
	a0 := r.Axis
	a1 := (a0 + 1) % 3
	a2 := (a0 + 2) % 3
	k := r.Min.E[a0]
	d = (k - in.Or.E[a0]) / in.Dir.E[a0]
	if d < dMin || d > dMax {
		return 0, r
	}
	e1 := in.Or.E[a1] + d*in.Dir.E[a1]
	e2 := in.Or.E[a2] + d*in.Dir.E[a2]
	if e1 < r.Min.E[a1] || e1 > r.Max.E[a1] || e2 < r.Min.E[a2] || e2 > r.Max.E[a2] {
		return 0, r
	}
	return d, r
}

// Box returns the Axis Aligned Bounding Box encompassing the Rect.
func (r *Rect) Box(t0, t1 float64) (box *AABB) {
	bias := geom.NewVec(0, 0, 0)
	bias.E[r.Axis] = 0.001
	return NewAABB(r.Min.Minus(bias), r.Max.Plus(bias))
}

// Bounce return the normal, light attenuation, and emittance color of the Rect at a point d distance along Ray in.
func (r *Rect) Bounce(in Ray, d float64) (norm geom.Unit, uv, p geom.Vec, m Material) {
	a0 := r.Axis
	a1 := (a0 + 1) % 3
	a2 := (a0 + 2) % 3
	e1 := in.Or.E[a1] + d*in.Dir.E[a1]
	e2 := in.Or.E[a2] + d*in.Dir.E[a2]
	u := (e1 - r.Min.E[a1]) / (r.Max.E[a1] - r.Min.E[a1])
	v := (e2 - r.Min.E[a2]) / (r.Max.E[a2] - r.Min.E[a2])
	p = in.At(d)
	norm = geom.NewUnit(0, 0, 0)
	norm.E[a0] = 1
	uv = geom.NewVec(u, v, 0)
	m = r.Mat
	return
}
