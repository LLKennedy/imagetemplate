package imagetemplate

import (
	fs "github.com/LLKennedy/imagetemplate/internal/filesystem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	newBuilder, err := NewBuilder()
	assert.NoError(t, err)
	img := newBuilder.GetCanvas().GetUnderlyingImage()
	assert.Equal(t, 0, img.Bounds().Size().X)
	assert.Equal(t, 0, img.Bounds().Size().Y)
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
