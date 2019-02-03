package geom

import (
	"math/rand"
)

// Unit is a 3D vector (x, y, z) with length 1.
type Unit Vec

// Dot returns the dot product of two unit vectors.
// Values above zero indicate vectors pointing in the same hemisphere.
// Values below zero indicate vectors pointing towards opposite hemispheres.
func (u Unit) Dot(u2 Unit) float64 {
	return Vec(u).Dot(Vec(u2))
}

// Inv inverts the unit vector.
func (u Unit) Inv() Unit {
	return Unit(Vec(u).Inv())
}

// RandUnit generates a random unit vector.
// BUG(Hunter): This isn't exactly uniform.
func RandUnit(rnd *rand.Rand) Unit {
	return Vec{2*rnd.Float64() - 1, 2*rnd.Float64() - 1, 2*rnd.Float64() - 1}.Unit()
}

// Scaled returns this unit vector scaled into a non-unit vector by n.
func (u Unit) Scaled(n float64) Vec {
	return Vec(u).Scaled(n)
}
