package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Ray is a light ray with an origin, direction, and time.
type Ray struct {
	Or  geom.Vec
	Dir geom.Unit
	T   float64
}

// NewRay returns a new ray with an origin, direction, and point in time.
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
