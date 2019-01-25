package trace

// List holds a list of Surfaces
type List struct {
	HH []HitBounder
}

// NewList creates a new list of HitBounders
func NewList(h ...HitBounder) *List {
	return &List{HH: h}
}

// Hit finds the first intersection (if any) between Ray r and any of the HitBounders in the List.
// If no intersection is found, t = 0.
func (l *List) Hit(r Ray, dMin, dMax float64) (d float64, b Bouncer) {
	closest := dMax
	for _, h := range l.HH {
		if hd, hb := h.Hit(r, dMin, closest); hd > 0 {
			closest, d = hd, hd
			b = hb
		}
	}
	return
}

// Add adds new HitBounders to this List
func (l *List) Add(h ...HitBounder) int {
	l.HH = append(l.HH, h...)
	return len(l.HH)
}

// Bounds returns the Axis Aligned Bounding Box encompassing all listed HitBounders
// between times t0 and t1
func (l *List) Bounds(t0, t1 float64) (bounds *AABB) {
	for _, h := range l.HH {
		bounds = h.Bounds(t0, t1).Plus(bounds)
	}
	if bounds == nil {
		panic("No Bounds defined")
	}
	return
}
