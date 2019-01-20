package tex

import (
	"math"
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
	u := p.X() - math.Floor(p.X())
	v := p.Y() - math.Floor(p.Y())
	w := p.Z() - math.Floor(p.Z())
	u = u * u * (3 - 2*u)
	v = v * v * (3 - 2*v)
	w = w * w * (3 - 2*w)
	i := int(math.Floor(p.X()))
	j := int(math.Floor(p.Y()))
	k := int(math.Floor(p.Z()))
	c := make([]float64, 8)
	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				x := permX[(i+di)&255]
				y := permY[(j+dj)&255]
				z := permZ[(k+dk)&255]
				c[4*di+2*dj+dk] = rndFloat[x^y^z]
			}
		}
	}
	return trilinear(c, u, v, w)
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

func trilinear(c []float64, u, v, w float64) (sum float64) {
	for i := 0.0; i < 2; i++ {
		for j := 0.0; j < 2; j++ {
			for k := 0.0; k < 2; k++ {
				xyz := c[4*int(i)+2*int(j)+int(k)]
				sum += (i*u + (1-i)*(1-u)) *
					(j*v + (1-j)*(1-v)) *
					(k*w + (1-k)*(1-w)) *
					xyz
			}
		}
	}
	return sum
}
