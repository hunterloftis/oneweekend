package mat

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/tex"
)

// Scatterer represents a material that scatters light.
type Scatterer interface {
	Scatter(in geom.Unit, n geom.Unit, p geom.Vec, u, v float64) (out geom.Unit, attenuation tex.Color, ok bool)
}
