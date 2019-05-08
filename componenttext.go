package imagetemplate

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/LLKennedy/gosysfonts"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// TextComponent implements the Component interface for text
type TextComponent struct {
	NamedPropertiesMap map[string][]string
	Content            string
	Start              image.Point
	Size               float64
	MaxWidth           int
	Alignment          TextAlignment
	PixelsPerInch      int //Should default to 72
	Font               *truetype.Font
	Colour             color.NRGBA
	reader             fileReader
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

// TextAlignment is a text alignment
type TextAlignment int

const (
	// TextAlignmentLeft aligns text left
	TextAlignmentLeft TextAlignment = iota
	// TextAlignmentRight aligns text right
	TextAlignmentRight
	// TextAlignmentCentre aligns text centrally
	TextAlignmentCentre
)

// Write draws text on the canvas
func (component TextComponent) Write(canvas Canvas) (Canvas, error) {
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
			case TextAlignmentLeft:
				alignmentOffset = 0
			case TextAlignmentRight:
				alignmentOffset = int(remainingWidth)
			case TextAlignmentCentre:
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
			return nil
		case "fontFile":
			stringVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("error converting %v to string", value)
			}
			if component.reader == nil {
				component.reader = ioutilFileReader{}
			}
			_, err := component.reader.ReadFile(stringVal) //set this variable again once it's working
			if err != nil {
				return err
			}
			rawFont, err := truetype.Parse(goregular.TTF)
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
			alignmentVal, isAlignment := value.(TextAlignment)
			stringVal, isString := value.(string)
			if !isAlignment && !isString {
				return fmt.Errorf("could not convert %v to text alignment or string", value)
			}
			if isAlignment {
				c.Alignment = alignmentVal
			} else {
				switch stringVal {
				case "left":
					c.Alignment = TextAlignmentLeft
				case "right":
					c.Alignment = TextAlignmentRight
				case "centre":
					c.Alignment = TextAlignmentCentre
				default:
					c.Alignment = TextAlignmentLeft
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
	props := make(NamedProperties)
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
	types := []propType{
		stringType,
		stringType,
		stringType,
	}
	var extractedVal interface{}
	validIndex := -1
	c.NamedPropertiesMap, extractedVal, validIndex, err = extractExclusiveProp(inputs, propNames, types, c.NamedPropertiesMap)
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
			c.reader = ioutilFileReader{}
		}
		_, err := c.reader.ReadFile(stringVal) //set this back once it works
		if err != nil {
			return component, props, err
		}
		rawFont, err := truetype.Parse(goregular.TTF)
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
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Alignment, "alignment", stringType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		alignmentString := newVal.(string)
		switch alignmentString {
		case "left":
			c.Alignment = TextAlignmentLeft
		case "right":
			c.Alignment = TextAlignmentRight
		case "centre":
			c.Alignment = TextAlignmentCentre
		default:
			c.Alignment = TextAlignmentLeft
		}
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.PixelsPerInch, "ppi", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.PixelsPerInch = newVal.(int)
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
	type invalidStruct struct {
		Message string
	}
	for key := range c.NamedPropertiesMap {
		props[key] = invalidStruct{Message: "Please replace me with real data"}
	}
	return c, props, nil
}
