package trace

import (
	"math"
	"math/rand"
	"sort"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

const costTraverse = 1
const costIntersect = 2

// BVH is a surface that contains other surfaces
// and organizes them into a bounding volume hierarchy.
// This improves performance of the Hit() function over simple lists for most sets of surfaces.
type BVH struct {
	left, right *BVH
	leaf        Surface
	bounds      *AABB
}

// NewBVH builds a new BVH containing surfaces in ss between times time0 and time1.
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

// Hit returns details of the intersection between r and this surface.
// If r does not intersect with this BVH, it returns nil.
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

// Bounds returns an axis-aligned bounding box that encloses
// all of the contained surfaces from time t0 to t1.
func (b *BVH) Bounds(t0, t1 float64) *AABB {
	return b.bounds
}

func split(t0, t1 float64, axis int, fraction float64, ss []Surface) (ll, rr []Surface) {
	ss2 := make([]Surface, len(ss))
	copy(ss2, ss)
	sort.Slice(ss2, func(i, j int) bool {
		m0 := ss2[i].Bounds(t0, t1).Mid()
		m1 := ss2[j].Bounds(t0, t1).Mid()
		return m0[axis] < m1[axis]
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

// AABB is an axis-aligned bounding box.
type AABB struct {
	min geom.Vec
	max geom.Vec
}

// NewAABB creates a new axis-aligned bounding box with the given min and max points.
func NewAABB(min, max geom.Vec) *AABB {
	return &AABB{min: min.Min(max), max: max.Max(min)}
}

// Hit returns whether or not r hits the box between distances dMin and dMax.
func (a *AABB) Hit(r Ray, dMin, dMax float64) bool {
	for i := 0; i < 3; i++ {
		invD := 1 / r.Dir[i]
		d0 := (a.min[i] - r.Or[i]) * invD
		d1 := (a.max[i] - r.Or[i]) * invD
		if invD < 0 {
			d0, d1 = d1, d0
		}
		if d0 > dMin {
			dMin = d0
		}
		if d1 < dMax {
			dMax = d1
		}
		if dMax <= dMin {
			return false
		}
	}
	return true
}

// Plus returns a new bounding box that encloses both this box and b.
// If b is nil, the new box will be equivalent to this box.
func (a *AABB) Plus(b *AABB) *AABB {
	if b == nil {
		return NewAABB(a.min, a.max)
	}
	return NewAABB(a.min.Min(b.min), a.max.Max(b.max))
}

// Corners returns the eight corners of this bounding box.
func (a *AABB) Corners() []geom.Vec {
	c := make([]geom.Vec, 0, 8)
	for i := 0.0; i < 2; i++ {
		for j := 0.0; j < 2; j++ {
			for k := 0.0; k < 2; k++ {
				x := i*a.min.X() + (1-i)*a.max.X()
				y := j*a.min.Y() + (1-j)*a.max.Y()
				z := k*a.min.Z() + (1-k)*a.max.Z()
				c = append(c, geom.Vec{x, y, z})
			}
		}
	}
	return c
}

// Extended returns an extended bounding box that also encloses v.
func (a *AABB) Extended(v geom.Vec) *AABB {
	return NewAABB(a.min.Min(v), a.max.Max(v))
}

// Min returns the minimum corner of this bounding box.
func (a *AABB) Min() geom.Vec {
	return a.min
}

// Max returns the maximum corner of this bounding box.
func (a *AABB) Max() geom.Vec {
	return a.max
}

// SurfaceArea returns the total surface area of this bounding box.
func (a *AABB) SurfaceArea() float64 {
	dims := a.max.Minus(a.min)
	front := dims.X() * dims.Y()
	side := dims.Z() * dims.Y()
	top := dims.X() * dims.Z()
	return (front + side + top) * 2
}

// Mid returns the mid-point of this bounding box.
func (a *AABB) Mid() geom.Vec {
	return a.min.Plus(a.max).Scaled(0.5)
}
