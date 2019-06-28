package cutils

import (
	"image/color"

	"github.com/LLKennedy/imagetemplate/v3/render"
)

// ParseColourStrings turns four strings representing a colour channel each into a color.NRGBA struct
func ParseColourStrings(red, green, blue, alpha string, inputProps map[string][]string) (color.NRGBA, map[string][]string, error) {
	colour := color.NRGBA{}
	var err error
	props, newVal, parseErr := render.ExtractSingleProp(red, "R", render.Uint8Type, inputProps)
	if parseErr != nil {
		err = CombineErrors(err, parseErr)
	} else if newVal != nil {
		colour.R = newVal.(uint8)
	}
	props, newVal, parseErr = render.ExtractSingleProp(green, "G", render.Uint8Type, props)
	if parseErr != nil {
		err = CombineErrors(err, parseErr)
	} else if newVal != nil {
		colour.G = newVal.(uint8)
	}
	props, newVal, parseErr = render.ExtractSingleProp(blue, "B", render.Uint8Type, props)
	if parseErr != nil {
		err = CombineErrors(err, parseErr)
	} else if newVal != nil {
		colour.B = newVal.(uint8)
	}
	props, newVal, parseErr = render.ExtractSingleProp(alpha, "A", render.Uint8Type, props)
	if parseErr != nil {
		err = CombineErrors(err, parseErr)
	} else if newVal != nil {
		colour.A = newVal.(uint8)
	}
	return colour, props, err
}
