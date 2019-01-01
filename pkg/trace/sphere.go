package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Sphere represents a spherical Surface
type Sphere struct {
}

// NewSphere creates a new Sphere with the given center and radius.
func NewSphere(center geom.Vec, radius float64) Sphere {
	return Sphere{}
}

// Hit finds the first intersection (if any) between Ray r and the Sphere's surface.
// If no intersection is found, t = 0.
func (s Sphere) Hit(r geom.Ray, tMin, tMax float64) (t, p geom.Vec, n geom.Unit) {
	return
}
