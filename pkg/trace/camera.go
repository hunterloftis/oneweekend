package trace

import (
	"math"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Camera originates Rays.
type Camera struct {
	lowerLeft  geom.Vec
	horizontal geom.Vec
	vertical   geom.Vec
	origin     geom.Vec
}

func NewCamera(lookFrom, lookAt geom.Vec, vup geom.Unit, vfov, aspect float64) (c Camera) {
	theta := vfov * math.Pi / 180
	halfH := math.Tan(theta / 2)
	halfW := aspect * halfH

	w := lookFrom.Minus(lookAt).Unit()
	u := vup.Cross(w.Vec).Unit()
	v := w.Cross(u.Vec).Unit()

	c.origin = lookFrom
	c.lowerLeft = c.origin.Minus(u.Scaled(halfW)).Minus(v.Scaled(halfH)).Minus(w.Vec)
	c.horizontal = u.Scaled(2 * halfW)
	c.vertical = v.Scaled(2 * halfH)
	return
}

// Ray returns a Ray passing through a given u, v coordinate.
func (c Camera) Ray(u, v float64) geom.Ray {
	return geom.NewRay(
		c.origin,
		c.lowerLeft.Plus((c.horizontal.Scaled(u)).Plus(c.vertical.Scaled(v))).Minus(c.origin).Unit(),
	)
}
