package imagetemplate

import (
	"fmt"
	"github.com/LLKennedy/gosysfonts"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

// TextComponent implements the Component interface for text
type TextComponent struct {
	NamedPropertiesMap map[string][]string
	Content            string
	Start              image.Point
	Size               int
	MaxWidth           int
	Font               *truetype.Font
	Colour             color.Color
}

type textFormat struct {
	Content  string `json:"content"`
	StartX   string `json:"startX"`
	StartY   string `json:"startY"`
	Size     string `json:"size"`
	MaxWidth string `json:"maxWidth"`
	Font     struct {
		FontName string `json:"fontName"`
		FontFile string `json:"fontFile"`
		FontURL  string `json:"fontURL"`
	} `json:"font"`
	Colour struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"colour"`
}

// Write draws text on the canvas
func (component TextComponent) Write(canvas Canvas) (Canvas, error) {
	c := canvas
	pool := gosysfonts.New()
	rawFont, err := pool.GetFont("Calibri")
	if err != nil {
		return canvas, err
	}
	face := truetype.NewFace(rawFont, &truetype.Options{Size: 14, Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64})
	c, err = c.Text(component.Content, component.Start, face, component.Colour, component.MaxWidth)
	if err != nil {
		return canvas, err
	}
	return c, nil
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
