package trace

import (
	"math"
	"sort"
)

// HitBounder represents a Hitter that also provides bounding box information.
type HitBounder interface {
	Hitter
	Bounds(t0, t1 float64) (bounds *AABB)
}

// BVH represents a bounding volume hierarchy
type BVH struct {
	left, right HitBounder
	bounds      *AABB
}

// NewBVH creates a new BVH representing HitBounders given as h between times time0 and time1
// When manually creating a new BVH, you should pass depth = 0.
func NewBVH(depth int, time0, time1 float64, h ...HitBounder) *BVH {
	b := BVH{}
	n := len(h)
	if n == 1 {
		b.left, b.right = h[0], h[0]
		b.bounds = h[0].Bounds(time0, time1)
		return &b
	}
	if n == 2 {
		b.left, b.right = h[0], h[1]
		b.bounds = h[0].Bounds(time0, time1).Plus(h[1].Bounds(time0, time1))
		return &b
	}
	axis := depth % 3
	sort.Slice(h, func(i, j int) bool {
		h0 := h[i].Bounds(time0, time1)
		h1 := h[j].Bounds(time0, time1)
		return h0.Min.E[axis] < h1.Min.E[axis]
	})
	split := int(math.Floor(float64(n)/2 + 1))
	b.left = NewBVH(depth+1, time0, time1, h[0:split]...)
	b.right = NewBVH(depth+1, time0, time1, h[split:]...)
	b.bounds = b.left.Bounds(time0, time1).Plus(b.right.Bounds(time0, time1))
	return &b
}

// Hit returns the distance d at which Ray r hits this BVH.
// If the Ray does not intersect with this BVH, d = 0.
// The returned Bouncer describes the ray's bouncing at this hit point.
func (b *BVH) Hit(r Ray, dMin, dMax float64) *Hit {
	if !b.bounds.Hit(r, dMin, dMax) {
		return nil
	}
	lHit := b.left.Hit(r, dMin, dMax)
	rHit := b.right.Hit(r, dMin, dMax)
	if lHit != nil && rHit != nil {
		if lHit.Dist < rHit.Dist {
			return lHit
		}
		return rHit
	}
	if lHit != nil {
		return lHit
	}
	if rHit != nil {
		return rHit
	}
	return nil
}

// Bounds returns a reference to an AABB encompassing the space of
// every HitBounder in this BVH from time t0 to t1.
func (b *BVH) Bounds(t0, t1 float64) *AABB {
	return b.bounds
}
