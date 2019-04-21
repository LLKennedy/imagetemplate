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
	SetCanvas(newCanvas Canvas) Builder
	GetComponents() []Component
	SetComponents(components []ToggleableComponent) Builder
	GetNamedPropertiesList() NamedProperties
	SetNamedProperties(properties NamedProperties) (Builder, error)
	ApplyComponents() (Builder, error)
	LoadComponentsFile(fileName string) (Builder, error)
	WriteToBMP() ([]byte, error)
}

// Template is the format of the JSON file used as a template for building images. See samples.json for examples, each element in the samples array is a complete and valid template object.
type Template struct {
	BaseImage struct {
		FileName string `json:"fileName"`
		Data     string `json:"data"`
		FileType string `json:"fileType"`
	} `json:"baseImage"`
	Components []ComponentTemplate `json:"components"`
}

// ComponentsElement represents a partial unmarshalled Component, with its properties left in raw form to be handled by each known type of Component.
type ComponentTemplate struct {
	Type        string               `json:"type"`
	Conditional ComponentConditional `json:"conditional"`
	Properties  json.RawMessage      `json:"properties"`
}

type ToggleableComponent struct {
	Conditional ComponentConditional
	Component   Component
}

// ImageBuilder uses golang's native Image package to implement the Builder interface
type ImageBuilder struct {
	Canvas          Canvas
	Components      []ToggleableComponent
	NamedProperties NamedProperties
}

// NewBuilder generates a new ImageBuilder with an internal canvas of the specified width and height, and optionally the specified starting colour. No provided colour will result in defaults for Image.
func NewBuilder(canvas Canvas, startingColour color.Color) (ImageBuilder, error) {
	if startingColour != nil {
		var err error
		canvas, err = canvas.Rectangle(image.Point{}, canvas.GetWidth(), canvas.GetHeight(), startingColour)
		if err != nil {
			return ImageBuilder{}, err
		}
	}
	return ImageBuilder{Canvas: canvas}, nil
}

// WriteToBMP outputs the contents of the builder to a BMP byte array
func (builder ImageBuilder) WriteToBMP() ([]byte, error) {
	var buf bytes.Buffer
	err := bmp.Encode(&buf, builder.Canvas.GetUnderlyingImage())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetCanvas returns the internal Canvas object
func (builder ImageBuilder) GetCanvas() Canvas {
	return builder.Canvas
}

// SetCanvas sets the internal Canvas object
func (builder ImageBuilder) SetCanvas(newCanvas Canvas) Builder {
	builder.Canvas = newCanvas
	return builder
}

// GetComponents gets the internal Component array
func (builder ImageBuilder) GetComponents() []Component {
	result := []Component{}
	for _, tComponent := range builder.Components {
		valid, err := tComponent.Conditional.Validate()
		if err != nil {
			continue
		}
		if valid {
			result = append(result, tComponent.Component)
		}
	}
	return result
}

// SetComponents sets the internal Component array
func (builder ImageBuilder) SetComponents(components []ToggleableComponent) Builder {
	builder.Components = components
	return builder
}

// GetNamedProperties returns the list of named properties in the builder object
func (builder ImageBuilder) GetNamedPropertiesList() NamedProperties {
	return builder.NamedProperties
}

// SetNamedProperties sets the values of names properties in all components and conditionals in the builder
func (builder ImageBuilder) SetNamedProperties(properties NamedProperties) (Builder, error) {
	b := builder
	for tIndex, tComponent := range b.Components {
		var err error
		tComponent.Component, err = tComponent.Component.SetNamedProperties(properties)
		if err != nil {
			return builder, err
		}
		for key, value := range properties {
			tComponent.Conditional, err = tComponent.Conditional.SetValue(key, value)
			if err != nil {
				return builder, err
			}
		}
		b.Components[tIndex] = tComponent
	}
	return builder, nil
}

// ApplyComponents iterates over the internal Component array, applying each in turn to the Canvas
func (builder ImageBuilder) ApplyComponents() (Builder, error) {
	b := builder
	for _, tComponent := range b.Components {
		valid, err := tComponent.Conditional.Validate()
		if err != nil {
			return builder, err
		}
		if valid {
			b.Canvas, err = tComponent.Component.Write(b.Canvas)
			if err != nil {
				return builder, err
			}
		}
	}
	return b, nil
}

// LoadComponentsFile sets the internal Component array based on the contents of the specified JSON file
func (builder ImageBuilder) LoadComponentsFile(fileName string) (Builder, error) {
	return builder, errors.New("Not implemented yet")
}
