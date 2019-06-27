package datetime

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/tools/godoc/vfs"
)

func (component *Component) delegatedSetProperties(name string, value interface{}) (err error) {
	switch name {
	case "time":
		stringVal, isStrings := value.([]string)
		timePointer, isTimePointer := value.(*time.Time)
		timeVal, isTime := value.(time.Time)
		if (!isStrings && !isTimePointer && !isTime) || (isStrings && len(stringVal) != 2) {
			return fmt.Errorf("error converting %v to []string, *time.Time or time.Time", value)
		}
		if isTime {
			component.Time = &timeVal
		}
		if isTimePointer {
			component.Time = timePointer
		}
		if isStrings {
			timeVal, err := time.Parse(stringVal[0], stringVal[1])
			if err != nil {
				return fmt.Errorf("cannot convert time string %v to time format %v", stringVal[1], stringVal[0])
			}
			component.Time = &timeVal
		}
		return nil
	case "timeFormat":
		stringVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("error converting %v to string", value)
		}
		component.TimeFormat = stringVal
		return nil
	case "fontName":
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
	case "fontFile":
		stringVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("error converting %v to string", value)
		}
		if component.fs == nil {
			component.fs = vfs.OS(".")
		}
		fontReader, err := component.fs.Open(stringVal)
		if err != nil {
			return err
		}
		defer fontReader.Close()
		fontData, err := ioutil.ReadAll(fontReader)
		if err != nil {
			return err
		}
		rawFont, err := truetype.Parse(fontData)
		if err != nil {
			return err
		}
		component.Font = rawFont
		return nil
	case "fontURL":
		return fmt.Errorf("fontURL not implemented")
	case "size":
		float64Val, ok := value.(float64)
		if !ok {
			return fmt.Errorf("error converting %v to float64", value)
		}
		component.Size = float64Val
		return nil
	case "alignment":
		alignmentVal, isAlignment := value.(Alignment)
		stringVal, isString := value.(string)
		if !isAlignment && !isString {
			return fmt.Errorf("could not convert %v to datetime alignment or string", value)
		}
		if isAlignment {
			component.Alignment = alignmentVal
			return nil
		}
		switch stringVal {
		case "left":
			component.Alignment = AlignmentLeft
			return nil
		case "right":
			component.Alignment = AlignmentRight
			return nil
		case "centre":
			component.Alignment = AlignmentCentre
			return nil
		default:
			component.Alignment = AlignmentLeft
			return nil
		}
	}
	if strings.Contains("RGBA", name) && len(name) == 1 {
		//Process colours
		colourVal, ok := value.(uint8)
		if !ok {
			return fmt.Errorf("error converting %v to uint8", value)
		}
		switch name {
		case "R":
			component.Colour.R = colourVal
			return nil
		case "G":
			component.Colour.G = colourVal
			return nil
		case "B":
			component.Colour.B = colourVal
			return nil
		case "A":
			component.Colour.A = colourVal
			return nil
		}
	}
	numberVal, ok := value.(int)
	if !ok {
		return fmt.Errorf("error converting %v to int", value)
	}
	switch name {
	case "startX":
		component.Start.X = numberVal
		return nil
	case "startY":
		component.Start.Y = numberVal
		return nil
	case "maxWidth":
		component.MaxWidth = numberVal
		return nil
	default:
		return fmt.Errorf("invalid component property in named property map: %v", name)
	}
}
