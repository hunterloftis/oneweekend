package trace

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Sphere represents a spherical Surface
type Sphere struct {
	center0, center1 geom.Vec
	t0, t1           float64
	rad              float64
	mat              Material
}

// NewSphere creates a new Sphere with the given center and radius.
func NewSphere(center geom.Vec, radius float64, m Material) *Sphere {
	return NewMovingSphere(center, center, 0, 1, radius, m)
}

// NewMovingSphere creates a new Sphere with two centers separated by times t0 and t1
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

// Hit finds the distance to the first intersection (if any) between Ray r and the Sphere's surface.
// If no intersection is found, d = 0.
func (s *Sphere) Hit(r Ray, dMin, dMax float64, _ *rand.Rand) *Hit {
	oc := r.Or.Minus(s.Center(r.T))
	a := r.Dir.Dot(r.Dir)
	b := oc.Dot(r.Dir.Vec)
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

// Center returns the center of the sphere at a given time t
func (s *Sphere) Center(t float64) geom.Vec {
	p := (t - s.t0) / (s.t1 - s.t0)
	offset := s.center1.Minus(s.center0).Scaled(p)
	return s.center0.Plus(offset)
}

// Bounds returns the Axis-Aligned Bounding Bounds of the sphere encompassing all times between t0 and t1
func (s *Sphere) Bounds(t0, t1 float64) AABB {
	bounds0 := NewAABB(
		s.Center(t0).Minus(geom.NewVec(s.rad, s.rad, s.rad)),
		s.Center(t0).Plus(geom.NewVec(s.rad, s.rad, s.rad)),
	)
	bounds1 := NewAABB(
		s.Center(t1).Minus(geom.NewVec(s.rad, s.rad, s.rad)),
		s.Center(t1).Plus(geom.NewVec(s.rad, s.rad, s.rad)),
	)
	return bounds0.Plus(bounds1)
}

// UV returns the u, v spherical-mapped coordinates of this Sphere at point p, time t.
func (s *Sphere) UV(p geom.Vec, t float64) (uv geom.Vec) {
	p2 := p.Minus(s.Center(t)).Scaled(1 / s.rad)
	phi := math.Atan2(p2.Z(), p2.X())
	theta := math.Asin(p2.Y())
	u := 1 - (phi+math.Pi)/(2*math.Pi)
	v := (theta + math.Pi/2) / math.Pi
	return geom.NewVec(u, v, 0)
}

// Rect represents an axis-aligned rectangle
type Rect struct {
	min, max geom.Vec
	axis     int
	mat      Material
}

// NewRect creates a new Rect with the given min and max points and a Material.
// One axis of min and max should be equal - that's the axis of the plane of the rectangle.
// Rects default to XY rects (with an axis of Z).
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

// Hit finds the distnace to the first intersection (if any) between Ray in and the Rect.
// If no intersection is found, d = 0.
func (r *Rect) Hit(in Ray, dMin, dMax float64, _ *rand.Rand) *Hit {
	a0 := r.axis
	a1 := (a0 + 1) % 3
	a2 := (a0 + 2) % 3
	k := r.min.E[a0]
	d := (k - in.Or.E[a0]) / in.Dir.E[a0]
	if d < dMin || d > dMax {
		return nil
	}
	e1 := in.Or.E[a1] + d*in.Dir.E[a1]
	e2 := in.Or.E[a2] + d*in.Dir.E[a2]
	if e1 < r.min.E[a1] || e1 > r.max.E[a1] || e2 < r.min.E[a2] || e2 > r.max.E[a2] {
		return nil
	}
	u := (e1 - r.min.E[a1]) / (r.max.E[a1] - r.min.E[a1])
	v := (e2 - r.min.E[a2]) / (r.max.E[a2] - r.min.E[a2])
	norm := geom.NewUnit(0, 0, 0)
	norm.E[a0] = 1
	return &Hit{
		Dist: d,
		UV:   geom.NewVec(u, v, 0),
		Pt:   in.At(d),
		Mat:  r.mat,
		Norm: norm,
	}
}

// Bounds returns the axis Aligned Bounding Box encompassing the Rect.
func (r *Rect) Bounds(t0, t1 float64) AABB {
	bias := geom.NewVec(0, 0, 0)
	bias.E[r.axis] = 0.001
	return NewAABB(r.min.Minus(bias), r.max.Plus(bias))
}

type Box struct {
	List
}

func NewBox(min, max geom.Vec, m Material) *Box {
	return &Box{List: *NewList(
		NewRect(geom.NewVec(min.X(), min.Y(), max.Z()), geom.NewVec(max.X(), max.Y(), max.Z()), m),
		NewFlip(NewRect(geom.NewVec(min.X(), min.Y(), min.Z()), geom.NewVec(max.X(), max.Y(), min.Z()), m)),
		NewRect(geom.NewVec(min.X(), max.Y(), min.Z()), geom.NewVec(max.X(), max.Y(), max.Z()), m),
		NewFlip(NewRect(geom.NewVec(min.X(), min.Y(), min.Z()), geom.NewVec(max.X(), min.Y(), max.Z()), m)),
		NewRect(geom.NewVec(max.X(), min.Y(), min.Z()), geom.NewVec(max.X(), max.Y(), max.Z()), m),
		NewFlip(NewRect(geom.NewVec(min.X(), min.Y(), min.Z()), geom.NewVec(min.X(), max.Y(), max.Z()), m)),
	)}
}
