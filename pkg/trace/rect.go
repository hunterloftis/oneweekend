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
func (r *Rect) Hit(in Ray, dMin, dMax float64) *Hit {
	a0 := r.Axis
	a1 := (a0 + 1) % 3
	a2 := (a0 + 2) % 3
	k := r.Min.E[a0]
	d := (k - in.Or.E[a0]) / in.Dir.E[a0]
	if d < dMin || d > dMax {
		return nil
	}
	e1 := in.Or.E[a1] + d*in.Dir.E[a1]
	e2 := in.Or.E[a2] + d*in.Dir.E[a2]
	if e1 < r.Min.E[a1] || e1 > r.Max.E[a1] || e2 < r.Min.E[a2] || e2 > r.Max.E[a2] {
		return nil
	}
	u := (e1 - r.Min.E[a1]) / (r.Max.E[a1] - r.Min.E[a1])
	v := (e2 - r.Min.E[a2]) / (r.Max.E[a2] - r.Min.E[a2])
	norm := geom.NewUnit(0, 0, 0)
	norm.E[a0] = 1
	return &Hit{
		Dist: d,
		UV:   geom.NewVec(u, v, 0),
		Pt:   in.At(d),
		Mat:  r.Mat,
		Norm: norm,
	}
}

// Bounds returns the Axis Aligned Bounding Box encompassing the Rect.
func (r *Rect) Bounds(t0, t1 float64) *AABB {
	bias := geom.NewVec(0, 0, 0)
	bias.E[r.Axis] = 0.001
	return NewAABB(r.Min.Minus(bias), r.Max.Plus(bias))
}
