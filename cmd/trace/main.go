package main

import (
	"math"
	"os"

	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/trace"
)

func main() {
	blue := trace.NewLambert(trace.NewColor(0, 0, 1))
	red := trace.NewLambert(trace.NewColor(1, 0, 0))
	r := math.Cos(math.Pi / 4)
	l := trace.NewList(
		trace.NewSphere(geom.NewVec(-r, 0, -1), r, blue),
		trace.NewSphere(geom.NewVec(r, 0, -1), r, red),
	)
	w := trace.NewWindow(400, 200)
	if err := w.WritePPM(os.Stdout, l, 100); err != nil {
		panic(err)
	}
}
