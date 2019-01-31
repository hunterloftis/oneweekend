package mat

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/tex"
)

// Dielectric describes a non-metallic material
type Dielectric struct {
	RefIdx float64
	nonEmitter
}

// NewDielectric creates a new material with the given index of refraction.
func NewDielectric(refractiveIndex float64) Dielectric {
	return Dielectric{RefIdx: refractiveIndex}
}

// Scatter reflects or refracts incoming light based on the ratio of indexes of refraction
func (d Dielectric) Scatter(in, n geom.Unit, _, _ geom.Vec) (out geom.Unit, attenuate tex.Color, ok bool) {
	var outNormal geom.Unit
	var ratio float64
	var cos float64

	if in.Dot(n) > 0 {
		outNormal = n.Inv()
		ratio = d.RefIdx
		cos = d.RefIdx * in.Dot(n) / in.Len()
	} else {
		outNormal = n
		ratio = 1 / d.RefIdx
		cos = -in.Dot(n) / in.Len()
	}

	out, refracted := refract(in, outNormal, ratio)
	if !refracted || schlick(cos, d.RefIdx) > rand.Float64() {
		out = reflect(in, n)
	}
	return out, tex.NewColor(1, 1, 1), true
}

func refract(u, n geom.Unit, ratio float64) (u2 geom.Unit, ok bool) {
	dt := u.Dot(n)
	disc := 1 - ratio*ratio*(1-dt*dt)
	if disc <= 0 {
		return u2, false
	}
	// TODO: compose this so it's more readable
	u2 = (u.Minus(n.Scaled(dt)).Scaled(ratio)).Minus(n.Scaled(math.Sqrt(disc))).Unit()
	return u2, true
}

func schlick(cos, refIdx float64) float64 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cos), 5)
}
