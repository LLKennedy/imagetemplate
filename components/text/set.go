package text

import (
	"fmt"

	"github.com/LLKennedy/imagetemplate/v3/internal/cutils"
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
		err = component.setSize(value)
	case "alignment":
		err = component.setAlignment(value)
	case "R", "G", "B", "A":
		err = component.setColour(name, value)
	case "startX", "startY":
		err = component.setStart(name, value)
	case "maxWidth":
		err = component.setMaxWidth(value)
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

func (component *Component) setSize(value interface{}) error {
	float64Val, ok := value.(float64)
	if !ok {
		return fmt.Errorf("error converting %v to float64", value)
	}
	component.Size = float64Val
	return nil
}

func (component *Component) setAlignment(value interface{}) error {
	alignmentVal, isAlignment := value.(Alignment)
	stringVal, isString := value.(string)
	if !isAlignment && !isString {
		return fmt.Errorf("could not convert %v to text alignment or string", value)
	}
	if isAlignment {
		component.Alignment = alignmentVal
		return nil
	}
	switch stringVal {
	case "left":
		component.Alignment = AlignmentLeft
	case "right":
		component.Alignment = AlignmentRight
	case "centre":
		component.Alignment = AlignmentCentre
	default:
		component.Alignment = AlignmentLeft
	}
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

func (component *Component) setStart(name string, value interface{}) error {
	if name == "startX" {
		numberVal, ok := value.(int)
		if !ok {
			return fmt.Errorf("error converting %v to int", value)
		}
		component.Start.X = numberVal
		return nil
	}
	numberVal, ok := value.(int)
	if !ok {
		return fmt.Errorf("error converting %v to int", value)
	}
	component.Start.Y = numberVal
	return nil
}

func (component *Component) setMaxWidth(value interface{}) error {
	numberVal, ok := value.(int)
	if !ok {
		return fmt.Errorf("error converting %v to int", value)
	}
	component.MaxWidth = numberVal
	return nil
}
