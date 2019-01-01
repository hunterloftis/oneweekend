package main

import (
	"os"

	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/trace"
)

func main() {
	l := trace.NewList(
		trace.NewSphere(geom.NewVec(0, 0, -1), 0.5),
		trace.NewSphere(geom.NewVec(0, -100.5, -1), 100),
	)
	f := trace.NewFrame(200, 100)
	if err := f.WritePPM(os.Stdout, l, 100); err != nil {
		panic(err)
	}
}
