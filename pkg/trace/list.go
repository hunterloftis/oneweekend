package trace

// List holds a list of Surfaces
type List struct {
	HH []Hitter
}

// NewList creates a new list of Hitters
func NewList(h ...Hitter) *List {
	return &List{HH: h}
}

// Hit finds the first intersection (if any) between Ray r and any of the Hitters in the List.
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

// Add adds new Hitters to this List
func (l *List) Add(h ...Hitter) int {
	l.HH = append(l.HH, h...)
	return len(l.HH)
}
