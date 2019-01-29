package trace

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

type RotateY struct {
	Child              HitBounder
	sinTheta, cosTheta float64
	bounds             *AABB
}

func NewRotateY(child HitBounder, angle float64) *RotateY {
	rads := angle * math.Pi / 180
	r := RotateY{
		Child:    child,
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

func (r *RotateY) Hit(in Ray, dMin, dMax float64) *Hit {
	in2 := NewRay(r.left(in.Or), geom.Unit{Vec: r.left(in.Dir.Vec)}, in.t)
	hit := r.Child.Hit(in2, dMin, dMax)
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
