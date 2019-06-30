package cutils

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"time"

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

// SetString turns an interface into a string and an error
func SetString(value interface{}) (string, error) {
	str, err := setType(value, setStringType)
	return str.(string), err
}

// SetInt turns an interface into a string and an error
func SetInt(value interface{}) (int, error) {
	i, err := setType(value, setIntType)
	return i.(int), err
}

// SetUint8 turns an interface into a uint8 and an error
func SetUint8(value interface{}) (uint8, error) {
	u, err := setType(value, setUint8Type)
	return u.(uint8), err
}

// SetFloat64 turns an interface into a float64 and an error
func SetFloat64(value interface{}) (float64, error) {
	f, err := setType(value, setFloat64Type)
	return f.(float64), err
}

// SetBool turns an interface into a bool and an error
func SetBool(value interface{}) (bool, error) {
	b, err := setType(value, setBoolType)
	return b.(bool), err
}

// SetTime turns an interface into a time and an error
func SetTime(value interface{}) (time.Time, error) {
	t, err := setType(value, setTimeType)
	return t.(time.Time), err
}

// SetTimePointer turns an interface into a time pointer and an error
func SetTimePointer(value interface{}) (*time.Time, error) {
	t, err := setType(value, setTimePointType)
	return t.(*time.Time), err
}

// SetImage turns an interface into an image and an error
func SetImage(value interface{}) (image.Image, error) {
	i, err := setType(value, setImageType)
	var ci image.Image
	ci, _ = i.(image.Image)
	return ci, err
}

// SetColour turns an interface into a colour and an error
func SetColour(value interface{}) (color.Color, error) {
	c, err := setType(value, setColourType)
	var cc color.Color
	cc, _ = c.(color.Color)
	return cc, err
}

type setTypes string

const (
	setStringType    setTypes = "string"
	setIntType       setTypes = "int"
	setUint8Type     setTypes = "uint8"
	setFloat64Type   setTypes = "float64"
	setBoolType      setTypes = "bool"
	setTimeType      setTypes = "time"
	setTimePointType setTypes = "time pointer"
	setImageType     setTypes = "image"
	setColourType    setTypes = "colour"
)

func setType(value interface{}, typeName setTypes) (converted interface{}, err error) {
	var ok bool
	switch typeName {
	case setStringType:
		if converted, ok = value.(string); !ok {
			converted = ""
		}
	case setIntType:
		if converted, ok = value.(int); !ok {
			converted = 0
		}
	case setUint8Type:
		if converted, ok = value.(uint8); !ok {
			converted = uint8(0)
		}
	case setFloat64Type:
		if converted, ok = value.(float64); !ok {
			converted = float64(0)
		}
	case setBoolType:
		if converted, ok = value.(bool); !ok {
			converted = false
		}
	case setTimeType:
		if converted, ok = value.(time.Time); !ok {
			converted = time.Time{}
		}
	case setTimePointType:
		converted, ok = value.(*time.Time)
	case setImageType:
		if converted, ok = value.(image.Image); !ok {
			var img image.Image
			converted = img
		}
	case setColourType:
		if converted, ok = value.(color.Color); !ok {
			var clr color.Color
			converted = clr
		}
	}
	if !ok {
		err = fmt.Errorf("error converting %v to %s", value, typeName)
	}
	return
}
