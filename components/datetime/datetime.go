package datetime

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"strings"
	"time"

	"github.com/LLKennedy/gosysfonts"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/tools/godoc/vfs"
)

// Component implements the Component interface for datetime
type Component struct {
	NamedPropertiesMap map[string][]string
	Time               *time.Time
	TimeFormat         string
	Start              image.Point
	Size               float64
	MaxWidth           int
	Alignment          Alignment
	Font               *truetype.Font
	Colour             color.NRGBA
	fs                 vfs.FileSystem
}

type datetimeFormat struct {
	Time       string       `json:"time"`
	TimeFormat string       `json:"timeFormat"`
	StartX     string       `json:"startX"`
	StartY     string       `json:"startY"`
	Size       string       `json:"size"`
	MaxWidth   string       `json:"maxWidth"`
	Alignment  string       `json:"alignment"`
	Font       fontFormat   `json:"font"`
	Colour     colourFormat `json:"colour"`
}

type fontFormat struct {
	FontName string `json:"fontName"`
	FontFile string `json:"fontFile"`
	FontURL  string `json:"fontURL"`
}

type colourFormat struct {
	Red   string `json:"R"`
	Green string `json:"G"`
	Blue  string `json:"B"`
	Alpha string `json:"A"`
}

// Alignment is a datetime alignment
type Alignment int

const (
	// AlignmentLeft aligns datetime left
	AlignmentLeft Alignment = iota
	// AlignmentRight aligns datetime right
	AlignmentRight
	// AlignmentCentre aligns datetime centrally
	AlignmentCentre
)

// Write draws datetime on the canvas
func (component Component) Write(canvas render.Canvas) (c render.Canvas, err error) {
	c = canvas
	defer func() {
		p := recover()
		if p != nil {
			err = fmt.Errorf("failed to write to canvas: %v", p)
		}
	}()
	fontSize := component.Size
	formattedTime := component.Time.Format(component.TimeFormat)
	fits := false
	tries := 0
	var face font.Face
	var alignmentOffset int
	for !fits && tries < 10 {
		fmt.Printf("new fontsize: %f", fontSize)
		tries++
		face = truetype.NewFace(component.Font, &truetype.Options{Size: fontSize, Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: canvas.GetPPI()})
		var realWidth int
		fits, realWidth = c.TryText(formattedTime, component.Start, face, component.Colour, component.MaxWidth)
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
		return canvas, fmt.Errorf("unable to fit datetime %s into maxWidth %d after %d tries", formattedTime, component.MaxWidth, tries)
	}
	c, err = c.Text(formattedTime, image.Pt(component.Start.X+alignmentOffset, component.Start.Y), face, component.Colour, component.MaxWidth)
	if err != nil {
		return canvas, err
	}
	return c, nil
}

// SetNamedProperties processes the named properties and sets them into the datetime properties
func (component Component) SetNamedProperties(properties render.NamedProperties) (render.Component, error) {
	c := component
	setFunc := func(name string, value interface{}) error {
		switch name {
		case "time":
			stringVal, isStrings := value.([]string)
			timePointer, isTimePointer := value.(*time.Time)
			timeVal, isTime := value.(time.Time)
			if (!isStrings && !isTimePointer && !isTime) || (isTime && len(stringVal) != 2) {
				return fmt.Errorf("error converting %v to []string, *time.Time or time.Time", value)
			}
			if isTime {
				c.Time = &timeVal
			}
			if isTimePointer {
				c.Time = timePointer
			}
			if isStrings {
				timeVal, err := time.Parse(stringVal[0], stringVal[1])
				if err != nil {
					return fmt.Errorf("cannot convert time string %v to time format %v", stringVal[1], stringVal[0])
				}
				c.Time = &timeVal
			}
			return nil
		case "timeFormat":
			stringVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("error converting %v to string", value)
			}
			c.TimeFormat = stringVal
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
			if component.fs == nil {
				component.fs = vfs.OS(".")
			}
			fontReader, err := component.fs.Open(stringVal)
			if err != nil {
				return err
			}
			defer fontReader.Close()
			fontData, err := ioutil.ReadAll(fontReader)
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
				return fmt.Errorf("could not convert %v to datetime alignment or string", value)
			}
			if isAlignment {
				c.Alignment = alignmentVal
				return nil
			} else {
				switch stringVal {
				case "left":
					c.Alignment = AlignmentLeft
					return nil
				case "right":
					c.Alignment = AlignmentRight
					return nil
				case "centre":
					c.Alignment = AlignmentCentre
					return nil
				default:
					c.Alignment = AlignmentLeft
					return nil
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

// GetJSONFormat returns the JSON structure of a datetime component
func (component Component) GetJSONFormat() interface{} {
	return &datetimeFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set datetime properties and fill the named properties map
func (component Component) VerifyAndSetJSONData(data interface{}) (render.Component, render.NamedProperties, error) {
	c := component
	props := make(render.NamedProperties)
	stringStruct, ok := data.(*datetimeFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	// Get named properties and assign each real property
	var newVal interface{}
	var err error
	// Deal with the font restrictions
	propData := []render.PropData{
		{
			InputValue: stringStruct.Font.FontName,
			PropName:   "fontName",
			Type:       render.StringType,
		},
		{
			InputValue: stringStruct.Font.FontFile,
			PropName:   "fontFile",
			Type:       render.StringType,
		},
		{
			InputValue: stringStruct.Font.FontURL,
			PropName:   "fontURL",
			Type:       render.StringType,
		},
	}
	var extractedVal interface{}
	var validIndex int
	c.NamedPropertiesMap, extractedVal, validIndex, err = render.ExtractExclusiveProp(propData, c.NamedPropertiesMap)
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
		if c.fs == nil {
			c.fs = vfs.OS(".")
		}
		fontReader, err := c.fs.Open(stringVal)
		if err != nil {
			return component, props, err
		}
		defer fontReader.Close()
		fontData, err := ioutil.ReadAll(fontReader)
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
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.Time, "time", render.TimeType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Time = newVal.(*time.Time)
	}
	c.NamedPropertiesMap, newVal, err = render.ExtractSingleProp(stringStruct.TimeFormat, "timeFormat", render.StringType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.TimeFormat = newVal.(string)
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
	for _, name := range []string{"datetime", "DateTime", "DATETIME", "Datetime", "Date/Time", "date/time", "date", "DATE", "Date"} {
		render.RegisterComponent(name, func(fs vfs.FileSystem) render.Component { return Component{fs: fs} })
	}
}
