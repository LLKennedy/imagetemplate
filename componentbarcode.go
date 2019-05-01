package imagetemplate

import (
	"fmt"
	"image"
)

// BarcodeComponent implements the Component interface for images
type BarcodeComponent struct {
	NamedPropertiesMap map[string][]string
	Image              image.Image
	TopLeft            image.Point
	Width              int
	Height             int
}

type barcodeFormat struct {
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

// Write draws a barcode on the canvas
func (component BarcodeComponent) Write(canvas Canvas) (Canvas, error) {
	return canvas, fmt.Errorf("Not implemented yet")
}

// SetNamedProperties proceses the named properties and sets them into the barcode properties
func (component BarcodeComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
	return component, fmt.Errorf("Not implemented yet")
}

// GetJSONFormat returns the JSON structure of a barcode component
func (component BarcodeComponent) GetJSONFormat() interface{} {
	return &barcodeFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set barcode properties and fill the named properties map
func (component BarcodeComponent) VerifyAndSetJSONData(interface{}) (Component, NamedProperties, error) {
	return component, nil, fmt.Errorf("Not implemented yet")
}
