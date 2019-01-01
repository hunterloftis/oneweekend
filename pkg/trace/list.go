package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// List holds a list of Surfaces
type List struct {
}

// NewList creates a new list of Surfaces
func NewList(s ...Surface) List {
	return List{}
}

// Hit finds the first intersection (if any) between Ray r and any of the Surfaces in the List.
// If no intersection is found, t = 0.
func (l List) Hit(r geom.Ray, tMin, tMax float64) (t, p geom.Vec, n geom.Unit) {
	return
}
