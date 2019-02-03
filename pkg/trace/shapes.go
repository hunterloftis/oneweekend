package trace

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Sphere is a spherical surface.
type Sphere struct {
	center0, center1 geom.Vec
	t0, t1           float64
	rad              float64
	mat              Material
}

// NewSphere creates a new sphere with the given center and radius.
func NewSphere(center geom.Vec, radius float64, m Material) *Sphere {
	return NewMovingSphere(center, center, 0, 1, radius, m)
}

// NewMovingSphere creates a new sphere that moves from center0 at time t0 to center1 at time t1.
func NewMovingSphere(center0, center1 geom.Vec, t0, t1, radius float64, m Material) *Sphere {
	return &Sphere{
		center0: center0,
		center1: center1,
		t0:      t0,
		t1:      t1,
		rad:     radius,
		mat:     m,
	}
}

// Hit returns details of the intersection between r and this surface.
// If r does not intersect with this BVH, it returns nil.
func (s *Sphere) Hit(r Ray, dMin, dMax float64, _ *rand.Rand) *Hit {
	oc := r.Or.Minus(s.Center(r.T))
	a := r.Dir.Dot(r.Dir)
	b := oc.Dot(geom.Vec(r.Dir))
	c := oc.Dot(oc) - s.rad*s.rad
	disc := b*b - a*c
	if disc <= 0 {
		return nil
	}
	sqrt := math.Sqrt(b*b - a*c)
	d := (-b - sqrt) / a
	if d <= dMin || d >= dMax {
		d = (-b + sqrt) / a
		if d <= dMin || d >= dMax {
			return nil
		}
	}
	p := r.At(d)
	return &Hit{
		Dist: d,
		Norm: p.Minus(s.Center(r.T)).Scaled(s.rad).Unit(),
		UV:   s.UV(p, r.T),
		Pt:   p,
		Mat:  s.mat,
	}
}

// Center returns the center of the sphere at time t.
func (s *Sphere) Center(t float64) geom.Vec {
	p := (t - s.t0) / (s.t1 - s.t0)
	offset := s.center1.Minus(s.center0).Scaled(p)
	return s.center0.Plus(offset)
}

// Bounds returns the axis-aligned bounding box enclosing this sphere
// between times t0 and t1.
func (s *Sphere) Bounds(t0, t1 float64) *AABB {
	bounds0 := NewAABB(
		s.Center(t0).Minus(geom.Vec{s.rad, s.rad, s.rad}),
		s.Center(t0).Plus(geom.Vec{s.rad, s.rad, s.rad}),
	)
	bounds1 := NewAABB(
		s.Center(t1).Minus(geom.Vec{s.rad, s.rad, s.rad}),
		s.Center(t1).Plus(geom.Vec{s.rad, s.rad, s.rad}),
	)
	return bounds0.Plus(bounds1)
}

// UV maps point p at time t to a uv coordinate.
// The uv coordinate is spherically mapped (lat/lon).
func (s *Sphere) UV(p geom.Vec, t float64) (uv geom.Vec) {
	p2 := p.Minus(s.Center(t)).Scaled(1 / s.rad)
	phi := math.Atan2(p2.Z(), p2.X())
	theta := math.Asin(p2.Y())
	u := 1 - (phi+math.Pi)/(2*math.Pi)
	v := (theta + math.Pi/2) / math.Pi
	return geom.Vec{u, v, 0}
}

// Rect is an axis-aligned rectangular surface.
type Rect struct {
	min, max geom.Vec
	axis     int
	mat      Material
}

// NewRect creates a new rectangular surface with material m from points min to max.
// One axis of min and max must be equal - that's the axis of the plane of the rectangle.
func NewRect(min, max geom.Vec, m Material) *Rect {
	r := Rect{
		min: min,
		max: max,
		mat: m,
	}
	r.axis = 2
	if min.X() == max.X() {
		r.axis = 0
	} else if min.Y() == max.Y() {
		r.axis = 1
	}
	return &r
}

// Hit returns details of the intersection between r and this surface.
// If r does not intersect with this BVH, it returns nil.
func (r *Rect) Hit(in Ray, dMin, dMax float64, _ *rand.Rand) *Hit {
	a0 := r.axis
	a1 := (a0 + 1) % 3
	a2 := (a0 + 2) % 3
	k := r.min[a0]
	d := (k - in.Or[a0]) / in.Dir[a0]
	if d < dMin || d > dMax {
		return nil
	}
	e1 := in.Or[a1] + d*in.Dir[a1]
	e2 := in.Or[a2] + d*in.Dir[a2]
	if e1 < r.min[a1] || e1 > r.max[a1] || e2 < r.min[a2] || e2 > r.max[a2] {
		return nil
	}
	u := (e1 - r.min[a1]) / (r.max[a1] - r.min[a1])
	v := (e2 - r.min[a2]) / (r.max[a2] - r.min[a2])
	norm := geom.Unit{0, 0, 0}
	norm[a0] = 1
	return &Hit{
		Dist: d,
		UV:   geom.Vec{u, v, 0},
		Pt:   in.At(d),
		Mat:  r.mat,
		Norm: norm,
	}
}

// Bounds returns an axis-aligned bounding box that encloses
// this rect from time t0 to t1.
func (r *Rect) Bounds(t0, t1 float64) *AABB {
	bias := geom.Vec{0, 0, 0}
	bias[r.axis] = 0.001
	return NewAABB(r.min.Minus(bias), r.max.Plus(bias))
}

// Box is an axis-aligned surface with six rectangular sides.
type Box struct {
	*List
}

// NewBox returns a new box that encloses the volume between points min and max.
func NewBox(min, max geom.Vec, m Material) *Box {
	return &Box{List: NewList(
		NewRect(geom.Vec{min.X(), min.Y(), max.Z()}, geom.Vec{max.X(), max.Y(), max.Z()}, m),
		NewFlip(NewRect(geom.Vec{min.X(), min.Y(), min.Z()}, geom.Vec{max.X(), max.Y(), min.Z()}, m)),
		NewRect(geom.Vec{min.X(), max.Y(), min.Z()}, geom.Vec{max.X(), max.Y(), max.Z()}, m),
		NewFlip(NewRect(geom.Vec{min.X(), min.Y(), min.Z()}, geom.Vec{max.X(), min.Y(), max.Z()}, m)),
		NewRect(geom.Vec{max.X(), min.Y(), min.Z()}, geom.Vec{max.X(), max.Y(), max.Z()}, m),
		NewFlip(NewRect(geom.Vec{min.X(), min.Y(), min.Z()}, geom.Vec{min.X(), max.Y(), max.Z()}, m)),
	)}
}
