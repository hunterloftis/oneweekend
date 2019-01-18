package main

import (
	"flag"
	"math/rand"
	"os"

	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/trace"
	"github.com/pkg/profile"
)

func main() {
	p := flag.Bool("profile", false, "generate a cpu profile")
	if flag.Parse(); *p {
		defer profile.Start().Stop()
	}
	w := trace.NewWindow(120, 80)
	if err := w.WritePPM(os.Stdout, scene(), 100); err != nil {
		panic(err)
	}
}

func scene() *trace.List {
	gray := trace.NewLambert(trace.NewColor(0.5, 0.5, 0.5))
	l := trace.NewList(
		trace.NewSphere(geom.NewVec(0, -1000, 0), 1000, gray),
		trace.NewSphere(geom.NewVec(0, 1, 0), 1, trace.NewDielectric(1.5)),
		trace.NewSphere(geom.NewVec(-4, 1, 0), 1, trace.NewLambert(trace.NewColor(0.4, 0.2, 0.1))),
		trace.NewSphere(geom.NewVec(4, 1, 0), 1, trace.NewMetal(trace.NewColor(0.7, 0.6, 0.5), 0)),
	)
	for a := -11.0; a < 11; a++ {
		for b := -11.0; b < 11; b++ {
			center := geom.NewVec(a+0.9*rand.Float64(), 0.2, b+0.9*rand.Float64())
			if center.Minus(geom.NewVec(4, 0.2, 0)).Len() <= 0.9 {
				continue
			}
			l.Add(trace.NewSphere(center, 0.2, mat()))
		}
	}
	return l
}

func mat() trace.Material {
	m := rand.Float64()
	if m < 0.8 {
		c := trace.NewColor(rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64())
		return trace.NewLambert(c)
	}
	if m < 0.95 {
		c := trace.NewColor(0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()))
		return trace.NewMetal(c, 0.5*rand.Float64())
	}
	return trace.NewDielectric(1.5)
}
