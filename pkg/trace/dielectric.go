package trace

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Dielectric describes a non-metallic material
type Dielectric struct {
	RefIdx float64
}

// NewDielectric creates a new material with the given index of refraction.
func NewDielectric(refractiveIndex float64) Dielectric {
	return Dielectric{RefIdx: refractiveIndex}
}

// Scatter reflects or refracts incoming light based on the ratio of indexes of refraction
func (d Dielectric) Scatter(in geom.Ray, p geom.Vec, n geom.Unit) (out geom.Ray, attenuation Color, ok bool) {
	attenuation = NewColor(1, 1, 1)

	outNormal := n
	ratio := 1 / d.RefIdx
	if in.Dir.Dot(n) > 0 {
		outNormal = n.Inv()
		ratio = d.RefIdx
	}

	r, refracted := refract(in.Dir, outNormal, ratio)
	if !refracted {
		r = reflect(in.Dir, n)
	}
	return geom.NewRay(p, r), attenuation, true
}

func refract(u, n geom.Unit, ratio float64) (u2 geom.Unit, ok bool) {
	dt := u.Dot(n)
	disc := 1 - ratio*ratio*(1-dt*dt)
	if disc <= 0 {
		return u2, false
	}
	// TODO: compose this so it's more readable
	u2 = (u.Minus(n.Scaled(dt)).Scaled(ratio)).Minus(n.Scaled(math.Sqrt(disc))).ToUnit()
	return u2, true
}
