package tex

import (
	"image"
	"io"

	_ "image/jpeg" // jpeg support
	_ "image/png"  // png support

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

// Image describes an image-mapped texture
type Image struct {
	Nx, Ny int
	Data   image.Image
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
		Nx:   bounds.Max.X,
		Ny:   bounds.Max.Y,
		Data: im,
	}
	return &i, nil
}

// Map maps a u, v coordinate in 3d space p to a Color on the image texture
func (i *Image) Map(uv, p geom.Vec) Color {
	x := int(uv.X() * float64(i.Nx))
	y := int((1 - uv.Y()) * float64(i.Ny))
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x > i.Nx-1 {
		x = i.Nx - 1
	}
	if y > i.Ny-1 {
		y = i.Ny - 1
	}
	c := i.Data.At(x, y)
	r, g, b, _ := c.RGBA()
	return NewColor(float64(r)/65535, float64(g)/65535, float64(b)/65535)
}
