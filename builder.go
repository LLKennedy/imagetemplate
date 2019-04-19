package imagetemplate

import (
	"bytes"
	// "errors"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

// Builder manipulates Canvas objects and outputs to a bitmap
type Builder interface {
	Canvas
	WriteToBMP() ([]byte, error)
}

// ImageBuilder uses golang's native Image package to implement the Builder interface
type ImageBuilder struct {
	Canvas Canvas
}

// NewBuilder generates a new ImageBuilder with an internal canvas of the specified width and height, and optionally the specified starting colour. No provided colour will result in defaults for Image.
func NewBuilder(width, height int, startingColour color.Color) (*ImageBuilder, error) {
	newCanvas, err := NewCanvas(width, height)
	if err != nil {
		return nil, err
	}
	if startingColour != nil {
		err = newCanvas.Rectangle(image.Point{X: 0, Y: 0}, width, height, startingColour)
		if err != nil {
			return nil, err
		}
	}
	return &ImageBuilder{Canvas: newCanvas}, nil
}

func (builder *ImageBuilder) WriteToBMP() ([]byte, error) {
	var buf bytes.Buffer
	err := bmp.Encode(&buf, builder.Canvas.GetUnderlyingImage())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//Canvas methods passed through to the internal canvas

// SetUnderlyingImage sets the underlying image in the canvas
func (builder *ImageBuilder) SetUnderlyingImage(newImage image.Image) {
	builder.Canvas.SetUnderlyingImage(newImage)
}

// GetUnderlyingImage gets the underlying image from the canvas
func (builder *ImageBuilder) GetUnderlyingImage() image.Image {
	return builder.Canvas.GetUnderlyingImage()
}

// Rectangle draws a rectangle on the canvas
func (builder *ImageBuilder) Rectangle(topLeft image.Point, width, height int, colour color.Color) error {
	return builder.Canvas.Rectangle(topLeft, width, height, colour)
}

// Circle draws a circle on the canvas
func (builder *ImageBuilder) Circle(centre image.Point, radius int, colour color.Color) error {
	return builder.Canvas.Circle(centre, radius, colour)
}

// Text draws text on the canvas
func (builder *ImageBuilder) Text(start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) error {
	return builder.Canvas.Text(start, typeFace, colour, fontSize, maxWidth, maxLines)
}
