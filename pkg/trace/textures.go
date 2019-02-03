package trace

import (
	"image"
	"io"
	"math"

	_ "image/jpeg" // jpeg format for Image
	_ "image/png"  // png format for Image

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Mapper maps a u, v coordinate in 3d space p to a Color
type Mapper interface {
	Map(uv, p geom.Vec) Color
}

// Checker represents an alternating checkered pattern of two sub-textures
type Checker struct {
	size      float64
	odd, even Mapper
}

// NewChecker returns a new checkered texture rendering sub-textures t0 and t1 in size squares.
func NewChecker(size float64, t0, t1 Mapper) Checker {
	return Checker{
		size: size,
		odd:  t0,
		even: t1,
	}
}

// Map maps a u, v coordinate in 3d space p to a Color
func (c Checker) Map(uv, p geom.Vec) Color {
	sines := math.Sin(c.size*p.X()) * math.Sin(c.size*p.Y()) * math.Sin(c.size*p.Z())
	if sines < 0 {
		return c.odd.Map(uv, p)
	}
	return c.even.Map(uv, p)
}

// Uniform represents a single, Uniform colored texture
type Uniform struct {
	c Color
}

// NewUniform returns a new Uniform texture
func NewUniform(r, g, b float64) Uniform {
	return Uniform{c: Color{r, g, b}}
}

// Map maps a u, v coordinate in 3d space p to a Color
func (u Uniform) Map(uv, p geom.Vec) Color {
	return u.c
}

type Bright struct {
	src   Mapper
	scale float64
}

func NewBright(src Mapper, scale float64) Bright {
	return Bright{src: src, scale: scale}
}

func (b Bright) Map(uv, p geom.Vec) Color {
	return b.src.Map(uv, p).Scaled(b.scale)
}

// Image describes an image-mapped texture
type Image struct {
	width, height int
	data          image.Image
}

// NewImage creates a new image texture by reading a png or jpeg from rc
func NewImage(rc io.ReadCloser) (*Image, error) {
	defer rc.Close()
	im, _, err := image.Decode(rc)
	if err != nil {
		return nil, err
	}
	bounds := im.Bounds()
	i := Image{
		width:  bounds.Max.X,
		height: bounds.Max.Y,
		data:   im,
	}
	return &i, nil
}

// Map maps a u, v coordinate in 3d space p to a Color on the image texture
func (i *Image) Map(uv, p geom.Vec) Color {
	x := int(uv.X() * float64(i.width))
	y := int((1 - uv.Y()) * float64(i.height))
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x > i.width-1 {
		x = i.width - 1
	}
	if y > i.height-1 {
		y = i.height - 1
	}
	c := i.data.At(x, y)
	r, g, b, _ := c.RGBA()
	return Color{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
}

// Noise represents a Perlin noise texture
type Noise struct {
	scale0, scale1 float64
	axis           int
}

// NewNoise returns a new noise texture with frequency scaled by scale
func NewNoise(scale0, scale1 float64, axis int) Noise {
	return Noise{scale0: scale0, scale1: scale1, axis: axis}
}

// Map maps a u, v coordinate in 3d space p to a Color
func (n Noise) Map(uv, p geom.Vec) Color {
	bright := 0.5 * (1 + math.Sin(n.scale0*p[n.axis]+10*turb(p.Scaled(n.scale1), 7)))
	return white.Scaled(bright)
}
