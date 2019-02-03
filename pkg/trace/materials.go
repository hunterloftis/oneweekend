package trace

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Dielectric is a transparent, non-metallic material.
// Glass, diamond, and water are all dielectrics.
type Dielectric struct {
	iRefract float64
	nonEmitter
}

// NewDielectric creates a new material with the given index of refraction.
func NewDielectric(iRefract float64) *Dielectric {
	return &Dielectric{iRefract: iRefract}
}

// Scatter reflects, refracts, and attenuates incoming light.
func (d *Dielectric) Scatter(in, n geom.Unit, _, _ geom.Vec, rnd *rand.Rand) (out geom.Unit, attenuate Color, ok bool) {
	var outNormal geom.Unit
	var ratio float64
	var cos float64

	if in.Dot(n) > 0 {
		outNormal = n.Inv()
		ratio = d.iRefract
		cos = d.iRefract * in.Dot(n)
	} else {
		outNormal = n
		ratio = 1 / d.iRefract
		cos = -in.Dot(n)
	}

	out, refracted := refract(in, outNormal, ratio)
	if !refracted || schlick(cos, d.iRefract) > rnd.Float64() {
		out = reflect(in, n)
	}
	return out, white, true
}

func refract(u, n geom.Unit, ratio float64) (u2 geom.Unit, ok bool) {
	dt := u.Dot(n)
	disc := 1 - ratio*ratio*(1-dt*dt)
	if disc <= 0 {
		return u2, false
	}
	// TODO: compose this so it's more readable
	u2 = (geom.Vec(u).Minus(n.Scaled(dt)).Scaled(ratio)).Minus(n.Scaled(math.Sqrt(disc))).Unit()
	return u2, true
}

func schlick(cos, iRefract float64) float64 {
	r0 := (1 - iRefract) / (1 + iRefract)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cos), 5)
}

// Isotropic is an isotropic volumetric material.
// It models materials with subsurface scattering like simple smoke and fog.
type Isotropic struct {
	nonEmitter
	texture Mapper
}

// NewIsotropic creates a new isotropic material that will scatter the given texture throughout a volume.
func NewIsotropic(texture Mapper) *Isotropic {
	return &Isotropic{texture: texture}
}

// Scatter scatters incoming light in random directions and attenuates it based on this material's texture.
func (i *Isotropic) Scatter(in, norm geom.Unit, uv, p geom.Vec, rnd *rand.Rand) (out geom.Unit, attenuate Color, ok bool) {
	return geom.RandUnit(rnd), i.texture.Map(uv, p), true
}

// Lambert describes a flat, diffuse material.
// Rubber and chalk are simple lambertian materials.
type Lambert struct {
	texture Mapper
	nonEmitter
}

// NewLambert creates a new Lambert material with the given texture.
func NewLambert(texture Mapper) *Lambert {
	return &Lambert{texture: texture}
}

// Scatter scatters incoming light rays in a hemisphere about the normal,
// attenuating them by the material's texture.
func (l *Lambert) Scatter(in, n geom.Unit, uv, p geom.Vec, rnd *rand.Rand) (out geom.Unit, attenuate Color, ok bool) {
	out = geom.Vec(n).Plus(geom.RandVecInSphere(rnd)).Unit()
	attenuate = l.texture.Map(uv, p)
	return out, attenuate, true
}

// Light is a material that emits light.
type Light struct {
	texture Mapper
}

// NewLight creates a new material that emits light based on the given texture.
func NewLight(texture Mapper) *Light {
	return &Light{texture: texture}
}

// Scatter returns false as this material emits light but doesn't scatter it.
func (l *Light) Scatter(in, norm geom.Unit, uv, p geom.Vec, _ *rand.Rand) (out geom.Unit, attenuate Color, ok bool) {
	return out, attenuate, false
}

// Emit returns the color emitted at coordinate uv and point p, based on the light's texture.
func (l *Light) Emit(uv, p geom.Vec) Color {
	return l.texture.Map(uv, p)
}

// Metal describes a reflective material.
type Metal struct {
	texture Mapper
	rough   float64
	nonEmitter
}

// NewMetal creates a new Metal material with a given texture and roughness.
func NewMetal(texture Mapper, roughness float64) *Metal {
	return &Metal{texture: texture, rough: roughness}
}

// Scatter reflects incoming light rays about the normal,
// randomly modulated by the metal's roughness.
func (m *Metal) Scatter(in, norm geom.Unit, uv, p geom.Vec, rnd *rand.Rand) (out geom.Unit, attenuate Color, ok bool) {
	r := reflect(in, norm)
	out = geom.Vec(r).Plus(geom.RandVecInSphere(rnd).Scaled(m.rough)).Unit()
	return out, m.texture.Map(uv, p), out.Dot(norm) > 0
}

// Reflect reflects this unit vector about a normal vector n.
func reflect(u, n geom.Unit) geom.Unit {
	return geom.Unit(geom.Vec(u).Minus(geom.Vec(n).Scaled(2 * u.Dot(n)))) // TODO: prove this is still a unit vector
}

type nonEmitter struct{}

func (n nonEmitter) Emit(uv, p geom.Vec) Color {
	return black
}
