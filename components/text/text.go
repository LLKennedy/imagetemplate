package text

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/LLKennedy/gosysfonts"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	fs "github.com/LLKennedy/imagetemplate/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/render"
)

// Component implements the Component interface for text
type Component struct {
	NamedPropertiesMap map[string][]string
	Content            string
	Start              image.Point
	Size               float64
	MaxWidth           int
	Alignment          Alignment
	PixelsPerInch      int //Should default to 72
	Font               *truetype.Font
	Colour             color.NRGBA
	reader             fs.FileReader
}

type textFormat struct {
	Content       string `json:"content"`
	StartX        string `json:"startX"`
	StartY        string `json:"startY"`
	Size          string `json:"size"`
	MaxWidth      string `json:"maxWidth"`
	Alignment     string `json:"alignment"`
	PixelsPerInch string `json:"ppi"`
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

// Alignment is a text alignment
type Alignment int

const (
	// AlignmentLeft aligns text left
	AlignmentLeft Alignment = iota
	// AlignmentRight aligns text right
	AlignmentRight
	// AlignmentCentre aligns text centrally
	AlignmentCentre
)

// Write draws text on the canvas
func (component Component) Write(canvas render.Canvas) (render.Canvas, error) {
	c := canvas
	fontSize := (component.Size / 72) * float64(component.PixelsPerInch) // one point in fonts is almost exactly 1/72nd of one inch
	fits := false
	tries := 0
	var face font.Face
	var alignmentOffset int
	for !fits && tries < 10 {
		tries++
		face = truetype.NewFace(component.Font, &truetype.Options{Size: fontSize, Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64})
		var realWidth int
		fits, realWidth = c.TryText(component.Content, component.Start, face, component.Colour, component.MaxWidth)
		if realWidth > component.MaxWidth {
			ratio := float64(component.MaxWidth) / float64(realWidth)
			fontSize = ratio * fontSize
		} else if realWidth < component.MaxWidth {
			remainingWidth := float64(component.MaxWidth) - float64(realWidth)
			switch component.Alignment {
			case AlignmentLeft:
				alignmentOffset = 0
			case AlignmentRight:
				alignmentOffset = int(remainingWidth)
			case AlignmentCentre:
				alignmentOffset = int(remainingWidth / 2)
			default:
				alignmentOffset = 0
			}
		}
	}
	if !fits {
		return canvas, fmt.Errorf("unable to fit text %v into maxWidth %d after %d tries", component.Content, component.MaxWidth, tries)
	}
	c, err := c.Text(component.Content, image.Pt(component.Start.X+alignmentOffset, component.Start.Y), face, component.Colour, component.MaxWidth)
	if err != nil {
		return canvas, err
	}
	return c, nil
}

// SetNamedProperties proceses the named properties and sets them into the text properties
func (component Component) SetNamedProperties(properties render.NamedProperties) (render.Component, error) {
	c := component
	setFunc := func(name string, value interface{}) error {
		switch name {
		case "content":
			stringVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("error converting %v to string", value)
			}
			c.Content = stringVal
			return nil
		case "fontName":
			stringVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("error converting %v to string", value)
			}
			pool := gosysfonts.New()
			rawFont, err := pool.GetFont(stringVal)
			if err != nil {
				return err
			}
			c.Font = rawFont
			return nil
		case "fontFile":
			stringVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("error converting %v to string", value)
			}
			if component.reader == nil {
				component.reader = fs.IoutilFileReader{}
			}
			fontData, err := component.reader.ReadFile(stringVal)
			if err != nil {
				return err
			}
			rawFont, err := truetype.Parse(fontData)
			if err != nil {
				return err
			}
			c.Font = rawFont
			return nil
		case "fontURL":
			return fmt.Errorf("fontURL not implemented")
		case "size":
			float64Val, ok := value.(float64)
			if !ok {
				return fmt.Errorf("error converting %v to float64", value)
			}
			c.Size = float64Val
			return nil
		case "alignment":
			alignmentVal, isAlignment := value.(Alignment)
			stringVal, isString := value.(string)
			if !isAlignment && !isString {
				return fmt.Errorf("could not convert %v to text alignment or string", value)
			}
			if isAlignment {
				c.Alignment = alignmentVal
			} else {
				switch stringVal {
				case "left":
					c.Alignment = AlignmentLeft
				case "right":
					c.Alignment = AlignmentRight
				case "centre":
					c.Alignment = AlignmentCentre
				default:
					c.Alignment = AlignmentLeft
				}
			}
		}
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
			default:
				//What? How did you get here?
				return fmt.Errorf("name was a string inside RGBA and Value was a valid uint8, but Name wasn't R, G, B, or A. Name was: %v", name)
			}
		}
		numberVal, ok := value.(int)
		if !ok {
			return fmt.Errorf("error converting %v to int", value)
		}
		switch name {
		case "startX":
			c.Start.X = numberVal
			return nil
		case "startY":
			c.Start.Y = numberVal
			return nil
		case "maxWidth":
			c.MaxWidth = numberVal
			return nil
		case "ppi":
			c.PixelsPerInch = numberVal
			if c.PixelsPerInch <= 0 {
				c.PixelsPerInch = 72
			}
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

// GetJSONFormat returns the JSON structure of a text component
func (component Component) GetJSONFormat() interface{} {
	return &textFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set text properties and fill the named properties map
func (component Component) VerifyAndSetJSONData(data interface{}) (render.Component, render.NamedProperties, error) {
	c := component
	props := make(render.NamedProperties)
	stringStruct, ok := data.(*textFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	// Get named properties and assign each real property
	var newVal interface{}
	var err error
	// Deal with the font restrictions
	inputs := []string{
		stringStruct.Font.FontName,
		stringStruct.Font.FontFile,
		stringStruct.Font.FontURL,
	}
	propNames := []string{
		"fontName",
		"fontFile",
		"fontURL",
	}
	types := []render.PropType{
		render.StringType,
		render.StringType,
		render.StringType,
	}
	var extractedVal interface{}
	validIndex := -1
	c.NamedPropertiesMap, extractedVal, validIndex, err = render.ExtractExclusiveProp(inputs, propNames, types, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	var fName, fFile, fURL interface{}
	switch validIndex {
	case 0:
		fName = extractedVal
	case 1:
		fFile = extractedVal
	case 2:
		fURL = extractedVal
	default:
		return component, props, fmt.Errorf("failed to extract font")
	}
	if fName != nil {
		stringVal := fName.(string)
		pool := gosysfonts.New()
		rawFont, err := pool.GetFont(stringVal)
		if err != nil {
			return component, props, err
		}
		c.Font = rawFont
	}
	if fFile != nil {
		stringVal := fFile.(string)
		if c.reader == nil {
			c.reader = fs.IoutilFileReader{}
		}
		fontData, err := c.reader.ReadFile(stringVal)
		if err != nil {
			return component, props, err
		}
		rawFont, err := truetype.Parse(fontData)
		if err != nil {
			return component, props, err
		}
		c.Font = rawFont
	}
	if fURL != nil {
		return component, props, fmt.Errorf("fontURL not implemented")
	}

	// All other props
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Content, "content", render.StringType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Content = newVal.(string)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.StartX, "startX", render.IntType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Start.X = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.StartY, "startY", render.IntType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Start.Y = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.MaxWidth, "maxWidth", render.IntType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.MaxWidth = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Size, "size", render.Float64Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Size = newVal.(float64)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Alignment, "alignment", render.StringType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		alignmentString := newVal.(string)
		switch alignmentString {
		case "left":
			c.Alignment = AlignmentLeft
		case "right":
			c.Alignment = AlignmentRight
		case "centre":
			c.Alignment = AlignmentCentre
		default:
			c.Alignment = AlignmentLeft
		}
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.PixelsPerInch, "ppi", render.IntType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.PixelsPerInch = newVal.(int)
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
	type invalidStruct struct {
		Message string
	}
	for key := range c.NamedPropertiesMap {
		props[key] = invalidStruct{Message: "Please replace me with real data"}
	}
	return c, props, nil
}