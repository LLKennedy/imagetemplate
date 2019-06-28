// Package cutils provides common parsing/conversion code for components to cut down on duplication
package cutils

import (
	"fmt"
	"image/color"
	"io/ioutil"

	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/golang/freetype/truetype"
	"golang.org/x/tools/godoc/vfs"
)

// CombineErrors combines two errors, maintaining a history of errors separated by newlines
func CombineErrors(history, latest error) error {
	switch {
	case history == nil:
		return latest
	case latest == nil:
		return history
	}
	return fmt.Errorf("%v\n%v", history, latest)
}

// LoadFontFile returns the font file found at the specified path on the specified file system
func LoadFontFile(fs vfs.FileSystem, fileName interface{}) (*truetype.Font, error) {
	path, ok := fileName.(string)
	if !ok {
		return nil, fmt.Errorf("error converting %v to string", fileName)
	}
	fontReader, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer fontReader.Close()
	fontData, err := ioutil.ReadAll(fontReader)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(fontData)
}

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
