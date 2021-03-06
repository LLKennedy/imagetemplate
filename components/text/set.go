package text

import (
	"fmt"

	"github.com/LLKennedy/imagetemplate/v3/cutils"
)

func (component *Component) delegatedSetProperties(name string, value interface{}) (err error) {
	switch name {
	case "content":
		err = component.setContent(value)
	case "fontName":
		err = component.setFontName(value)
	case "fontFile":
		err = component.setFontFile(value)
	case "fontURL":
		err = component.setFontURL(value)
	case "size":
		component.Size, err = cutils.SetFloat64(value)
	case "alignment":
		err = component.setTextAlignment(value)
	case "R", "G", "B", "A":
		err = component.setColour(name, value)
	case "startX", "startY":
		err = component.setStart(name, value)
	case "maxWidth":
		component.MaxWidth, err = cutils.SetInt(value)
	default:
		err = fmt.Errorf("invalid component property in named property map: %v", name)
	}
	return
}

func (component *Component) setContent(value interface{}) error {
	stringVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("error converting %v to string", value)
	}
	component.Content = stringVal
	return nil
}

func (component *Component) setFontName(value interface{}) error {
	stringVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("error converting %v to string", value)
	}
	rawFont, err := component.getFontPool().GetFont(stringVal)
	if err != nil {
		return err
	}
	component.Font = rawFont
	return nil
}

func (component *Component) setFontFile(value interface{}) (err error) {
	component.Font, err = cutils.LoadFontFile(component.getFileSystem(), value)
	return
}

func (component *Component) setFontURL(value interface{}) error {
	return fmt.Errorf("fontURL not implemented")
}

func (component *Component) setTextAlignment(value interface{}) error {
	alignmentVal, isTextAlignment := value.(cutils.TextAlignment)
	stringVal, isString := value.(string)
	if !isTextAlignment && !isString {
		return fmt.Errorf("could not convert %v to text alignment or string", value)
	}
	if isTextAlignment {
		component.TextAlignment = alignmentVal
		return nil
	}
	component.TextAlignment = cutils.StringToAlignment(stringVal)
	return nil
}

func (component *Component) setColour(name string, value interface{}) error {
	//Process colours
	colourVal, ok := value.(uint8)
	if !ok {
		return fmt.Errorf("error converting %v to uint8", value)
	}
	switch name {
	case "R":
		component.Colour.R = colourVal
	case "G":
		component.Colour.G = colourVal
	case "B":
		component.Colour.B = colourVal
	default:
		component.Colour.A = colourVal
	}
	return nil
}

func (component *Component) setStart(name string, value interface{}) (err error) {
	if name == "startX" {
		component.Start.X, err = cutils.SetInt(value)
		return err
	}
	component.Start.Y, err = cutils.SetInt(value)
	return err
}
