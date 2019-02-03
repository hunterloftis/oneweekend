/*
Package trace renders 3D scenes into 2D images by CPU ray-tracing.

Basics

A render requires just four things: a window, a surface, a camera, and materials.

The window defines the render's dimensions and outputs the image.
The surface is the scene that rays are traced against.
The camera determines your point of view and projects rays at the surface.

Most renders will need at least two materials -
one to emit light, and one to scatter or reflect light.
*/
package trace

import (
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// bias is a small number that is used in place of zero to avoid
// issues with numerical precision.
const bias = 0.001

// Surface is a bounded object in 3D space that can be hit by a Ray.
type Surface interface {
	Hit(r Ray, dMin, dMax float64, rnd *rand.Rand) *Hit
	Bounds(t0, t1 float64) (bounds *AABB)
}

// Material determines how light scatters and emits when it hits a Surface.
type Material interface {
	Scatter(in, norm geom.Unit, uv, p geom.Vec, rnd *rand.Rand) (out geom.Unit, attenuate Color, ok bool)
	Emit(uv, p geom.Vec) Color
}

// Hit records the details of a Ray->Surface intersection.
type Hit struct {
	Dist float64
	Norm geom.Unit
	UV   geom.Vec
	Pt   geom.Vec
	Mat  Material
}
