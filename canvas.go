package imagetemplate

import (
	"errors"
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
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
)

// Canvas holds the image struct and associated properties
type Canvas interface {
	SetUnderlyingImage(newImage image.Image) Canvas
	GetUnderlyingImage() image.Image
	GetWidth() int
	GetHeight() int
	Rectangle(topLeft image.Point, width, height int, colour color.Color) (Canvas, error)
	Circle(centre image.Point, radius int, colour color.Color) (Canvas, error)
	Text(text string, start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) (Canvas, error)
	SubImage(start image.Point, subImage image.Image) (Canvas, error)
	Barcode(codeType BarcodeType, content []byte, extra BarcodeExtraData, start image.Point, width, height int, dataColour color.Color, bgColour color.Color) (Canvas, error)
}

// ImageCanvas uses golang's native Image package to implement the Canvas interface
type ImageCanvas struct {
	Image draw.Image
}

// NewCanvas generates a new canvas of the given width and height
func NewCanvas(width, height int) (ImageCanvas, error) {
	if width <= 0 && height <= 0 {
		return ImageCanvas{}, errors.New("Invalid width and height")
	} else if width <= 0 {
		return ImageCanvas{}, errors.New("Invalid width")
	} else if height <= 0 {
		return ImageCanvas{}, errors.New("Invalid height")
	}
	return ImageCanvas{
		Image: image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: width, Y: height},
		}),
	}, nil
}

// SetUnderlyingImage sets the internal Image property to the given object
func (canvas ImageCanvas) SetUnderlyingImage(newImage image.Image) Canvas {
	drawImage, ok := newImage.(draw.Image)
	if !ok {
		bounds := newImage.Bounds()
		drawImage = image.NewNRGBA(bounds)
		draw.Draw(drawImage, bounds, newImage, bounds.Min, draw.Src)
	}
	canvas.Image = drawImage
	return canvas
}

// GetUnderlyingImage gets the internal Image property
func (canvas ImageCanvas) GetUnderlyingImage() image.Image {
	return canvas.Image
}

// GetWidth returns the width of the underlying Image
func (canvas ImageCanvas) GetWidth() int {
	return canvas.Image.Bounds().Size().X
}

// GetWidth returns the width of the underlying Image
func (canvas ImageCanvas) GetHeight() int {
	return canvas.Image.Bounds().Size().Y
}

// Rectangle draws a rectangle of a specific colour on the canvas
func (canvas ImageCanvas) Rectangle(topLeft image.Point, width, height int, colour color.Color) (Canvas, error) {
	c := canvas
	colourPlane := image.Uniform{C: colour}
	if width <= 0 && height <= 0 {
		return canvas, errors.New("Invalid width and height")
	} else if width <= 0 {
		return canvas, errors.New("Invalid width")
	} else if height <= 0 {
		return canvas, errors.New("Invalid height")
	}
	rect := image.Rectangle{
		Min: image.Point{
			X: topLeft.X,
			Y: topLeft.Y,
		},
		Max: image.Point{
			X: topLeft.X + width,
			Y: topLeft.Y + height,
		},
	}
	draw.Draw(c.Image, rect, &colourPlane, topLeft, draw.Over)
	return c, nil
}

// Circle draws a circle of a specific colour on the canvas
func (canvas ImageCanvas) Circle(centre image.Point, radius int, colour color.Color) (Canvas, error) {
	c := canvas
	colourPlane := image.Uniform{C: colour}
	mask := &circle{p: centre, r: radius}
	draw.DrawMask(c.Image, mask.Bounds(), &colourPlane, image.ZP, mask, mask.Bounds().Min, draw.Over)
	return c, nil
}

// Text draws text on the canvas
func (canvas ImageCanvas) Text(text string, start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) (Canvas, error) {
	c := canvas
	drawer := &font.Drawer{
		Dst:  c.Image,
		Face: typeFace,
		Src:  image.NewUniform(colour),
	}
	drawer.DrawString(text)
	return c, nil
}

// SubImage draws another image on the canvas
func (canvas ImageCanvas) SubImage(start image.Point, subImage image.Image) (Canvas, error) {
	c := canvas
	draw.Draw(c.Image, subImage.Bounds(), subImage, start, draw.Over)
	return c, nil
}

// BarcodeType wraps the barcode types into a single enum
type BarcodeType string

const (
	// BarcodeTypeAztec           is an alias for an imported barcode type
	BarcodeTypeAztec = barcode.TypeAztec
	// BarcodeTypeCodabar         is an alias for an imported barcode type
	BarcodeTypeCodabar = barcode.TypeCodabar
	// BarcodeTypeCode128         is an alias for an imported barcode type
	BarcodeTypeCode128 = barcode.TypeCode128
	// BarcodeTypeCode39          is an alias for an imported barcode type
	BarcodeTypeCode39 = barcode.TypeCode39
	// BarcodeTypeCode93          is an alias for an imported barcode type
	BarcodeTypeCode93 = barcode.TypeCode93
	// BarcodeTypeDataMatrix      is an alias for an imported barcode type
	BarcodeTypeDataMatrix = barcode.TypeDataMatrix
	// BarcodeTypeEAN8            is an alias for an imported barcode type
	BarcodeTypeEAN8 = barcode.TypeEAN8
	// BarcodeTypeEAN13           is an alias for an imported barcode type
	BarcodeTypeEAN13 = barcode.TypeEAN13
	// BarcodeTypePDF             is an alias for an imported barcode type
	BarcodeTypePDF = barcode.TypePDF
	// BarcodeTypeQR              is an alias for an imported barcode type
	BarcodeTypeQR = barcode.TypeQR
	// BarcodeType2of5            is an alias for an imported barcode type
	BarcodeType2of5 = barcode.Type2of5
	// BarcodeType2of5Interleaved is an alias for an imported barcode type
	BarcodeType2of5Interleaved = barcode.Type2of5Interleaved
)

// BarcodeExtraData contains additional data required for some barcode formats, leave any fields not named for the type in use alone
type BarcodeExtraData struct {
	// AztecMinECCPercent       is required for aztec barcodes
	AztecMinECCPercent int
	// AztecUserSpecifiedLayers is required for aztec barcodes
	AztecUserSpecifiedLayers int
	// Code39IncludeChecksum    is required for code39 barcodes
	Code39IncludeChecksum bool
	// Code39FullAsciiMode      is required for code39 barcodes
	Code39FullAsciiMode bool
	// Code93IncludeChecksum    is required for code93 barcodes
	Code93IncludeChecksum bool
	// Code93FullAsciiMode      is required for code93 barcodes
	Code93FullAsciiMode bool
	// PDFSecurityLevel         is required for pdf417 barcodes
	PDFSecurityLevel byte
	// QRLevel                  is required for qr barcodes
	QRLevel qr.ErrorCorrectionLevel
	// QRMode                   is required for qr barcodes
	QRMode qr.Encoding
}

// Barcode draws a barcode on the canvas
func (canvas ImageCanvas) Barcode(codeType BarcodeType, content []byte, extra BarcodeExtraData, start image.Point, width, height int, dataColour color.Color, backgroundColour color.Color) (Canvas, error) {
	c := canvas
	var encodedBarcode barcode.Barcode
	var err error
	switch codeType {
	case BarcodeTypeAztec:
		encodedBarcode, err = aztec.Encode(content, extra.AztecMinECCPercent, extra.AztecUserSpecifiedLayers)
		if err != nil {
			return canvas, err
		}
	case BarcodeTypeCodabar:
		encodedBarcode, err = codabar.Encode(string(content))
		if err != nil {
			return canvas, err
		}
	case BarcodeTypeCode128:
		encodedBarcode, err = code128.Encode(string(content))
		if err != nil {
			return canvas, err
		}
	case BarcodeTypeCode39:
		encodedBarcode, err = code39.Encode(string(content), extra.Code39IncludeChecksum, extra.Code39FullAsciiMode)
		if err != nil {
			return canvas, err
		}
	case BarcodeTypeCode93:
		encodedBarcode, err = code93.Encode(string(content), extra.Code93IncludeChecksum, extra.Code93FullAsciiMode)
		if err != nil {
			return canvas, err
		}
	case BarcodeTypeDataMatrix:
		encodedBarcode, err = datamatrix.Encode(string(content))
		if err != nil {
			return canvas, err
		}
	case BarcodeTypeEAN8:
		if len(content) != 8 {
			return canvas, errors.New("EAN8 Barcode requires 8 characters")
		}
		encodedBarcode, err = ean.Encode(string(content))
		if err != nil {
			return canvas, err
		}
	case BarcodeTypeEAN13:
		if len(content) != 13 {
			return canvas, errors.New("EAN13 Barcode requires 13 characters")
		}
		encodedBarcode, err = ean.Encode(string(content))
		if err != nil {
			return canvas, err
		}
	case BarcodeTypePDF:
		encodedBarcode, err = pdf417.Encode(string(content), extra.PDFSecurityLevel)
		if err != nil {
			return canvas, err
		}
	case BarcodeTypeQR:
		encodedBarcode, err = qr.Encode(string(content), extra.QRLevel, extra.QRMode)
		if err != nil {
			return canvas, err
		}
	case BarcodeType2of5:
		encodedBarcode, err = twooffive.Encode(string(content), false)
		if err != nil {
			return canvas, err
		}
	case BarcodeType2of5Interleaved:
		encodedBarcode, err = twooffive.Encode(string(content), true)
		if err != nil {
			return canvas, err
		}
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

// Steal the circle example code from https://blog.golang.org/go-imagedraw-package
type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
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
