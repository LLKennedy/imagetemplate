// Package imagetemplate defines a template for drawing custom images from pre-defined components, and provides to tools to load and implement that template.
package imagetemplate

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg" // jpeg imported for image decoding
	_ "image/png"  // png imported for image decoding
	"strconv"
	"strings"

	_ "github.com/LLKennedy/imagetemplate/v2/components/barcode"   // add barcode component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v2/components/circle"    // add circle component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v2/components/datetime"  // add datetime component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v2/components/image"     // add image component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v2/components/rectangle" // add rectangle component to registry by default
	_ "github.com/LLKennedy/imagetemplate/v2/components/text"      // add text component to registry by default
	fs "github.com/LLKennedy/imagetemplate/v2/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/v2/render"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff" // tiff imported for image decoding
)

// Builder manipulates Canvas objects and outputs to a bitmap
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
	BaseImage  BaseImage           `json:"baseImage"`
	Components []ComponentTemplate `json:"components"`
}

// BaseImage is the template format of the base image settings
type BaseImage struct {
	FileName   string     `json:"fileName"`
	Data       string     `json:"data"`
	BaseColour BaseColour `json:"baseColour"`
	BaseWidth  string     `json:"width"`
	BaseHeight string     `json:"height"`
}

// BaseColour is the template format of the base colour settings
type BaseColour struct {
	Red   string `json:"R"`
	Green string `json:"G"`
	Blue  string `json:"B"`
	Alpha string `json:"A"`
}

// ComponentTemplate is a partial unmarshalled Component, with its properties left in raw form to be handled by each known type of Component.
type ComponentTemplate struct {
	Type        string                      `json:"type"`
	Conditional render.ComponentConditional `json:"conditional"`
	Properties  json.RawMessage             `json:"properties"`
}

// ToggleableComponent is a component with its conditional
type ToggleableComponent struct {
	Conditional render.ComponentConditional
	Component   render.Component
}

// ImageBuilder uses golang's native Image package to implement the Builder interface
type ImageBuilder struct {
	Canvas          render.Canvas
	Components      []ToggleableComponent
	NamedProperties render.NamedProperties
	reader          fs.FileReader
}

// NewBuilder generates a new ImageBuilder with an internal canvas of the specified width and height, and optionally the specified starting colour. No provided colour will result in defaults for Image.
func NewBuilder() ImageBuilder {
	return ImageBuilder{reader: fs.IoutilFileReader{}}
}

// WriteToBMP outputs the contents of the builder to a BMP byte array
func (builder ImageBuilder) WriteToBMP() ([]byte, error) {
	var buf bytes.Buffer
	err := bmp.Encode(&buf, builder.GetCanvas().GetUnderlyingImage())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// LoadComponentsFile sets the internal Component array based on the contents of the specified JSON file
func (builder ImageBuilder) LoadComponentsFile(fileName string) (Builder, error) {
	b := builder
	// Load initial data into template object
	fileData, err := builder.reader.ReadFile(fileName)
	if err != nil {
		return builder, err
	}
	return b.LoadComponentsData(fileData)
}

// LoadComponentsData sets the internal component array based on the contents of the specified JSON data
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

func (builder ImageBuilder) setBackgroundImage(template Template) (ImageBuilder, error) {
	b := builder
	// Check the state of the optional and required properties
	dataSet := template.BaseImage.Data != ""
	fileSet := template.BaseImage.FileName != ""
	baseColourSet := template.BaseImage.BaseWidth != "" && template.BaseImage.BaseHeight != "" && (template.BaseImage.BaseColour.Red != "" || template.BaseImage.BaseColour.Green != "" || template.BaseImage.BaseColour.Blue != "" || template.BaseImage.BaseColour.Alpha != "")
	if (dataSet && fileSet) || (dataSet && baseColourSet) || (fileSet && baseColourSet) {
		return builder, fmt.Errorf("cannot load base image from file and load from data string and generate from base colour, specify only data or fileName or base colour")
	}
	if !dataSet && !fileSet && !baseColourSet {
		return builder.SetCanvas(builder.GetCanvas()).(ImageBuilder), nil
	}
	// Get image data from string or file
	var imageData []byte
	var err error
	var baseImage image.Image
	if baseColourSet {
		width64, err := strconv.ParseInt(template.BaseImage.BaseWidth, 10, 64) //Use ParseInt instead of Atoi for compatibility with go 1.7
		if err != nil {
			return builder, err
		}
		width := int(width64)
		height64, err := strconv.ParseInt(template.BaseImage.BaseHeight, 10, 64) //Use ParseInt instead of Atoi for compatibility with go 1.7
		if err != nil {
			return builder, err
		}
		height := int(height64)
		red64, err := strconv.ParseUint(template.BaseImage.BaseColour.Red, 0, 8)
		if err != nil {
			return builder, err
		}
		red := uint8(red64)
		green64, err := strconv.ParseUint(template.BaseImage.BaseColour.Green, 0, 8)
		if err != nil {
			return builder, err
		}
		green := uint8(green64)
		blue64, err := strconv.ParseUint(template.BaseImage.BaseColour.Blue, 0, 8)
		if err != nil {
			return builder, err
		}
		blue := uint8(blue64)
		alpha64, err := strconv.ParseUint(template.BaseImage.BaseColour.Alpha, 0, 8)
		if err != nil {
			return builder, err
		}
		alpha := uint8(alpha64)
		var rectImage image.Image
		var colourPlane image.Image
		rectangle := image.Rect(0, 0, width, height)
		rectImage = image.NewNRGBA(rectangle)
		colourPlane = image.NewUniform(color.NRGBA{R: red, G: green, B: blue, A: alpha})
		draw.Draw(rectImage.(draw.Image), rectangle, colourPlane, image.Point{X: 0, Y: 0}, draw.Over)
		baseImage = rectImage
	} else {
		if dataSet {
			sReader := strings.NewReader(template.BaseImage.Data)
			decoder := base64.NewDecoder(base64.StdEncoding, sReader)
			_, err = decoder.Read(imageData)
			if err != nil {
				return builder, err
			}
		} else {
			imageData, err = builder.reader.ReadFile(template.BaseImage.FileName)
			if err != nil {
				return builder, err
			}
		}
		// Decode image data
		imageBuffer := bytes.NewBuffer(imageData)
		baseImage, _, err = image.Decode(imageBuffer)
		if err != nil {
			return builder, err
		}
		if ycbcr, ok := baseImage.(*image.YCbCr); ok {
			var newImage draw.Image
			newImage = image.NewNRGBA(ycbcr.Rect)
			draw.Draw(newImage, ycbcr.Rect, ycbcr, ycbcr.Bounds().Min, draw.Over)
			baseImage = newImage
		}
	}
	if b.Canvas == nil {
		// No current canvas, uses loaded image as canvas
		var drawImage draw.Image
		drawImage = image.NewNRGBA(baseImage.Bounds())
		draw.Draw(drawImage, baseImage.Bounds(), baseImage, baseImage.Bounds().Min, draw.Over)
		b = b.SetCanvas(render.ImageCanvas{Image: drawImage}).(ImageBuilder)
		return b, nil
	}
	// Check if resizing is necessary
	currentHeight, currentWidth := baseImage.Bounds().Size().Y, baseImage.Bounds().Size().X
	targetHeight, targetWidth := b.GetCanvas().GetHeight(), b.GetCanvas().GetWidth()
	if targetHeight != currentHeight || targetWidth != currentWidth {
		// Compare aspect ratios
		targetAspect := float64(targetWidth) / float64(targetHeight)
		currentAspect := float64(currentWidth) / float64(currentHeight)
		var resizedWidth, resizedHeight int
		if targetAspect == currentAspect {
			// Identical apsect ratios
			resizedWidth = targetWidth
			resizedHeight = targetHeight
		} else if targetAspect < currentAspect {
			// Fit wide image into thin frame
			resizedHeight = targetHeight
		} else {
			// Fit thin image into wide frame
			resizedWidth = targetWidth
		}
		baseImage = imaging.Resize(baseImage, resizedWidth, resizedHeight, imaging.Lanczos)
	}
	canvas, err := b.GetCanvas().DrawImage(image.Point{X: 0, Y: 0}, baseImage)
	if err != nil {
		return builder, err
	}
	b = b.SetCanvas(canvas).(ImageBuilder)
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

// GetCanvas returns the internal Canvas object
func (builder ImageBuilder) GetCanvas() render.Canvas {
	if builder.Canvas == nil {
		return render.ImageCanvas{}
	}
	return builder.Canvas
}

// SetCanvas sets the internal Canvas object
func (builder ImageBuilder) SetCanvas(newCanvas render.Canvas) Builder {
	builder.Canvas = newCanvas
	return builder
}

// GetComponents gets the internal Component array
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

// SetComponents sets the internal Component array
func (builder ImageBuilder) SetComponents(components []ToggleableComponent) Builder {
	builder.Components = components
	return builder
}

// GetNamedPropertiesList returns the list of named properties in the builder object
func (builder ImageBuilder) GetNamedPropertiesList() render.NamedProperties {
	return builder.NamedProperties
}

// SetNamedProperties sets the values of names properties in all components and conditionals in the builder
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

// ApplyComponents iterates over the internal Component array, applying each in turn to the Canvas
func (builder ImageBuilder) ApplyComponents() (Builder, error) {
	b := builder
	for _, tComponent := range b.Components {
		if tComponent.Conditional.Name == "" {
			var err error
			b.Canvas, err = tComponent.Component.Write(b.Canvas)
			if err != nil {
				return builder, err
			}
		} else {
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
	}
	return b, nil
}
