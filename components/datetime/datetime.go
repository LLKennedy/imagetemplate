// Package datetime is a text-based time component with customisable content, size, colour, location and time format.
package datetime

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"runtime/debug"
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
	/*
		NamedPropertiesMap maps user/application variables to properties of the component.
		This field is filled automatically by VerifyAndSetJSONData, then used in
		SetNamedProperties to determine whether a variable being passed in is relevant to this
		component.

		For example, map[string][]string{"expiry": []string{"time"}} would indicate that
		the user specified variable "expiry" will fill the Time property.
	*/
	NamedPropertiesMap map[string][]string
	// Time is the timestamp to render.
	Time *time.Time
	// TimeFormat is the format with which to parse a string-based time input.
	TimeFormat string
	// Start is the coordinates of the dot relative to the top-left corner of the canvas.
	Start image.Point
	// Size is the size of the text in points.
	Size float64
	// MaxWidth is the maximum number of horizontal pixels the dot can move before scaling text.
	MaxWidth int
	// Alignment aligns text to the left, right or centre.
	Alignment Alignment
	// Font is the typeface to use.
	Font *truetype.Font
	// Colour is the colour of the text.
	Colour color.NRGBA
	// fs is the file system.
	fs vfs.FileSystem
	// fontPool is the pool of available fonts.
	fontPool gosysfonts.Pool
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
			err = fmt.Errorf("failed to write to canvas: %v\n%s", p, debug.Stack())
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
			if (!isStrings && !isTimePointer && !isTime) || (isStrings && len(stringVal) != 2) {
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
			rawFont, err := c.getFontPool().GetFont(stringVal)
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
			}
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
	startTime := time.Now()
	c := component
	props := make(render.NamedProperties)
	stringStruct, ok := data.(*datetimeFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	return c.parseJSONFormat(stringStruct, startTime, props)
}

func (component Component) parseJSONFormat(stringStruct *datetimeFormat, startTime time.Time, props render.NamedProperties) (c Component, foundProps render.NamedProperties, err error) {
	c = component
	// Get named properties and assign each real property
	c, err = c.parseFont(stringStruct, err)
	c, err = c.parseTime(stringStruct, startTime, err)
	c, err = c.parseStart(stringStruct, err)
	c, err = c.parseMaxWidth(stringStruct, err)
	c, err = c.parseSize(stringStruct, err)
	c, err = c.parseAlignment(stringStruct, err)
	c, err = c.parseColour(stringStruct, err)

	// Fill discovered properties with real data
	for key := range c.NamedPropertiesMap {
		props[key] = struct {
			Message string
		}{Message: "Please replace me with real data"}
	}

	// Return original component on error
	if err != nil {
		c = component
	}
	return c, props, err
}

func combineErrors(history error, latest error) error {
	if history == nil {
		return latest
	}
	return fmt.Errorf("%v\n%v", history, latest)
}

func (component Component) parseFont(stringStruct *datetimeFormat, history error) (c Component, err error) {
	err = history
	c = component
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
	var parseErr error
	c.NamedPropertiesMap, extractedVal, validIndex, parseErr = render.ExtractExclusiveProp(propData, component.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	if extractedVal != nil {
		switch validIndex {
		case 0:
			c, err = c.parseFontName(extractedVal.(string), history)
		case 1:
			c, err = c.parseFontFile(extractedVal.(string), history)
		case 2:
			c, err = c.parseFontURL(extractedVal.(string), history)
		}
	}
	return
}

func (component Component) parseFontName(name string, history error) (c Component, err error) {
	err = history
	c = component
	rawFont, parseErr := c.getFontPool().GetFont(name)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	c.Font = rawFont
	return
}

func (component Component) parseFontURL(url string, history error) (c Component, err error) {
	err = history
	c = component
	err = combineErrors(err, fmt.Errorf("fontURL not implemented"))
	return
}

func (component Component) parseFontFile(path string, history error) (c Component, err error) {
	err = history
	c = component
	if c.fs == nil {
		c.fs = vfs.OS(".")
	}
	fontReader, parseErr := c.fs.Open(path)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	defer fontReader.Close()
	fontData, parseErr := ioutil.ReadAll(fontReader)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	rawFont, parseErr := truetype.Parse(fontData)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	c.Font = rawFont
	return
}

func (component Component) getFontPool() gosysfonts.Pool {
	if component.fontPool == nil {
		return gosysfonts.New()
	}
	return component.fontPool
}

func (component Component) parseTime(stringStruct *datetimeFormat, startTime time.Time, history error) (c Component, err error) {
	err = history
	c = component
	// TODO: rewrite this logic to handle standalone time vs passed in time vs passed in string time vs hard-coded string time etc.
	props, newVal, parseErr := render.ExtractSingleProp(stringStruct.Time, "time", render.TimeType, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
	} else {
		c.NamedPropertiesMap = props
		if newVal != nil {
			timeVal := startTime.Add(newVal.(time.Duration))
			c.Time = &timeVal
		}
	}
	props, newVal, parseErr = render.ExtractSingleProp(stringStruct.TimeFormat, "timeFormat", render.StringType, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	c.NamedPropertiesMap = props
	if newVal != nil {
		c.TimeFormat = newVal.(string)
	}
	return
}

func (component Component) parseStart(stringStruct *datetimeFormat, history error) (c Component, err error) {
	err = history
	c = component
	props, newVal, parseErr := render.ExtractSingleProp(stringStruct.StartX, "startX", render.IntType, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
	} else {
		c.NamedPropertiesMap = props
		if newVal != nil {
			c.Start.X = newVal.(int)
		}
	}
	props, newVal, parseErr = render.ExtractSingleProp(stringStruct.StartY, "startY", render.IntType, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	c.NamedPropertiesMap = props
	if newVal != nil {
		c.Start.Y = newVal.(int)
	}
	return
}

func (component Component) parseMaxWidth(stringStruct *datetimeFormat, history error) (c Component, err error) {
	err = history
	c = component
	props, newVal, parseErr := render.ExtractSingleProp(stringStruct.MaxWidth, "maxWidth", render.IntType, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	c.NamedPropertiesMap = props
	if newVal != nil {
		c.MaxWidth = newVal.(int)
	}
	return
}

func (component Component) parseSize(stringStruct *datetimeFormat, history error) (c Component, err error) {
	err = history
	c = component
	props, newVal, parseErr := render.ExtractSingleProp(stringStruct.Size, "size", render.Float64Type, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	c.NamedPropertiesMap = props
	if newVal != nil {
		c.Size = newVal.(float64)
	}
	return
}

func (component Component) parseAlignment(stringStruct *datetimeFormat, history error) (c Component, err error) {
	err = history
	c = component
	props, newVal, parseErr := render.ExtractSingleProp(stringStruct.Alignment, "alignment", render.StringType, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	c.NamedPropertiesMap = props
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
	return
}

func (component Component) parseColour(stringStruct *datetimeFormat, history error) (c Component, err error) {
	err = history
	c = component
	props, newVal, parseErr := render.ExtractSingleProp(stringStruct.Colour.Red, "R", render.Uint8Type, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
	} else {
		c.NamedPropertiesMap = props
		if newVal != nil {
			c.Colour.R = newVal.(uint8)
		}
	}
	props, newVal, parseErr = render.ExtractSingleProp(stringStruct.Colour.Green, "G", render.Uint8Type, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
	} else {
		c.NamedPropertiesMap = props
		if newVal != nil {
			c.Colour.G = newVal.(uint8)
		}
	}
	props, newVal, parseErr = render.ExtractSingleProp(stringStruct.Colour.Blue, "B", render.Uint8Type, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
	} else {
		c.NamedPropertiesMap = props
		if newVal != nil {
			c.Colour.B = newVal.(uint8)
		}
	}
	props, newVal, parseErr = render.ExtractSingleProp(stringStruct.Colour.Alpha, "A", render.Uint8Type, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
	} else {
		c.NamedPropertiesMap = props
		if newVal != nil {
			c.Colour.A = newVal.(uint8)
		}
	}
	return
}

func init() {
	for _, name := range []string{"datetime", "DateTime", "DATETIME", "Datetime", "Date/Time", "date/time", "date", "DATE", "Date"} {
		render.RegisterComponent(name, func(fs vfs.FileSystem) render.Component { return Component{fs: fs, fontPool: gosysfonts.New()} })
	}
}
