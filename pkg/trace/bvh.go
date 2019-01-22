package trace

import (
	"math"
	"sort"
)

// BVH represents a bounding volume hierarchy
type BVH struct {
	left, right HitBoxer
	box         *AABB
}

// NewBVH creates a new BVH representing HitBoxers given as h between times time0 and time1
// When manually creating a new BVH, you should pass depth = 0.
func NewBVH(depth int, time0, time1 float64, h ...HitBoxer) *BVH {
	b := BVH{}
	n := len(h)
	if n == 1 {
		b.left, b.right = h[0], h[0]
		b.box = h[0].Box(time0, time1)
		return &b
	}
	if n == 2 {
		b.left, b.right = h[0], h[1]
		b.box = h[0].Box(time0, time1).Plus(h[1].Box(time0, time1))
		return &b
	}
	axis := depth % 3
	sort.Slice(h, func(i, j int) bool {
		h0 := h[i].Box(time0, time1)
		h1 := h[j].Box(time0, time1)
		return h0.Min.E[axis] < h1.Min.E[axis]
	})
	split := int(math.Floor(float64(n)/2 + 1))
	b.left = NewBVH(depth+1, time0, time1, h[0:split]...)
	b.right = NewBVH(depth+1, time0, time1, h[split:]...)
	b.box = b.left.Box(time0, time1).Plus(b.right.Box(time0, time1))
	return &b
}

// Hit returns the distance d at which Ray r hits this BVH.
// If the Ray does not intersect with this BVH, d = 0.
// The returned Bouncer describes the ray's bouncing at this hit point.
func (b *BVH) Hit(r Ray, dMin, dMax float64) (d float64, bo Bouncer) {
	if !b.box.Hit(r, dMin, dMax) {
		return 0, nil
	}
	lDist, lBounce := b.left.Hit(r, dMin, dMax)
	rDist, rBounce := b.right.Hit(r, dMin, dMax)
	if lDist > 0 && rDist > 0 {
		if lDist < rDist {
			return lDist, lBounce
		}
		return rDist, rBounce
	}
	if lDist > 0 {
		return lDist, lBounce
	}
	if rDist > 0 {
		return rDist, rBounce
	}
	return 0, nil
}

// Box returns a reference to an AABB encompassing the space of
// every HitBoxer in this BVH from time t0 to t1.
func (b *BVH) Box(t0, t1 float64) *AABB {
	return b.box
}
