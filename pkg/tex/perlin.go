package tex

import (
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

var (
	rndFloat []float64
	permX    []int
	permY    []int
	permZ    []int
)

func init() {
	rndFloat = perlinGen()
	permX = perlinGenPerm()
	permY = perlinGenPerm()
	permZ = perlinGenPerm()
}

func Perlin(p geom.Vec) float64 {
	// u := p.X() - math.Floor(p.X())
	// v := p.Y() - math.Floor(p.Y())
	// w := p.Z() - math.Floor(p.Z())
	i := int(4*p.X()) & 255
	j := int(4*p.Y()) & 255
	k := int(4*p.Z()) & 255
	return rndFloat[permX[i]^permY[j]^permZ[k]]
}

func perlinGen() []float64 {
	p := make([]float64, 256)
	for i := 0; i < 256; i++ {
		p[i] = rand.Float64()
	}
	return p
}

func perlinPermute(p []int, n int) {
	for i := n - 1; i > 0; i-- {
		target := int(rand.Float64() * float64(i+1))
		p[i], p[target] = p[target], p[i]
	}
}

func perlinGenPerm() []int {
	p := make([]int, 256)
	for i := 0; i < 256; i++ {
		p[i] = i
	}
	perlinPermute(p, 256)
	return p
}
