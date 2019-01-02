package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Lambert describes a diffuse material.
type Lambert struct {
	Albedo Color
}

func NewLambert(albedo Color) Lambert {
	return Lambert{Albedo: albedo}
}

// Scatter scatters incoming light rays in a hemisphere about the normal.
func (l Lambert) Scatter(in geom.Ray, p geom.Vec, n geom.Unit) (out geom.Ray, attenuation Color, ok bool) {
	target := p.Plus(n.Vec).Plus(geom.RandVecInSphere())
	out = geom.NewRay(p, target.Minus(p).ToUnit())
	return out, l.Albedo, true
}
