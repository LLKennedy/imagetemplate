// Package text is a simple text component with customisable content, size, colour, location and font.
package text

import (
	"fmt"
	"image"
	"image/color"
	"runtime/debug"

	"github.com/LLKennedy/gosysfonts"
	"github.com/LLKennedy/imagetemplate/v3/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/tools/godoc/vfs"
)

// Component implements the Component interface for text.
type Component struct {
	/*
		NamedPropertiesMap maps user/application variables to properties of the component.
		This field is filled automatically by VerifyAndSetJSONData, then used in
		SetNamedProperties to determine whether a variable being passed in is relevant to this
		component.

		For example, map[string][]string{"username": []string{"content"}} would indicate that
		the user specified variable "username" will fill the Content property.
	*/
	NamedPropertiesMap map[string][]string
	// Content is the text to render.
	Content string
	// Start is the coordinates of the dot relative to the top-left corner of the canvas.
	Start image.Point
	// Size is the size of the text in points.
	Size float64
	// MaxWidth is the maximum number of horizontal pixels the dot can move before scaling text.
	MaxWidth int
	// cutils.TextAlignment aligns text to the left, right or centre.
	TextAlignment cutils.TextAlignment
	// Font is the typeface to use.
	Font *truetype.Font
	// Colour is the colour of the text.
	Colour color.NRGBA
	// fs is the file system.
	fs vfs.FileSystem
	// fontPool is the pool of available fonts.
	fontPool gosysfonts.Pool
}

type textFormat struct {
	Content       string `json:"content"`
	StartX        string `json:"startX"`
	StartY        string `json:"startY"`
	Size          string `json:"size"`
	MaxWidth      string `json:"maxWidth"`
	TextAlignment string `json:"alignment"`
	Font          struct {
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

// Write draws text on the canvas.
func (component Component) Write(canvas render.Canvas) (c render.Canvas, err error) {
	c = canvas
	defer func() {
		p := recover()
		if p != nil {
			err = fmt.Errorf("failed to write to canvas: %v\n%s", p, debug.Stack())
		}
	}()
	fontSize := component.Size
	fits := false
	tries := 0
	var face font.Face
	var alignmentOffset int
	for !fits && tries < 10 {
		tries++
		face = truetype.NewFace(component.Font, &truetype.Options{Size: fontSize, Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: canvas.GetPPI()})
		var realWidth int
		fits, realWidth = c.TryText(component.Content, component.Start, font.Face(face), component.Colour, component.MaxWidth)
		fontSize, alignmentOffset = cutils.ScaleFontsToWidth(fontSize, realWidth, component.MaxWidth, component.TextAlignment)
	}
	if !fits {
		return canvas, fmt.Errorf("unable to fit text %v into maxWidth %d after %d tries", component.Content, component.MaxWidth, tries)
	}
	c, err = c.Text(component.Content, image.Pt(component.Start.X+alignmentOffset, component.Start.Y), font.Face(face), component.Colour, component.MaxWidth)
	if err != nil {
		return canvas, err
	}
	return c, nil
}

// SetNamedProperties processes the named properties and sets them into the text properties.
func (component Component) SetNamedProperties(properties render.NamedProperties) (render.Component, error) {
	c := component
	var err error
	c.NamedPropertiesMap, err = render.StandardSetNamedProperties(properties, component.NamedPropertiesMap, (&c).delegatedSetProperties)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a text component.
func (component Component) GetJSONFormat() interface{} {
	return &textFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set text properties and fill the named properties map.
func (component Component) VerifyAndSetJSONData(data interface{}) (render.Component, render.NamedProperties, error) {
	c := component
	props := make(render.NamedProperties)
	stringStruct, ok := data.(*textFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	return c.parseJSONFormat(stringStruct, props)
}

func (component Component) getFontPool() gosysfonts.Pool {
	if component.fontPool == nil {
		return gosysfonts.New()
	}
	return component.fontPool
}

func (component Component) getFileSystem() vfs.FileSystem {
	if component.fs == nil {
		return vfs.OS(".")
	}
	return component.fs
}

func init() {
	for _, name := range []string{"text", "Text", "TEXT", "words", "Words", "WORDS", "writing", "Writing", "WRITING"} {
		render.RegisterComponent(name, func(fs vfs.FileSystem) render.Component { return Component{fs: fs, fontPool: gosysfonts.New()} })
	}
}
