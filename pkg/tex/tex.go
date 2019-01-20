package tex

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
)

type Mapper interface {
	Map(u, v float64, p geom.Vec) Color
}
