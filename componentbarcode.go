package imagetemplate

import (
	"fmt"
	"github.com/boombuler/barcode/qr"
	"image"
	"image/color"
	"strings"
)

// BarcodeComponent implements the Component interface for images
type BarcodeComponent struct {
	NamedPropertiesMap map[string][]string
	Content            string
	Type               BarcodeType
	TopLeft            image.Point
	Width              int
	Height             int
	DataColour         color.NRGBA
	BackgroundColour   color.NRGBA
	Extra              BarcodeExtraData
}

type barcodeFormat struct {
	Content    string `json:"content"`
	Type       string `json:"barcodeType"`
	TopLeftX   string `json:"topLeftX"`
	TopLeftY   string `json:"topLeftY"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	DataColour struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"dataColour"`
	BackgroundColour struct {
		Red   string `json:"R"`
		Green string `json:"G"`
		Blue  string `json:"B"`
		Alpha string `json:"A"`
	} `json:"backgroundColour"`
}

// Write draws a barcode on the canvas
func (component BarcodeComponent) Write(canvas Canvas) (Canvas, error) {
	c := canvas
	var err error
	c, err = c.Barcode(component.Type, []byte(component.Content), component.Extra, component.TopLeft, component.Width, component.Height, component.DataColour, component.BackgroundColour)
	if err != nil {
		return canvas, err
	}
	return c, nil
}

// SetNamedProperties proceses the named properties and sets them into the barcode properties
func (component BarcodeComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
	c := component
	setFunc := func(name string, value interface{}) error {
		switch name {
		case "content":
			stringVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("error converting %v to string", value)
			}
			c.Content = stringVal
			return nil
		case "barcodeType":
			stringVal, ok := value.(BarcodeType)
			if !ok {
				return fmt.Errorf("error converting %v to barcode type", value)
			}
			c.Extra = BarcodeExtraData{}
			switch stringVal {
			case BarcodeType2of5:
			case BarcodeType2of5Interleaved:
			case BarcodeTypeAztec:
				c.Extra.AztecMinECCPercent = 50      //TODO: get a beter value for this, or set it from the file
				c.Extra.AztecUserSpecifiedLayers = 4 //TODO: get a better value for this, or set it from the file
			case BarcodeTypeCodabar:
			case BarcodeTypeCode128:
			case BarcodeTypeCode39:
				c.Extra.Code39IncludeChecksum = true
				c.Extra.Code39FullASCIIMode = true
			case BarcodeTypeCode93:
				c.Extra.Code93IncludeChecksum = true
				c.Extra.Code93FullASCIIMode = true
			case BarcodeTypeDataMatrix:
			case BarcodeTypeEAN13:
			case BarcodeTypeEAN8:
			case BarcodeTypePDF:
				c.Extra.PDFSecurityLevel = 4 //TODO: get a better value for this, or set it from the file
			case BarcodeTypeQR:
				c.Extra.QRLevel = qr.Q
				c.Extra.QRMode = qr.Unicode
			default:
				return fmt.Errorf("unsupported barcode type %v", stringVal)
			}
			c.Type = stringVal
			return nil
		}
		if strings.Contains("dRdGdBdAbRbGbBbA", name) && len(name) == 2 {
			//Process colours
			colourVal, ok := value.(uint8)
			if !ok {
				return fmt.Errorf("error converting %v to uint8", value)
			}
			switch name {
			case "dR":
				c.DataColour.R = colourVal
				return nil
			case "dG":
				c.DataColour.G = colourVal
				return nil
			case "dB":
				c.DataColour.B = colourVal
				return nil
			case "dA":
				c.DataColour.A = colourVal
				return nil
			case "bR":
				c.BackgroundColour.R = colourVal
				return nil
			case "bG":
				c.BackgroundColour.G = colourVal
				return nil
			case "bB":
				c.BackgroundColour.B = colourVal
				return nil
			case "bA":
				c.BackgroundColour.A = colourVal
				return nil
			default:
				//What? How did you get here?
				return fmt.Errorf("name was a string inside RGBA and Value was a valid uint8, but Name wasn't R, G, B, or A. Name was: %v", name)
			}
		}
		numberVal, ok := value.(int)
		if !ok {
			return fmt.Errorf("error converting %v to int", value)
		}
		switch name {
		case "topLeftX":
			c.TopLeft.X = numberVal
			return nil
		case "topLeftY":
			c.TopLeft.Y = numberVal
			return nil
		case "width":
			c.Width = numberVal
			return nil
		case "height":
			c.Height = numberVal
			return nil
		default:
			return fmt.Errorf("invalid component property in named property map: %v", name)
		}
	}
	var err error
	c.NamedPropertiesMap, err = StandardSetNamedProperties(properties, component.NamedPropertiesMap, setFunc)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a barcode component
func (component BarcodeComponent) GetJSONFormat() interface{} {
	return &barcodeFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set barcode properties and fill the named properties map
func (component BarcodeComponent) VerifyAndSetJSONData(data interface{}) (Component, NamedProperties, error) {
	c := component
	var props NamedProperties
	stringStruct, ok := data.(*barcodeFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	// Get named properties and assign each real property
	var newVal interface{}
	var err error
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Content, "content", stringType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Content = newVal.(string)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.TopLeftX, "topLeftX", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.TopLeft.X = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.TopLeftY, "topLeftY", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.TopLeft.Y = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Width, "width", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Width = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.Height, "height", intType, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.Height = newVal.(int)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.DataColour.Red, "dR", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.DataColour.R = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.DataColour.Green, "dG", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.DataColour.G = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.DataColour.Blue, "dB", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.DataColour.B = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.DataColour.Alpha, "dA", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.DataColour.A = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.BackgroundColour.Red, "bR", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.BackgroundColour.R = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.BackgroundColour.Green, "bG", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.BackgroundColour.G = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.BackgroundColour.Blue, "bB", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.BackgroundColour.B = newVal.(uint8)
	}
	c.NamedPropertiesMap, newVal, err = extractSingleProp(stringStruct.BackgroundColour.Alpha, "bA", uint8Type, c.NamedPropertiesMap)
	if err != nil {
		return component, props, err
	}
	if newVal != nil {
		c.BackgroundColour.A = newVal.(uint8)
	}
	return c, props, nil
}
