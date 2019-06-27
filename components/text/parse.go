package text

import (
	"fmt"
	"io/ioutil"

	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/golang/freetype/truetype"
	"golang.org/x/tools/godoc/vfs"
)

func (component Component) parseJSONFormat(stringStruct *textFormat, props render.NamedProperties) (c Component, foundProps render.NamedProperties, err error) {
	c = component
	// Get named properties and assign each real property
	c, err = c.parseFont(stringStruct, err)
	c, err = c.parseContent(stringStruct, err)
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

func (component Component) parseFont(stringStruct *textFormat, history error) (c Component, err error) {
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

func (component Component) parseContent(stringStruct *textFormat, history error) (c Component, err error) {
	err = history
	c = component
	props, newVal, parseErr := render.ExtractSingleProp(stringStruct.Content, "content", render.StringType, c.NamedPropertiesMap)
	if parseErr != nil {
		err = combineErrors(err, parseErr)
		return
	}
	c.NamedPropertiesMap = props
	if newVal != nil {
		c.Content = newVal.(string)
	}
	return
}

func (component Component) parseStart(stringStruct *textFormat, history error) (c Component, err error) {
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

func (component Component) parseMaxWidth(stringStruct *textFormat, history error) (c Component, err error) {
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

func (component Component) parseSize(stringStruct *textFormat, history error) (c Component, err error) {
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

func (component Component) parseAlignment(stringStruct *textFormat, history error) (c Component, err error) {
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

func (component Component) parseColour(stringStruct *textFormat, history error) (c Component, err error) {
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
