package text

import (
	"github.com/LLKennedy/imagetemplate/v3/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
)

func (component Component) parseJSONFormat(stringStruct *textFormat, props render.NamedProperties) (c Component, foundProps render.NamedProperties, err error) {
	c = component
	var parseErr error
	// Get named properties and assign each real property
	c.Font, c.NamedPropertiesMap, parseErr = cutils.ParseFont(stringStruct.Font.FontName, stringStruct.Font.FontFile, stringStruct.Font.FontURL, cutils.ParseFontOptions{Props: c.NamedPropertiesMap, FileSystem: c.getFileSystem(), FontPool: c.getFontPool()})
	err = cutils.CombineErrors(err, parseErr)
	c.Content, c.NamedPropertiesMap, parseErr = cutils.ExtractString(stringStruct.Content, "content", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.Start, c.NamedPropertiesMap, parseErr = cutils.ParsePoint(stringStruct.StartX, stringStruct.StartY, "startX", "startY", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.MaxWidth, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.MaxWidth, "maxWidth", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.Size, c.NamedPropertiesMap, parseErr = cutils.ExtractFloat(stringStruct.Size, "size", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.TextAlignment, c.NamedPropertiesMap, parseErr = cutils.ExtractTextAlignment(stringStruct.TextAlignment, "alignment", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
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
