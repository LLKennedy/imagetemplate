package datetime

import (
	"time"

	"github.com/LLKennedy/imagetemplate/v3/internal/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
)

func (component Component) parseJSONFormat(stringStruct *datetimeFormat, startTime time.Time, props render.NamedProperties) (c Component, foundProps render.NamedProperties, err error) {
	c = component
	var parseErr error
	// Get named properties and assign each real property
	c.Font, c.NamedPropertiesMap, parseErr = cutils.ParseFont(stringStruct.Font.FontName, stringStruct.Font.FontFile, stringStruct.Font.FontURL, cutils.ParseFontOptions{Props: c.NamedPropertiesMap, FileSystem: c.getFileSystem(), FontPool: c.getFontPool()})
	err = cutils.CombineErrors(err, parseErr)
	c, err = c.parseTime(stringStruct, startTime, err)
	c.Start, c.NamedPropertiesMap, parseErr = cutils.ParsePoint(stringStruct.StartX, stringStruct.StartY, "startX", "startY", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.MaxWidth, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.MaxWidth, "maxWidth", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.Size, c.NamedPropertiesMap, parseErr = cutils.ExtractFloat(stringStruct.Size, "size", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
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

func (component Component) parseTime(stringStruct *datetimeFormat, startTime time.Time, history error) (c Component, err error) {
	err = history
	c = component
	// TODO: rewrite this logic to handle standalone time vs passed in time vs passed in string time vs hard-coded string time etc.
	props, newVal, parseErr := render.ExtractSingleProp(stringStruct.Time, "time", render.TimeType, c.NamedPropertiesMap)
	if parseErr != nil {
		err = cutils.CombineErrors(err, parseErr)
	} else {
		c.NamedPropertiesMap = props
		if newVal != nil {
			timeVal := startTime.Add(newVal.(time.Duration))
			c.Time = &timeVal
		}
	}
	c.TimeFormat, c.NamedPropertiesMap, parseErr = cutils.ExtractString(stringStruct.TimeFormat, "timeFormat", c.NamedPropertiesMap)
	if parseErr != nil {
		err = cutils.CombineErrors(err, parseErr)
	}
	return
}

func (component Component) parseAlignment(stringStruct *datetimeFormat, history error) (c Component, err error) {
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
