package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// List holds a list of Surfaces
type List struct {
	HH []Hitter
}

// NewList creates a new list of Hitters
func NewList(s ...Hitter) List {
	return List{HH: s}
}

// Hit finds the first intersection (if any) between Ray r and any of the Hitters in the List.
// If no intersection is found, t = 0.
func (l List) Hit(r geom.Ray, tMin, tMax float64) (t float64, s Surfacer) {
	closest := tMax
	for _, h := range l.HH {
		if ht, hs := h.Hit(r, tMin, closest); ht > 0 {
			closest, t = ht, ht
			s = hs
		}
	}
	return
}
