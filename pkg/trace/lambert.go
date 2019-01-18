package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Lambert describes a diffuse material.
type Lambert struct {
	Albedo Color
}

// NewLambert creates a new Lambert material with the given color.
func NewLambert(albedo Color) Lambert {
	return Lambert{Albedo: albedo}
}

// Scatter scatters incoming light rays in a hemisphere about the normal.
func (l Lambert) Scatter(in geom.Unit, n geom.Unit) (out geom.Unit, attenuation Color, ok bool) {
	out = n.Plus(geom.RandVecInSphere()).Unit()
	return out, l.Albedo, true
}
