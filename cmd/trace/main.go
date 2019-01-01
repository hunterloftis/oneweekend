package main

import (
	"fmt"

	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/img"
)

func main() {
	nx := 200
	ny := 100
	fmt.Println("P3")
	fmt.Println(nx, ny)
	fmt.Println("255")

	lowerLeft := geom.NewVec(-2, -1, -1)
	horizontal := geom.NewVec(4, 0, 0)
	vertical := geom.NewVec(0, 2, 0)
	origin := geom.NewVec(0, 0, 0)

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u := float64(i) / float64(nx)
			v := float64(j) / float64(ny)
			r := geom.NewRay(
				origin,
				lowerLeft.Plus((horizontal.Scaled(u)).Plus(vertical.Scaled(v))).ToUnit(),
			)
			col := color(r)
			ir := int(255.99 * col.R())
			ig := int(255.99 * col.G())
			ib := int(255.99 * col.B())
			fmt.Println(ir, ig, ib)
		}
	}
}

func color(r geom.Ray) img.Color {
	t := 0.5 * (r.Dir.Y() + 1.0)
	white := img.NewColor(1, 1, 1).Scaled(1 - t)
	blue := img.NewColor(0.5, 0.7, 1).Scaled(t)
	return white.Plus(blue)
}
