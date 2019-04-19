package imagetemplate

import (
	"errors"
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

func NewBuilder(width, height int) *ImageBuilder {
	return &ImageBuilder{Canvas: NewCanvas(width, height)}
}

func (builder *ImageBuilder) WriteToBMP() ([]byte, error) {
	return nil, errors.New("Not implemented yet")
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
