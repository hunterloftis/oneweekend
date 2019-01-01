package main

import (
	"fmt"

	"github.com/hunterloftis/oneweekend/pkg/img"
)

func main() {
	nx := 200
	ny := 100
	fmt.Println("P3")
	fmt.Println(nx, ny)
	fmt.Println("255")
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			col := img.NewColor(float64(i)/float64(nx), float64(j)/float64(ny), 0.2)
			ir := int(255.99 * col.R())
			ig := int(255.99 * col.G())
			ib := int(255.99 * col.B())
			fmt.Println(ir, ig, ib)
		}
	}
}
