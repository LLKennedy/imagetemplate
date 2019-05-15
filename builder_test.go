package imagetemplate

import (
	"fmt"
	fs "github.com/LLKennedy/imagetemplate/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/render"
	"github.com/stretchr/testify/assert"
	"image/color"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	dimensions := []int{-10000, -500, -33, -1, 0, 1, 33, 500, 10000}
	for _, x := range dimensions {
		for _, y := range dimensions {
			t.Run(fmt.Sprintf("Creating builder with dimensions %d by %d", x, y), func(t *testing.T) {
				var newBuilder Builder
				newCanvas, err := render.NewCanvas(x, y)
				if x <= 0 && y <= 0 && err != nil {
					assert.Nil(t, newBuilder)
					assert.EqualError(t, err, "invalid width and height")
					return
				}
				if x <= 0 && err != nil {
					assert.Nil(t, newBuilder)
					assert.EqualError(t, err, "invalid width")
					return
				}
				if y <= 0 && err != nil {
					assert.Nil(t, newBuilder)
					assert.EqualError(t, err, "invalid height")
					return
				}
				if err != nil {
					t.Fatalf("%v", err)
				}
				newBuilder, err = NewBuilder(newCanvas, nil)
				if err != nil {
					t.Fatalf("%v", err)
				}
				assert.NotNil(t, newBuilder)
				realBuilder, ok := newBuilder.(ImageBuilder)
				assert.True(t, ok)
				assert.NotNil(t, realBuilder)
				imageBounds := realBuilder.Canvas.GetUnderlyingImage().Bounds()
				assert.Equal(t, imageBounds.Size().X, x)
				assert.Equal(t, imageBounds.Size().Y, y)

			})
		}
	}
	specifiedColour := color.NRGBA{R: 123, G: 231, B: 132, A: 213}
	width := 50
	height := 50
	newCanvas, err := render.NewCanvas(width, height)
	assert.NoError(t, err)
	newBuilder, err := NewBuilder(newCanvas, &specifiedColour)
	assert.NoError(t, err)
	img := newBuilder.Canvas.GetUnderlyingImage()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			t.Run(fmt.Sprintf("Testing set pixel at %d, %d", x, y), func(t *testing.T) {
				setColour := img.At(x, y)
				assert.Equal(t, specifiedColour, setColour)
			})
		}
	}
}

func TestLoadComponentsData(t *testing.T) {
	t.Run("circles", func(t *testing.T) {
		sampleData := `{
			"baseImage": {
				"fileName": "baseone.bmp"
			},
			"components": [
				{
					"type": "circle",
					"properties": {
						"centreX": "145",
						"centreY": "300",
						"radius": "25",
						"colour": {
							"R": "0",
							"G": "255",
							"B": "0",
							"A": "255"
						}
					}
				},
				{
					"type": "circle",
					"properties": {
						"centreX": "380",
						"centreY": "154",
						"radius": "81",
						"colour": {
							"R": "0",
							"G": "0",
							"B": "255",
							"A": "255"
						}
					}
				},
				{
					"type": "circle",
					"properties": {
						"centreX": "297",
						"centreY": "185",
						"radius": "48",
						"colour": {
							"R": "255",
							"G": "127",
							"B": "39",
							"A": "255"
						}
					}
				},
				{
					"type": "circle",
					"properties": {
						"centreX": "133",
						"centreY": "388",
						"radius": "80",
						"colour": {
							"R": "127",
							"G": "127",
							"B": "127",
							"A": "255"
						}
					}
				},
				{
					"type": "circle",
					"properties": {
						"centreX": "350",
						"centreY": "390",
						"radius": "80",
						"colour": {
							"R": "255",
							"G": "174",
							"B": "201",
							"A": "255"
						}
					}
				},
				{
					"type": "circle",
					"properties": {
						"centreX": "105",
						"centreY": "105",
						"radius": "50",
						"colour": {
							"R": "255",
							"G": "0",
							"B": "0",
							"A": "255"
						}
					}
				}
			]
		}`
		var newBuilder Builder
		reader := fs.MockReader{Files: make(map[string]fs.MockFile)}
		reader.Files["baseone.bmp"] = fs.MockFile{Err: nil, Data: []byte{0x42, 0x4d, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x36, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xcc, 0x48, 0x3f, 0xff, 0xff, 0xff, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0xff, 0xff, 0xff, 0x24, 0x1c, 0xed, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0xf2, 0xff, 0x00}}
		newBuilder = ImageBuilder{reader: reader}
		newBuilder, err := newBuilder.LoadComponentsData([]byte(sampleData))
		assert.NoError(t, err)
		newBuilder, err = newBuilder.ApplyComponents()
		assert.NoError(t, err)
		//TODO: output and check results
	})
}
