package render

import (
	"errors"
	"image"
	"image/color"
	"image/draw"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/aztec"
	"github.com/boombuler/barcode/codabar"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/code93"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/pdf417"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/twooffive"
)

// BarcodeType wraps the barcode types into a single enum.
type BarcodeType string

const (
	// BarcodeTypeAztec           is an alias for an imported barcode type
	BarcodeTypeAztec BarcodeType = barcode.TypeAztec
	// BarcodeTypeCodabar         is an alias for an imported barcode type
	BarcodeTypeCodabar BarcodeType = barcode.TypeCodabar
	// BarcodeTypeCode128         is an alias for an imported barcode type
	BarcodeTypeCode128 BarcodeType = barcode.TypeCode128
	// BarcodeTypeCode39          is an alias for an imported barcode type
	BarcodeTypeCode39 BarcodeType = barcode.TypeCode39
	// BarcodeTypeCode93          is an alias for an imported barcode type
	BarcodeTypeCode93 BarcodeType = barcode.TypeCode93
	// BarcodeTypeDataMatrix      is an alias for an imported barcode type
	BarcodeTypeDataMatrix BarcodeType = barcode.TypeDataMatrix
	// BarcodeTypeEAN8            is an alias for an imported barcode type
	BarcodeTypeEAN8 BarcodeType = barcode.TypeEAN8
	// BarcodeTypeEAN13           is an alias for an imported barcode type
	BarcodeTypeEAN13 BarcodeType = barcode.TypeEAN13
	// BarcodeTypePDF             is an alias for an imported barcode type
	BarcodeTypePDF BarcodeType = barcode.TypePDF
	// BarcodeTypeQR              is an alias for an imported barcode type
	BarcodeTypeQR BarcodeType = barcode.TypeQR
	// BarcodeType2of5            is an alias for an imported barcode type
	BarcodeType2of5 BarcodeType = barcode.Type2of5
	// BarcodeType2of5Interleaved is an alias for an imported barcode type
	BarcodeType2of5Interleaved BarcodeType = barcode.Type2of5Interleaved
)

// ToBarcodeType attempts to convert a barcode type string to a defined BarcodeType constant.
func ToBarcodeType(raw string) (BarcodeType, error) {
	switch raw {
	case string(BarcodeTypeAztec):
		return BarcodeTypeAztec, nil
	case string(BarcodeTypeCodabar):
		return BarcodeTypeCodabar, nil
	case string(BarcodeTypeCode128):
		return BarcodeTypeCode128, nil
	case string(BarcodeTypeCode39):
		return BarcodeTypeCode39, nil
	case string(BarcodeTypeCode93):
		return BarcodeTypeCode93, nil
	case string(BarcodeTypeDataMatrix):
		return BarcodeTypeDataMatrix, nil
	case string(BarcodeTypeEAN8):
		return BarcodeTypeEAN8, nil
	case string(BarcodeTypeEAN13):
		return BarcodeTypeEAN13, nil
	case string(BarcodeTypePDF):
		return BarcodeTypePDF, nil
	case string(BarcodeTypeQR):
		return BarcodeTypeQR, nil
	case string(BarcodeType2of5):
		return BarcodeType2of5, nil
	case string(BarcodeType2of5Interleaved):
		return BarcodeType2of5Interleaved, nil
	default:
		return BarcodeType(""), errors.New("barcode type does not match defined constants")
	}
}

// BarcodeExtraData contains additional data required for some barcode formats, leave any fields not named for the type in use alone.
type BarcodeExtraData struct {
	// AztecMinECCPercent is required for aztec barcodes
	AztecMinECCPercent int
	// AztecUserSpecifiedLayers is required for aztec barcodes
	AztecUserSpecifiedLayers int
	// Code39IncludeChecksum is required for code39 barcodes
	Code39IncludeChecksum bool
	// Code39FullASCIIMode is required for code39 barcodes
	Code39FullASCIIMode bool
	// Code93IncludeChecksum is required for code93 barcodes
	Code93IncludeChecksum bool
	// Code93FullASCIIMode is required for code93 barcodes
	Code93FullASCIIMode bool
	// PDFSecurityLevel is required for pdf417 barcodes
	PDFSecurityLevel byte
	// QRLevel is required for qr barcodes
	QRLevel qr.ErrorCorrectionLevel
	// QRMode is required for qr barcodes
	QRMode qr.Encoding
}

// Barcode draws a barcode on the canvas.
func (canvas ImageCanvas) Barcode(codeType BarcodeType, content []byte, extra BarcodeExtraData, start image.Point, width, height int, dataColour color.Color, backgroundColour color.Color) (Canvas, error) {
	c := canvas
	if c.Image == nil {
		return canvas, errors.New("no image set for canvas to draw on")
	}
	var encodedBarcode barcode.Barcode
	var err error
	switch codeType {
	case BarcodeTypeAztec:
		encodedBarcode, err = aztec.Encode(content, extra.AztecMinECCPercent, extra.AztecUserSpecifiedLayers)
	case BarcodeTypeCodabar:
		encodedBarcode, err = codabar.Encode(string(content))
	case BarcodeTypeCode128:
		encodedBarcode, err = code128.Encode(string(content))
	case BarcodeTypeCode39:
		encodedBarcode, err = code39.Encode(string(content), extra.Code39IncludeChecksum, extra.Code39FullASCIIMode)
	case BarcodeTypeCode93:
		encodedBarcode, err = code93.Encode(string(content), extra.Code93IncludeChecksum, extra.Code93FullASCIIMode)
		if err != nil && err.Error() == "invalid data!" {
			err = errors.New("invalid data")
		}
	case BarcodeTypeDataMatrix:
		encodedBarcode, err = datamatrix.Encode(string(content))
	case BarcodeTypeEAN8:
		if len(content) != 8 {
			err = errors.New("EAN8 Barcode requires 8 characters")
		} else {
			encodedBarcode, err = ean.Encode(string(content))
		}
	case BarcodeTypeEAN13:
		if len(content) != 13 {
			err = errors.New("EAN13 Barcode requires 13 characters")
		} else {
			encodedBarcode, err = ean.Encode(string(content))
		}
	case BarcodeTypePDF:
		encodedBarcode, err = pdf417.Encode(string(content), extra.PDFSecurityLevel)
	case BarcodeTypeQR:
		encodedBarcode, err = qr.Encode(string(content), extra.QRLevel, extra.QRMode)
	case BarcodeType2of5:
		encodedBarcode, err = twooffive.Encode(string(content), false)
	case BarcodeType2of5Interleaved:
		encodedBarcode, err = twooffive.Encode(string(content), true)
	}
	if err != nil {
		return canvas, err
	}
	encodedBarcode, err = barcode.Scale(encodedBarcode, width, height)
	if err != nil {
		return canvas, err
	}

	if dataColour == nil {
		dataColour = color.Black
	}
	if backgroundColour == nil {
		backgroundColour = color.White
	}
	boundRect := encodedBarcode.Bounds()
	draw.DrawMask(c.Image, image.Rect(start.X, start.Y, start.X+width, start.Y+height), image.NewUniform(backgroundColour), image.ZP, blackAndWhiteMask{bw: encodedBarcode, bColour: color.Transparent, wColour: color.Opaque}, boundRect.Min, draw.Over)
	draw.DrawMask(c.Image, image.Rect(start.X, start.Y, start.X+width, start.Y+height), image.NewUniform(dataColour), image.ZP, blackAndWhiteMask{bw: encodedBarcode, bColour: color.Opaque, wColour: color.Transparent}, boundRect.Min, draw.Over)
	return c, nil
}

type blackAndWhiteMask struct {
	bw      image.Image
	bColour color.Alpha16
	wColour color.Alpha16
}

func (m blackAndWhiteMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (m blackAndWhiteMask) Bounds() image.Rectangle {
	return m.bw.Bounds()
}

func (m blackAndWhiteMask) At(x, y int) color.Color {
	if m.bw.At(x, y) == color.Black {
		return m.bColour
	}
	return m.wColour
}
