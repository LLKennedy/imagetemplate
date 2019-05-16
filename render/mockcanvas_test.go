package render

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
)

func TestMockCanvas(t *testing.T) {
	m := MockCanvas{
		FixedSetUnderlyingImage: &MockCanvas{},
		FixedGetUnderlyingImage: nil,
		FixedGetWidth: 50,
		FixedGetHeight: 25,
		FixedRectangleError: nil,
		FixedCircleError: nil,
		FixedTextError: nil,
		FixedTryTextBool: true,
		FixedTryTextInt: 75,
		FixedDrawImageError: nil,
		FixedBarcodeError: nil,
		FixedPixelsPerInch: 80,
	}
	assert.Equal(t, m.FixedSetUnderlyingImage, m.SetUnderlyingImage(nil))
	assert.Equal(t, m.FixedGetUnderlyingImage, m.GetUnderlyingImage())
	assert.Equal(t, m.FixedGetWidth, m.GetWidth())
	assert.Equal(t, m.FixedGetHeight, m.GetHeight())
	mm, err := m.Rectangle(image.Pt(0,0), 30, 25, color.Black)
	assert.Equal(t, m, mm)
	assert.Equal(t, m.FixedRectangleError, err)
	mm, err = m.Circle(image.Pt(0,0), 6, color.White)
	assert.Equal(t, m, mm)
	assert.Equal(t, m.FixedCircleError, err)
	mm, err = m.Text("some text", image.Pt(0,0), nil, color.White, 50)
	assert.Equal(t, m, mm)
	assert.Equal(t, m.FixedTextError, err)
	fits, width := m.TryText("some text", image.Pt(0,0), nil, color.White, 50)
	assert.Equal(t, m.FixedTryTextBool, fits)
	assert.Equal(t, m.FixedTryTextInt, width)
	mm, err = m.DrawImage(image.Pt(0,0), nil)
	assert.Equal(t, m, mm)
	assert.Equal(t, m.FixedDrawImageError, err)
	mm, err = m.Barcode(BarcodeTypeAztec, []byte{}, BarcodeExtraData{}, image.Pt(0,0), 50, 20, color.White, color.Black)
	assert.Equal(t, m, mm)
	assert.Equal(t, m.FixedBarcodeError, err)
	ppi := m.GetPPI()
	assert.Equal(t, m.FixedPixelsPerInch, ppi)
}