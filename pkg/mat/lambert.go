package mat

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/tex"
)

// Lambert describes a diffuse material.
type Lambert struct {
	Albedo tex.Mapper
}

// NewLambert creates a new Lambert material with the given color.
func NewLambert(albedo tex.Mapper) Lambert {
	return Lambert{Albedo: albedo}
}

// Scatter scatters incoming light rays in a hemisphere about the normal.
func (l Lambert) Scatter(in, n geom.Unit, p geom.Vec) (out geom.Unit, attenuation tex.Color, ok bool) {
	out = n.Plus(geom.RandVecInSphere()).Unit()
	attenuation = l.Albedo.Map(0, 0, p)
	return out, attenuation, true
}