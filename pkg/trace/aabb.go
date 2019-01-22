package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// AABB Represents an axis-aligned bounding box
type AABB struct {
	Min geom.Vec
	Max geom.Vec
}

// NewAABB creates a new AABB with the given min and max positions
func NewAABB(min, max geom.Vec) *AABB {
	return &AABB{Min: min, Max: max}
}

// Hit returns whether or not a given ray hits the AABB between dMin and dMax distances
func (a *AABB) Hit(r Ray, dMin, dMax float64) bool {
	for i := 0; i < 3; i++ {
		invD := 1 / r.Dir.E[i]
		d0 := (a.Min.E[i] - r.Or.E[i]) * invD
		d1 := (a.Max.E[i] - r.Or.E[i]) * invD
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

// Plus returns a new AABB that encompasses the space of both a and b.
// If b is nil, the new AABB will be equivalent to a.
func (a *AABB) Plus(b *AABB) *AABB {
	if b == nil {
		return NewAABB(a.Min, a.Max)
	}
	return NewAABB(a.Min.Min(b.Min), a.Max.Max(b.Max))
}
