package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Metal describes a reflective material
type Metal struct {
	Albedo Color
	Rough  float64
}

// NewMetal creates a new Metal material with a given color and roughness.
func NewMetal(albedo Color, roughness float64) Metal {
	return Metal{Albedo: albedo, Rough: roughness}
}

// Scatter reflects incoming light rays about the normal.
func (m Metal) Scatter(in geom.Unit, n geom.Unit) (out geom.Unit, attenuation Color, ok bool) {
	r := reflect(in, n)
	out = r.Plus(geom.RandVecInSphere().Scaled(m.Rough)).Unit()
	return out, m.Albedo, out.Dot(n) > 0
}

// Reflect reflects this unit vector about a normal vector n.
func reflect(u, n geom.Unit) geom.Unit {
	return geom.Unit{Vec: u.Minus(n.Scaled(2 * u.Dot(n)))}
}
