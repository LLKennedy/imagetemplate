package render

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

func TestBlankCanvas(t *testing.T) {
	blankCanvas := ImageCanvas{}
	t.Run("circle", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.Circle(image.ZP, 10, color.White)
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.EqualError(t, err, "no image set for canvas to draw on")
	})
	t.Run("barcode", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.Barcode(BarcodeTypeAztec, []byte{}, BarcodeExtraData{}, image.ZP, 0, 0, color.White, color.Black)
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.EqualError(t, err, "no image set for canvas to draw on")
	})
	t.Run("draw image", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.DrawImage(image.ZP, image.NewNRGBA(image.Rect(0, 0, 10, 10)))
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.EqualError(t, err, "no image set for canvas to draw on")
	})
	t.Run("get height", func(t *testing.T) {
		height := blankCanvas.GetHeight()
		assert.Equal(t, 0, height)
	})
	t.Run("get image func", func(t *testing.T) {
		img := blankCanvas.GetUnderlyingImage()
		assert.Nil(t, img)
	})
	t.Run("get width", func(t *testing.T) {
		width := blankCanvas.GetWidth()
		assert.Equal(t, 0, width)
	})
	t.Run("set ppi", func(t *testing.T) {
		ppi := float64(352)
		modifiedCanvas := blankCanvas.SetPPI(ppi)
		assert.Equal(t, ppi, modifiedCanvas.(ImageCanvas).pixelsPerInch)
	})
	t.Run("get ppi", func(t *testing.T) {
		newPPI := float64(718)
		blankCanvas.pixelsPerInch = newPPI
		ppi := blankCanvas.GetPPI()
		assert.Equal(t, newPPI, ppi)
	})
	t.Run("rectangle", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.Rectangle(image.ZP, 10, 10, color.White)
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.EqualError(t, err, "no image set for canvas to draw on")
	})
	t.Run("text", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.Text("hello", image.ZP, nil, color.White, 100)
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.EqualError(t, err, "no image set for canvas to draw on")
	})
	t.Run("get image", func(t *testing.T) {
		img := blankCanvas.Image
		assert.Nil(t, img)
	})
	t.Run("set image", func(t *testing.T) {
		modifiedCanvas := blankCanvas.SetUnderlyingImage(image.NewNRGBA(image.Rect(0, 0, 10, 10)))
		assert.NotEqual(t, blankCanvas, modifiedCanvas)
		assert.NotNil(t, modifiedCanvas.GetUnderlyingImage())
	})
}

func TestNewCanvas(t *testing.T) {
	t.Run("invalid width and height", func(t *testing.T) {
		testVals := []int{-9999999, -100, -50, -1, 0}
		for _, x := range testVals {
			for _, y := range testVals {
				t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
					newCanvas, err := NewCanvas(x, y)
					assert.Nil(t, newCanvas.Image)
					assert.EqualError(t, err, "invalid width and height")
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
				assert.EqualError(t, err, "invalid width")
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
				assert.EqualError(t, err, "invalid height")
			})
		}
	})
	t.Run("valid input", func(t *testing.T) {
		testVals := []int{1, 50, 100}
		for _, x := range testVals {
			for _, y := range testVals {
				t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
					newCanvas, err := NewCanvas(x, y)
					assert.NoError(t, err)
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
					assert.EqualError(t, err, "invalid width and height")
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
				assert.EqualError(t, err, "invalid width")
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
				assert.EqualError(t, err, "invalid height")
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
				assert.NoError(t, err)
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
				assert.NoError(t, err)
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
				assert.NoError(t, err)
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
		assert.EqualError(t, err, "invalid radius")
		modifiedCanvas, err = newCanvas.Circle(image.Pt(3, 3), 0, color.NRGBA{G: 255, A: 255})
		assert.Equal(t, newCanvas, modifiedCanvas)
		assert.EqualError(t, err, "invalid radius")
	})
	t.Run("valid circle", func(t *testing.T) {
		modifiedCanvas, err := newCanvas.Circle(image.Pt(3, 3), 2, color.NRGBA{G: 255, A: 255})
		assert.NoError(t, err)
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

func TestText(t *testing.T) {
	var newCanvas Canvas
	newCanvas, _ = NewCanvas(30, 30)
	newCanvas = newCanvas.SetUnderlyingImage(image.NewNRGBA(newCanvas.GetUnderlyingImage().Bounds()))
	var regFont font.Face
	ttFont, _ := truetype.Parse(goregular.TTF)
	regFont = truetype.NewFace(ttFont, &truetype.Options{Size: 14, Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64})
	t.Run("invalid maxWidth", func(t *testing.T) {
		modifiedCanvas, err := newCanvas.Text("a", image.Pt(2, 10), regFont, color.White, -1)
		assert.Equal(t, newCanvas, modifiedCanvas)
		assert.EqualError(t, err, "invalid maxWidth")
	})
	t.Run("valid text", func(t *testing.T) {
		purple := color.NRGBA{R: 130, B: 200, A: 255}
		modifiedCanvas, err := newCanvas.Text("test", image.Pt(2, 15), regFont, purple, 28)
		assert.NoError(t, err)
		assert.NotNil(t, modifiedCanvas)
		var data bytes.Buffer
		err = bmp.Encode(&data, modifiedCanvas.GetUnderlyingImage())
		assert.NoError(t, err)
		assert.Equal(t, testTextOutput, data.Bytes())
	})
	t.Run("oversized text", func(t *testing.T) {
		purple := color.NRGBA{R: 130, B: 200, A: 255}
		_, err := newCanvas.Text("test this much longer string which definitely won't fit", image.Pt(2, 15), regFont, purple, 28)
		assert.EqualError(t, err, "resultant drawn text was longer than maxWidth")
	})
}

func TestTryText(t *testing.T) {
	var newCanvas Canvas
	newCanvas, _ = NewCanvas(30, 30)
	var regFont font.Face
	ttFont, _ := truetype.Parse(goregular.TTF)
	regFont = truetype.NewFace(ttFont, &truetype.Options{Size: 14, Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64})
	newCanvas = newCanvas.SetUnderlyingImage(image.NewNRGBA(newCanvas.GetUnderlyingImage().Bounds()))
	t.Run("invalid maxWidth", func(t *testing.T) {
		fits, width := newCanvas.TryText("a", image.Pt(2, 10), regFont, color.White, -1)
		assert.False(t, fits)
		assert.Equal(t, -1, width)
	})
	t.Run("invalid image", func(t *testing.T) {
		modifiedCanvas := newCanvas.SetUnderlyingImage(nil)
		fits, width := modifiedCanvas.TryText("a", image.Pt(2, 10), regFont, color.White, 10)
		assert.False(t, fits)
		assert.Equal(t, -2, width)
	})
	t.Run("valid text", func(t *testing.T) {
		purple := color.NRGBA{R: 130, B: 200, A: 255}
		fits, width := newCanvas.TryText("test", image.Pt(2, 15), regFont, purple, 28)
		assert.True(t, fits)
		assert.Equal(t, 23, width) //no maths to calculate this, I just ran the function
	})
	t.Run("oversized text", func(t *testing.T) {
		purple := color.NRGBA{R: 130, B: 200, A: 255}
		fits, width := newCanvas.TryText("test this much longer string which definitely won't fit", image.Pt(2, 15), regFont, purple, 28)
		assert.False(t, fits)
		assert.Equal(t, 325, width) //no maths to calculate this, I just ran the function
	})
}

func TestDrawImage(t *testing.T) {
	var newCanvas, otherCanvas Canvas
	sx := 30
	sy := 50
	w := 20
	h := 20
	newCanvas, _ = NewCanvas(100, 100)
	newCanvas, _ = newCanvas.Rectangle(image.ZP, 100, 100, color.NRGBA{B: 255, A: 255})
	otherCanvas, _ = NewCanvas(20, 20)
	otherCanvas, _ = otherCanvas.Rectangle(image.ZP, w, h, color.NRGBA{G: 255, A: 255})
	modifiedCanvas, err := newCanvas.DrawImage(image.Pt(sx, sy), otherCanvas.GetUnderlyingImage())
	assert.NoError(t, err)
	result := modifiedCanvas.GetUnderlyingImage()
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			t.Run(fmt.Sprintf("x=%d,y=%d", x, y), func(t *testing.T) {
				if x >= sx && x < sx+w && y >= sy && y < sy+h {
					assert.Equal(t, color.NRGBA{G: 255, A: 255}, result.At(x, y))
				} else {
					assert.Equal(t, color.NRGBA{B: 255, A: 255}, result.At(x, y))
				}
			})
		}
	}
}

func TestBarcode(t *testing.T) {
	for i := 0; i < 1; i++ {
		t.Run(fmt.Sprintf("random colours, run number %d", i), func(t *testing.T) {
			type bEncoder func([]byte, BarcodeExtraData) (barcode.Barcode, error)
			type testBarcode struct {
				name             string
				codeType         BarcodeType
				content          []byte
				extra            BarcodeExtraData
				start            image.Point
				width, height    int
				dataColour       color.Color
				backgroundColour color.Color
				err              error
				refData          []byte
			}
			tests := []testBarcode{
				testBarcode{
					name:             "qr",
					codeType:         BarcodeTypeQR,
					content:          []byte("www.github.com/LLKennedy/imagetemplate"),
					extra:            BarcodeExtraData{QRLevel: qr.Q, QRMode: qr.Unicode},
					start:            image.ZP,
					width:            130,
					height:           130,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcodeQR,
				},
				testBarcode{
					name:             "bad qr",
					codeType:         BarcodeTypeQR,
					content:          []byte("test"),
					extra:            BarcodeExtraData{QRMode: 1},
					start:            image.ZP,
					width:            130,
					height:           130,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf(`"test" can not be encoded as Numeric`),
				},
				testBarcode{
					name:             "aztec",
					codeType:         BarcodeTypeAztec,
					content:          []byte("www.github.com/LLKennedy/imagetemplate"),
					extra:            BarcodeExtraData{AztecMinECCPercent: 50, AztecUserSpecifiedLayers: 4},
					start:            image.ZP,
					width:            130,
					height:           130,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcodeAztec,
				},
				testBarcode{
					name:             "bad aztec",
					codeType:         BarcodeTypeAztec,
					content:          []byte("www.github.com/LLKennedy/imagetemplate"),
					extra:            BarcodeExtraData{AztecMinECCPercent: 50, AztecUserSpecifiedLayers: 150},
					start:            image.ZP,
					width:            130,
					height:           130,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Illegal value 150 for layers"),
				},
				testBarcode{
					name:             "pdf",
					codeType:         BarcodeTypePDF,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{PDFSecurityLevel: 4},
					start:            image.ZP,
					width:            130,
					height:           130,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcodePDF,
				},
				testBarcode{
					name:             "bad pdf",
					codeType:         BarcodeTypePDF,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{PDFSecurityLevel: 150},
					start:            image.ZP,
					width:            130,
					height:           130,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
				testBarcode{
					name:             "datamatrix",
					codeType:         BarcodeTypeDataMatrix,
					content:          []byte("https://www.github.com/LLKennedy/imagetemplate"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           130,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcodeDataMatrix,
				},
				testBarcode{
					name:             "bad datamatrix",
					codeType:         BarcodeTypeDataMatrix,
					content:          []byte("https://www.github.com/LLKennedy/imagetemplate"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           130,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
				testBarcode{
					name:             "nine of three",
					codeType:         BarcodeTypeCode93,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{Code93IncludeChecksum: true, Code93FullASCIIMode: true},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcode93,
				},
				testBarcode{
					name:             "bad nine of three",
					codeType:         BarcodeTypeCode93,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{Code93IncludeChecksum: true, Code93FullASCIIMode: true},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
				testBarcode{
					name:             "two of five",
					codeType:         BarcodeType2of5,
					content:          []byte("12345678"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcode25,
				},
				testBarcode{
					name:             "bad two of five",
					codeType:         BarcodeType2of5,
					content:          []byte("12345678"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
				testBarcode{
					name:             "two of five interleaved",
					codeType:         BarcodeType2of5Interleaved,
					content:          []byte("12345678"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcode25i,
				},
				testBarcode{
					name:             "bad two of five interleaved",
					codeType:         BarcodeType2of5Interleaved,
					content:          []byte("12345678"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
				testBarcode{
					name:             "codabar",
					codeType:         BarcodeTypeCodabar,
					content:          []byte("B123456D"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcodeCodabar,
				},
				testBarcode{
					name:             "bad codabar",
					codeType:         BarcodeTypeCodabar,
					content:          []byte("B123456D"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
				testBarcode{
					name:             "code128",
					codeType:         BarcodeTypeCode128,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcodeCode128,
				},
				testBarcode{
					name:             "bad code128",
					codeType:         BarcodeTypeCode128,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
				testBarcode{
					name:             "ean13",
					codeType:         BarcodeTypeEAN13,
					content:          []byte("5901234123457"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcodeEan13,
				},
				testBarcode{
					name:             "bad ean13",
					codeType:         BarcodeTypeEAN13,
					content:          []byte("5901234123457"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
				testBarcode{
					name:             "ean8",
					codeType:         BarcodeTypeEAN8,
					content:          []byte("11223344"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcodeEan8,
				},
				testBarcode{
					name:             "bad ean8",
					codeType:         BarcodeTypeEAN8,
					content:          []byte("11223344"),
					extra:            BarcodeExtraData{},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
				testBarcode{
					name:             "three of nine",
					codeType:         BarcodeTypeCode39,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{Code39IncludeChecksum: true, Code39FullASCIIMode: true},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					refData:          barcode39,
				},
				testBarcode{
					name:             "bad three of nine",
					codeType:         BarcodeTypeCode39,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{Code39IncludeChecksum: true, Code39FullASCIIMode: true},
					start:            image.ZP,
					width:            130,
					height:           65,
					dataColour:       color.Black,
					backgroundColour: color.White,
					err:              fmt.Errorf("Invalid security level 150"),
				},
			}
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					freshCanvas, _ := NewCanvas(test.width, test.height)
					modifiedCanvas, err := freshCanvas.Barcode(test.codeType, test.content, test.extra, test.start, test.width, test.height, test.dataColour, test.backgroundColour)
					if test.err == nil {
						assert.NoError(t, err)
						if err != nil {
							t.Fatal(err)
						}
						var imageBytes bytes.Buffer
						err = bmp.Encode(&imageBytes, modifiedCanvas.GetUnderlyingImage())
						assert.NoError(t, err)
						assert.Equal(t, test.refData, imageBytes.Bytes())
					} else {
						assert.EqualError(t, err, test.err.Error())
					}
				})
			}
		})
	}
}
