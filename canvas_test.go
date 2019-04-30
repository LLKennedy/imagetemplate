package imagetemplate

import (
	"bytes"
	"fmt"
	// "github.com/boombuler/barcode"
	// "github.com/boombuler/barcode/aztec"
	// "github.com/boombuler/barcode/codabar"
	// "github.com/boombuler/barcode/code128"
	// "github.com/boombuler/barcode/code39"
	// "github.com/boombuler/barcode/code93"
	// "github.com/boombuler/barcode/datamatrix"
	// "github.com/boombuler/barcode/ean"
	// "github.com/boombuler/barcode/pdf417"
	// "github.com/boombuler/barcode/qr"
	// "math/rand"
	// "github.com/boombuler/barcode/twooffive"
	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestBlankCanvas(t *testing.T) {
	blankCanvas := ImageCanvas{}
	t.Run("circle", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.Circle(image.ZP, 10, color.White)
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "No image set for canvas to draw on", err.Error())
		}
	})
	t.Run("barcode", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.Barcode(BarcodeTypeAztec, []byte{}, BarcodeExtraData{}, image.ZP, 0, 0, color.White, color.Black)
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "No image set for canvas to draw on", err.Error())
		}
	})
	t.Run("draw image", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.DrawImage(image.ZP, image.NewNRGBA(image.Rect(0, 0, 10, 10)))
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "No image set for canvas to draw on", err.Error())
		}
	})
	t.Run("get height", func(t *testing.T) {
		height := blankCanvas.GetHeight()
		assert.Equal(t, -1, height)
	})
	t.Run("get image func", func(t *testing.T) {
		img := blankCanvas.GetUnderlyingImage()
		assert.Nil(t, img)
	})
	t.Run("get width", func(t *testing.T) {
		width := blankCanvas.GetWidth()
		assert.Equal(t, -1, width)
	})
	t.Run("rectangle", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.Rectangle(image.ZP, 10, 10, color.White)
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "No image set for canvas to draw on", err.Error())
		}
	})
	t.Run("text", func(t *testing.T) {
		modifiedCanvas, err := blankCanvas.Text("hello", image.ZP, nil, color.White, 100)
		assert.Equal(t, blankCanvas, modifiedCanvas)
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "No image set for canvas to draw on", err.Error())
		}
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
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "Invalid maxWidth", err.Error())
		}
	})
	t.Run("valid text", func(t *testing.T) {
		purple := color.NRGBA{R: 130, B: 200, A: 255}
		modifiedCanvas, err := newCanvas.Text("test", image.Pt(2, 15), regFont, purple, 28)
		assert.Nil(t, err)
		assert.NotNil(t, modifiedCanvas)
		var data bytes.Buffer
		err = bmp.Encode(&data, modifiedCanvas.GetUnderlyingImage())
		assert.Nil(t, err)
		assert.Equal(t, testTextOutput, data.Bytes())
	})
	t.Run("oversized text", func(t *testing.T) {
		purple := color.NRGBA{R: 130, B: 200, A: 255}
		_, err := newCanvas.Text("test this much longer string which definitely won't fit", image.Pt(2, 15), regFont, purple, 28)
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "Resultant drawn text was longer than maxWidth", err.Error())
		}
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
	assert.Nil(t, err)
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

/*
func TestBarcode(t *testing.T) {
	for i := 0; i < 1; i++ {
		t.Run(fmt.Sprintf("random colours, run number %d", i), func(t *testing.T) {
			var grandCanvas Canvas
			grandWidth := 520
			grandHeight := grandWidth / 2
			squareSize := grandWidth / 4
			flatHeight := grandWidth / 8
			grandCanvas, err := NewCanvas(grandWidth, grandHeight)
			assert.NoError(t, err)
			colours := make([]color.Color, 24)
			for c := range colours {
				newBytes := make([]byte, 4)
				rand.Read(newBytes)
				colours[c] = color.NRGBA{R: uint8(newBytes[0]), G: uint8(newBytes[1]), B: uint8(newBytes[2]), A: uint8(newBytes[3])}
			}
			colourNum := -1
			nextColour := func() color.Color {
				colourNum++
				return colours[colourNum]
			}
			type bEncoder func([]byte, BarcodeExtraData) (barcode.Barcode, error)
			type testBarcode struct {
				name             string
				encodeFunc       bEncoder
				codeType         BarcodeType
				content          []byte
				extra            BarcodeExtraData
				start            image.Point
				width, height    int
				dataColour       color.Color
				backgroundColour color.Color
			}
			tests := []testBarcode{
				testBarcode{
					name: "qr",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return qr.Encode(string(content), extra.QRLevel, extra.QRMode)
					},
					codeType:         BarcodeTypeQR,
					content:          []byte("www.github.com/LLKennedy/imagetemplate"),
					extra:            BarcodeExtraData{QRLevel: qr.Q, QRMode: qr.Unicode},
					start:            image.ZP,
					width:            squareSize,
					height:           squareSize,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
				testBarcode{
					name: "aztec",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return aztec.Encode(content, extra.AztecMinECCPercent, extra.AztecUserSpecifiedLayers)
					},
					codeType:         BarcodeTypeAztec,
					content:          []byte("www.github.com/LLKennedy/imagetemplate"),
					extra:            BarcodeExtraData{AztecMinECCPercent: 50, AztecUserSpecifiedLayers: 4},
					start:            image.Point{X: squareSize, Y: 0},
					width:            squareSize,
					height:           squareSize,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
				testBarcode{
					name: "pdf",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return pdf417.Encode(string(content), extra.PDFSecurityLevel)
					},
					codeType:         BarcodeTypePDF,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{PDFSecurityLevel: 4},
					start:            image.Point{X: squareSize * 2, Y: 0},
					width:            squareSize,
					height:           squareSize,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
				testBarcode{
					name: "datamatrix",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return datamatrix.Encode(string(content))
					},
					codeType:         BarcodeTypeDataMatrix,
					content:          []byte("https://www.github.com/LLKennedy/imagetemplate"),
					extra:            BarcodeExtraData{},
					start:            image.Point{X: squareSize * 3, Y: 0},
					width:            squareSize,
					height:           squareSize,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
				testBarcode{
					name: "nine of three",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return code93.Encode(string(content), extra.Code93IncludeChecksum, extra.Code93FullAsciiMode)
					},
					codeType:         BarcodeTypeCode93,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{Code93IncludeChecksum: true, Code93FullAsciiMode: true},
					start:            image.Point{X: 0, Y: squareSize},
					width:            squareSize,
					height:           flatHeight,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
				// testBarcode{
				// 	name: "two of five",
				// 	encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
				// 		return twooffive.Encode(string(content), false)
				// 	},
				// 	codeType:         BarcodeType2of5,
				// 	content:          []byte("12345678"),
				// 	extra:            BarcodeExtraData{},
				// 	start:            image.Point{X: squareSize, Y: squareSize},
				// 	width:            squareSize,
				// 	height:           flatHeight,
				// 	dataColour:       nextColour(),
				// 	backgroundColour: nextColour(),
				// },
				// testBarcode{
				// 	name: "two of five interleaved",
				// 	encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
				// 		return twooffive.Encode(string(content), true)
				// 	},
				// 	codeType:         BarcodeType2of5,
				// 	content:          []byte("12345678"),
				// 	extra:            BarcodeExtraData{},
				// 	start:            image.Point{X: squareSize*2, Y: squareSize},
				// 	width:            squareSize,
				// 	height:           flatHeight,
				// 	dataColour:       nextColour(),
				// 	backgroundColour: nextColour(),
				// },
				testBarcode{
					name: "codabar",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return codabar.Encode(string(content))
					},
					codeType:         BarcodeTypeCodabar,
					content:          []byte("B123456D"),
					extra:            BarcodeExtraData{},
					start:            image.Point{X: squareSize * 3, Y: squareSize},
					width:            squareSize,
					height:           flatHeight,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
				testBarcode{
					name: "code128",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return code128.Encode(string(content))
					},
					codeType:         BarcodeTypeCode128,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{},
					start:            image.Point{X: 0, Y: flatHeight * 3},
					width:            squareSize,
					height:           flatHeight,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
				testBarcode{
					name: "ean13",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return ean.Encode(string(content))
					},
					codeType:         BarcodeTypeEAN13,
					content:          []byte("5901234123457"),
					extra:            BarcodeExtraData{},
					start:            image.Point{X: squareSize, Y: flatHeight * 3},
					width:            squareSize,
					height:           flatHeight,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
				testBarcode{
					name: "ean8",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return ean.Encode(string(content))
					},
					codeType:         BarcodeTypeEAN8,
					content:          []byte("11223344"),
					extra:            BarcodeExtraData{},
					start:            image.Point{X: squareSize * 2, Y: flatHeight * 3},
					width:            squareSize,
					height:           flatHeight,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
				testBarcode{
					name: "three of nine",
					encodeFunc: func(content []byte, extra BarcodeExtraData) (barcode.Barcode, error) {
						return code39.Encode(string(content), extra.Code39IncludeChecksum, extra.Code39FullAsciiMode)
					},
					codeType:         BarcodeTypeCode39,
					content:          []byte("Luke"),
					extra:            BarcodeExtraData{Code39IncludeChecksum: true, Code39FullAsciiMode: true},
					start:            image.Point{X: squareSize * 3, Y: flatHeight * 3},
					width:            squareSize,
					height:           flatHeight,
					dataColour:       nextColour(),
					backgroundColour: nextColour(),
				},
			}
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					grandCanvas, err = NewCanvas(grandWidth, grandHeight)
					rawBarcode, err := test.encodeFunc(test.content, test.extra)
					assert.NoError(t, err)
					if err != nil {
						t.Fatal(err)
					}
					rawBarcode, err = barcode.Scale(rawBarcode, test.width, test.height)
					assert.NoError(t, err)
					if err != nil {
						t.Fatal(err)
					}
					grandCanvas, err = grandCanvas.Barcode(test.codeType, test.content, test.extra, test.start, test.width, test.height, test.dataColour, test.backgroundColour)
					assert.NoError(t, err)
					underlyingImage := grandCanvas.GetUnderlyingImage()
					for x := 0; x < test.width; x++ {
						for y := 0; y < test.height; y++ {
							t.Run(fmt.Sprintf("x=%d,y=%d", x+test.start.X, y+test.start.Y), func(t *testing.T) {
								refPixel := rawBarcode.At(x, y)
								realPixel := underlyingImage.At(x+test.start.X, y+test.start.Y)
								if refPixel == color.White {
									assert.Equal(t, test.backgroundColour, realPixel)
								} else if refPixel == color.Black {
									assert.Equal(t, test.dataColour, realPixel)
								}
							})
						}
					}
				})
			}
		})
	}
}
*/

