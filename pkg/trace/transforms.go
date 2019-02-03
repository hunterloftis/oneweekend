package trace

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Translate is a surface that translates a child surface.
type Translate struct {
	child  Surface
	offset geom.Vec
}

// NewTranslate returns a new surface that translates child by offset.
func NewTranslate(child Surface, offset geom.Vec) *Translate {
	return &Translate{child: child, offset: offset}
}

// Hit returns details of the intersection between r and this surface.
// If r does not intersect with this surface, it returns nil.
func (t *Translate) Hit(r Ray, dMin, dMax float64, rnd *rand.Rand) *Hit {
	r2 := NewRay(r.Or.Minus(t.offset), r.Dir, r.T)
	hit := t.child.Hit(r2, dMin, dMax, rnd)
	if hit != nil {
		hit.Pt = hit.Pt.Plus(t.offset)
	}
	return hit
}

// Bounds returns an axis-aligned bounding box that encloses
// this surface from time t0 to t1.
func (t *Translate) Bounds(t0, t1 float64) *AABB {
	b := t.child.Bounds(t0, t1)
	return NewAABB(b.Min().Plus(t.offset), b.Max().Plus(t.offset))
}

// RotateY is a surface that rotates a child surface on the Y axis.
type RotateY struct {
	child              Surface
	sinTheta, cosTheta float64
	bounds             *AABB
}

// NewRotateY returns a new surface that rotates child by angle on the Y axis.
func NewRotateY(child Surface, angle float64) *RotateY {
	rads := angle * math.Pi / 180
	r := RotateY{
		child:    child,
		sinTheta: math.Sin(rads),
		cosTheta: math.Cos(rads),
		bounds:   child.Bounds(0, 1),
	}
	for _, p := range r.bounds.Corners() {
		p2 := r.right(p)
		r.bounds = r.bounds.Extended(p2)
	}
	return &r
}

// Hit returns details of the intersection between r and this surface.
// If r does not intersect with this surface, it returns nil.
func (r *RotateY) Hit(in Ray, dMin, dMax float64, rnd *rand.Rand) *Hit {
	in2 := NewRay(r.left(in.Or), geom.Unit(r.left(geom.Vec(in.Dir))), in.T)
	hit := r.child.Hit(in2, dMin, dMax, rnd)
	if hit != nil {
		hit.Norm = geom.Unit(r.right(geom.Vec(hit.Norm)))
		hit.Pt = r.right(hit.Pt)
	}
	return hit
}

// Bounds returns an axis-aligned bounding box that encloses
// this surface from time t0 to t1.
func (r *RotateY) Bounds(t0, t1 float64) *AABB {
	return r.bounds
}

func (r *RotateY) right(dir geom.Vec) geom.Vec {
	x := r.cosTheta*dir.X() + r.sinTheta*dir.Z()
	z := -r.sinTheta*dir.X() + r.cosTheta*dir.Z()
	return geom.Vec{x, dir.Y(), z}
}

func (r *RotateY) left(dir geom.Vec) geom.Vec {
	x := r.cosTheta*dir.X() - r.sinTheta*dir.Z()
	z := r.sinTheta*dir.X() + r.cosTheta*dir.Z()
	return geom.Vec{x, dir.Y(), z}
}

// Flip is a surface that inverts the normals of a child surface.
type Flip struct {
	Surface
}

// NewFlip returns a new surface that inverts the normals of child.
func NewFlip(child Surface) *Flip {
	return &Flip{Surface: child}
}

// Hit returns details of the intersection between r and this surface.
// If r does not intersect with this surface, it returns nil.
// It modifies the Hit record by inverting the normals of the original intersection.
func (f *Flip) Hit(in Ray, dMin, dMax float64, rnd *rand.Rand) *Hit {
	hit := f.Surface.Hit(in, dMin, dMax, rnd)
	if hit != nil {
		hit.Norm = hit.Norm.Inv()
	}
	return hit
}
