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

func (l *Light) Scatter(in, norm geom.Unit, uv, p geom.Vec) (out geom.Unit, attenuate tex.Color, ok bool) {
	return out, attenuate, false
}

func (l *Light) Emit(uv, p geom.Vec) tex.Color {
	return l.Mapper.Map(uv, p)
}
