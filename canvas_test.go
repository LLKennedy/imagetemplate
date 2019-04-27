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
								t.Run(fmt.Sprintf("pX=%d,pY=%d", pX, pY), func(t *testing.T) {
									assert.Equal(t, color.NRGBA{}, newCanvas.Image.At(pX, pY))
								})
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

func TestRectangle(t *testing.T) {
	newCanvas, _ := NewCanvas(5, 5)
	t.Run("invalid width and height", func(t *testing.T) {
		testVals := []int{-9999999, -100, -50, -1, 0}
		for _, x := range testVals {
			for _, y := range testVals {
				t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
					modifiedCanvas, err := newCanvas.Rectangle(image.ZP, x, y, color.Black)
					assert.Equal(t, newCanvas, modifiedCanvas)
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
				modifiedCanvas, err := newCanvas.Rectangle(image.ZP, x, y, color.Black)
				assert.Equal(t, newCanvas, modifiedCanvas)
				assert.Equal(t, err.Error(), "Invalid width")
			})
		}
	})
	t.Run("invalid width", func(t *testing.T) {
		testVals := []int{-9999999, -100, -50, -1, 0}
		x := 100
		for _, y := range testVals {
			t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
				modifiedCanvas, err := newCanvas.Rectangle(image.ZP, x, y, color.Black)
				assert.Equal(t, newCanvas, modifiedCanvas)
				assert.Equal(t, err.Error(), "Invalid height")
			})
		}
	})
	t.Run("ouside range", func(t *testing.T) {
		tests := [][]int{
			[]int{-100, -100, 50, 50},
			[]int{-100, 100, 50, 50},
			[]int{100, -100, 50, 50},
			[]int{100, 100, 50, 50},
		}
		for _, test := range tests {
			t.Run(fmt.Sprintf("x=%d,y=%d,w=%d,h=%d", test[0], test[1], test[2], test[3]), func(t *testing.T) {
				modifiedCanvas, err := newCanvas.Rectangle(image.Pt(test[0], test[1]), test[2], test[3], color.Black)
				assert.Nil(t, err)
				underlyingImage := modifiedCanvas.GetUnderlyingImage()
				for pX := 0; pX < 10; pX++ {
					for pY := 0; pY < 10; pY++ {
						t.Run(fmt.Sprintf("pX=%d,pY=%d", pX, pY), func(t *testing.T) {
							assert.Equal(t, color.NRGBA{}, underlyingImage.At(pX, pY))
						})
					}
				}
			})
		}
	})
	t.Run("ouside range", func(t *testing.T) {
		tests := [][]int{
			[]int{-100, -100, 50, 50},
			[]int{-100, 100, 50, 50},
			[]int{100, -100, 50, 50},
			[]int{100, 100, 50, 50},
		}
		for _, test := range tests {
			t.Run(fmt.Sprintf("x=%d,y=%d,w=%d,h=%d", test[0], test[1], test[2], test[3]), func(t *testing.T) {
				modifiedCanvas, err := newCanvas.Rectangle(image.Pt(test[0], test[1]), test[2], test[3], color.NRGBA{R: 255, A: 255})
				assert.Nil(t, err)
				underlyingImage := modifiedCanvas.GetUnderlyingImage()
				for pX := 0; pX < 10; pX++ {
					for pY := 0; pY < 10; pY++ {
						t.Run(fmt.Sprintf("pX=%d,pY=%d", pX, pY), func(t *testing.T) {
							assert.Equal(t, color.NRGBA{}, underlyingImage.At(pX, pY))
						})
					}
				}
			})
		}
	})
	t.Run("visible draws", func(t *testing.T) {
		tests := [][]int{
			[]int{-4, -4, 5, 5},
			[]int{-4, 4, 5, 5},
			[]int{4, -4, 5, 5},
			[]int{4, 4, 5, 5},
			[]int{0, 0, 5, 5},
			[]int{0, 2, 5, 5},
			[]int{2, 0, 5, 5},
		}
		for _, test := range tests {
			x, y, w, h := test[0], test[1], test[2], test[3]
			t.Run(fmt.Sprintf("x=%d,y=%d,w=%d,h=%d", x, y, w, h), func(t *testing.T) {
				modifiedCanvas, err := newCanvas.Rectangle(image.Pt(x, y), w, h, color.NRGBA{R: 255, A: 255})
				assert.Nil(t, err)
				underlyingImage := modifiedCanvas.GetUnderlyingImage()
				for pX := 0; pX < 10; pX++ {
					for pY := 0; pY < 10; pY++ {
						t.Run(fmt.Sprintf("pX=%d,pY=%d", pX, pY), func(t *testing.T) {
							if pX >= x && pX < x+w && pY >= y && pY < y+h {
								assert.Equal(t, color.NRGBA{R: 255, A: 255}, underlyingImage.At(pX, pY))
							} else {
								assert.Equal(t, color.NRGBA{}, underlyingImage.At(pX, pY))
							}
						})
					}
				}
			})
			newCanvas = newCanvas.SetUnderlyingImage(image.NewNRGBA(image.Rect(0, 0, 10, 10))).(ImageCanvas)
		}
	})
}

func TestCircle(t *testing.T) {
	newCanvas, _ := NewCanvas(10, 10)
	t.Run("invalid radius", func(t *testing.T) {
		modifiedCanvas, err := newCanvas.Circle(image.Pt(3, 3), -1, color.NRGBA{G: 255, A: 255})
		assert.Equal(t, newCanvas, modifiedCanvas)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "Invalid radius")
		modifiedCanvas, err = newCanvas.Circle(image.Pt(3, 3), 0, color.NRGBA{G: 255, A: 255})
		assert.Equal(t, newCanvas, modifiedCanvas)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "Invalid radius")
	})
	t.Run("valid circle", func(t *testing.T) {
		modifiedCanvas, err := newCanvas.Circle(image.Pt(3, 3), 2, color.NRGBA{G: 255, A: 255})
		assert.Nil(t, err)
		testCircle := &circle{p: image.Pt(3, 3), r: 2}
		assert.Equal(t, color.AlphaModel, testCircle.ColorModel())
		underlyingImage := modifiedCanvas.GetUnderlyingImage()
		tMask := color.Alpha{0}
		oMask := color.Alpha{255}
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
					currentPixel := testCircle.At(x, y)
					if currentPixel == tMask {
						assert.Equal(t, color.NRGBA{}, underlyingImage.At(x, y))
					} else if currentPixel == oMask {
						assert.Equal(t, color.NRGBA{G: 255, A: 255}, underlyingImage.At(x, y))
					} else {
						assert.Fail(t, fmt.Sprintf("Invalid value for current pixel: %v", currentPixel))
					}
				})
			}
		}
	})
}
