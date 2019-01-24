package mat

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/tex"
)

// Material represents a material that scatters light.
type Material interface {
	Scatter(in geom.Unit, n geom.Unit, u, v float64, p geom.Vec) (out geom.Unit, attenuation tex.Color, ok bool)
	Emit(u, v float64, p geom.Vec) tex.Color
}

type nonEmitter struct{}

func (n nonEmitter) Emit(u, v float64, p geom.Vec) tex.Color {
	return tex.NewColor(0, 0, 0)
}
