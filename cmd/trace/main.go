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
	w := trace.NewWindow(898, 898)
	if err := w.WritePPM(os.Stdout, final(), 5000); err != nil {
		panic(err)
	}
}

func final() *trace.List {
	rand.Seed(10)
	nb := 20
	w := 100.0
	ns := 1000
	list := trace.NewList()
	boxList := trace.NewList()
	boxList2 := trace.NewList()
	white := mat.NewLambert(tex.NewConstant(tex.NewColor(0.73, 0.73, 0.73)))
	ground := mat.NewLambert(tex.NewConstant(tex.NewColor(0.48, 0.83, 0.53)))
	light := mat.NewLight(tex.NewConstant(tex.NewColor(7, 7, 7)))
	for i := 0; i < nb; i++ {
		for j := 0; j < nb; j++ {
			min := geom.NewVec(-1000+float64(i)*w, 0, -1000+float64(j)*w)
			max := geom.NewVec(w, 1+99*rand.Float64(), w).Plus(min)
			boxList.Add(trace.NewBox(min, max, ground))
		}
	}
	list.Add(trace.NewBVH(0, 0, 1, boxList.HH...))
	list.Add(trace.NewRect(geom.NewVec(123, 554, 147), geom.NewVec(423, 554, 412), light))
	center := geom.NewVec(400, 400, 200)
	list.Add(trace.NewMovingSphere(center, center.Plus(geom.NewVec(30, 0, 0)), 0, 1, 50, mat.NewLambert(tex.NewConstant(tex.NewColor(0.7, 0.3, 0.1)))))
	list.Add(trace.NewSphere(geom.NewVec(260, 150, 45), 50, mat.NewDielectric(1.5)))
	list.Add(trace.NewSphere(geom.NewVec(0, 150, 145), 50, mat.NewMetal(tex.NewColor(0.8, 0.8, 0.9), 1)))
	boundary := trace.NewSphere(geom.NewVec(360, 150, 145), 70, mat.NewDielectric(1.5))
	list.Add(boundary)
	list.Add(trace.NewVolume(boundary, 0.2, mat.NewIsotropic(tex.NewConstant(tex.NewColor(0.2, 0.4, 0.9)))))
	boundary = trace.NewSphere(geom.NewVec(0, 0, 0), 5000, mat.NewDielectric(1.5))
	list.Add(trace.NewVolume(boundary, 0.0001, mat.NewIsotropic(tex.NewConstant(tex.NewColor(1, 1, 1)))))
	f, err := os.Open("images/earthmap.jpg")
	if err != nil {
		panic(err)
	}
	earth, err := tex.NewImage(f)
	if err != nil {
		panic(err)
	}
	list.Add(trace.NewSphere(geom.NewVec(400, 200, 400), 100, mat.NewLambert(earth)))
	perlin := tex.NewNoise(0.1)
	list.Add(trace.NewSphere(geom.NewVec(220, 280, 300), 80, mat.NewLambert(perlin)))
	for j := 0; j < ns; j++ {
		boxList2.Add(trace.NewSphere(geom.NewVec(165*rand.Float64(), 165*rand.Float64(), 165*rand.Float64()), 10, white))
	}
	list.Add(trace.NewTranslate(trace.NewRotateY(trace.NewBVH(0, 0, 1, boxList2.HH...), 15), geom.NewVec(-100, 270, 395)))
	return list
}

func cornellSmoke() *trace.List {
	green := mat.NewLambert(tex.NewConstant(tex.NewColor(0.12, 0.45, 0.15)))
	red := mat.NewLambert(tex.NewConstant(tex.NewColor(0.65, 0.05, 0.05)))
	light := mat.NewLight(tex.NewConstant(tex.NewColor(7, 7, 7)))
	white := mat.NewLambert(tex.NewConstant(tex.NewColor(0.73, 0.73, 0.73)))
	smoke := mat.NewIsotropic(tex.NewConstant(tex.NewColor(0, 0, 0)))
	fog := mat.NewIsotropic(tex.NewConstant(tex.NewColor(1, 1, 1)))
	b1 := trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.NewVec(0, 0, 0), geom.NewVec(165, 165, 165), white), -18), geom.NewVec(130, 0, 65))
	b2 := trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.NewVec(0, 0, 0), geom.NewVec(165, 330, 165), white), 15), geom.NewVec(265, 0, 295))
	return trace.NewList(
		trace.NewFlipped(trace.NewRect(geom.NewVec(555, 0, 0), geom.NewVec(555, 555, 555), green)),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(0, 555, 555), red),
		trace.NewRect(geom.NewVec(113, 554, 127), geom.NewVec(443, 554, 432), light),
		trace.NewFlipped(trace.NewRect(geom.NewVec(0, 555, 0), geom.NewVec(555, 555, 555), white)),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(555, 0, 555), white),
		trace.NewFlipped(trace.NewRect(geom.NewVec(0, 0, 555), geom.NewVec(555, 555, 555), white)),
		trace.NewVolume(b1, 0.01, fog),
		trace.NewVolume(b2, 0.01, smoke),
	)
}

func cornell() *trace.List {
	green := mat.NewLambert(tex.NewConstant(tex.NewColor(0.12, 0.45, 0.15)))
	red := mat.NewLambert(tex.NewConstant(tex.NewColor(0.65, 0.05, 0.05)))
	light := mat.NewLight(tex.NewConstant(tex.NewColor(15, 15, 15)))
	white := mat.NewLambert(tex.NewConstant(tex.NewColor(0.73, 0.73, 0.73)))
	return trace.NewList(
		trace.NewFlipped(trace.NewRect(geom.NewVec(555, 0, 0), geom.NewVec(555, 555, 555), green)),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(0, 555, 555), red),
		trace.NewRect(geom.NewVec(213, 554, 227), geom.NewVec(343, 554, 332), light),
		trace.NewFlipped(trace.NewRect(geom.NewVec(0, 555, 0), geom.NewVec(555, 555, 555), white)),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(555, 0, 555), white),
		trace.NewFlipped(trace.NewRect(geom.NewVec(0, 0, 555), geom.NewVec(555, 555, 555), white)),
		trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.NewVec(0, 0, 0), geom.NewVec(165, 165, 165), white), -18), geom.NewVec(130, 0, 65)),
		trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.NewVec(0, 0, 0), geom.NewVec(165, 330, 165), white), 15), geom.NewVec(265, 0, 295)),
	)
}

func simpleLight() *trace.List {
	perlin := tex.NewNoise(4)
	return trace.NewList(
		trace.NewSphere(geom.NewVec(0, -1000, 0), 1000, mat.NewLambert(perlin)),
		trace.NewSphere(geom.NewVec(0, 2, 0), 2, mat.NewLambert(perlin)),
		trace.NewSphere(geom.NewVec(0, 7, 0), 2, mat.NewLight(tex.NewConstant(tex.NewColor(4, 4, 4)))),
		trace.NewRect(geom.NewVec(3, 1, -2), geom.NewVec(5, 3, -2), mat.NewLight(tex.NewConstant(tex.NewColor(4, 4, 4)))),
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
		tex.NewConstant(tex.NewColor(0.2, 0.3, 0.1)),
		tex.NewConstant(tex.NewColor(0.9, 0.9, 0.9)),
	)
	l := trace.NewList(
		trace.NewSphere(geom.NewVec(0, -1000, 0), 1000, mat.NewLambert(check)),
		trace.NewSphere(geom.NewVec(0, 1, 0), 1, mat.NewDielectric(1.5)),
		trace.NewSphere(geom.NewVec(-4, 1, 0), 1, mat.NewLambert(tex.NewConstant(tex.NewColor(0.4, 0.2, 0.1)))),
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

func randMat() (trace.Material, bool) {
	m := rand.Float64()
	if m < 0.8 {
		c := tex.NewColor(rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64())
		return mat.NewLambert(tex.NewConstant(c)), true
	}
	if m < 0.95 {
		c := tex.NewColor(0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()))
		return mat.NewMetal(c, 0.5*rand.Float64()), false
	}
	return mat.NewDielectric(1.5), false
}
