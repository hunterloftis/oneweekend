package mat

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/tex"
)

// Metal describes a reflective material
type Metal struct {
	Albedo tex.Color
	Rough  float64
	nonEmitter
}

// NewMetal creates a new Metal material with a given color and roughness.
func NewMetal(albedo tex.Color, roughness float64) Metal {
	return Metal{Albedo: albedo, Rough: roughness}
}

// Scatter reflects incoming light rays about the normal.
func (m Metal) Scatter(in, norm geom.Unit, _, _ geom.Vec) (out geom.Unit, attenuate tex.Color, ok bool) {
	r := reflect(in, norm)
	out = r.Plus(geom.RandVecInSphere().Scaled(m.Rough)).Unit()
	return out, m.Albedo, out.Dot(norm) > 0
}

// Reflect reflects this unit vector about a normal vector n.
func reflect(u, n geom.Unit) geom.Unit {
	return geom.Unit{Vec: u.Minus(n.Scaled(2 * u.Dot(n)))}
}
