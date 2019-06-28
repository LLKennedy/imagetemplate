// Package circle is a simple circle component with customisable size, colour, location and radius.
package circle

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/LLKennedy/imagetemplate/v3/internal/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"golang.org/x/tools/godoc/vfs"
)

// Component implements the Component interface for circles.
type Component struct {
	/*
		NamedPropertiesMap maps user/application variables to properties of the component.
		This field is filled automatically by VerifyAndSetJSONData, then used in
		SetNamedProperties to determine whether a variable being passed in is relevant to this
		component.

		For example, map[string][]string{"circleSize": []string{"radius"}} would indicate that
		the user specified variable "circleSize" will fill the Radius property.
	*/
	NamedPropertiesMap map[string][]string
	/*
		Centre is the coordinates of the centre of the circle relative to the top-left corner
		of the canvas.
	*/
	Centre image.Point
	// Radius is the radius of the circle.
	Radius int
	// Colour is the colour of the circle.
	Colour color.NRGBA
}

type circleFormat struct {
	CentreX string       `json:"centreX"`
	CentreY string       `json:"centreY"`
	Radius  string       `json:"radius"`
	Colour  colourFormat `json:"colour"`
}

type colourFormat struct {
	Red   string `json:"R"`
	Green string `json:"G"`
	Blue  string `json:"B"`
	Alpha string `json:"A"`
}

// Write draws a circle on the canvas.
func (component Component) Write(canvas render.Canvas) (render.Canvas, error) {
	if len(component.NamedPropertiesMap) != 0 {
		return canvas, fmt.Errorf("cannot draw circle, not all named properties are set: %v", component.NamedPropertiesMap)
	}
	return canvas.Circle(component.Centre, component.Radius, component.Colour)
}

// SetNamedProperties processes the named properties and sets them into the circle properties.
func (component Component) SetNamedProperties(properties render.NamedProperties) (render.Component, error) {
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
	c.NamedPropertiesMap, err = render.StandardSetNamedProperties(properties, component.NamedPropertiesMap, setFunc)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a circle component.
func (component Component) GetJSONFormat() interface{} {
	return &circleFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set circle properties and fill the named properties map.
func (component Component) VerifyAndSetJSONData(data interface{}) (render.Component, render.NamedProperties, error) {
	c := component
	props := make(render.NamedProperties)
	stringStruct, ok := data.(*circleFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	// Get named properties and assign each real property
	var newVal interface{}
	var err error
	c.Centre.X, c.NamedPropertiesMap, err = cutils.ExtractInt(stringStruct.CentreX, "centreX", c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	c.Centre.Y, c.NamedPropertiesMap, err = cutils.ExtractInt(stringStruct.CentreY, "centreY", c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	c.Radius, c.NamedPropertiesMap, err = cutils.ExtractInt(stringStruct.Radius, "radius", c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Colour.Red, "R", render.Uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Colour.R = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Colour.Green, "G", render.Uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Colour.G = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Colour.Blue, "B", render.Uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Colour.B = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Colour.Alpha, "A", render.Uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Colour.A = newVal.(uint8)
	}
	for key := range c.NamedPropertiesMap {
		props[key] = struct {
			Message string
		}{Message: "Please replace me with real data"}
	}
	return c, props, nil
}

func init() {
	for _, name := range []string{"circle", "Circle", "CIRCLE"} {
		render.RegisterComponent(name, func(vfs.FileSystem) render.Component { return Component{} })
	}
}
