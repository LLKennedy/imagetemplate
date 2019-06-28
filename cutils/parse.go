package cutils

import (
	"fmt"
	"image"
	"image/color"

	"github.com/LLKennedy/gosysfonts"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/golang/freetype/truetype"
	"golang.org/x/tools/godoc/vfs"
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

// ParseFontOptions are the optional parameters for ParseFont
type ParseFontOptions struct {
	// Props is the initial named properties map
	Props map[string][]string
	// FileSystem is the vfs FileSystem to use
	FileSystem vfs.FileSystem
	// FontPool is the gosysfonts Pool to use
	FontPool gosysfonts.Pool
}

// ParseFont turns a font name, file path or url into a truetype font
func ParseFont(fontName, fileName, url string, opts ParseFontOptions) (*truetype.Font, map[string][]string, error) {
	if opts.Props == nil {
		opts.Props = map[string][]string{}
	}
	if opts.FileSystem == nil {
		opts.FileSystem = vfs.OS(".")
	}
	if opts.FontPool == nil {
		opts.FontPool = gosysfonts.New()
	}
	propData := []render.PropData{
		{
			InputValue: fontName,
			PropName:   "fontName",
			Type:       render.StringType,
		},
		{
			InputValue: fileName,
			PropName:   "fontFile",
			Type:       render.StringType,
		},
		{
			InputValue: url,
			PropName:   "fontURL",
			Type:       render.StringType,
		},
	}
	props, extractedVal, validIndex, err := render.ExtractExclusiveProp(propData, opts.Props)
	if err != nil {
		return nil, opts.Props, err
	}
	stringVal := extractedVal.(string)
	var font *truetype.Font
	if extractedVal != nil {
		switch validIndex {
		case 0:
			font, err = opts.FontPool.GetFont(stringVal)
		case 1:
			font, err = LoadFontFile(opts.FileSystem, stringVal)
		case 2:
			err = fmt.Errorf("fontURL not implemented")
		}
	}
	return font, props, err
}

// ParsePoint turns an X and Y string into an image.Point
func ParsePoint(x, y, xName, yName string, props map[string][]string) (image.Point, map[string][]string, error) {
	pX, newProps, xErr := ExtractInt(x, xName, props)
	pY, newProps, yErr := ExtractInt(y, yName, newProps)
	if xErr != nil || yErr != nil {
		newProps = props
	}
	return image.Pt(pX, pY), newProps, CombineErrors(xErr, yErr)
}
