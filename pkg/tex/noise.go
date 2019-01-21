package tex

import "github.com/hunterloftis/oneweekend/pkg/geom"

type Noise struct {
	Scale float64
}

func NewNoise(scale float64) Noise {
	return Noise{Scale: scale}
}

func (n Noise) Map(u, v float64, p geom.Vec) Color {
	// bright := 0.5*Perlin(p.Scaled(n.Scale)) + 0.5
	bright := turb(p.Scaled(n.Scale), 7)
	return NewColor(1, 1, 1).Scaled(bright)
}
