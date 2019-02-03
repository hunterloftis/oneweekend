package trace

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

type Translate struct {
	child  Surface
	offset geom.Vec
}

func NewTranslate(child Surface, offset geom.Vec) *Translate {
	return &Translate{child: child, offset: offset}
}

func (t *Translate) Hit(r Ray, dMin, dMax float64, rnd *rand.Rand) *Hit {
	r2 := NewRay(r.Or.Minus(t.offset), r.Dir, r.T)
	hit := t.child.Hit(r2, dMin, dMax, rnd)
	if hit != nil {
		hit.Pt = hit.Pt.Plus(t.offset)
	}
	return hit
}

func (t *Translate) Bounds(t0, t1 float64) *AABB {
	b := t.child.Bounds(t0, t1)
	return NewAABB(b.Min().Plus(t.offset), b.Max().Plus(t.offset))
}

type RotateY struct {
	child              Surface
	sinTheta, cosTheta float64
	bounds             *AABB
}

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

func (r *RotateY) Hit(in Ray, dMin, dMax float64, rnd *rand.Rand) *Hit {
	in2 := NewRay(r.left(in.Or), geom.Unit{Vec: r.left(in.Dir.Vec)}, in.T)
	hit := r.child.Hit(in2, dMin, dMax, rnd)
	if hit != nil {
		hit.Norm = geom.Unit{Vec: r.right(hit.Norm.Vec)}
		hit.Pt = r.right(hit.Pt)
	}
	return hit
}

func (r *RotateY) Bounds(t0, t1 float64) *AABB {
	return r.bounds
}

func (r *RotateY) right(dir geom.Vec) geom.Vec {
	x := r.cosTheta*dir.X() + r.sinTheta*dir.Z()
	z := -r.sinTheta*dir.X() + r.cosTheta*dir.Z()
	return geom.NewVec(x, dir.Y(), z)
}

func (r *RotateY) left(dir geom.Vec) geom.Vec {
	x := r.cosTheta*dir.X() - r.sinTheta*dir.Z()
	z := r.sinTheta*dir.X() + r.cosTheta*dir.Z()
	return geom.NewVec(x, dir.Y(), z)
}

// Flip overrides the Bounce method of a Surface to invert surface normals.
type Flip struct {
	Surface
}

// NewFlip creates a new Flip instance.
func NewFlip(child Surface) *Flip {
	return &Flip{Surface: child}
}

// Hit finds the distnace to the first intersection (if any) between Ray in and the child Surface.
// If no intersection is found, d = 0.
// Instead of returning the child, it returns the Flip instance as the Bouncer target.
func (f *Flip) Hit(in Ray, dMin, dMax float64, rnd *rand.Rand) *Hit {
	hit := f.Surface.Hit(in, dMin, dMax, rnd)
	if hit != nil {
		hit.Norm = hit.Norm.Inv()
	}
	return hit
}
