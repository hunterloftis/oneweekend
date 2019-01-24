package main

import (
	"flag"
	"math/rand"
	"os"

	"github.com/hunterloftis/oneweekend/pkg/geom"
	"github.com/hunterloftis/oneweekend/pkg/mat"
	"github.com/hunterloftis/oneweekend/pkg/tex"
	"github.com/hunterloftis/oneweekend/pkg/trace"
	"github.com/pkg/profile"
)

func main() {
	p := flag.Bool("profile", false, "generate a cpu profile")
	if flag.Parse(); *p {
		defer profile.Start().Stop()
	}
	w := trace.NewWindow(500, 500)
	if err := w.WritePPM(os.Stdout, cornell(), 100); err != nil {
		panic(err)
	}
}

func cornell() *trace.List {
	green := mat.NewLambert(tex.NewSolid(tex.NewColor(0.12, 0.45, 0.15)))
	red := mat.NewLambert(tex.NewSolid(tex.NewColor(0.65, 0.05, 0.05)))
	light := mat.NewLight(tex.NewSolid(tex.NewColor(15, 15, 15)))
	white := mat.NewLambert(tex.NewSolid(tex.NewColor(0.73, 0.73, 0.73)))
	return trace.NewList(
		trace.NewRect(geom.NewVec(555, 0, 0), geom.NewVec(555, 555, 555), green),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(0, 555, 555), red),
		trace.NewRect(geom.NewVec(213, 554, 227), geom.NewVec(343, 554, 332), light),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(555, 0, 555), white),
		trace.NewRect(geom.NewVec(0, 0, 555), geom.NewVec(555, 555, 555), white),
	)
}

func simpleLight() *trace.List {
	perlin := tex.NewNoise(4)
	return trace.NewList(
		trace.NewSphere(geom.NewVec(0, -1000, 0), 1000, mat.NewLambert(perlin)),
		trace.NewSphere(geom.NewVec(0, 2, 0), 2, mat.NewLambert(perlin)),
		trace.NewSphere(geom.NewVec(0, 7, 0), 2, mat.NewLight(tex.NewSolid(tex.NewColor(4, 4, 4)))),
		trace.NewRect(geom.NewVec(3, 1, -2), geom.NewVec(5, 3, -2), mat.NewLight(tex.NewSolid(tex.NewColor(4, 4, 4)))),
	)
}

func earth() *trace.Sphere {
	f, err := os.Open("images/earthmap.jpg")
	if err != nil {
		panic(err)
	}
	t, err := tex.NewImage(f)
	if err != nil {
		panic(err)
	}
	return trace.NewSphere(geom.NewVec(0, 0, 0), 1, mat.NewLambert(t))
}

func twoPerlinSpheres() *trace.List {
	perlin := tex.NewNoise(5)
	return trace.NewList(
		trace.NewSphere(geom.NewVec(0, -1000, 0), 1000, mat.NewLambert(perlin)),
		trace.NewSphere(geom.NewVec(0, 2, 0), 2, mat.NewLambert(perlin)),
	)
}

func cover() *trace.BVH {
	check := tex.NewChecker(10,
		tex.NewSolid(tex.NewColor(0.2, 0.3, 0.1)),
		tex.NewSolid(tex.NewColor(0.9, 0.9, 0.9)),
	)
	l := trace.NewList(
		trace.NewSphere(geom.NewVec(0, -1000, 0), 1000, mat.NewLambert(check)),
		trace.NewSphere(geom.NewVec(0, 1, 0), 1, mat.NewDielectric(1.5)),
		trace.NewSphere(geom.NewVec(-4, 1, 0), 1, mat.NewLambert(tex.NewSolid(tex.NewColor(0.4, 0.2, 0.1)))),
		trace.NewSphere(geom.NewVec(4, 1, 0), 1, mat.NewMetal(tex.NewColor(0.7, 0.6, 0.5), 0)),
	)
	for a := -11.0; a < 11; a++ {
		for b := -11.0; b < 11; b++ {
			center := geom.NewVec(a+0.9*rand.Float64(), 0.2, b+0.9*rand.Float64())
			if center.Minus(geom.NewVec(4, 0.2, 0)).Len() <= 0.9 {
				continue
			}
			m, move := randMat()
			if move {
				l.Add(trace.NewMovingSphere(center, center.Plus(geom.NewVec(0, 0.5*rand.Float64(), 0)), 0, 1, 0.2, m))
			} else {
				l.Add(trace.NewSphere(center, 0.2, m))
			}
		}
	}
	return trace.NewBVH(0, 0, 1, l.HH...)
}

func randMat() (mat.Material, bool) {
	m := rand.Float64()
	if m < 0.8 {
		c := tex.NewColor(rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64())
		return mat.NewLambert(tex.NewSolid(c)), true
	}
	if m < 0.95 {
		c := tex.NewColor(0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()))
		return mat.NewMetal(c, 0.5*rand.Float64()), false
	}
	return mat.NewDielectric(1.5), false
}
