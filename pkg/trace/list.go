package trace

// List holds a list of Surfaces
type List struct {
	HH []HitBoxer
}

// NewList creates a new list of HitBoxers
func NewList(h ...HitBoxer) *List {
	return &List{HH: h}
}

// Hit finds the first intersection (if any) between Ray r and any of the HitBoxers in the List.
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

// Add adds new HitBoxers to this List
func (l *List) Add(h ...HitBoxer) int {
	l.HH = append(l.HH, h...)
	return len(l.HH)
}

// Box returns the Axis Aligned Bounding Box encompassing all listed HitBoxers
// between times t0 and t1
func (l *List) Box(t0, t1 float64) (box *AABB) {
	for _, h := range l.HH {
		box = h.Box(t0, t1).Plus(box)
	}
	if box == nil {
		panic("No Boxes defined")
	}
	return box
}
