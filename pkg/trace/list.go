package trace

import "math/rand"

// List contains a list of surfaces.
type List struct {
	ss []Surface
}

// NewList creates a new list containing ss.
func NewList(ss ...Surface) *List {
	return &List{ss: ss}
}

// Hit returns details of the intersection between r and this surface.
// If r does not intersect with this BVH, it returns nil.
func (l *List) Hit(r Ray, dMin, dMax float64, rnd *rand.Rand) (nearest *Hit) {
	for _, h := range l.ss {
		if hit := h.Hit(r, dMin, dMax, rnd); hit != nil {
			dMax = hit.Dist
			nearest = hit
		}
	}
	return
}

// Add adds new surfaces to the list.
func (l *List) Add(ss ...Surface) int {
	l.ss = append(l.ss, ss...)
	return len(l.ss)
}

// Bounds returns an axis-aligned bounding box that encloses
// all of the contained surfaces from time t0 to t1.
func (l *List) Bounds(t0, t1 float64) (bounds *AABB) {
	for _, h := range l.ss {
		bounds = h.Bounds(t0, t1).Plus(bounds)
	}
	return
}

// Surfaces returns all the surfaces the list contains.
func (l *List) Surfaces() []Surface {
	return l.ss
}
