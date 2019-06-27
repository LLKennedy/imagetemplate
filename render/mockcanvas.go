package render

import (
	"image"
	"image/color"

	"github.com/stretchr/testify/mock"
	"golang.org/x/image/font"
)

// MockCanvas is a mock implementation of the Canvas interface for testing purposes.
type MockCanvas struct {
	mock.Mock
}

// SetUnderlyingImage returns the preset value(s).
func (m *MockCanvas) SetUnderlyingImage(newImage image.Image) Canvas {
	args := m.Called(newImage)
	return args.Get(0).(Canvas)
}

// GetUnderlyingImage returns the preset value(s).
func (m *MockCanvas) GetUnderlyingImage() image.Image {
	args := m.Called()
	return args.Get(0).(image.Image)
}

// GetWidth returns the preset value(s).
func (m *MockCanvas) GetWidth() int {
	args := m.Called()
	return args.Int(0)
}

// GetHeight returns the preset value(s).
func (m *MockCanvas) GetHeight() int {
	args := m.Called()
	return args.Int(0)
}

// GetPPI returns the preset value(s).
func (m *MockCanvas) GetPPI() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

// SetPPI returns the preset value(s).
func (m *MockCanvas) SetPPI(ppi float64) Canvas {
	args := m.Called(ppi)
	return args.Get(0).(Canvas)
}

// Rectangle returns the preset value(s).
func (m *MockCanvas) Rectangle(topLeft image.Point, width, height int, colour color.Color) (Canvas, error) {
	args := m.Called(topLeft, width, height, colour)
	return args.Get(0).(Canvas), args.Error(1)
}

// Circle returns the preset value(s).
func (m *MockCanvas) Circle(centre image.Point, radius int, colour color.Color) (Canvas, error) {
	args := m.Called(centre, radius, colour)
	return args.Get(0).(Canvas), args.Error(1)
}

// Text returns the preset value(s).
func (m *MockCanvas) Text(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (Canvas, error) {
	args := m.Called(text, start, typeFace, colour, maxWidth)
	return args.Get(0).(Canvas), args.Error(1)
}

// TryText returns the preset value(s).
func (m *MockCanvas) TryText(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (bool, int) {
	args := m.Called(text, start, typeFace, colour, maxWidth)
	return args.Bool(0), args.Int(1)
}

// DrawImage returns the preset value(s).
func (m *MockCanvas) DrawImage(start image.Point, subImage image.Image) (Canvas, error) {
	args := m.Called(start, subImage)
	return args.Get(0).(Canvas), args.Error(1)
}

// Barcode returns the preset value(s).
func (m *MockCanvas) Barcode(codeType BarcodeType, content []byte, extra BarcodeExtraData, start image.Point, width, height int, dataColour color.Color, bgColour color.Color) (Canvas, error) {
	args := m.Called(codeType, content, extra, start, width, height, dataColour, bgColour)
	return args.Get(0).(Canvas), args.Error(1)
}
