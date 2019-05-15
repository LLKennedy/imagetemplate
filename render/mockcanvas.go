package render

import (
	"golang.org/x/image/font"
	"image"
	"image/color"
)

// MockCanvas is a mock implementation of the Canvas interface for testing purposes
type MockCanvas struct {
	FixedSetUnderlyingImage Canvas
	FixedGetUnderlyingImage image.Image
	FixedGetWidth           int
	FixedGetHeight          int
	FixedRectangleError     error
	FixedCircleError        error
	FixedTextError          error
	FixedTryTextBool        bool
	FixedTryTextInt         int
	FixedDrawImageError     error
	FixedBarcodeError       error
	FixedPixelsPerInch      int
}

// SetUnderlyingImage returns the preset value(s)
func (m MockCanvas) SetUnderlyingImage(newImage image.Image) Canvas {
	return m.FixedSetUnderlyingImage
}

// GetUnderlyingImage returns the preset value(s)
func (m MockCanvas) GetUnderlyingImage() image.Image {
	return m.FixedGetUnderlyingImage
}

// GetWidth returns the preset value(s)
func (m MockCanvas) GetWidth() int {
	return m.FixedGetWidth
}

// GetHeight returns the preset value(s)
func (m MockCanvas) GetHeight() int {
	return m.FixedGetHeight
}

// GetPPI returns the preset value(s)
func (m MockCanvas) GetPPI() int {
	return m.FixedPixelsPerInch
}

// Rectangle returns the preset value(s)
func (m MockCanvas) Rectangle(topLeft image.Point, width, height int, colour color.Color) (Canvas, error) {
	return m, m.FixedRectangleError
}

// Circle returns the preset value(s)
func (m MockCanvas) Circle(centre image.Point, radius int, colour color.Color) (Canvas, error) {
	return m, m.FixedCircleError
}

// Text returns the preset value(s)
func (m MockCanvas) Text(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (Canvas, error) {
	return m, m.FixedTextError
}

// TryText returns the preset value(s)
func (m MockCanvas) TryText(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (bool, int) {
	return m.FixedTryTextBool, m.FixedTryTextInt
}

// DrawImage returns the preset value(s)
func (m MockCanvas) DrawImage(start image.Point, subImage image.Image) (Canvas, error) {
	return m, m.FixedDrawImageError
}

// Barcode returns the preset value(s)
func (m MockCanvas) Barcode(codeType BarcodeType, content []byte, extra BarcodeExtraData, start image.Point, width, height int, dataColour color.Color, bgColour color.Color) (Canvas, error) {
	return m, m.FixedBarcodeError
}
