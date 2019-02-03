package trace

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Volume is a volumetric surface.
// It can intersect with rays not only on its surface boundary, but also within its enclosing volume.
// Fog and smoke are volumetric surfaces.
type Volume struct {
	boundary Surface
	density  float64
	phase    Material
}

// NewVolume returns a new volumetric surface with a boundary, density, and material.
// Most materials are not designed for use with volumetric surfaces,
// so usually phase will be Isotropic.
func NewVolume(boundary Surface, density float64, phase Material) *Volume {
	return &Volume{
		boundary: boundary,
		density:  density,
		phase:    phase,
	}
}

// Hit returns details of the intersection between r and this volume.
// If r does not intersect, it returns nil.
// Intersections are not deterministic for volumes, and it will hit or not,
// and at various distances, based on the volume's density and a random factor.
func (v *Volume) Hit(r Ray, dMin, dMax float64, rnd *rand.Rand) *Hit {
	hit1 := v.boundary.Hit(r, -math.MaxFloat64, math.MaxFloat64, rnd)
	if hit1 == nil {
		return nil
	}
	hit2 := v.boundary.Hit(r, hit1.Dist+bias, math.MaxFloat64, rnd)
	if hit2 == nil {
		return nil
	}
	if hit1.Dist < dMin {
		hit1.Dist = dMin
	}
	if hit2.Dist > dMax {
		hit2.Dist = dMax
	}
	if hit1.Dist > hit2.Dist {
		return nil
	}
	dHit := -(1 / v.density) * math.Log(rnd.Float64())
	d := hit1.Dist + dHit
	if d >= hit2.Dist {
		return nil
	}
	return &Hit{
		Dist: d,
		Norm: geom.Unit{1, 0, 0},
		UV:   geom.Vec{0, 0, 0},
		Pt:   r.At(d),
		Mat:  v.phase,
	}
}

// Bounds returns an axis-aligned bounding box that encloses
// this volume from time t0 to t1.
func (v *Volume) Bounds(t0, t1 float64) *AABB {
	return v.boundary.Bounds(t0, t1)
}
