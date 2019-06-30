// Package rectangle is a simple rectangle component with customisable size, colour and location.
package rectangle

import (
	"fmt"
	"image"
	"image/color"

	"github.com/LLKennedy/imagetemplate/v3/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"golang.org/x/tools/godoc/vfs"
)

// Component implements the Component interface for rectangles.
type Component struct {
	/*
		NamedPropertiesMap maps user/application variables to properties of the component.
		This field is filled automatically by VerifyAndSetJSONData, then used in
		SetNamedProperties to determine whether a variable being passed in is relevant to this
		component.

		For example, map[string][]string{"squareSize": []string{"width", "height"}} would
		indicate that the user specified variable "squareSize" will fill the Width and Height
		properties.
	*/
	NamedPropertiesMap map[string][]string
	/*
		TopLeft is the coordinates of the top-left corner of the rectangle relative to the
		top-left corner of the canvas.
	*/
	TopLeft image.Point
	// Width is the width of the rectangle.
	Width int
	// Height is the height of the rectangle.
	Height int
	// Colour is the colour of the rectangle.
	Colour color.NRGBA
}

type rectangleFormat struct {
	TopLeftX string `json:"topLeftX"`
	TopLeftY string `json:"topLeftY"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Colour   struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"colour"`
}

// Write draws a rectangle on the canvas.
func (component Component) Write(canvas render.Canvas) (render.Canvas, error) {
	if len(component.NamedPropertiesMap) != 0 {
		return canvas, fmt.Errorf("cannot draw rectangle, not all named properties are set: %v", component.NamedPropertiesMap)
	}
	return canvas.Rectangle(component.TopLeft, component.Width, component.Height, component.Colour)
}

// SetNamedProperties processes the named properties and sets them into the rectangle properties.
func (component Component) SetNamedProperties(properties render.NamedProperties) (render.Component, error) {
	c := component
	var err error
	c.NamedPropertiesMap, err = render.StandardSetNamedProperties(properties, component.NamedPropertiesMap, (&c).delegatedSetProperties)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a rectangle component.
func (component Component) GetJSONFormat() interface{} {
	return &rectangleFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set rectangle properties and fill the named properties map.
func (component Component) VerifyAndSetJSONData(data interface{}) (render.Component, render.NamedProperties, error) {
	c := component
	props := make(render.NamedProperties)
	stringStruct, ok := data.(*rectangleFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	// Get named properties and assign each real property
	var newVal interface{}
	var err error
	c.TopLeft.X, c.NamedPropertiesMap, err = cutils.ExtractInt(stringStruct.TopLeftX, "topLeftX", c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	c.TopLeft.Y, c.NamedPropertiesMap, err = cutils.ExtractInt(stringStruct.TopLeftY, "topLeftY", c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	c.Width, c.NamedPropertiesMap, err = cutils.ExtractInt(stringStruct.Width, "width", c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	c.Height, c.NamedPropertiesMap, err = cutils.ExtractInt(stringStruct.Height, "height", c.NamedPropertiesMap)
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
	for _, name := range []string{"rect", "RECT", "Rect", "rectangle", "Rectangle", "RECTANGLE"} {
		render.RegisterComponent(name, func(vfs.FileSystem) render.Component { return Component{} })
	}
}
