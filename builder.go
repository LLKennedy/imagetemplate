// Package imagetemplate defines a template for drawing custom images from pre-defined components, and provides to tools to load and implement that template.
package imagetemplate

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
)

// Builder manipulates Canvas objects and outputs to a bitmap
type Builder interface {
	GetCanvas() Canvas
	SetCanvas(newCanvas Canvas)
	GetComponents() []Component
	SetComponents(components []Component)
	GetNamedProperties() NamedProperties
	SetNamedProperties(properties NamedProperties)
	ApplyComponents() error
	LoadComponentsFile(fileName string) error
	WriteToBMP() ([]byte, error)
}

// Template is the format of the JSON file used as a template for building images. See samples.json for examples, each element in the samples array is a complete and valid template object.
type Template struct {
	BaseImage struct {
		FileName string `json:"fileName"`
		Data     string `json:"data"`
		FileType string `json:"fileType"`
	} `json:"baseImage"`
	Components []ComponentsElement `json:"components"`
}

// ComponentsElement represents a partial unmarshalled Component, with its properties left in raw form to be handled by each known type of Component.
type ComponentsElement struct {
	Type        string               `json:"type"`
	Conditional ComponentConditional `json:"conditional"`
	Properties  json.RawMessage      `json:"properties"`
}

// ImageBuilder uses golang's native Image package to implement the Builder interface
type ImageBuilder struct {
	Canvas          Canvas
	Components      []Component
	NamedProperties NamedProperties
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

// WriteToBMP outputs the contents of the builder to a BMP byte array
func (builder *ImageBuilder) WriteToBMP() ([]byte, error) {
	var buf bytes.Buffer
	err := bmp.Encode(&buf, builder.Canvas.GetUnderlyingImage())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetCanvas returns the internal Canvas object
func (builder *ImageBuilder) GetCanvas() Canvas {
	return builder.Canvas
}

// SetCanvas sets the internal Canvas object
func (builder *ImageBuilder) SetCanvas(newCanvas Canvas) {
	builder.Canvas = newCanvas
}

// GetComponents gets the internal Component array
func (builder *ImageBuilder) GetComponents() []Component {
	return builder.Components
}

// SetComponents sets the internal Component array
func (builder *ImageBuilder) SetComponents(components []Component) {
	builder.Components = components
}

// GetNamedProperties returns the list of named properties in the builder object
func (builder *ImageBuilder) GetNamedProperties() NamedProperties {
	return builder.NamedProperties
}

// SetNamedProperties sets the list of named properties in the builder object
func (builder *ImageBuilder) SetNamedProperties(properties NamedProperties) {
	builder.NamedProperties = properties
}

// ApplyComponents iterates over the internal Component array, applying each in turn to the Canvas
func (builder *ImageBuilder) ApplyComponents() error {
	return errors.New("Not implemented yet")
}

// LoadComponentsFile sets the internal Component array based on the contents of the specified JSON file
func (builder *ImageBuilder) LoadComponentsFile(fileName string) error {
	return errors.New("Not implemented yet")
}
