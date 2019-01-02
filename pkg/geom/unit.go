package geom

// Unit represents a unit vector (length 1)
type Unit struct {
	Vec
}

// Dot returns the dot product of two unit vectors.
// Values above zero indicate vectors pointing in the same hemisphere.
// Values below zero indicate vectors pointing towards opposite hemispheres.
// TODO: check to see if that's generally true, or only for unit vectors
func (u Unit) Dot(u2 Unit) float64 {
	return u.Vec.Dot(u2.Vec)
}

// Inv inverts the unit vector.
func (u Unit) Inv() Unit {
	return Unit{Vec: u.Vec.Inv()}
}
