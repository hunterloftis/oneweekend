package trace

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Window gathers the results of ray traces on a width x height grid.
type Window struct {
	width, height int
}

// NewWindow creates a new Window with specific dimensions
func NewWindow(width, height int) *Window {
	return &Window{width: width, height: height}
}

func (wi *Window) Aspect() float64 {
	return float64(wi.width) / float64(wi.height)
}

type result struct {
	row    int
	pixels string
}

// WritePPM traces each pixel in the Window and writes the results to w in PPM format
func (wi Window) WritePPM(w io.Writer, cam *Camera, s Surface, samples int) error {
	if _, err := fmt.Fprintln(w, "P3"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, wi.width, wi.height); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "255"); err != nil {
		return err
	}

	worker := func(jobs <-chan int, results chan<- result, rnd *rand.Rand) {
		for y := range jobs {
			res := result{row: y, pixels: ""}
			for x := 0; x < wi.width; x++ {
				c := black
				for n := 0; n < samples; n++ {
					u := (float64(x) + rnd.Float64()) / float64(wi.width)
					v := 1 - (float64(y)+rnd.Float64())/float64(wi.height)
					r := cam.Ray(u, v, rnd)
					c = c.Plus(color(r, s, 0, rnd))
				}
				c = c.Scaled(1 / float64(samples)).Gamma(2)
				ir := int(math.Min(255, 255.99*c.R()))
				ig := int(math.Min(255, 255.99*c.G()))
				ib := int(math.Min(255, 255.99*c.B()))
				res.pixels += fmt.Sprintln(ir, ig, ib)
			}
			results <- res
		}
	}

	workers := runtime.NumCPU() + 1
	jobs := make(chan int, wi.height)
	results := make(chan result, workers+1)
	pending := make(map[int]string, 0)
	cursor := 0

	for w := 0; w < workers; w++ {
		go worker(jobs, results, rand.New(rand.NewSource(time.Now().Unix())))
	}
	for y := 0; y < wi.height; y++ {
		jobs <- y
	}
	close(jobs)
	for y := 0; y < wi.height; y++ {
		r := <-results
		pending[r.row] = r.pixels
		for pending[cursor] != "" {
			fmt.Fprint(w, pending[cursor])
			delete(pending, cursor)
			cursor++
		}
	}

	return nil
}

func color(r Ray, h Surface, depth int, rnd *rand.Rand) Color {
	if depth > 50 {
		return black
	}
	if hit := h.Hit(r, bias, math.MaxFloat64, rnd); hit != nil {
		emit := hit.Mat.Emit(hit.UV, hit.Pt)
		out, attenuate, ok := hit.Mat.Scatter(r.Dir, hit.Norm, hit.UV, hit.Pt, rnd)
		if !ok {
			return emit
		}
		indirect := color(NewRay(hit.Pt, out, r.T), h, depth+1, rnd).Times(attenuate)
		return emit.Plus(indirect)
	}
	return black
}

// Camera originates Rays.
type Camera struct {
	lowerLeft    geom.Vec
	horizontal   geom.Vec
	vertical     geom.Vec
	origin       geom.Vec
	u, v, w      geom.Unit
	lensRadius   float64
	time0, time1 float64
}

// NewCamera creates a new Camera
// TODO: this argument list is getting pretty ridiculous
func NewCamera(lookFrom, lookAt geom.Vec, vup geom.Unit, vfov, aspect, aperture, focus, t0, t1 float64) *Camera {
	theta := vfov * math.Pi / 180
	halfH := math.Tan(theta / 2)
	halfW := aspect * halfH

	c := Camera{}
	c.w = lookFrom.Minus(lookAt).Unit()
	c.u = geom.Vec(vup).Cross(geom.Vec(c.w)).Unit()
	c.v = geom.Vec(c.w).Cross(geom.Vec(c.u)).Unit()

	width := c.u.Scaled(halfW * focus)
	height := c.v.Scaled(halfH * focus)
	dist := c.w.Scaled(focus)

	c.time0 = t0
	c.time1 = t1
	c.lensRadius = aperture / 2
	c.origin = lookFrom
	c.lowerLeft = c.origin.Minus(width).Minus(height).Minus(dist)
	c.horizontal = width.Scaled(2)
	c.vertical = height.Scaled(2)
	return &c
}

// Ray returns a Ray passing through a given s, t coordinate.
func (c *Camera) Ray(s, t float64, rnd *rand.Rand) Ray {
	rd := geom.RandVecInDisk(rnd).Scaled(c.lensRadius)
	offset := c.u.Scaled(rd.X()).Plus(c.v.Scaled(rd.Y()))
	source := c.origin.Plus(offset)
	dest := c.lowerLeft.Plus(c.horizontal.Scaled(s).Plus(c.vertical.Scaled(t)))
	time := c.time0 + (c.time1-c.time0)*rnd.Float64()
	return NewRay(source, dest.Minus(source).Unit(), time)
}
