package imagetemplate

import (
	"errors"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
)

// Canvas holds the image struct and associated properties
type Canvas interface {
	SetUnderlyingImage(newImage draw.Image)
	GetUnderlyingImage() image.Image
	GetWidth() int
	GetHeight() int
	Rectangle(topLeft image.Point, width, height int, colour color.Color) error
	Circle(centre image.Point, radius int, colour color.Color) error
	Text(start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) error
	SubImage(start image.Point, subImage image.Image) error
}

// ImageCanvas uses golang's native Image package to implement the Canvas interface
type ImageCanvas struct {
	Image draw.Image
}

// NewCanvas generates a new canvas of the given width and height
func NewCanvas(width, height int) (*ImageCanvas, error) {
	if width <= 0 && height <= 0 {
		return nil, errors.New("Invalid width and height")
	} else if width <= 0 {
		return nil, errors.New("Invalid width")
	} else if height <= 0 {
		return nil, errors.New("Invalid height")
	}
	return &ImageCanvas{
		Image: image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: width, Y: height},
		}),
	}, nil
}

// SetUnderlyingImage sets the internal Image property to the given object
func (canvas *ImageCanvas) SetUnderlyingImage(newImage draw.Image) {
	canvas.Image = newImage
}

// GetUnderlyingImage gets the internal Image property
func (canvas *ImageCanvas) GetUnderlyingImage() image.Image {
	return canvas.Image
}

// GetWidth returns the width of the underlying Image
func (canvas *ImageCanvas) GetWidth() int {
	return canvas.Image.Bounds().Size().X
}

// GetWidth returns the width of the underlying Image
func (canvas *ImageCanvas) GetHeight() int {
	return canvas.Image.Bounds().Size().Y
}

// Rectangle draws a rectangle of a specific colour on the canvas
func (canvas *ImageCanvas) Rectangle(topLeft image.Point, width, height int, colour color.Color) error {
	colourPlane := image.Uniform{C: colour}
	if width <= 0 && height <= 0 {
		return errors.New("Invalid width and height")
	} else if width <= 0 {
		return errors.New("Invalid width")
	} else if height <= 0 {
		return errors.New("Invalid height")
	}
	rect := image.Rectangle{
		Min: image.Point{
			X: topLeft.X,
			Y: topLeft.Y,
		},
		Max: image.Point{
			X: topLeft.X + width,
			Y: topLeft.Y + height,
		},
	}
	draw.Draw(canvas.Image, rect, &colourPlane, topLeft, draw.Over)
	return nil
}

// CircleAlignment dicates where precisely a circle is centred
type CircleAlignment int

// Circle draws a circle of a specific colour on the canvas
func (canvas *ImageCanvas) Circle(centre image.Point, radius int, colour color.Color) error {
	colourPlane := image.Uniform{C: colour}
	mask := &circle{p: centre, r: radius}
	draw.DrawMask(canvas.Image, mask.Bounds(), &colourPlane, image.ZP, mask, mask.Bounds().Min, draw.Over)
	return nil
}

// Text draws text on the canvas
func (canvas *ImageCanvas) Text(start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) error {
	return errors.New("Not implemented yet")
}

// SubImage draws another image on the canvas
func (canvas *ImageCanvas) SubImage(start image.Point, subImage image.Image) error {
	return errors.New("Not implemented yet")
}

// Steal the circle example code from https://blog.golang.org/go-imagedraw-package
type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
