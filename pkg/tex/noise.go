package tex

import "github.com/hunterloftis/oneweekend/pkg/geom"

type Noise struct{}

func NewNoise() Noise {
	return Noise{}
}

func (n Noise) Map(u, v float64, p geom.Vec) Color {
	return NewColor(1, 1, 1).Scaled(Perlin(p))
}
