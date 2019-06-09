package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg" // jpeg imported for image decoding
	_ "image/png"  // png imported for image decoding
	"io"
	"strings"

	fs "github.com/LLKennedy/imagetemplate/v2/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/v2/render"
	"github.com/disintegration/imaging"
	_ "golang.org/x/image/bmp"  // bmp imported for image decoding
	_ "golang.org/x/image/tiff" // tiff imported for image decoding
)

// Component implements the Component interface for images
type Component struct {
	NamedPropertiesMap map[string][]string
	Image              image.Image
	TopLeft            image.Point
	Width              int
	Height             int
	reader             fs.FileReader
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
func (component Component) Write(canvas render.Canvas) (render.Canvas, error) {
	if len(component.NamedPropertiesMap) != 0 {
		return canvas, fmt.Errorf("cannot draw image, not all named properties are set: %v", component.NamedPropertiesMap)
	}
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
func (component Component) SetNamedProperties(properties render.NamedProperties) (render.Component, error) {
	c := component
	setFunc := func(name string, value interface{}) error {
		switch name {
		case "data":
			bytesVal, isBytes := value.([]byte)
			stringVal, isString := value.(string)
			readerVal, isReader := value.(io.Reader)
			if !isBytes && !isString && !isReader {
				return fmt.Errorf("error converting %v to []byte, string or io.Reader", value)
			}
			var reader io.Reader
			if isBytes {
				reader = bytes.NewBuffer(bytesVal)
			} else if isString {
				stringReader := strings.NewReader(stringVal)
				reader = base64.NewDecoder(base64.StdEncoding, stringReader)
			} else if isReader {
				reader = readerVal
			}
			img, _, err := image.Decode(reader)
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
				component.reader = fs.IoutilFileReader{}
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
	c.NamedPropertiesMap, err = render.StandardSetNamedProperties(properties, component.NamedPropertiesMap, setFunc)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a image component
func (component Component) GetJSONFormat() interface{} {
	return &imageFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set image properties and fill the named properties map
func (component Component) VerifyAndSetJSONData(data interface{}) (render.Component, render.NamedProperties, error) {
	c := component
	props := make(render.NamedProperties)
	stringStruct, ok := data.(*imageFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	// Get named properties and assign each real property
	var newVal interface{}
	var err error
	// Deal with the file/data restrictions
	propData := []render.PropData{
		{
			InputValue: stringStruct.FileName,
			PropName:   "fileName",
			Type:       render.StringType,
		},
		{
			InputValue: stringStruct.Data,
			PropName:   "data",
			Type:       render.StringType,
		},
	}
	var extractedVal interface{}
	validIndex := -1
	c.NamedPropertiesMap, extractedVal, validIndex, err = render.ExtractExclusiveProp(propData, c.NamedPropertiesMap)
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
			component.reader = fs.IoutilFileReader{}
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
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.TopLeftX, "topLeftX", render.IntType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.TopLeft.X = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.TopLeftY, "topLeftY", render.IntType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.TopLeft.Y = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Width, "width", render.IntType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Width = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Height, "height", render.IntType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Height = newVal.(int)
	}

	for key := range c.NamedPropertiesMap {
		props[key] = struct {
			Message string
		}{Message: "Please replace me with real data"}
	}
	return c, props, nil
}

func init() {
	for _, name := range []string{"image", "img", "photo", "Image", "IMG", "Photo", "picture", "Picture", "IMAGE", "PHOTO", "PICTURE"} {
		render.RegisterComponent(name, func() render.Component { return Component{} })
	}
}
