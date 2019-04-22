package imagetemplate

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

// CircleComponent implements the Component interface for circles
type CircleComponent struct {
	NamedPropertiesMap map[string][]string
	Centre             image.Point
	Radius             int
	Colour             color.NRGBA
}

type circleFormat struct {
	CentreX string `json:"centreX"`
	CentreY string `json:"centreY"`
	Radius  string `json:"radius"`
	Colour  struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"colour"`
}

// Write draws a circle on the canvas
func (component CircleComponent) Write(canvas Canvas) (Canvas, error) {
	if len(component.NamedPropertiesMap) != 0 {
		return canvas, fmt.Errorf("Cannot draw circle, not all named properties are set: %v", component.NamedPropertiesMap)
	}
	return canvas.Circle(component.Centre, component.Radius, component.Colour)
}

// SetNamedProperties processes the named properties and sets them into the circle properties
func (component CircleComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
	c := component
	setFunc := func(name string, value interface{}) error {
		if strings.Contains("RGBA", name) && len(name) == 1 {
			//Process colours
			colourVal, ok := value.(uint8)
			if !ok {
				return fmt.Errorf("Error converting %v to uint8", value)
			}
			switch name {
			case "R":
				c.Colour.R = colourVal
				return nil
			case "G":
				c.Colour.G = colourVal
				return nil
			case "B":
				c.Colour.B = colourVal
				return nil
			case "A":
				c.Colour.A = colourVal
				return nil
			default:
				//What? How did you get here?
				return fmt.Errorf("Name was a string inside RGBA and Value was a valid uint8, but Name wasn't R, G, B, or A. Name was: %v", name)
			}
		}
		numberVal, ok := value.(int)
		if !ok {
			return fmt.Errorf("Error converting %v to int", value)
		}
		switch name {
		case "centreX":
			c.Centre.X = numberVal
			return nil
		case "centreY":
			c.Centre.Y = numberVal
			return nil
		case "radius":
			c.Radius = numberVal
			return nil
		default:
			return fmt.Errorf("Invalid component property in named property map: %v", name)
		}
	}
	var err error
	c.NamedPropertiesMap, err = StandardSetNamedProperties(properties, component.NamedPropertiesMap, setFunc)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a circle component
func (component CircleComponent) GetJSONFormat() interface{} {
	return &circleFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set circle properties and fill the named properties map
func (component CircleComponent) VerifyAndSetJSONData(data interface{}) (Component, NamedProperties, error) {
	c := component
	var props NamedProperties
	stringStruct, ok := data.(*circleFormat)
	if !ok {
		return component, props, fmt.Errorf("Failed to convert returned data to component properties")
	}
	// Get named properties and assign each real property
	var newVal interface{}
	var err error
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.CentreX, "centreX", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Centre.X = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.CentreY, "centreY", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Centre.Y = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Radius, "radius", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Radius = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Colour.Red, "R", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Colour.R = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Colour.Green, "G", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Colour.G = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Colour.Blue, "B", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Colour.B = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Colour.Alpha, "A", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Colour.A = newVal.(uint8)
	}
	return c, props, nil
}

// RectangleComponent implements the Component interface for rectangles
type RectangleComponent struct {
	NamedPropertiesMap map[string][]string
	TopLeft            image.Point
	Width              int
	Height             int
	Colour             color.Color
}

type rectangleFormat struct {
	CentreX string `json:"centreX"`
	CentreY string `json:"centreY"`
	Radius  string `json:"radius"`
	Colour  struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"colour"`
}

// Write draws a rectangle on the canvas
func (component RectangleComponent) Write(canvas Canvas) (Canvas, error) {
	return canvas, fmt.Errorf("Not implemented yet")
}

// SetNamedProperties proceses the named properties and sets them into the rectangle properties
func (component RectangleComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
	return component, fmt.Errorf("Not implemented yet")
}

// GetJSONFormat returns the JSON structure of a rectangle component
func (component RectangleComponent) GetJSONFormat() interface{} {
	return &rectangleFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set rectangle properties and fill the named properties map
func (component RectangleComponent) VerifyAndSetJSONData(interface{}) (Component, NamedProperties, error) {
	return component, nil, fmt.Errorf("Not implemented yet")
}

// ImageComponent implements the Component interface for images
type ImageComponent struct {
	NamedPropertiesMap map[string][]string
	Image              image.Image
	TopLeft            image.Point
	Width              int
	Height             int
}

type imageFormat struct {
	CentreX string `json:"centreX"`
	CentreY string `json:"centreY"`
	Radius  string `json:"radius"`
	Colour  struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"colour"`
}

// Write draws an image on the canvas
func (component ImageComponent) Write(canvas Canvas) (Canvas, error) {
	return canvas, fmt.Errorf("Not implemented yet")
}

// SetNamedProperties proceses the named properties and sets them into the image properties
func (component ImageComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
	return component, fmt.Errorf("Not implemented yet")
}

// GetJSONFormat returns the JSON structure of a image component
func (component ImageComponent) GetJSONFormat() interface{} {
	return &imageFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set image properties and fill the named properties map
func (component ImageComponent) VerifyAndSetJSONData(interface{}) (Component, NamedProperties, error) {
	return component, nil, fmt.Errorf("Not implemented yet")
}

// TextComponent implements the Component interface for text
type TextComponent struct {
	NamedPropertiesMap map[string][]string
	TopLeft            image.Point
	Width              int
	Height             int
	Colour             color.Color
}

type textFormat struct {
	CentreX string `json:"centreX"`
	CentreY string `json:"centreY"`
	Radius  string `json:"radius"`
	Colour  struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"colour"`
}

// Write draws text on the canvas
func (component TextComponent) Write(canvas Canvas) (Canvas, error) {
	return canvas, fmt.Errorf("Not implemented yet")
}

// SetNamedProperties proceses the named properties and sets them into the text properties
func (component TextComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
	return component, fmt.Errorf("Not implemented yet")
}

// GetJSONFormat returns the JSON structure of a text component
func (component TextComponent) GetJSONFormat() interface{} {
	return &textFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set text properties and fill the named properties map
func (component TextComponent) VerifyAndSetJSONData(interface{}) (Component, NamedProperties, error) {
	return component, nil, fmt.Errorf("Not implemented yet")
}
