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

func NewCamera(vfov, aspect float64) (c Camera) {
	theta := vfov * math.Pi / 180
	halfH := math.Tan(theta / 2)
	halfW := aspect * halfH
	c.lowerLeft = geom.NewVec(-halfW, -halfH, -1)
	c.horizontal = geom.NewVec(2*halfW, 0, 0)
	c.vertical = geom.NewVec(0, 2*halfH, 0)
	c.origin = geom.NewVec(0, 0, 0)
	return
}

// Ray returns a Ray passing through a given u, v coordinate.
func (c Camera) Ray(u, v float64) geom.Ray {
	return geom.NewRay(
		c.origin,
		c.lowerLeft.Plus((c.horizontal.Scaled(u)).Plus(c.vertical.Scaled(v))).ToUnit(),
	)
}
