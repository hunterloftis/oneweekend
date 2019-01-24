package mat

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/tex"
)

type nonEmitter struct{}

func (n nonEmitter) Emit(uv, p geom.Vec) tex.Color {
	return tex.NewColor(0, 0, 0)
}
