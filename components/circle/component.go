package imagetemplate

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"github.com/LLKennedy/imagetemplate/render"
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
func (component CircleComponent) Write(canvas render.Canvas) (Canvas, error) {
	if len(component.NamedPropertiesMap) != 0 {
		return canvas, fmt.Errorf("cannot draw circle, not all named properties are set: %v", component.NamedPropertiesMap)
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
				return fmt.Errorf("error converting %v to uint8", value)
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
			}
		}
		numberVal, ok := value.(int)
		if !ok {
			return fmt.Errorf("error converting %v to int", value)
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

// GetJSONFormat returns the JSON structure of a circle component
func (component CircleComponent) GetJSONFormat() interface{} {
	return &circleFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set circle properties and fill the named properties map
func (component CircleComponent) VerifyAndSetJSONData(data interface{}) (Component, NamedProperties, error) {
	c := component
	props := make(NamedProperties)
	stringStruct, ok := data.(*circleFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
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
	type invalidStruct struct {
		Message string
	}
	for key := range c.NamedPropertiesMap {
		props[key] = invalidStruct{Message: "Please replace me with real data"}
	}
	return c, props, nil
}
