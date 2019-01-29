package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

type Translate struct {
	Child  HitBounder
	Offset geom.Vec
}

func NewTranslate(child HitBounder, offset geom.Vec) *Translate {
	return &Translate{Child: child, Offset: offset}
}

func (t *Translate) Hit(r Ray, dMin, dMax float64) *Hit {
	r2 := NewRay(r.Or.Minus(t.Offset), r.Dir, r.t)
	hit := t.Child.Hit(r2, dMin, dMax)
	if hit != nil {
		hit.Pt = hit.Pt.Plus(t.Offset)
	}
	return hit
}

func (t *Translate) Bounds(t0, t1 float64) *AABB {
	b := t.Child.Bounds(t0, t1)
	return NewAABB(b.Min.Plus(t.Offset), b.Max.Plus(t.Offset))
}
