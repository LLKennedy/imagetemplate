package barcode

import (
	"fmt"

	"github.com/LLKennedy/imagetemplate/v3/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/boombuler/barcode/qr"
)

func (component *Component) delegatedSetProperties(name string, value interface{}) (err error) {
	switch name {
	case "content":
		component.Content, err = cutils.SetString(value)
	case "barcodeType":
		stringVal, ok := value.(render.BarcodeType)
		if !ok {
			return fmt.Errorf("error converting %v to barcode type", value)
		}
		component.Extra = render.BarcodeExtraData{}
		switch stringVal {
		case render.BarcodeType2of5:
		case render.BarcodeType2of5Interleaved:
		case render.BarcodeTypeAztec:
			component.Extra.AztecMinECCPercent = 50      //TODO: get a beter value for this, or set it from the file
			component.Extra.AztecUserSpecifiedLayers = 4 //TODO: get a better value for this, or set it from the file
		case render.BarcodeTypeCodabar:
		case render.BarcodeTypeCode128:
		case render.BarcodeTypeCode39:
			component.Extra.Code39IncludeChecksum = true
			component.Extra.Code39FullASCIIMode = true
		case render.BarcodeTypeCode93:
			component.Extra.Code93IncludeChecksum = true
			component.Extra.Code93FullASCIIMode = true
		case render.BarcodeTypeDataMatrix:
		case render.BarcodeTypeEAN13:
		case render.BarcodeTypeEAN8:
		case render.BarcodeTypePDF:
			component.Extra.PDFSecurityLevel = 4 //TODO: get a better value for this, or set it from the file
		case render.BarcodeTypeQR:
			component.Extra.QRLevel = qr.Q
			component.Extra.QRMode = qr.Unicode
		}
		component.Type = stringVal
		return nil
	case "dR", "dG", "dB", "dA", "bR", "bG", "bB", "bA":
		err = component.setColour(name, value)
	case "topLeftX", "topLeftY":
		err = component.setTopLeft(name, value)
	case "width":
		component.Width, err = cutils.SetInt(value)
	case "height":
		component.Height, err = cutils.SetInt(value)
	default:
		return fmt.Errorf("invalid component property in named property map: %v", name)
	}
	return
}

func (component *Component) setColour(name string, value interface{}) error {
	//Process colours
	colourVal, ok := value.(uint8)
	if !ok {
		return fmt.Errorf("error converting %v to uint8", value)
	}
	switch name {
	case "dR":
		component.DataColour.R = colourVal
	case "dG":
		component.DataColour.G = colourVal
	case "dB":
		component.DataColour.B = colourVal
	case "dA":
		component.DataColour.A = colourVal
	case "bR":
		component.BackgroundColour.R = colourVal
	case "bG":
		component.BackgroundColour.G = colourVal
	case "bB":
		component.BackgroundColour.B = colourVal
	case "bA":
		component.BackgroundColour.A = colourVal
	}
	return nil
}

func (component *Component) setTopLeft(name string, value interface{}) (err error) {
	if name == "topLeftX" {
		component.TopLeft.X, err = cutils.SetInt(value)
		return err
	}
	component.TopLeft.Y, err = cutils.SetInt(value)
	return err
}
