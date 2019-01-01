package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// List holds a list of Surfaces
type List struct {
	SS []Surface
}

// NewList creates a new list of Surfaces
func NewList(s ...Surface) List {
	return List{SS: s}
}

// Hit finds the first intersection (if any) between Ray r and any of the Surfaces in the List.
// If no intersection is found, t = 0.
func (l List) Hit(r geom.Ray, tMin, tMax float64) (t float64, p geom.Vec, n geom.Unit) {
	closest := tMax
	for _, s := range l.SS {
		if st, sp, sn := s.Hit(r, tMin, closest); st > 0 {
			closest, t = st, st
			p = sp
			n = sn
		}
	}
	return
}
