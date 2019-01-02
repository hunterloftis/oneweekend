package main

import (
	"os"

	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/trace"
)

func main() {
	blue := trace.NewLambert(trace.NewColor(0.1, 0.2, 0.5))
	yellow := trace.NewLambert(trace.NewColor(0.8, 0.8, 0.0))
	bronze := trace.NewMetal(trace.NewColor(0.8, 0.6, 0.2), 1)
	glass := trace.NewDielectric(1.5)
	l := trace.NewList(
		trace.NewSphere(geom.NewVec(0, 0, -1), 0.5, blue),
		trace.NewSphere(geom.NewVec(0, -100.5, -1), 100, yellow),
		trace.NewSphere(geom.NewVec(1, 0, -1), 0.5, bronze),
		trace.NewSphere(geom.NewVec(-1, 0, -1), 0.5, glass),
	)
	f := trace.NewFrame(400, 200)
	if err := f.WritePPM(os.Stdout, l, 100); err != nil {
		panic(err)
	}
}
