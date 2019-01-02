package geom

// Unit represents a unit vector (length 1)
type Unit struct {
	Vec
}

// Reflect reflects this unit vector about a normal vector n.
// TODO: verify that these operations on unit vectors always return a unit vector
func (u Unit) Reflect(n Unit) Unit {
	return Unit{Vec: u.Minus(n.Scaled(2 * u.Dot(n)))}
}

// Dot returns the dot product of two unit vectors.
// Values above zero indicate vectors pointing in the same hemisphere.
// Values below zero indicate vectors pointing towards opposite hemispheres.
// TODO: check to see if that's generally true, or only for unit vectors
func (u Unit) Dot(u2 Unit) float64 {
	return u.E[0]*u2.E[0] + u.E[1]*u2.E[1] + u.E[2]*u2.E[2]
}
