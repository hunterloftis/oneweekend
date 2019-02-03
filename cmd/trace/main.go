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
	w := trace.NewWindow(250, 250)
	cam, scene := final(w.Aspect())
	if err := w.WritePPM(os.Stdout, cam, scene, 250); err != nil {
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
			min := geom.Vec{-1000 + float64(i)*w, 0, -1000 + float64(j)*w}
			max := geom.Vec{w, 1 + 99*rand.Float64(), w}.Plus(min)
			list.Add(trace.NewBox(min, max, ground))
		}
	}
	list.Add(trace.NewRect(geom.Vec{123, 554, 147}, geom.Vec{423, 554, 412}, light))
	center := geom.Vec{400, 400, 200}
	list.Add(trace.NewMovingSphere(center, center.Plus(geom.Vec{30, 0, 0}), 0, 1, 50, trace.NewLambert(trace.NewUniform(0.7, 0.3, 0.1))))
	list.Add(trace.NewSphere(geom.Vec{260, 150, 45}, 50, trace.NewDielectric(1.5)))
	boundary := trace.NewSphere(geom.Vec{0, 0, 0}, 5000, trace.NewDielectric(1.5))
	list.Add(trace.NewVolume(boundary, 0.0001, trace.NewIso(trace.NewUniform(1, 1, 1))))
	f, err := os.Open("images/earthmap.jpg")
	if err != nil {
		panic(err)
	}
	earth, err := trace.NewImage(f)
	if err != nil {
		panic(err)
	}
	list.Add(trace.NewSphere(geom.Vec{100, 150, 100}, 50, trace.NewMetal(trace.NewBright(earth, 1.5), 0.1)))
	list.Add(trace.NewSphere(geom.Vec{400, 200, 400}, 100, trace.NewLight(trace.NewBright(earth, 3))))
	perlin := trace.NewNoise(0.1, 0.1, 0)
	list.Add(trace.NewSphere(geom.Vec{220, 280, 300}, 80, trace.NewLambert(perlin)))
	from := geom.Vec{478, 278, -600}
	at := geom.Vec{278, 278, 0}
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.Unit{0, 1, 0}, 40, aspect, 0, focus, 0, 1)
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
	light := trace.NewLight(trace.NewUniform(7, 7, 7))
	for i := 0; i < nb; i++ {
		for j := 0; j < nb; j++ {
			min := geom.Vec{-1000 + float64(i)*w, 0, -1000 + float64(j)*w}
			max := geom.Vec{w, 1 + 99*rand.Float64(), w}.Plus(min)
			boxList.Add(trace.NewBox(min, max, ground))
		}
	}
	list.Add(trace.NewBVH(0, 1, boxList.Surfaces()...))
	list.Add(trace.NewRect(geom.Vec{123, 554, 147}, geom.Vec{423, 554, 412}, light))
	center := geom.Vec{400, 400, 200}
	list.Add(trace.NewMovingSphere(center, center.Plus(geom.Vec{30, 0, 0}), 0, 1, 50, trace.NewLambert(trace.NewUniform(0.7, 0.3, 0.1))))
	list.Add(trace.NewSphere(geom.Vec{260, 150, 45}, 50, trace.NewDielectric(1.5)))
	list.Add(trace.NewSphere(geom.Vec{0, 150, 145}, 50, trace.NewMetal(trace.NewUniform(0.8, 0.8, 0.9), 1)))
	boundary := trace.NewSphere(geom.Vec{360, 150, 145}, 70, trace.NewDielectric(1.5))
	list.Add(boundary)
	list.Add(trace.NewVolume(boundary, 0.2, trace.NewIso(trace.NewUniform(0.2, 0.4, 0.9))))
	boundary = trace.NewSphere(geom.Vec{0, 0, 0}, 5000, trace.NewDielectric(1.5))
	list.Add(trace.NewVolume(boundary, 0.0001, trace.NewIso(trace.NewUniform(1, 1, 1))))
	f, err := os.Open("images/earthmap.jpg")
	if err != nil {
		panic(err)
	}
	earth, err := trace.NewImage(f)
	if err != nil {
		panic(err)
	}
	list.Add(trace.NewSphere(geom.Vec{400, 200, 400}, 100, trace.NewLambert(earth)))
	perlin := trace.NewNoise(0.1, 0.1, 0)
	list.Add(trace.NewSphere(geom.Vec{220, 280, 300}, 80, trace.NewLambert(perlin)))
	for j := 0; j < ns; j++ {
		boxList2.Add(trace.NewSphere(geom.Vec{165 * rand.Float64(), 165 * rand.Float64(), 165 * rand.Float64()}, 10, white))
	}
	list.Add(trace.NewTranslate(trace.NewRotateY(trace.NewBVH(0, 1, boxList2.Surfaces()...), 15), geom.Vec{-100, 270, 395}))
	from := geom.Vec{478, 278, -600}
	at := geom.Vec{278, 278, 0}
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.Unit{0, 1, 0}, 40, aspect, 0, focus, 0, 1)
	return cam, list
}

func cornellSmoke(aspect float64) (trace.Camera, *trace.List) {
	green := trace.NewLambert(trace.NewUniform(0.12, 0.45, 0.15))
	red := trace.NewLambert(trace.NewUniform(0.65, 0.05, 0.05))
	light := trace.NewLight(trace.NewUniform(7, 7, 7))
	white := trace.NewLambert(trace.NewUniform(0.73, 0.73, 0.73))
	smoke := trace.NewIso(trace.NewUniform(0, 0, 0))
	fog := trace.NewIso(trace.NewUniform(1, 1, 1))
	b1 := trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.Vec{0, 0, 0}, geom.Vec{165, 165, 165}, white), -18), geom.Vec{130, 0, 65})
	b2 := trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.Vec{0, 0, 0}, geom.Vec{165, 330, 165}, white), 15), geom.Vec{265, 0, 295})
	from := geom.Vec{278, 278, -800}
	at := geom.Vec{278, 278, 0}
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.Unit{0, 1, 0}, 40, aspect, 0, focus, 0, 1)
	return cam, trace.NewList(
		trace.NewFlip(trace.NewRect(geom.Vec{555, 0, 0}, geom.Vec{555, 555, 555}, green)),
		trace.NewRect(geom.Vec{0, 0, 0}, geom.Vec{0, 555, 555}, red),
		trace.NewRect(geom.Vec{113, 554, 127}, geom.Vec{443, 554, 432}, light),
		trace.NewFlip(trace.NewRect(geom.Vec{0, 555, 0}, geom.Vec{555, 555, 555}, white)),
		trace.NewRect(geom.Vec{0, 0, 0}, geom.Vec{555, 0, 555}, white),
		trace.NewFlip(trace.NewRect(geom.Vec{0, 0, 555}, geom.Vec{555, 555, 555}, white)),
		trace.NewVolume(b1, 0.01, fog),
		trace.NewVolume(b2, 0.01, smoke),
	)
}

func cornell(aspect float64) (trace.Camera, *trace.List) {
	green := trace.NewLambert(trace.NewUniform(0.12, 0.45, 0.15))
	red := trace.NewLambert(trace.NewUniform(0.65, 0.05, 0.05))
	light := trace.NewLight(trace.NewUniform(15, 15, 15))
	white := trace.NewLambert(trace.NewUniform(0.73, 0.73, 0.73))
	from := geom.Vec{278, 278, -800}
	at := geom.Vec{278, 278, 0}
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.Unit{0, 1, 0}, 40, aspect, 0, focus, 0, 1)
	return cam, trace.NewList(
		trace.NewFlip(trace.NewRect(geom.Vec{555, 0, 0}, geom.Vec{555, 555, 555}, green)),
		trace.NewRect(geom.Vec{0, 0, 0}, geom.Vec{0, 555, 555}, red),
		trace.NewRect(geom.Vec{213, 554, 227}, geom.Vec{343, 554, 332}, light),
		trace.NewFlip(trace.NewRect(geom.Vec{0, 555, 0}, geom.Vec{555, 555, 555}, white)),
		trace.NewRect(geom.Vec{0, 0, 0}, geom.Vec{555, 0, 555}, white),
		trace.NewFlip(trace.NewRect(geom.Vec{0, 0, 555}, geom.Vec{555, 555, 555}, white)),
		trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.Vec{0, 0, 0}, geom.Vec{165, 165, 165}, white), -18), geom.Vec{130, 0, 65}),
		trace.NewTranslate(trace.NewRotateY(trace.NewBox(geom.Vec{0, 0, 0}, geom.Vec{165, 330, 165}, white), 15), geom.Vec{265, 0, 295}),
	)
}

func simpleLight(aspect float64) (trace.Camera, *trace.List) {
	perlin := trace.NewNoise(4, 1, 2)
	from := geom.Vec{25, 4, 6}
	at := geom.Vec{0, 2, 0}
	focus := 10.0
	cam := trace.NewCamera(from, at, geom.Unit{0, 1, 0}, 20, aspect, 0, focus, 0, 1)
	return cam, trace.NewList(
		trace.NewSphere(geom.Vec{0, -1000, 0}, 1000, trace.NewLambert(perlin)),
		trace.NewSphere(geom.Vec{0, 2, 0}, 2, trace.NewLambert(perlin)),
		trace.NewSphere(geom.Vec{0, 7, 0}, 2, trace.NewLight(trace.NewUniform(4, 4, 4))),
		trace.NewRect(geom.Vec{3, 1, -2}, geom.Vec{5, 3, -2}, trace.NewLight(trace.NewUniform(4, 4, 4))),
	)
}
