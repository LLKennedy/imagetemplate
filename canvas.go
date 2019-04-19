package imagetemplate

import (
	"errors"
	// _ "golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	// _ "image/gif"
	// _ "image/jpeg"
	// _ "image/png"
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

// Circle draws a circle of a specific colour on the canvas
func (canvas *ImageCanvas) Circle(centre image.Point, radius int, colour color.Color) error {
	return errors.New("Not implemented yet")
}

// Text draws text on the canvas
func (canvas *ImageCanvas) Text(start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) error {
	return errors.New("Not implemented yet")
}
