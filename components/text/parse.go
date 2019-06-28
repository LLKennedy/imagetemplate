package text

import (
	"github.com/LLKennedy/imagetemplate/v3/internal/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
)

func (component Component) parseJSONFormat(stringStruct *textFormat, props render.NamedProperties) (c Component, foundProps render.NamedProperties, err error) {
	c = component
	var parseErr error
	// Get named properties and assign each real property
	c.Font, c.NamedPropertiesMap, parseErr = cutils.ParseFont(stringStruct.Font.FontName, stringStruct.Font.FontFile, stringStruct.Font.FontURL, cutils.ParseFontOptions{Props: c.NamedPropertiesMap, FileSystem: c.getFileSystem(), FontPool: c.getFontPool()})
	err = cutils.CombineErrors(err, parseErr)
	c, err = c.parseContent(stringStruct, err)
	c.Start, c.NamedPropertiesMap, parseErr = cutils.ParsePoint(stringStruct.StartX, stringStruct.StartY, "startX", "startY", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c, err = c.parseMaxWidth(stringStruct, err)
	c, err = c.parseSize(stringStruct, err)
	c, err = c.parseAlignment(stringStruct, err)
	c.Colour, c.NamedPropertiesMap, parseErr = cutils.ParseColourStrings(stringStruct.Colour.Red, stringStruct.Colour.Green, stringStruct.Colour.Blue, stringStruct.Colour.Alpha, c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
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

func (component Component) parseContent(stringStruct *textFormat, history error) (c Component, err error) {
	err = history
	c = component
	var parseErr error
	c.Content, c.NamedPropertiesMap, parseErr = cutils.ExtractString(stringStruct.Content, "content", c.NamedPropertiesMap)
	if parseErr != nil {
		err = cutils.CombineErrors(err, parseErr)
	}
	return
}

func (component Component) parseMaxWidth(stringStruct *textFormat, history error) (c Component, err error) {
	c = component
	c.MaxWidth, c.NamedPropertiesMap, err = cutils.ExtractInt(stringStruct.MaxWidth, "maxWidth", c.NamedPropertiesMap)
	err = cutils.CombineErrors(history, err)
	return
}

func (component Component) parseSize(stringStruct *textFormat, history error) (c Component, err error) {
	err = history
	c = component
	props, newVal, parseErr := render.ExtractSingleProp(stringStruct.Size, "size", render.Float64Type, c.NamedPropertiesMap)
	if parseErr != nil {
		err = cutils.CombineErrors(err, parseErr)
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
	var alignmentString string
	var parseErr error
	alignmentString, c.NamedPropertiesMap, parseErr = cutils.ExtractString(stringStruct.Alignment, "alignment", c.NamedPropertiesMap)
	if parseErr != nil {
		err = cutils.CombineErrors(err, parseErr)
		return
	}
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
	return
}
