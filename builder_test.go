package imagetemplate

import (
	"fmt"
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
				newCanvas, err := NewCanvas(x, y)
				if x <= 0 && y <= 0 && err != nil {
					assert.Nil(t, newBuilder)
					assert.Equal(t, "Invalid width and height", err.Error())
					return
				}
				if x <= 0 && err != nil {
					assert.Nil(t, newBuilder)
					assert.Equal(t, "Invalid width", err.Error())
					return
				}
				if y <= 0 && err != nil {
					assert.Nil(t, newBuilder)
					assert.Equal(t, "Invalid height", err.Error())
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
	newCanvas, err := NewCanvas(width, height)
	assert.Nil(t, err)
	newBuilder, err := NewBuilder(newCanvas, &specifiedColour)
	assert.Nil(t, err)
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
		newCanvas, err := NewCanvas(500, 500)
		if err != nil {
			t.Fatalf("canvas returned error: %v", err)
		}
		sampleData := `{
			"baseImage": {
				"width": "500",
				"height": "500",
				"baseColour": {
					"R": "220", 
					"G": "220", 
					"B": "220",
					"A": "255"
				}
			},
			"components": [
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
		newBuilder, err = NewBuilder(newCanvas, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		assert.Nil(t, err)
		newBuilder, err = newBuilder.LoadComponentsData([]byte(sampleData))
		assert.Nil(t, err)
	})
}
