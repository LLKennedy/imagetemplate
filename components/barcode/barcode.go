// Package barcode is a component for rendering barcodes with customisable content and colour for both background and data channels.
package barcode

import (
	"fmt"
	"image"
	"image/color"

	"github.com/LLKennedy/imagetemplate/v3/render"
	"golang.org/x/tools/godoc/vfs"
)

// Component implements the Component interface for images.
type Component struct {
	/*
		NamedPropertiesMap maps user/application variables to properties of the component.
		This field is filled automatically by VerifyAndSetJSONData, then used in
		SetNamedProperties to determine whether a variable being passed in is relevant to this
		component.

		For example, map[string][]string{"websiteURL": []string{"content"}} would indicate that
		the user specified variable "websiteURL" will fill the Content property.
	*/
	NamedPropertiesMap map[string][]string
	// Content is the data which will be encoded as a barcode.
	Content string
	// Type is the sort of barcode to encode, such as QR, PDF, or two of five.
	Type render.BarcodeType
	/*
		TopLeft is the coordinates of the top-left corner of the rendered barcode (including
		background) relative to the top-left corner of the canvas.
	*/
	TopLeft image.Point
	// Width is the width of the barcode (including background).
	Width int
	// Height is the height of the barcode (including background).
	Height int
	// DataColour is the colour which will fill the data channel.
	DataColour color.NRGBA
	// BackgroundColour is the colour which will fill the background channel.
	BackgroundColour color.NRGBA
	// Extra is additional information required by certain barcode types.
	Extra render.BarcodeExtraData
}

type barcodeFormat struct {
	Content    string `json:"content"`
	Type       string `json:"barcodeType"`
	TopLeftX   string `json:"topLeftX"`
	TopLeftY   string `json:"topLeftY"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	DataColour struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"dataColour"`
	BackgroundColour struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"backgroundColour"`
}

// Write draws a barcode on the canvas.
func (component Component) Write(canvas render.Canvas) (render.Canvas, error) {
	if len(component.NamedPropertiesMap) != 0 {
		return canvas, fmt.Errorf("cannot draw barcode, not all named properties are set: %v", component.NamedPropertiesMap)
	}
	c := canvas
	var err error
	c, err = c.Barcode(component.Type, []byte(component.Content), component.Extra, component.TopLeft, component.Width, component.Height, component.DataColour, component.BackgroundColour)
	if err != nil {
		return canvas, err
	}
	return c, nil
}

// SetNamedProperties processes the named properties and sets them into the barcode properties.
func (component Component) SetNamedProperties(properties render.NamedProperties) (render.Component, error) {
	c := component
	var err error
	c.NamedPropertiesMap, err = render.StandardSetNamedProperties(properties, component.NamedPropertiesMap, (&c).delegatedSetProperties)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a barcode component.
func (component Component) GetJSONFormat() interface{} {
	return &barcodeFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set barcode properties and fill the named properties map.
func (component Component) VerifyAndSetJSONData(data interface{}) (render.Component, render.NamedProperties, error) {
	c := component
	props := make(render.NamedProperties)
	stringStruct, ok := data.(*barcodeFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	return c.parseJSONFormat(stringStruct, props)
}

func init() {
	for _, name := range []string{"barcode", "bar", "code", "Barcode", "BARCODE", "BAR", "Bar Code", "bar code"} {
		render.RegisterComponent(name, func(vfs.FileSystem) render.Component { return Component{} })
	}
}
