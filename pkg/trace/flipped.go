package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// HitBoxBouncer implements HitBoxer and Bouncer.
type HitBoxBouncer interface {
	HitBoxer
	Bouncer
}

// Flipped overrides the Bounce method of a HitBoxBouncer to invert surface normals.
type Flipped struct {
	HitBoxBouncer
}

// NewFlipped creates a new Flipped instance.
func NewFlipped(child HitBoxBouncer) *Flipped {
	return &Flipped{HitBoxBouncer: child}
}

// Hit finds the distnace to the first intersection (if any) between Ray in and the child HitBoxBouncer.
// If no intersection is found, d = 0.
// Instead of returning the child, it returns the Flipped instance as the Bouncer target.
func (f *Flipped) Hit(in Ray, dMin, dMax float64) (d float64, b Bouncer) {
	d, _ = f.HitBoxBouncer.Hit(in, dMin, dMax)
	return d, f
}

// Bounce returns the normal, uv coords, point of impact, and material encountered by Ray r at distance d.
func (f *Flipped) Bounce(r Ray, d float64) (norm geom.Unit, uv, p geom.Vec, m Material) {
	norm2, uv, p, mat := f.HitBoxBouncer.Bounce(r, d)
	return norm2.Inv(), uv, p, mat
}
