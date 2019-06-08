package render

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"testing"
)

func TestMockCanvas(t *testing.T) {
	m := new(MockCanvas)
	m.On("SetUnderlyingImage", image.Image(nil)).Return(m)
	assert.Equal(t, m, m.SetUnderlyingImage(image.Image(nil)))
	m.On("GetUnderlyingImage").Return(image.NewNRGBA(image.Rect(0, 0, 1, 1)))
	assert.Equal(t, image.NewNRGBA(image.Rect(0, 0, 1, 1)), m.GetUnderlyingImage())
	m.On("GetWidth").Return(30)
	assert.Equal(t, 30, m.GetWidth())
	m.On("GetHeight").Return(40)
	assert.Equal(t, 40, m.GetHeight())
	m.On("Rectangle", image.Pt(0, 0), 30, 25, color.Black).Return(m, fmt.Errorf("some error"))
	c, err := m.Rectangle(image.Pt(0, 0), 30, 25, color.Black)
	assert.Equal(t, m, c)
	assert.EqualError(t, err, "some error")
	m.On("Circle", image.Pt(0, 0), 6, color.White).Return(m, fmt.Errorf("some error"))
	c, err = m.Circle(image.Pt(0, 0), 6, color.White)
	assert.Equal(t, m, c)
	assert.EqualError(t, err, "some error")
	m.On("Text", "test", image.Pt(0, 0), nil, color.White, 50).Return(m, fmt.Errorf("some error"))
	c, err = m.Text("test", image.Pt(0, 0), nil, color.White, 50)
	assert.Equal(t, m, c)
	assert.EqualError(t, err, "some error")
	m.On("TryText", "test", image.Pt(0, 0), nil, color.White, 50).Return(true, 100)
	works, size := m.TryText("test", image.Pt(0, 0), nil, color.White, 50)
	assert.True(t, works)
	assert.Equal(t, size, 100)
	m.On("DrawImage", image.Pt(0, 0), nil).Return(m, fmt.Errorf("some error"))
	c, err = m.DrawImage(image.Pt(0, 0), nil)
	assert.Equal(t, m, c)
	assert.EqualError(t, err, "some error")
	barcodeType := BarcodeTypeAztec
	m.On("Barcode", barcodeType, []byte{}, BarcodeExtraData{}, image.Pt(0, 0), 50, 20, color.White, color.Black).Return(m, fmt.Errorf("some error"))
	c, err = m.Barcode(barcodeType, []byte{}, BarcodeExtraData{}, image.Pt(0, 0), 50, 20, color.White, color.Black)
	assert.Equal(t, m, c)
	assert.EqualError(t, err, "some error")
	m.On("SetPPI", float64(100)).Return(m)
	assert.Equal(t, m, m.SetPPI(float64(100)))
	m.On("GetPPI").Return(float64(80))
	assert.Equal(t, float64(80), m.GetPPI())
	m.AssertExpectations(t)
}
