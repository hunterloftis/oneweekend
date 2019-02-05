package trace

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

type result struct {
	row    int
	pixels string
}

// Window gathers the results of ray traces in a width x height grid.
type Window struct {
	width, height int
}

// NewWindow creates a new Window with dimensions width and height.
func NewWindow(width, height int) *Window {
	return &Window{width: width, height: height}
}

// WritePPM traces each pixel in the Window and writes the results to w in PPM format.
func (wi *Window) WritePPM(w io.Writer, cam *Camera, s Surface, samples int) error {
	if _, err := fmt.Fprint(w, "P3\n", wi.width, wi.height, "\n255\n"); err != nil {
		return err
	}

	// create worker goroutines and one job per image row.
	aspect := float64(wi.width) / float64(wi.height)
	nw := runtime.NumCPU() + 1
	jobs := make(chan int, wi.height)
	results := make(chan result, nw*2)
	worker := func(rnd *rand.Rand) {
		for y := range jobs {
			var px strings.Builder
			for x := 0; x < wi.width; x++ {
				c := black
				for n := 0; n < samples; n++ {
					u := (float64(x) + rnd.Float64()) / float64(wi.width)
					v := (float64(y) + rnd.Float64()) / float64(wi.height)
					r := cam.Ray(u, v, aspect, rnd)
					c = color(r, s, 0, rnd).Plus(c)
				}
				c = c.Scaled(1 / float64(samples)).Gamma(2)
				r, g, b := c.RGBInt()
				fmt.Fprintln(&px, r, g, b)
			}
			results <- result{row: y, pixels: px.String()}
		}
	}
	for w := 0; w < nw; w++ {
		go worker(rand.New(rand.NewSource(time.Now().Unix())))
	}
	for y := 0; y < wi.height; y++ {
		jobs <- y
	}
	close(jobs)

	// buffer results and write them in order.
	cursor := 0
	pending := make(map[int]string, 0)
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

// color recursively traces rays into s, starting with r.
// It returns a color which is not deterministic,
// but is just one random path that r could take.
func color(r Ray, s Surface, depth int, rnd *rand.Rand) Color {
	if depth >= 50 {
		return black
	}
	hit := s.Hit(r, bias, math.MaxFloat64, rnd)
	if hit == nil {
		return black
	}
	emit := hit.Mat.Emit(hit.UV, hit.Pt)
	out, attenuate, ok := hit.Mat.Scatter(r.Dir, hit.Norm, hit.UV, hit.Pt, rnd)
	if !ok {
		return emit
	}
	indirect := color(NewRay(hit.Pt, out, r.T), s, depth+1, rnd).Times(attenuate)
	return emit.Plus(indirect)
}

// Camera generates rays from a given point of view.
type Camera struct {
	vertical     geom.Vec
	origin       geom.Vec
	u, v, w      geom.Unit
	lensRadius   float64
	time0, time1 float64
	halfH        float64
	focus        float64
	dist         geom.Vec
	height       geom.Vec
}

// NewCamera returns a new Camera.
// TODO: this argument list is getting pretty ridiculous
func NewCamera(lookFrom, lookAt geom.Vec, vup geom.Unit, vfov, aperture, focus, t0, t1 float64) *Camera {
	theta := vfov * math.Pi / 180
	c := Camera{
		focus:      focus,
		origin:     lookFrom,
		time0:      t0,
		time1:      t1,
		lensRadius: aperture / 2,
		halfH:      math.Tan(theta / 2),
		w:          lookFrom.Minus(lookAt).Unit(),
	}
	c.u = geom.Vec(vup).Cross(geom.Vec(c.w)).Unit()
	c.v = geom.Vec(c.w).Cross(geom.Vec(c.u)).Unit()
	c.height = c.v.Scaled(c.halfH * c.focus)
	c.dist = c.w.Scaled(c.focus)
	c.vertical = c.height.Scaled(2)
	return &c
}

// Ray returns a ray from this camera passing through a given s, t coordinate.
// aspect is the ratio of width/height.
func (c *Camera) Ray(s, t, aspect float64, rnd *rand.Rand) Ray {
	time := c.time0 + (c.time1-c.time0)*rnd.Float64()

	// start the ray at a random point on the disc (the aperture).
	rd := geom.RandVecInDisk(rnd).Scaled(c.lensRadius)
	offset := c.u.Scaled(rd.X()).Plus(c.v.Scaled(rd.Y()))
	source := c.origin.Plus(offset)

	// stretch the horizontal vectors by the aspect ratio.
	width := c.u.Scaled(c.halfH * c.focus * aspect)
	upperLeft := c.origin.Minus(width).Plus(c.height).Minus(c.dist)
	horizontal := width.Scaled(2)

	// create the ray by starting in the corner and adding horizontal and vertical components.
	dest := upperLeft.Plus(horizontal.Scaled(s).Minus(c.vertical.Scaled(t)))
	return NewRay(source, dest.Minus(source).Unit(), time)
}
