package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

// Lambert describes a lambertian material attenuated by an albedo
type Lambert struct {
	Albedo Color
}

// Scatter scatters incoming light rays in a lambertian pattern
func (l Lambert) Scatter(in geom.Ray, p geom.Vec, n geom.Unit) (out geom.Ray, attenuation Color) {
	target := p.Plus(n.Vec).Plus(geom.RandVecInSphere())
	out = geom.NewRay(p, target.Minus(p).ToUnit())
	return out, l.Albedo
}
