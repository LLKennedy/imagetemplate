// Package scaffold parses JSON data, matches it to known template components and controls rendering of the resultant image.
package scaffold

import (
	"bytes"
	"encoding/json"
	_ "image/jpeg" // jpeg imported for image decoding
	_ "image/png"  // png imported for image decoding
	"io/ioutil"

	_ "github.com/LLKennedy/imagetemplate/v3/components/barcode"   // add barcode component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v3/components/circle"    // add circle component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v3/components/datetime"  // add datetime component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v3/components/image"     // add image component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v3/components/rectangle" // add rectangle component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v3/components/text"      // add text component to registry by default
	"github.com/LLKennedy/imagetemplate/v3/render"

	"golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff" // tiff imported for image decoding
	"golang.org/x/tools/godoc/vfs"
)

// Builder manipulates Canvas objects and outputs to a bitmap.
type Builder interface {
	GetCanvas() render.Canvas
	SetCanvas(newCanvas render.Canvas) Builder
	GetComponents() []render.Component
	SetComponents(components []ToggleableComponent) Builder
	GetNamedPropertiesList() render.NamedProperties
	SetNamedProperties(properties render.NamedProperties) (Builder, error)
	ApplyComponents() (Builder, error)
	LoadComponentsFile(fileName string) (Builder, error)
	LoadComponentsData(fileData []byte) (Builder, error)
	WriteToBMP() ([]byte, error)
}

// Template is the format of the JSON file used as a template for building images. See samples.json for examples, each element in the samples array is a complete and valid template object.
type Template struct {
	// BaseImage is the bottom layer of the canvas, on which to draw everything else.
	BaseImage BaseImage `json:"baseImage"`
	// Components are all other elements to be rendered.
	Components []ComponentTemplate `json:"components"`
}

// BaseImage is the template format of the base image settings.
type BaseImage struct {
	// FileName is the file to load the base image from.
	FileName string `json:"fileName"`
	// Data is the base64-encoded data to load the base image from.
	Data string `json:"data"`
	// BaseColour is the pure colour to use as a base image.
	BaseColour BaseColour `json:"baseColour"`
	// BaseWidth is the width to use for pure colour.
	BaseWidth string `json:"width"`
	// BaseHeight is the height to use for pure colour.
	BaseHeight string `json:"height"`
	// PPI is the pixels per inch to set in the canvas.
	PPI string `json:"ppi"`
}

// BaseColour is the template format of the base colour settings.
type BaseColour struct {
	// Red is the red channel.
	Red string `json:"R"`
	// Green is the green channel.
	Green string `json:"G"`
	// Blue is the blue channel.
	Blue string `json:"B"`
	// Alpha is the alpha channel.
	Alpha string `json:"A"`
}

// ComponentTemplate is a partial unmarshalled Component, with its properties left in raw form to be handled by each known type of Component.
type ComponentTemplate struct {
	// Type is the type of the component, such as Rectangle or Barcode.
	Type string `json:"type"`
	// Conditional is the condition(s) on which the component will render.
	Conditional render.ComponentConditional `json:"conditional"`
	// Properties are the raw, unprocessed JSON data for the component to parse
	Properties json.RawMessage `json:"properties"`
}

// ToggleableComponent is a component with its conditional.
type ToggleableComponent struct {
	// Conditional is the condition(s) on which the component will render.
	Conditional render.ComponentConditional
	// Component is the component to render.
	Component render.Component
}

// ImageBuilder uses golang's native Image package to implement the Builder interface.
type ImageBuilder struct {
	// Canvas is the canvas on which the image is drawn.
	Canvas render.Canvas
	// Components are the components to render.
	Components []ToggleableComponent
	// NamedProperties are the user/application defined variables
	NamedProperties render.NamedProperties
	// fs is the file system
	fs vfs.FileSystem
}

// NewBuilder generates a new ImageBuilder with an internal canvas of the specified width and height, and optionally the specified starting colour. No provided colour will result in defaults for Image.
func NewBuilder(fs vfs.FileSystem) Builder {
	if fs == nil {
		fs = vfs.OS(".")
	}
	return ImageBuilder{fs: fs}
}

// WriteToBMP outputs the contents of the builder to a BMP byte array.
func (builder ImageBuilder) WriteToBMP() ([]byte, error) {
	var buf bytes.Buffer
	err := bmp.Encode(&buf, builder.GetCanvas().GetUnderlyingImage())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// LoadComponentsFile sets the internal Component array based on the contents of the specified JSON file.
func (builder ImageBuilder) LoadComponentsFile(fileName string) (Builder, error) {
	b := builder
	// Load initial data into template object
	file, err := builder.fs.Open(fileName)
	if err != nil {
		return builder, err
	}
	defer file.Close()
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return builder, err
	}
	return b.LoadComponentsData(fileData)
}

// LoadComponentsData sets the internal component array based on the contents of the specified JSON data.
func (builder ImageBuilder) LoadComponentsData(fileData []byte) (Builder, error) {
	b := builder
	var template Template
	err := json.Unmarshal(fileData, &template)
	if err != nil {
		return builder, err
	}

	// Parse background image info
	b, err = b.setBackgroundImage(template)
	if err != nil {
		return builder, err
	}

	// Try each known component type to fit the properties
	b.Components, b.NamedProperties, err = parseComponents(template.Components)
	if err != nil {
		return builder, err
	}

	return b, nil
}

func parseComponents(templates []ComponentTemplate) ([]ToggleableComponent, render.NamedProperties, error) {
	var results []ToggleableComponent
	namedProperties := render.NamedProperties{}
	for _, template := range templates {
		//Handle conditional first
		result := ToggleableComponent{}
		tempProperties := namedProperties
		result.Conditional = template.Conditional
		for key, value := range result.Conditional.GetNamedPropertiesList() {
			tempProperties[key] = value
		}
		newComponent, err := render.Decode(template.Type)
		if err != nil {
			return results, namedProperties, err
		}
		// Get JSON struct to parse into
		shape := newComponent.GetJSONFormat()
		err = json.Unmarshal(template.Properties, shape)
		if err != nil {
			// Invalid JSON
			// FIXME: output more helpful expected vs. actual data here for debugging templates
			return results, namedProperties, err
		}
		// Set real properties from JSON struct
		newComponent, compNamedProps, err := newComponent.VerifyAndSetJSONData(shape)
		if err != nil {
			// Invalid data
			return results, namedProperties, err
		}
		for key, value := range compNamedProps {
			tempProperties[key] = value
		}
		result.Component = newComponent
		results = append(results, result)
		namedProperties = tempProperties
	}
	return results, namedProperties, nil
}

// GetCanvas returns the internal Canvas object.
func (builder ImageBuilder) GetCanvas() render.Canvas {
	if builder.Canvas == nil {
		return render.ImageCanvas{}
	}
	return builder.Canvas
}

// SetCanvas sets the internal Canvas object.
func (builder ImageBuilder) SetCanvas(newCanvas render.Canvas) Builder {
	builder.Canvas = newCanvas
	return builder
}

// GetComponents gets the internal Component array.
func (builder ImageBuilder) GetComponents() []render.Component {
	result := []render.Component{}
	for _, tComponent := range builder.Components {
		valid, err := tComponent.Conditional.Validate()
		if valid && err == nil {
			result = append(result, tComponent.Component)
		}
	}
	return result
}

// SetComponents sets the internal Component array.
func (builder ImageBuilder) SetComponents(components []ToggleableComponent) Builder {
	builder.Components = components
	return builder
}

// GetNamedPropertiesList returns the list of named properties in the builder object.
func (builder ImageBuilder) GetNamedPropertiesList() render.NamedProperties {
	return builder.NamedProperties
}

// SetNamedProperties sets the values of names properties in all components and conditionals in the builder.
func (builder ImageBuilder) SetNamedProperties(properties render.NamedProperties) (Builder, error) {
	b := builder
	b.Components = builder.Components[:]
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
	return b, nil
}

// ApplyComponents iterates over the internal Component array, applying each in turn to the Canvas.
func (builder ImageBuilder) ApplyComponents() (Builder, error) {
	b := builder
	for _, tComponent := range b.Components {
		if tComponent.Conditional.Name == "" {
			var err error
			b.Canvas, err = tComponent.Component.Write(b.GetCanvas())
			if err != nil {
				return builder, err
			}
		} else {
			valid, err := tComponent.Conditional.Validate()
			if err != nil {
				return builder, err
			}
			if valid {
				b.Canvas, err = tComponent.Component.Write(b.GetCanvas())
				if err != nil {
					return builder, err
				}
			}
		}
	}
	return b, nil
}
