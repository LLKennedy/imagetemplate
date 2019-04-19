package imagetemplate

import (
	"bytes"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
)

// Builder manipulates Canvas objects and outputs to a bitmap
type Builder interface {
	GetCanvas() Canvas
	SetCanvas(newCanvas Canvas)
	GetComponents() []Component
	SetComponents([]Component)
	WriteToBMP() ([]byte, error)
}

// ImageBuilder uses golang's native Image package to implement the Builder interface
type ImageBuilder struct {
	Canvas     Canvas
	Components []Component
}

// NewBuilder generates a new ImageBuilder with an internal canvas of the specified width and height, and optionally the specified starting colour. No provided colour will result in defaults for Image.
func NewBuilder(canvas Canvas, startingColour color.Color) (*ImageBuilder, error) {
	if startingColour != nil {
		err := canvas.Rectangle(image.Point{}, canvas.GetWidth(), canvas.GetHeight(), startingColour)
		if err != nil {
			return nil, err
		}
	}
	return &ImageBuilder{Canvas: canvas}, nil
}

func (builder *ImageBuilder) WriteToBMP() ([]byte, error) {
	var buf bytes.Buffer
	err := bmp.Encode(&buf, builder.Canvas.GetUnderlyingImage())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (builder *ImageBuilder) GetCanvas() Canvas {
	return builder.Canvas
}

func (builder *ImageBuilder) SetCanvas(newCanvas Canvas) {
	builder.Canvas = newCanvas
}

func (builder *ImageBuilder) GetComponents() []Component {
	return builder.Components
}

func (builder *ImageBuilder) SetComponents(components []Component) {
	builder.Components = components
}
