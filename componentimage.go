package imagetemplate

import (
	"fmt"
	"image"
)

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
	return canvas, fmt.Errorf("not implemented yet")
}

// SetNamedProperties proceses the named properties and sets them into the image properties
func (component ImageComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
	return component, fmt.Errorf("not implemented yet")
}

// GetJSONFormat returns the JSON structure of a image component
func (component ImageComponent) GetJSONFormat() interface{} {
	return &imageFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set image properties and fill the named properties map
func (component ImageComponent) VerifyAndSetJSONData(interface{}) (Component, NamedProperties, error) {
	return component, nil, fmt.Errorf("not implemented yet")
}
