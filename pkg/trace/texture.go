package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

type Texture interface {
	At(u, v float64, p geom.Vec) Color
}
