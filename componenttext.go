package imagetemplate

import (
	"fmt"
	"github.com/LLKennedy/gosysfonts"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"strings"
)

// TextComponent implements the Component interface for text
type TextComponent struct {
	NamedPropertiesMap map[string][]string
	Content            string
	Start              image.Point
	Size               float64
	MaxWidth           int
	Font               *truetype.Font
	Colour             color.NRGBA
	reader fileReader
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
	fontSize := component.Size
	fits := false
	tries := 0
	var face font.Face
	for !fits && tries < 10 {
		tries++
		face = truetype.NewFace(component.Font, &truetype.Options{Size: fontSize, Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64})
		var realWidth int
		fits, realWidth = c.TryText(component.Content, component.Start, face, component.Colour, component.MaxWidth)
		if realWidth > component.MaxWidth {
			ratio := float64(component.MaxWidth) / float64(realWidth)
			fontSize = ratio * fontSize
		}
	}
	if !fits {
		return canvas, fmt.Errorf("unable to fit text %v into maxWidth %d after %d tries", component.Content, component.MaxWidth, tries)
	}
	c, err := c.Text(component.Content, component.Start, face, component.Colour, component.MaxWidth)
	if err != nil {
		return canvas, err
	}
	return c, nil
}

// SetNamedProperties proceses the named properties and sets them into the text properties
func (component TextComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
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
		case "fontFile":
			stringVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("error converting %v to string", value)
			}
			if component.reader == nil {
				component.reader = ioutilFileReader{}
			}
			ttfBytes, err := component.reader.ReadFile(stringVal)
			if err != nil {
				return err
			}
			rawFont, err := truetype.Parse(ttfBytes)
			if err != nil {
				return err
			}
			c.Font = rawFont
		case "fontURL":
			return fmt.Errorf("fontURL not implemented")
		case "size":
			float64Val, ok := value.(float64)
			if !ok {
				return fmt.Errorf("error converting %v to float64", value)
			}
			c.Size = float64Val
		}
		if strings.Contains("RGBA", name) && len(name) == 1 {
			//Process colours
			colourVal, ok := value.(uint8)
			if !ok {
				return fmt.Errorf("Error converting %v to uint8", value)
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
				return fmt.Errorf("Name was a string inside RGBA and Value was a valid uint8, but Name wasn't R, G, B, or A. Name was: %v", name)
			}
		}
		numberVal, ok := value.(int)
		if !ok {
			return fmt.Errorf("Error converting %v to int", value)
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
		default:
			return fmt.Errorf("Invalid component property in named property map: %v", name)
		}
	}
	var err error
	c.NamedPropertiesMap, err = StandardSetNamedProperties(properties, component.NamedPropertiesMap, setFunc)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a text component
func (component TextComponent) GetJSONFormat() interface{} {
	return &textFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set text properties and fill the named properties map
func (component TextComponent) VerifyAndSetJSONData(data interface{}) (Component, NamedProperties, error) {
	c := component
	var props NamedProperties
	stringStruct, ok := data.(*textFormat)
	if !ok {
		return component, props, fmt.Errorf("Failed to convert returned data to component properties")
	}
	// Get named properties and assign each real property
	var newVal interface{}
	var err error
	// Deal with the font restrictions
	var fName, fFile, fURL interface{}
	c.NamedPropertiesMap, fName, err = extractSingleProp(stringStruct.Font.FontName, "fontName", stringType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	c.NamedPropertiesMap, fFile, err = extractSingleProp(stringStruct.Font.FontFile, "fontFile", stringType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	c.NamedPropertiesMap, fURL, err = extractSingleProp(stringStruct.Font.FontURL, "fontURL", stringType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	trueCount := 0
	if fName != nil {trueCount++}
	if fFile != nil {trueCount++}
	if fURL != nil {trueCount++}
	if trueCount != 1 {
		return component, props, fmt.Errorf("exactly one of fontName, fontFile and fontURL must be set")
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
			c.reader = ioutilFileReader{}
		}
		ttfBytes, err := c.reader.ReadFile(stringVal)
		if err != nil {
			return component, props, err
		}
		rawFont, err := truetype.Parse(ttfBytes)
		if err != nil {
			return component, props, err
		}
		c.Font = rawFont
	}
	if fURL != nil {
		return component, props, fmt.Errorf("fontURL not implemented")
	}

	// All other props
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Content, "content", stringType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Content = newVal.(string)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.StartX, "startX", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Start.X = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.StartY, "startY", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Start.Y = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.MaxWidth, "maxWidth", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.MaxWidth = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Size, "size", float64Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Size = newVal.(float64)
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
	return c, props, nil
}
