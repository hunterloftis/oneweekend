package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Ray represents a light ray that exists both in space and in time
type Ray struct {
	Or  geom.Vec
	Dir geom.Unit
	T   float64
}

// NewRay creates a new ray given an origin, direction, and point in time
func NewRay(origin geom.Vec, direction geom.Unit, time float64) Ray {
	return Ray{
		Or:  origin,
		Dir: direction,
		T:   time,
	}
}

// At returns the position of the ray at distance d.
func (r Ray) At(d float64) geom.Vec {
	return r.Or.Plus(r.Dir.Scaled(d))
}
