package imagetemplate

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/disintegration/imaging"
	_ "golang.org/x/image/bmp"  // bmp imported for image decoding
	_ "golang.org/x/image/tiff" // tiff imported for image decoding
	"image"
	_ "image/jpeg" // jpeg imported for image decoding
	_ "image/png"  // png imported for image decoding
	"strings"
)

// ImageComponent implements the Component interface for images
type ImageComponent struct {
	NamedPropertiesMap map[string][]string
	Image              image.Image
	TopLeft            image.Point
	Width              int
	Height             int
	reader             fileReader
}

type imageFormat struct {
	TopLeftX string `json:"topLeftX"`
	TopLeftY string `json:"topLeftY"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	FileName string `json:"fileName"`
	Data     string `json:"data"`
}

// Write draws an image on the canvas
func (component ImageComponent) Write(canvas Canvas) (Canvas, error) {
	c := canvas
	var err error
	scaledImage := imaging.Resize(component.Image, component.Width, component.Height, imaging.Lanczos)
	c, err = c.DrawImage(component.TopLeft, scaledImage)
	if err != nil {
		return canvas, err
	}
	return c, nil
}

// SetNamedProperties proceses the named properties and sets them into the image properties
func (component ImageComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
	c := component
	setFunc := func(name string, value interface{}) error {
		switch name {
		case "data":
			bytesVal, ok := value.([]byte)
			if !ok {
				return fmt.Errorf("error converting %v to []byte", value)
			}
			buf := bytes.NewBuffer(bytesVal)
			img, _, err := image.Decode(buf)
			if err != nil {
				return err
			}
			c.Image = img
			return nil
		case "fileName":
			stringVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("error converting %v to string", value)
			}
			if component.reader == nil {
				component.reader = ioutilFileReader{}
			}
			bytesVal, err := component.reader.ReadFile(stringVal)
			if err != nil {
				return err
			}
			buf := bytes.NewBuffer(bytesVal)
			img, _, err := image.Decode(buf)
			if err != nil {
				return err
			}
			c.Image = img
			return nil
		}
		numberVal, ok := value.(int)
		if !ok {
			return fmt.Errorf("error converting %v to int", value)
		}
		switch name {
		case "topLeftX":
			c.TopLeft.X = numberVal
			return nil
		case "topLeftY":
			c.TopLeft.Y = numberVal
			return nil
		case "width":
			c.Width = numberVal
			return nil
		case "height":
			c.Height = numberVal
			return nil
		default:
			return fmt.Errorf("invalid component property in named property map: %v", name)
		}
	}
	var err error
	c.NamedPropertiesMap, err = StandardSetNamedProperties(properties, component.NamedPropertiesMap, setFunc)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a image component
func (component ImageComponent) GetJSONFormat() interface{} {
	return &imageFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set image properties and fill the named properties map
func (component ImageComponent) VerifyAndSetJSONData(data interface{}) (Component, NamedProperties, error) {
	c := component
	var props NamedProperties
	stringStruct, ok := data.(*imageFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	// Get named properties and assign each real property
	var newVal interface{}
	var err error
	// Deal with the file/data restrictions
	inputs := []string{
		stringStruct.FileName,
		stringStruct.Data,
	}
	propNames := []string{
		"fileName",
		"data",
	}
	types := []propType{
		stringType,
		stringType,
	}
	var extractedVal interface{}
	validIndex := -1
	c.NamedPropertiesMap, extractedVal, validIndex, err = extractExclusiveProp(inputs, propNames, types, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	var file, fdata interface{}
	switch validIndex {
	case 0:
		file = extractedVal
	case 1:
		fdata = extractedVal
	default:
		return component, props, fmt.Errorf("failed to extract image file")
	}
	if fdata != nil {
		base64Val := fdata.(string)
		r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Val))
		img, _, err := image.Decode(r)
		if err != nil {
			return component, props, err
		}
		c.Image = img
	}
	if file != nil {
		stringVal := file.(string)
		if component.reader == nil {
			component.reader = ioutilFileReader{}
		}
		bytesVal, err := component.reader.ReadFile(stringVal)
		if err != nil {
			return component, props, err
		}
		buf := bytes.NewBuffer(bytesVal)
		img, _, err := image.Decode(buf)
		if err != nil {
			return component, props, err
		}
		c.Image = img
	}

	// All other props
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.TopLeftX, "topLeftX", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.TopLeft.X = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.TopLeftY, "topLeftY", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.TopLeft.Y = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Width, "width", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Width = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Height, "height", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Height = newVal.(int)
	}
	return c, props, nil
}
