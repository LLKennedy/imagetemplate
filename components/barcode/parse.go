package barcode

import (
	"fmt"

	"github.com/LLKennedy/imagetemplate/v3/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
)

func (component Component) parseJSONFormat(stringStruct *barcodeFormat, props render.NamedProperties) (c Component, foundProps render.NamedProperties, err error) {
	c = component
	var parseErr error

	// Get named properties and assign each real property
	var typeString string
	typeString, c.NamedPropertiesMap, err = cutils.ExtractString(stringStruct.Type, "barcodeType", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	if typeString != "" {
		c.Type, err = render.ToBarcodeType(typeString)
		if err != nil {
			return component, props, fmt.Errorf("for barcode type %s: %v", typeString, err)
		}
	}
	c.Content, c.NamedPropertiesMap, parseErr = cutils.ExtractString(stringStruct.Content, "content", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.TopLeft, c.NamedPropertiesMap, parseErr = cutils.ParsePoint(stringStruct.TopLeftX, stringStruct.TopLeftY, "topLeftX", "topLeftY", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.Width, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.Width, "width", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.Height, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.Height, "height", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.DataColour, c.NamedPropertiesMap, parseErr = cutils.ParseColourStrings(cutils.ColourStrings{R: stringStruct.DataColour.Red, G: stringStruct.DataColour.Green, B: stringStruct.DataColour.Blue, A: stringStruct.DataColour.Alpha}, "d", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.BackgroundColour, c.NamedPropertiesMap, parseErr = cutils.ParseColourStrings(cutils.ColourStrings{R: stringStruct.BackgroundColour.Red, G: stringStruct.BackgroundColour.Green, B: stringStruct.BackgroundColour.Blue, A: stringStruct.BackgroundColour.Alpha}, "b", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)

	for key := range c.NamedPropertiesMap {
		props[key] = struct{ Message string }{Message: "Please replace me with real data"}
	}

	// Return original component on error
	if err != nil {
		c = component
	}
	return c, props, err
}
