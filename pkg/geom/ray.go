package geom

// Ray defines a ray with an origin and a direction
type Ray struct {
	Or  Vec
	Dir Unit
}

// NewRay creates a new ray given an origin and a direction
func NewRay(origin Vec, direction Unit) Ray {
	return Ray{Or: origin, Dir: direction}
}

// At returns the ray at point t
func (r Ray) At(t float64) Vec {
	return r.Or.Plus(r.Dir.Scaled(t))
}
