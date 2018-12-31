package main

import "fmt"

func main() {
	nx := 200
	ny := 100
	fmt.Println("P3")
	fmt.Println(nx, ny)
	fmt.Println("255")
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			r := float64(i) / float64(nx)
			g := float64(j) / float64(ny)
			b := 0.2
			ir := int(255.99 * r)
			ig := int(255.99 * g)
			ib := int(255.99 * b)
			fmt.Println(ir, ig, ib)
		}
	}
}
