package datetime

import (
	"time"

	"github.com/LLKennedy/imagetemplate/v3/internal/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
)

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

func (component Component) parseFont(stringStruct *datetimeFormat, history error) (c Component, err error) {
	c = component
	c.Font, c.NamedPropertiesMap, err = cutils.ParseFont(stringStruct.Font.FontName, stringStruct.Font.FontFile, stringStruct.Font.FontURL, cutils.ParseFontOptions{Props: c.NamedPropertiesMap, FileSystem: c.getFileSystem(), FontPool: c.getFontPool()})
	err = cutils.CombineErrors(history, err)
	return
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

func (component Component) parseStart(stringStruct *datetimeFormat, history error) (c Component, err error) {
	c = component
	var parseErr error
	c.Start.X, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.StartX, "startX", c.NamedPropertiesMap)
	err = cutils.CombineErrors(history, parseErr)
	c.Start.Y, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.StartY, "startY", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	return
}

func (component Component) parseMaxWidth(stringStruct *datetimeFormat, history error) (c Component, err error) {
	c = component
	c.MaxWidth, c.NamedPropertiesMap, err = cutils.ExtractInt(stringStruct.MaxWidth, "maxWidth", c.NamedPropertiesMap)
	err = cutils.CombineErrors(history, err)
	return
}

func (component Component) parseSize(stringStruct *datetimeFormat, history error) (c Component, err error) {
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

func (component Component) parseColour(stringStruct *datetimeFormat, history error) (c Component, err error) {
	err = history
	c = component
	var parseErr error
	c.Colour, c.NamedPropertiesMap, parseErr = cutils.ParseColourStrings(stringStruct.Colour.Red, stringStruct.Colour.Green, stringStruct.Colour.Blue, stringStruct.Colour.Alpha, c.NamedPropertiesMap)
	if parseErr != nil {
		err = cutils.CombineErrors(history, parseErr)
	}
	return
}
