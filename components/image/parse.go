package image

import (
	"encoding/base64"
	"image"
	"strings"

	"github.com/LLKennedy/imagetemplate/v3/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"golang.org/x/tools/godoc/vfs"
)

func (component Component) parseJSONFormat(stringStruct *imageFormat, props render.NamedProperties) (c Component, foundProps render.NamedProperties, err error) {
	c = component
	var parseErr error
	// Deal with the file/data restrictions
	c, parseErr = c.parseImageFile(stringStruct.FileName, stringStruct.Data, props)
	err = cutils.CombineErrors(err, parseErr)
	c.TopLeft.X, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.TopLeftX, "topLeftX", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.TopLeft.Y, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.TopLeftY, "topLeftY", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.Width, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.Width, "width", c.NamedPropertiesMap)
	err = cutils.CombineErrors(err, parseErr)
	c.Height, c.NamedPropertiesMap, parseErr = cutils.ExtractInt(stringStruct.Height, "height", c.NamedPropertiesMap)
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

func (component Component) parseImageFile(filename, data string, props render.NamedProperties) (c Component, err error) {
	c = component
	propData := []render.PropData{
		{
			InputValue: filename,
			PropName:   "fileName",
			Type:       render.StringType,
		},
		{
			InputValue: data,
			PropName:   "data",
			Type:       render.StringType,
		},
	}
	var extractedVal interface{}
	var validIndex int
	c.NamedPropertiesMap, extractedVal, validIndex, err = render.ExtractExclusiveProp(propData, c.NamedPropertiesMap)
	if err != nil {
		c = component
	} else if extractedVal != nil {
		switch validIndex {
		case 0:
			stringVal := extractedVal.(string)
			if c.fs == nil {
				c.fs = vfs.OS(".")
			}
			bytesVal, err := c.fs.Open(stringVal)
			if err != nil {
				return component, err
			}
			defer bytesVal.Close()
			img, _, err := image.Decode(bytesVal)
			if err != nil {
				return component, err
			}
			c.Image = img
		case 1:
			base64Val := extractedVal.(string)
			r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Val))
			img, _, err := image.Decode(r)
			if err != nil {
				return component, err
			}
			c.Image = img
		}
	}
	return
}
