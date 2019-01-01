package trace

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Sphere represents a spherical Surface
type Sphere struct {
	Center geom.Vec
	Rad    float64
}

// NewSphere creates a new Sphere with the given center and radius.
func NewSphere(center geom.Vec, radius float64) Sphere {
	return Sphere{Center: center, Rad: radius}
}

// Hit finds the first intersection (if any) between Ray r and the Sphere's surface.
// If no intersection is found, t = 0.
func (s Sphere) Hit(r geom.Ray, tMin, tMax float64) (t float64, p geom.Vec, n geom.Unit) {
	oc := r.Or.Minus(s.Center)
	a := r.Dir.Dot(r.Dir.Vec)
	b := oc.Dot(r.Dir.Vec)
	c := oc.Dot(oc) - s.Rad*s.Rad
	disc := b*b - a*c
	if disc <= 0 {
		return 0, p, n
	}
	sqrt := math.Sqrt(b*b - a*c)
	t = (-b - sqrt) / a
	if t > tMin && t < tMax {
		p = r.At(t)
		return t, p, p.Minus(s.Center).ToUnit()
	}
	t = (-b + sqrt) / a
	if t > tMin && t < tMax {
		p = r.At(t)
		return t, p, p.Minus(s.Center).ToUnit()
	}
	return 0, p, n
}
