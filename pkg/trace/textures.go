package trace

import (
	"image"
	"io"
	"math"

	_ "image/jpeg" // jpeg format for Image
	_ "image/png"  // png format for Image

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Mapper maps a uv coordinate at point p to a Color.
type Mapper interface {
	Map(uv, p geom.Vec) Color
}

// Checker is an alternating checkered pattern of two sub-textures.
type Checker struct {
	size      float64
	odd, even Mapper
}

// NewChecker returns a new checkered texture rendering sub-textures t0 and t1 in size cubes.
func NewChecker(size float64, t0, t1 Mapper) *Checker {
	return &Checker{
		size: size,
		odd:  t0,
		even: t1,
	}
}

// Map maps a uv coordinate at point p to a Color in cubes based on the checker's size.
func (c *Checker) Map(uv, p geom.Vec) Color {
	sines := math.Sin(c.size*p.X()) * math.Sin(c.size*p.Y()) * math.Sin(c.size*p.Z())
	if sines < 0 {
		return c.odd.Map(uv, p)
	}
	return c.even.Map(uv, p)
}

// Uniform is a uniformly colored texture.
type Uniform struct {
	c Color
}

// NewUniform returns a new uniform texture with the given RGB color.
func NewUniform(r, g, b float64) *Uniform {
	return &Uniform{c: Color{r, g, b}}
}

// Map returns the texture's color.
func (u *Uniform) Map(uv, p geom.Vec) Color {
	return u.c
}

// Bright is a texture that modifies the brightness of a source texture.
type Bright struct {
	src   Mapper
	scale float64
}

// NewBright returns a new texture that scales the brightness of src.
func NewBright(src Mapper, scale float64) *Bright {
	return &Bright{src: src, scale: scale}
}

// Map maps a uv coordinate at point p to the source texture's scaled color.
func (b *Bright) Map(uv, p geom.Vec) Color {
	return b.src.Map(uv, p).Scaled(b.scale)
}

// Image is an image-mapped texture.
type Image struct {
	width, height int
	data          image.Image
}

// NewImage creates a new image texture by reading a png or jpeg from rc.
// If it encounters an error while reading, it returns a nil Image.
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

// Map maps a uv coordinate to a color on the image.
func (i *Image) Map(uv, _ geom.Vec) Color {
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

// Noise is a Perlin noise texture.
type Noise struct {
	scale0, scale1 float64
	axis           int
}

// NewNoise returns a new noise texture.
// scale0 is the overall scale.
// scale1 is the turbulence scale.
// axis is the axis (0, 1, or 2) by which the perlin noise is modulated.
func NewNoise(scale0, scale1 float64, axis int) *Noise {
	return &Noise{scale0: scale0, scale1: scale1, axis: axis}
}

// Map maps point p to a color via perlin noise.
func (n *Noise) Map(_, p geom.Vec) Color {
	bright := 0.5 * (1 + math.Sin(n.scale0*p[n.axis]+10*turb(p.Scaled(n.scale1), 7)))
	return white.Scaled(bright)
}
