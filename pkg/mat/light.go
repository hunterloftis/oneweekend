package mat

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/tex"
)

type Light struct {
	Mapper tex.Mapper
}

func NewLight(m tex.Mapper) *Light {
	return &Light{Mapper: m}
}

func (l *Light) Scatter(in, n geom.Unit, u, v float64, p geom.Vec) (out geom.Unit, attenuate tex.Color, ok bool) {
	return out, attenuate, false
}

func (l *Light) Emit(u, v float64, p geom.Vec) tex.Color {
	return l.Mapper.Map(u, v, p)
}
