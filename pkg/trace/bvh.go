package trace

import (
	"math"
	"sort"
)

type BVH struct {
	left, right HitBoxer
	box         *AABB
}

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

func (b *BVH) Box(t0, t1 float64) *AABB {
	return b.box
}
