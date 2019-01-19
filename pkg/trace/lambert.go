package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Lambert describes a diffuse material.
type Lambert struct {
	Albedo Texture
}

// NewLambert creates a new Lambert material with the given color.
func NewLambert(albedo Texture) Lambert {
	return Lambert{Albedo: albedo}
}

// Scatter scatters incoming light rays in a hemisphere about the normal.
func (l Lambert) Scatter(in, n geom.Unit, p geom.Vec) (out geom.Unit, attenuation Color, ok bool) {
	out = n.Plus(geom.RandVecInSphere()).Unit()
	attenuation = l.Albedo.At(0, 0, p)
	return out, attenuation, true
}
