package trace

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Sphere represents a spherical Surface
type Sphere struct {
	Center geom.Vec
	Rad    float64
	Mat    Material
}

// NewSphere creates a new Sphere with the given center and radius.
func NewSphere(center geom.Vec, radius float64, m Material) Sphere {
	return Sphere{Center: center, Rad: radius, Mat: m}
}

// Hit finds the first intersection (if any) between Ray r and the Sphere's surface.
// If no intersection is found, t = 0.
func (s Sphere) Hit(r geom.Ray, tMin, tMax float64) (t float64, surf Surfacer) {
	oc := r.Or.Minus(s.Center)
	a := r.Dir.Dot(r.Dir.Vec)
	b := oc.Dot(r.Dir.Vec)
	c := oc.Dot(oc) - s.Rad*s.Rad
	disc := b*b - a*c
	if disc <= 0 {
		return 0, s
	}
	sqrt := math.Sqrt(b*b - a*c)
	t = (-b - sqrt) / a
	if t > tMin && t < tMax {
		return t, s
	}
	t = (-b + sqrt) / a
	if t > tMin && t < tMax {
		return t, s
	}
	return 0, s
}

// Surface returns the normal and material at point p on the Sphere
func (s Sphere) Surface(p geom.Vec) (n geom.Unit, m Material) {
	return p.Minus(s.Center).ToUnit(), s.Mat
}
