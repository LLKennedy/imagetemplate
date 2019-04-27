package imagetemplate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestNewCanvas(t *testing.T) {
	t.Run("invalid width and height", func(t *testing.T) {
		testVals := []int{-9999999, -100, -50, -1, 0}
		for _, x := range testVals {
			for _, y := range testVals {
				t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
					newCanvas, err := NewCanvas(x, y)
					assert.Nil(t, newCanvas.Image)
					assert.Equal(t, err.Error(), "Invalid width and height")
				})
			}
		}
	})
	t.Run("invalid width", func(t *testing.T) {
		testVals := []int{-9999999, -100, -50, -1, 0}
		y := 100
		for _, x := range testVals {
			t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
				newCanvas, err := NewCanvas(x, y)
				assert.Nil(t, newCanvas.Image)
				assert.Equal(t, err.Error(), "Invalid width")
			})
		}
	})
	t.Run("invalid width", func(t *testing.T) {
		testVals := []int{-9999999, -100, -50, -1, 0}
		x := 100
		for _, y := range testVals {
			t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
				newCanvas, err := NewCanvas(x, y)
				assert.Nil(t, newCanvas.Image)
				assert.Equal(t, err.Error(), "Invalid height")
			})
		}
	})
	t.Run("valid input", func(t *testing.T) {
		testVals := []int{1, 50, 100}
		for _, x := range testVals {
			for _, y := range testVals {
				t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
					newCanvas, err := NewCanvas(x, y)
					assert.Nil(t, err)
					assert.NotNil(t, newCanvas.Image)
					if newCanvas.Image != nil {
						bounds := newCanvas.Image.Bounds()
						width := bounds.Max.X - bounds.Min.X
						height := bounds.Max.Y - bounds.Min.Y
						assert.Equal(t, x, width)
						assert.Equal(t, y, height)
						assert.Equal(t, color.NRGBAModel, newCanvas.Image.ColorModel())
						for pX := 0; pX < x; pX++ {
							for pY := 0; pY < y; pY++ {
								assert.Equal(t, color.NRGBA{}, newCanvas.Image.At(pX, pY))
							}
						}
					}
				})
			}
		}
	})
}

func TestSetUnderlyingImage(t *testing.T) {
	t.Run("using draw.Image", func(t *testing.T) {
		newCanvas, _ := NewCanvas(100, 100)
		otherImage := image.NewGray16(image.Rect(0, 0, 100, 100))
		for x := 0; x < 100; x++ {
			for y := 0; y < 100; y++ {
				assert.NotEqual(t, newCanvas.Image.At(x, y), otherImage.At(x, y))
			}
		}
		modifiedCanvas := newCanvas.SetUnderlyingImage(otherImage)
		//Check original isn't changed
		for x := 0; x < 100; x++ {
			for y := 0; y < 100; y++ {
				assert.NotEqual(t, newCanvas.Image.At(x, y), otherImage.At(x, y))
			}
		}
		//Check new version has changed
		modifiedImage := modifiedCanvas.GetUnderlyingImage()
		for x := 0; x < 100; x++ {
			for y := 0; y < 100; y++ {
				assert.Equal(t, modifiedImage.At(x, y), otherImage.At(x, y))
			}
		}
	})
	t.Run("using image.Image", func(t *testing.T) {
		newCanvas, _ := NewCanvas(150, 150)
		var otherImage image.Image
		otherImage = image.Rect(0, 0, 150, 150)
		drawImage, canConvert := otherImage.(draw.Image)
		assert.Nil(t, drawImage)
		assert.False(t, canConvert)
		for x := 0; x < 100; x++ {
			for y := 0; y < 100; y++ {
				t.Run(fmt.Sprintf("first test x=%d,y=%d", x, y), func(t *testing.T) {
					assert.NotEqual(t, newCanvas.Image.At(x, y), otherImage.At(x, y))
				})
			}
		}
		modifiedCanvas := newCanvas.SetUnderlyingImage(otherImage)
		//Check original isn't changed
		for x := 0; x < 100; x++ {
			for y := 0; y < 100; y++ {
				t.Run(fmt.Sprintf("second test x=%d,y=%d", x, y), func(t *testing.T) {
					assert.NotEqual(t, newCanvas.Image.At(x, y), otherImage.At(x, y))
				})
			}
		}
		//Check new version has changed
		modifiedImage := modifiedCanvas.GetUnderlyingImage()
		equivalentNRGBA := image.NewNRGBA(image.Rect(0, 0, 150, 150))
		for x := 0; x < 100; x++ {
			for y := 0; y < 100; y++ {
				//Parse otherImage (type Alpha16) through NRGBA to match exactly
				equivalentNRGBA.Set(x, y, otherImage.At(x, y))
				t.Run(fmt.Sprintf("third test x=%d,y=%d", x, y), func(t *testing.T) {
					assert.Equal(t, equivalentNRGBA.At(x, y), modifiedImage.At(x, y))
				})
			}
		}
	})
}

func TestGetUnderlyingImage(t *testing.T) {
	newCanvas, _ := NewCanvas(200, 200)
	otherImage := image.NewGray(image.Rect(0, 0, 200, 200))
	initialImage := newCanvas.GetUnderlyingImage()
	assert.NotEqual(t, initialImage, otherImage)
	modifiedCanvas := newCanvas.SetUnderlyingImage(otherImage)
	retrievedImage := modifiedCanvas.GetUnderlyingImage()
	assert.Equal(t, retrievedImage, otherImage)
}

func TestGetWidthHeight(t *testing.T) {
	testVals := []int{1, 5, 10, 44, 103, 314, 1000}
	for _, x := range testVals {
		for _, y := range testVals {
			t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
				newCanvas, _ := NewCanvas(x, y)
				assert.Equal(t, newCanvas.GetWidth(), x)
				assert.Equal(t, newCanvas.GetHeight(), y)
			})
		}
	}
}
