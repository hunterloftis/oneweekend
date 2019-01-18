package trace

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Camera originates Rays.
type Camera struct {
	lowerLeft    geom.Vec
	horizontal   geom.Vec
	vertical     geom.Vec
	origin       geom.Vec
	u, v, w      geom.Unit
	lensRadius   float64
	time0, time1 float64
}

// NewCamera creates a new Camera
// TODO: this argument list is getting pretty ridiculous
func NewCamera(lookFrom, lookAt geom.Vec, vup geom.Unit, vfov, aspect, aperture, focus, t0, t1 float64) (c Camera) {
	theta := vfov * math.Pi / 180
	halfH := math.Tan(theta / 2)
	halfW := aspect * halfH

	c.w = lookFrom.Minus(lookAt).Unit()
	c.u = vup.Cross(c.w.Vec).Unit()
	c.v = c.w.Cross(c.u.Vec).Unit()

	width := c.u.Scaled(halfW * focus)
	height := c.v.Scaled(halfH * focus)
	dist := c.w.Scaled(focus)

	c.time0 = t0
	c.time1 = t1
	c.lensRadius = aperture / 2
	c.origin = lookFrom
	c.lowerLeft = c.origin.Minus(width).Minus(height).Minus(dist)
	c.horizontal = width.Scaled(2)
	c.vertical = height.Scaled(2)
	return
}

// Ray returns a Ray passing through a given s, t coordinate.
func (c Camera) Ray(s, t float64) Ray {
	rd := geom.RandVecInDisk().Scaled(c.lensRadius)
	offset := c.u.Scaled(rd.X()).Plus(c.v.Scaled(rd.Y()))
	source := c.origin.Plus(offset)
	dest := c.lowerLeft.Plus(c.horizontal.Scaled(s).Plus(c.vertical.Scaled(t)))
	time := c.time0 + (c.time1-c.time0)*rand.Float64()
	return NewRay(source, dest.Minus(source).Unit(), time)
}
