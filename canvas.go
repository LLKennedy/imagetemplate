package imagetemplate

import (
	"errors"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

// Canvas holds the image struct and associated properties
type Canvas interface {
	SetUnderlyingImage(newImage image.Image)
	GetUnderlyingImage() image.Image
	Rectangle(topLeft image.Point, width, height int, colour color.Color) error
	Circle(centre image.Point, radius int, colour color.Color) error
	Text(start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) error
}

// ImageCanvas uses golang's native Image package to implement the Canvas interface
type ImageCanvas struct {
	Image image.Image
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
func (canvas *ImageCanvas) SetUnderlyingImage(newImage image.Image) {
	canvas.Image = newImage
}

// GetUnderlyingImage gets the internal Image property
func (canvas *ImageCanvas) GetUnderlyingImage() image.Image {
	return canvas.Image
}

// Rectangle draws a rectangle of a specific colour on the canvas
func (canvas *ImageCanvas) Rectangle(topLeft image.Point, width, height int, colour color.Color) error {
	return errors.New("Not implemented yet")
}

// Circle draws a circle of a specific colour on the canvas
func (canvas *ImageCanvas) Circle(centre image.Point, radius int, colour color.Color) error {
	return errors.New("Not implemented yet")
}

// Text draws text on the canvas
func (canvas *ImageCanvas) Text(start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) error {
	return errors.New("Not implemented yet")
}
