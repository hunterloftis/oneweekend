package trace

import "math/rand"

// List holds a list of Surfaces
type List struct {
	ss []Surface
}

// NewList creates a new list of HitBounders
func NewList(h ...Surface) *List {
	return &List{ss: h}
}

// Hit finds the first intersection (if any) between Ray r and any of the HitBounders in the List.
// If no intersection is found, t = 0.
func (l *List) Hit(r Ray, dMin, dMax float64, rnd *rand.Rand) (nearest *Hit) {
	for _, h := range l.ss {
		if hit := h.Hit(r, dMin, dMax, rnd); hit != nil {
			dMax = hit.Dist
			nearest = hit
		}
	}
	return
}

// Add adds new HitBounders to this List
func (l *List) Add(ss ...Surface) int {
	l.ss = append(l.ss, ss...)
	return len(l.ss)
}

// Bounds returns the Axis Aligned Bounding Box encompassing all listed HitBounders
// between times t0 and t1
func (l *List) Bounds(t0, t1 float64) (bounds AABB) {
	for _, h := range l.ss {
		bounds = h.Bounds(t0, t1).Plus(bounds)
	}
	return
}

// Surfaces returns all the Surfaces contained by this list.
func (l *List) Surfaces() []Surface {
	return l.ss
}
