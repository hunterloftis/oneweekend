package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

var (
	lowerLeft  = geom.NewVec(-2, -1, -1)
	horizontal = geom.NewVec(4, 0, 0)
	vertical   = geom.NewVec(0, 2, 0)
	origin     = geom.NewVec(0, 0, 0)
)

// Camera originates Rays.
type Camera struct {
}

// Ray returns a Ray passing through a given u, v coordinate.
func (c Camera) Ray(u, v float64) geom.Ray {
	return geom.NewRay(
		origin,
		lowerLeft.Plus((horizontal.Scaled(u)).Plus(vertical.Scaled(v))).ToUnit(),
	)
}
