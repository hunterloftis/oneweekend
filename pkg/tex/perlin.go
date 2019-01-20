package tex

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

var (
	rndUnit []geom.Unit
	permX   []int
	permY   []int
	permZ   []int
)

func init() {
	rndUnit = generate()
	permX = generatePerm()
	permY = generatePerm()
	permZ = generatePerm()
}

func Perlin(p geom.Vec) float64 {
	i := int(math.Floor(p.X()))
	j := int(math.Floor(p.Y()))
	k := int(math.Floor(p.Z()))
	u := p.X() - math.Floor(p.X())
	v := p.Y() - math.Floor(p.Y())
	w := p.Z() - math.Floor(p.Z())
	c := make([]geom.Unit, 8)
	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				x := permX[(i+di)&255]
				y := permY[(j+dj)&255]
				z := permZ[(k+dk)&255]
				c[4*di+2*dj+dk] = rndUnit[x^y^z]
			}
		}
	}
	return interp(c, u, v, w)
}

func generate() []geom.Unit {
	p := make([]geom.Unit, 256)
	for i := 0; i < 256; i++ {
		p[i] = geom.RandUnit()
	}
	return p
}

func permute(p []int, n int) {
	for i := n - 1; i > 0; i-- {
		target := int(rand.Float64() * float64(i+1))
		p[i], p[target] = p[target], p[i]
	}
}

func generatePerm() []int {
	p := make([]int, 256)
	for i := 0; i < 256; i++ {
		p[i] = i
	}
	permute(p, 256)
	return p
}

func interp(c []geom.Unit, u, v, w float64) (sum float64) {
	uu := u * u * (3 - 2*u)
	vv := v * v * (3 - 2*v)
	ww := w * w * (3 - 2*w)
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				weightV := geom.NewVec(u-float64(i), v-float64(j), w-float64(k))
				xyz := c[4*i+2*j+k]
				sum += (float64(i)*uu + (1-float64(i))*(1-uu)) *
					(float64(j)*vv + (1-float64(j))*(1-vv)) *
					(float64(k)*ww + (1-float64(k))*(1-ww)) *
					xyz.Vec.Dot(weightV)
			}
		}
	}
	return sum
}
