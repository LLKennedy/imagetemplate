package imagetemplate

import (
	"image"
	"image/color"
	"golang.org/x/image/font"
)

type mockCanvas struct {
	FixedSetUnderlyingImage Canvas
	FixedGetUnderlyingImage image.Image
	FixedGetWidth int
	FixedGetHeight int
	FixedRectangleError error
	FixedCircleError error
	FixedTextError error
	FixedTryTextBool bool
	FixedTryTextInt int
	FixedDrawImageError error
	FixedBarcodeError error
	FixedPixelsPerInch int
}

func (m mockCanvas) SetUnderlyingImage(newImage image.Image) Canvas {
	return m.FixedSetUnderlyingImage
}
func (m mockCanvas) GetUnderlyingImage() image.Image {
	return m.FixedGetUnderlyingImage
}
func (m mockCanvas) GetWidth() int {
	return m.FixedGetWidth
}
func (m mockCanvas) GetHeight() int {
	return m.FixedGetHeight
}
func (m mockCanvas) GetPPI() int {
	return m.FixedPixelsPerInch
}
func (m mockCanvas) Rectangle(topLeft image.Point, width, height int, colour color.Color) (Canvas, error) {
	return m, m.FixedRectangleError
}
func (m mockCanvas) Circle(centre image.Point, radius int, colour color.Color) (Canvas, error) {
	return m, m.FixedCircleError
}
func (m mockCanvas) Text(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (Canvas, error) {
	return m, m.FixedTextError
}
func (m mockCanvas) TryText(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (bool, int) {
	return m.FixedTryTextBool, m.FixedTryTextInt
}
func (m mockCanvas) DrawImage(start image.Point, subImage image.Image) (Canvas, error) {
	return m, m.FixedDrawImageError
}
func (m mockCanvas) Barcode(codeType BarcodeType, content []byte, extra BarcodeExtraData, start image.Point, width, height int, dataColour color.Color, bgColour color.Color) (Canvas, error) {
	return m, m.FixedBarcodeError
}