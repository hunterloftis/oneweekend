package trace

import (
	"math"
	"math/rand"
	"sort"
)

const costTraverse = 1
const costIntersect = 2

// BVH represents a bounding volume hierarchy
type BVH struct {
	left, right *BVH
	leaf        Surface
	bounds      *AABB
}

// NewBVH creates a new BVH representing Surfaces given as h between times time0 and time1.
func NewBVH(time0, time1 float64, ss ...Surface) *BVH {
	return newBVHNode(0, time0, time1, ss...)
}

func newBVHNode(depth int, time0, time1 float64, ss ...Surface) *BVH {
	b := BVH{}
	if len(ss) < 4 {
		b.leaf = NewList(ss...)
		b.bounds = b.leaf.Bounds(time0, time1)
		return &b
	}

	cheapest := float64(len(ss)) * costIntersect
	var left, right []Surface
	for axis := 0; axis < 3; axis++ {
		for fraction := 0.33; fraction < 0.7; fraction += 0.17 {
			ll, rr := split(time0, time1, axis, fraction, ss)
			if c := cost(time0, time1, ll, rr); c < cheapest {
				cheapest = c
				left, right = ll, rr
			}
		}
	}
	if left == nil {
		b.leaf = NewList(ss...)
		b.bounds = b.leaf.Bounds(time0, time1)
		return &b
	}

	b.left = newBVHNode(depth+1, time0, time1, left...)
	b.right = newBVHNode(depth+1, time0, time1, right...)
	b.bounds = b.left.Bounds(time0, time1).Plus(b.right.Bounds(time0, time1))
	return &b
}

// Hit returns the distance d at which Ray r hits this BVH.
// If the Ray does not intersect with this BVH, d = 0.
// The returned Bouncer describes the ray's bouncing at this hit point.
func (b *BVH) Hit(r Ray, dMin, dMax float64, rnd *rand.Rand) *Hit {
	if !b.bounds.Hit(r, dMin, dMax) {
		return nil
	}
	if b.leaf != nil {
		return b.leaf.Hit(r, dMin, dMax, rnd)
	}
	lHit := b.left.Hit(r, dMin, dMax, rnd)
	rHit := b.right.Hit(r, dMin, dMax, rnd)
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

// Bounds returns a reference to an *AABB encompassing the space of
// every Surface in this BVH from time t0 to t1.
func (b *BVH) Bounds(t0, t1 float64) *AABB {
	return b.bounds
}

func split(t0, t1 float64, axis int, fraction float64, ss []Surface) (ll, rr []Surface) {
	ss2 := make([]Surface, len(ss))
	copy(ss2, ss)
	sort.Slice(ss2, func(i, j int) bool {
		m0 := ss2[i].Bounds(t0, t1).Mid()
		m1 := ss2[j].Bounds(t0, t1).Mid()
		return m0.E[axis] < m1.E[axis]
	})
	split := int(math.Floor(float64(len(ss2))*fraction + 1))
	return ss2[0:split], ss2[split:]
}

func cost(t0, t1 float64, ll, rr []Surface) float64 {
	sa := boundsAround(t0, t1, append(ll, rr...)).SurfaceArea()
	saL := boundsAround(t0, t1, ll).SurfaceArea()
	saR := boundsAround(t0, t1, rr).SurfaceArea()
	nL := float64(len(ll))
	nR := float64(len(rr))
	costLeft := (saL / sa) * nL
	costRight := (saR / sa) * nR
	return costTraverse + costIntersect*(costLeft+costRight)
}

func boundsAround(t0, t1 float64, ss []Surface) (b *AABB) {
	if len(ss) == 0 {
		panic("boundsAround called with zero Surfaces")
	}
	for _, s := range ss {
		b = s.Bounds(t0, t1).Plus(b)
	}
	return
}
