package imagetemplate

import (
	"fmt"
	"image"
	"image/color"
)

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
