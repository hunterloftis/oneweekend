package tex

import (
	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Mapper maps a u, v coordinate in 3d space p to a Color
type Mapper interface {
	Map(u, v float64, p geom.Vec) Color
}
