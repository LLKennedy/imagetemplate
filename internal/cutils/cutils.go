// Package cutils provides common parsing/conversion code for components to cut down on duplication
package cutils

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/tools/godoc/vfs"
)

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
