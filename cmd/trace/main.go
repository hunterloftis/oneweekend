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
		wd, _ := os.Getwd()
		defer profile.Start(profile.ProfilePath(wd)).Stop()
	}
	w := trace.NewWindow(500, 500)
	cam, scene := final(w.Aspect())
	if err := w.WritePPM(os.Stdout, cam, scene, 500); err != nil {
		panic(err)
	}
}

func custom(aspect float64) (trace.Camera, *trace.BVH) {
	rand.Seed(10)
	nb := 20
	w := 100.0
	list := trace.NewList()
	ground := trace.NewLambert(trace.NewUniform(0.48, 0.83, 0.53))
	light := trace.NewLight(trace.NewUniform(1, 1, 1))
	for i := 0; i < nb; i++ {
		for j := 0; j < nb; j++ {
			min := geom.NewVec(-1000+float64(i)*w, 0, -1000+float64(j)*w)
			max := geom.NewVec(w, 1+99*rand.Float64(), w).Plus(min)
			list.Add(trace.NewBox(min, max, ground))
		}
	}
	list.Add(trace.NewRect(geom.NewVec(123, 554, 147), geom.NewVec(423, 554, 412), light))
	center := geom.NewVec(400, 400, 200)
	list.Add(trace.NewMovingSphere(center, center.Plus(geom.NewVec(30, 0, 0)), 0, 1, 50, trace.NewLambert(trace.NewUniform(0.7, 0.3, 0.1))))
	list.Add(trace.NewSphere(geom.NewVec(260, 150, 45), 50, trace.NewDielectric(1.5)))
	boundary := trace.NewSphere(geom.NewVec(0, 0, 0), 5000, trace.NewDielectric(1.5))
	list.Add(trace.NewVolume(boundary, 0.0001, trace.NewIso(trace.NewUniform(1, 1, 1))))
	f, err := os.Open("images/earthmap.jpg")
	if err != nil {
		panic(err)
	}
	earth, err := trace.NewImage(f)
	if err != nil {
		panic(err)
	}
	list.Add(trace.NewSphere(geom.NewVec(100, 150, 100), 50, trace.NewMetal(trace.NewBright(earth, 1.5), 0.1)))
	list.Add(trace.NewSphere(geom.NewVec(400, 200, 400), 100, trace.NewLight(trace.NewBright(earth, 3))))
	perlin := trace.NewNoise(0.1, 0.1, 0)
	list.Add(trace.NewSphere(geom.NewVec(220, 280, 300), 80, trace.NewLambert(perlin)))
	from := geom.NewVec(478, 278, -600)
	at := geom.NewVec(278, 278, 0)
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.NewUnit(0, 1, 0), 40, aspect, 0, focus, 0, 1)
	return cam, trace.NewBVH(0, 1, list.Surfaces()...)
}

func final(aspect float64) (trace.Camera, *trace.List) {
	rand.Seed(10)
	nb := 20
	w := 100.0
	ns := 1000
	list := trace.NewList()
	boxList := trace.NewList()
	boxList2 := trace.NewList()
	white := trace.NewLambert(trace.NewUniform(0.73, 0.73, 0.73))
	ground := trace.NewLambert(trace.NewUniform(0.48, 0.83, 0.53))
	light := trace.NewLight(trace.NewUniform(2, 2, 2))
	for i := 0; i < nb; i++ {
		for j := 0; j < nb; j++ {
			min := geom.NewVec(-1000+float64(i)*w, 0, -1000+float64(j)*w)
			max := geom.NewVec(w, 1+99*rand.Float64(), w).Plus(min)
			boxList.Add(trace.NewBox(min, max, ground))
		}
	}
	list.Add(trace.NewBVH(0, 1, boxList.Surfaces()...))
	list.Add(trace.NewRect(geom.NewVec(123, 554, 147), geom.NewVec(423, 554, 412), light))
	center := geom.NewVec(400, 400, 200)
	list.Add(trace.NewMovingSphere(center, center.Plus(geom.NewVec(30, 0, 0)), 0, 1, 50, trace.NewLambert(trace.NewUniform(0.7, 0.3, 0.1))))
	list.Add(trace.NewSphere(geom.NewVec(260, 150, 45), 50, trace.NewDielectric(1.5)))
	list.Add(trace.NewSphere(geom.NewVec(0, 150, 145), 50, trace.NewMetal(trace.NewUniform(0.8, 0.8, 0.9), 1)))
	boundary := trace.NewSphere(geom.NewVec(360, 150, 145), 70, trace.NewDielectric(1.5))
	list.Add(boundary)
	list.Add(trace.NewVolume(boundary, 0.2, trace.NewIso(trace.NewUniform(0.2, 0.4, 0.9))))
	boundary = trace.NewSphere(geom.NewVec(0, 0, 0), 5000, trace.NewDielectric(1.5))
	list.Add(trace.NewVolume(boundary, 0.0001, trace.NewIso(trace.NewUniform(1, 1, 1))))
	f, err := os.Open("images/earthmap.jpg")
	if err != nil {
		panic(err)
	}
	earth, err := trace.NewImage(f)
	if err != nil {
		panic(err)
	}
	list.Add(trace.NewSphere(geom.NewVec(400, 200, 400), 100, trace.NewLambert(earth)))
	perlin := trace.NewNoise(0.1, 0.1, 0)
	list.Add(trace.NewSphere(geom.NewVec(220, 280, 300), 80, trace.NewLambert(perlin)))
	for j := 0; j < ns; j++ {
		boxList2.Add(trace.NewSphere(geom.NewVec(165*rand.Float64(), 165*rand.Float64(), 165*rand.Float64()), 10, white))
	}
	list.Add(trace.NewTranslate(trace.NewRotateY(trace.NewBVH(0, 1, boxList2.Surfaces()...), 15), geom.NewVec(-100, 270, 395)))
	from := geom.NewVec(478, 278, -600)
	at := geom.NewVec(278, 278, 0)
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.NewUnit(0, 1, 0), 40, aspect, 0, focus, 0, 1)
	return cam, list
}

func cornellSmoke(aspect float64) (trace.Camera, *trace.List) {
	green := trace.NewLambert(trace.NewUniform(0.12, 0.45, 0.15))
	red := trace.NewLambert(trace.NewUniform(0.65, 0.05, 0.05))
	light := trace.NewLight(trace.NewUniform(7, 7, 7))
	white := trace.NewLambert(trace.NewUniform(0.73, 0.73, 0.73))
	smoke := trace.NewIso(trace.NewUniform(0, 0, 0))
	fog := trace.NewIso(trace.NewUniform(1, 1, 1))
	b1 := trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.NewVec(0, 0, 0), geom.NewVec(165, 165, 165), white), -18), geom.NewVec(130, 0, 65))
	b2 := trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.NewVec(0, 0, 0), geom.NewVec(165, 330, 165), white), 15), geom.NewVec(265, 0, 295))
	from := geom.NewVec(278, 278, -800)
	at := geom.NewVec(278, 278, 0)
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.NewUnit(0, 1, 0), 40, aspect, 0, focus, 0, 1)
	return cam, trace.NewList(
		trace.NewFlip(trace.NewRect(geom.NewVec(555, 0, 0), geom.NewVec(555, 555, 555), green)),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(0, 555, 555), red),
		trace.NewRect(geom.NewVec(113, 554, 127), geom.NewVec(443, 554, 432), light),
		trace.NewFlip(trace.NewRect(geom.NewVec(0, 555, 0), geom.NewVec(555, 555, 555), white)),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(555, 0, 555), white),
		trace.NewFlip(trace.NewRect(geom.NewVec(0, 0, 555), geom.NewVec(555, 555, 555), white)),
		trace.NewVolume(b1, 0.01, fog),
		trace.NewVolume(b2, 0.01, smoke),
	)
}

func cornell(aspect float64) (trace.Camera, *trace.List) {
	green := trace.NewLambert(trace.NewUniform(0.12, 0.45, 0.15))
	red := trace.NewLambert(trace.NewUniform(0.65, 0.05, 0.05))
	light := trace.NewLight(trace.NewUniform(15, 15, 15))
	white := trace.NewLambert(trace.NewUniform(0.73, 0.73, 0.73))
	from := geom.NewVec(278, 278, -800)
	at := geom.NewVec(278, 278, 0)
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.NewUnit(0, 1, 0), 40, aspect, 0, focus, 0, 1)
	return cam, trace.NewList(
		trace.NewFlip(trace.NewRect(geom.NewVec(555, 0, 0), geom.NewVec(555, 555, 555), green)),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(0, 555, 555), red),
		trace.NewRect(geom.NewVec(213, 554, 227), geom.NewVec(343, 554, 332), light),
		trace.NewFlip(trace.NewRect(geom.NewVec(0, 555, 0), geom.NewVec(555, 555, 555), white)),
		trace.NewRect(geom.NewVec(0, 0, 0), geom.NewVec(555, 0, 555), white),
		trace.NewFlip(trace.NewRect(geom.NewVec(0, 0, 555), geom.NewVec(555, 555, 555), white)),
		trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.NewVec(0, 0, 0), geom.NewVec(165, 165, 165), white), -18), geom.NewVec(130, 0, 65)),
		trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.NewVec(0, 0, 0), geom.NewVec(165, 330, 165), white), 15), geom.NewVec(265, 0, 295)),
	)
}

func simpleLight(aspect float64) (trace.Camera, *trace.List) {
	perlin := trace.NewNoise(4, 1, 2)
	from := geom.NewVec(25, 4, 6)
	at := geom.NewVec(0, 2, 0)
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.NewUnit(0, 1, 0), 20, aspect, 0, focus, 0, 1)
	return cam, trace.NewList(
		trace.NewSphere(geom.NewVec(0, -1000, 0), 1000, trace.NewLambert(perlin)),
		trace.NewSphere(geom.NewVec(0, 2, 0), 2, trace.NewLambert(perlin)),
		trace.NewSphere(geom.NewVec(0, 7, 0), 2, trace.NewLight(trace.NewUniform(4, 4, 4))),
		trace.NewRect(geom.NewVec(3, 1, -2), geom.NewVec(5, 3, -2), trace.NewLight(trace.NewUniform(4, 4, 4))),
	)
}
