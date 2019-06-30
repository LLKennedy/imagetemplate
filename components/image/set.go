package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/LLKennedy/imagetemplate/v3/cutils"
)

func (component *Component) delegatedSetProperties(name string, value interface{}) (err error) {
	switch name {
	case "data":
		err = component.setData(value)
	case "fileName":
		err = component.setFileName(value)
	case "topLeftX":
		component.TopLeft.X, err = cutils.SetInt(value)
	case "topLeftY":
		component.TopLeft.Y, err = cutils.SetInt(value)
	case "width":
		component.Width, err = cutils.SetInt(value)
	case "height":
		component.Height, err = cutils.SetInt(value)
	default:
		err = fmt.Errorf("invalid component property in named property map: %v", name)
	}
	return
}

func (component *Component) setData(value interface{}) error {
	bytesVal, isBytes := value.([]byte)
	stringVal, isString := value.(string)
	readerVal, isReader := value.(io.Reader)
	if !isBytes && !isString && !isReader {
		return fmt.Errorf("error converting %v to []byte, string or io.Reader", value)
	}
	var reader io.Reader
	if isBytes {
		reader = bytes.NewBuffer(bytesVal)
	} else if isString {
		stringReader := strings.NewReader(stringVal)
		reader = base64.NewDecoder(base64.StdEncoding, stringReader)
	} else if isReader {
		reader = readerVal
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}
	component.Image = img
	return nil
}

func (component *Component) setFileName(value interface{}) error {
	stringVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("error converting %v to string", value)
	}
	bytesVal, err := component.getFileSystem().Open(stringVal)
	if err != nil {
		return err
	}
	defer bytesVal.Close()
	img, _, err := image.Decode(bytesVal)
	if err != nil {
		return err
	}
	component.Image = img
	return nil
}
