package geom

import (
	"fmt"
	"io"
	"math"
	"math/rand"
)

// Vec represents a 3-element vector
type Vec [3]float64

// RandVecInSphere creates a random Vec within a unit sphere
// TODO: I don't like rejection methods. Isn't there a way to generate 2 angles and accomplish the same thing reliably?
func RandVecInSphere(rnd *rand.Rand) Vec {
	for {
		v := Vec{rnd.Float64(), rnd.Float64(), rnd.Float64()}.Scaled(2).Minus(Vec{1, 1, 1})
		if v.LenSq() < 1 {
			return v
		}
	}
}

// RandVecInDisk creates a random Vec within a unit disk
// TODO: more rejection methods :/
func RandVecInDisk(rnd *rand.Rand) Vec {
	xy := Vec{1, 1, 0}
	for {
		v := Vec{rnd.Float64(), rnd.Float64(), 0}.Scaled(2).Minus(xy)
		if v.Dot(v) < 1 {
			return v
		}
	}
}

// X returns the first element
func (v Vec) X() float64 {
	return v[0]
}

// Y returns the second element
func (v Vec) Y() float64 {
	return v[1]
}

// Z returns the third element
func (v Vec) Z() float64 {
	return v[2]
}

// Inv returns this vector's inverse as a new vector
func (v Vec) Inv() Vec {
	return Vec{-v[0], -v[1], -v[2]}
}

// Len returns the vector's length
func (v Vec) Len() float64 {
	return math.Sqrt(v.LenSq())
}

// LenSq returns the square of the vector's length
func (v Vec) LenSq() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

// Unit converts this vector to a unit vector
func (v Vec) Unit() (u Unit) {
	k := 1.0 / v.Len()
	u[0] = v[0] * k
	u[1] = v[1] * k
	u[2] = v[2] * k
	return
}

// IStream streams in space-separated vector values from a Reader
func (v Vec) IStream(r io.Reader) error {
	_, err := fmt.Fscan(r, v[0], v[1], v[2])
	return err
}

// OStream writes space-separated vector values to a Writer
func (v Vec) OStream(w io.Writer) error {
	_, err := fmt.Fprint(w, v[0], v[1], v[2])
	return err
}

// Plus returns the sum of two vectors
func (v Vec) Plus(v2 Vec) Vec {
	return Vec{v[0] + v2[0], v[1] + v2[1], v[2] + v2[2]}
}

// Minus returns the difference of two vectors
func (v Vec) Minus(v2 Vec) Vec {
	return Vec{v[0] - v2[0], v[1] - v2[1], v[2] - v2[2]}
}

// Times returns the multiplication of two vectors
func (v Vec) Times(v2 Vec) Vec {
	return Vec{v[0] * v2[0], v[1] * v2[1], v[2] * v2[2]}
}

// Div returns the division of two vectors
func (v Vec) Div(v2 Vec) Vec {
	return Vec{v[0] / v2[0], v[1] / v2[1], v[2] / v2[2]}
}

// Scaled returns a vector scaled by a scalar
func (v Vec) Scaled(n float64) Vec {
	return Vec{v[0] * n, v[1] * n, v[2] * n}
}

// Dot returns the dot product of two vectors
func (v Vec) Dot(v2 Vec) float64 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

// Cross returns the cross product of two vectors
func (v Vec) Cross(v2 Vec) Vec {
	return Vec{
		v[1]*v2[2] - v[2]*v2[1],
		v[2]*v2[0] - v[0]*v2[2],
		v[0]*v2[1] - v[1]*v2[0],
	}
}

// Min returns a new Vector using the lowest elements of two vectors
func (v Vec) Min(v2 Vec) Vec {
	for i := 0; i < 3; i++ {
		if v2[i] < v[i] {
			v[i] = v2[i]
		}
	}
	return v
}

// Max returns a new Vector using the highest element of two vectors
func (v Vec) Max(v2 Vec) Vec {
	for i := 0; i < 3; i++ {
		if v2[i] > v[i] {
			v[i] = v2[i]
		}
	}
	return v
}
