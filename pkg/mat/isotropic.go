package mat

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/tex"
)

type Isotropic struct {
	nonEmitter
	albedo tex.Mapper
}

func NewIsotropic(albedo tex.Mapper) *Isotropic {
	return &Isotropic{albedo: albedo}
}

func (i *Isotropic) Scatter(in, norm geom.Unit, uv, p geom.Vec) (out geom.Unit, attenuate tex.Color, ok bool) {
	return geom.RandUnit(), i.albedo.Map(uv, p), true
}
