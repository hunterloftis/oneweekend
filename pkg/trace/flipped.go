package trace

// Flipped overrides the Bounce method of a HitBounder to invert surface normals.
type Flipped struct {
	HitBounder
}

// NewFlipped creates a new Flipped instance.
func NewFlipped(child HitBounder) *Flipped {
	return &Flipped{HitBounder: child}
}

// Hit finds the distnace to the first intersection (if any) between Ray in and the child HitBounder.
// If no intersection is found, d = 0.
// Instead of returning the child, it returns the Flipped instance as the Bouncer target.
func (f *Flipped) Hit(in Ray, dMin, dMax float64) *Hit {
	hit := f.HitBounder.Hit(in, dMin, dMax)
	if hit != nil {
		hit.Norm = hit.Norm.Inv()
	}
	return hit
}
