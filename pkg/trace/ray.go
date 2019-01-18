package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Ray represents a light ray that exists both in space and in time
type Ray struct {
	Or  geom.Vec
	Dir geom.Unit
	t   float64
}

// NewRay creates a new ray given an origin, direction, and point in time
func NewRay(origin geom.Vec, direction geom.Unit, time float64) Ray {
	return Ray{
		Or:  origin,
		Dir: direction,
		t:   time,
	}
}

// At returns the ray at point t
func (r Ray) At(t float64) geom.Vec {
	return r.Or.Plus(r.Dir.Scaled(t))
}
