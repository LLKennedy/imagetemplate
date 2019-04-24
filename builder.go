// Package imagetemplate defines a template for drawing custom images from pre-defined components, and provides to tools to load and implement that template.
package imagetemplate

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	_ "image/jpeg" // jpeg imported for image decoding
	_ "image/png"  // png imported for image decoding
	"io/ioutil"
	"strconv"
	"strings"
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
	LoadComponentsData(fileData []byte) (Builder, error)
	WriteToBMP() ([]byte, error)
}

// Template is the format of the JSON file used as a template for building images. See samples.json for examples, each element in the samples array is a complete and valid template object.
type Template struct {
	BaseImage struct {
		FileName   string `json:"fileName"`
		Data       string `json:"data"`
		BaseColour struct {
			Red   string `json:"R"`
			Green string `json:"G"`
			Blue  string `json:"B"`
			Alpha string `json:"A"`
		} `json:"baseColour"`
		BaseWidth  string `json:"width"`
		BaseHeight string `json:"height"`
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

// LoadComponentsFile sets the internal Component array based on the contents of the specified JSON file
func (builder ImageBuilder) LoadComponentsFile(fileName string) (Builder, error) {
	b := builder
	// Load initial data into template object
	fileData, err := ioutil.ReadFile(fileName)
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
	b.Canvas, err = setBackgroundImage(b.Canvas, template)
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

func setBackgroundImage(canvas Canvas, template Template) (Canvas, error) {
	c := canvas
	// Check the state of the optional and required properties
	dataSet := template.BaseImage.Data != ""
	fileSet := template.BaseImage.FileName != ""
	baseColourSet := template.BaseImage.BaseColour.Red != ""
	if (dataSet && fileSet) || (dataSet && baseColourSet) || fileSet && baseColourSet {
		return canvas, fmt.Errorf("Cannot load base image from file and load from data string and generate from base colour, specify only data or fileName or base colour")
	}
	if !dataSet && !fileSet && !baseColourSet {
		return c, nil
	}
	// Get image data from string or file
	var imageData []byte
	var err error
	var baseImage image.Image
	if baseColourSet {
		width, err := strconv.Atoi(template.BaseImage.BaseWidth)
		if err != nil {
			return canvas, err
		}
		height, err := strconv.Atoi(template.BaseImage.BaseWidth)
		if err != nil {
			return canvas, err
		}
		red, err := strconv.ParseUint(template.BaseImage.BaseColour.Red, 0, 8)
		if err != nil {
			return canvas, err
		}
		green, err := strconv.ParseUint(template.BaseImage.BaseColour.Green, 0, 8)
		if err != nil {
			return canvas, err
		}
		blue, err := strconv.ParseUint(template.BaseImage.BaseColour.Blue, 0, 8)
		if err != nil {
			return canvas, err
		}
		alpha, err := strconv.ParseUint(template.BaseImage.BaseColour.Alpha, 0, 8)
		if err != nil {
			return canvas, err
		}
		var rectImage image.Image
		var colourPlane image.Image
		rectangle := image.Rect(0, 0, width, height)
		rectImage = image.NewNRGBA(rectangle)
		colourPlane = image.NewUniform(color.NRGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(alpha)})
		draw.Draw(rectImage.(draw.Image), rectangle, colourPlane, image.Point{X: 0, Y: 0}, draw.Over)
		baseImage = rectImage
	} else {
		if dataSet {
			sReader := strings.NewReader(template.BaseImage.Data)
			decoder := base64.NewDecoder(base64.RawStdEncoding, sReader)
			_, err = decoder.Read(imageData)
			if err != nil {
				return canvas, err
			}
		} else {
			imageData, err = ioutil.ReadFile(template.BaseImage.FileName)
			if err != nil {
				return canvas, err
			}
		}
		// Decode image data
		imageBuffer := bytes.NewBuffer(imageData)
		baseImage, _, err = image.Decode(imageBuffer)
		if err != nil {
			return canvas, err
		}
	}
	if c == nil {
		// No current image, use loaded image instead
		drawImg, ok := baseImage.(draw.Image)
		if !ok {
			return canvas, fmt.Errorf("Could not create write-access Image from image data")
		}
		c = ImageCanvas{Image: drawImg}
		return c, nil
	}
	// Check if resizing is necessary
	currentHeight, currentWidth := baseImage.Bounds().Size().Y, baseImage.Bounds().Size().X
	targetHeight, targetWidth := c.GetHeight(), c.GetWidth()
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
	c.SubImage(image.Point{X: 0, Y: 0}, baseImage)
	return c, nil
}

func parseComponents(templates []ComponentTemplate) ([]ToggleableComponent, NamedProperties, error) {
	var results []ToggleableComponent
	namedProperties := NamedProperties{}
	for tCount, template := range templates {
		//Handle conditional first
		result := ToggleableComponent{}
		tempProperties := namedProperties
		result.Conditional = template.Conditional
		for key, value := range result.Conditional.GetNamedPropertiesList() {
			tempProperties[key] = value
		}
		var typeRange []string
		switch template.Type {
		case "circle", "Circle", "rectangle", "Rectangle", "rect", "Rect", "image", "Image", "photo", "Photo", "text", "Text", "words", "Words":
			typeRange = []string{template.Type}
		default:
			typeRange = []string{"circle", "Circle", "rectange", "Rectangle", "rect", "Rect", "image", "Image", "photo", "Photo", "text", "Text", "words", "Words"}
		}
		for _, compType := range typeRange {
			var newComponent Component
			switch compType {
			case "circle", "Circle":
				newComponent = CircleComponent{}
			case "rectangle", "Rectangle", "rect", "Rect":
				newComponent = RectangleComponent{}
			case "image", "Image", "photo", "Photo":
				newComponent = ImageComponent{}
			case "text", "Text", "words", "Words":
				newComponent = TextComponent{}
			}
			// Get JSON struct to parse into
			shape := newComponent.GetJSONFormat()
			err := json.Unmarshal(template.Properties, shape)
			if err != nil {
				// Invalid JSON
				return results, namedProperties, err
			}
			// Set real properties from JSON struct
			newComponent, compNamedProps, err := newComponent.VerifyAndSetJSONData(shape)
			if err != nil {
				// Didn't match this type
				continue
			}
			for key, value := range compNamedProps {
				tempProperties[key] = value
			}
			result.Component = newComponent
			results = append(results, result)
			namedProperties = tempProperties
		}
		if len(results) <= tCount {
			// Failed to find a matching type
			return results, namedProperties, fmt.Errorf("Failed to find type matching component with user-specified type %v", template.Type)
		}
	}
	return results, namedProperties, nil
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
