package imagetemplate

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	fs "github.com/LLKennedy/imagetemplate/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/render"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	newBuilder, err := NewBuilder()
	assert.NoError(t, err)
	img := newBuilder.GetCanvas().GetUnderlyingImage()
	assert.Equal(t, 0, img.Bounds().Size().X)
	assert.Equal(t, 0, img.Bounds().Size().Y)
}

type fakeImage struct {
	at         color.Color
	bounds     image.Rectangle
	colorModel color.Model
}

func (i *fakeImage) At(x, y int) color.Color {
	return i.at
}

func (i *fakeImage) Bounds() image.Rectangle {
	return i.bounds
}

func (i *fakeImage) ColorModel() color.Model {
	return i.colorModel
}

func TestWriteToBMP(t *testing.T) {
	t.Run("nil image", func(t *testing.T) {
		newBuilder := ImageBuilder{}
		data, err := newBuilder.WriteToBMP()
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x42, 0x4d, 0x36, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x36, 0x0, 0x0, 0x0, 0x28, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x18, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, data)
	})
	t.Run("bad image", func(t *testing.T) {
		newBuilder := Builder(ImageBuilder{})
		img := fakeImage{bounds: image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(-1, -1)}}
		canvas := render.MockCanvas{
			FixedGetUnderlyingImage: &img,
		}
		newBuilder = newBuilder.SetCanvas(canvas)
		data, err := newBuilder.WriteToBMP()
		assert.EqualError(t, err, "bmp: negative bounds")
		assert.Nil(t, data)
	})
}

func TestLoadComponentsFile(t *testing.T) {
	builder := ImageBuilder{}
	t.Run("file error", func(t *testing.T) {
		builder.reader = fs.MockReader{Files: map[string]fs.MockFile{"myfile.json": fs.MockFile{Err: fmt.Errorf("some file error")}}}
		newBuilder, err := builder.LoadComponentsFile("myfile.json")
		assert.Equal(t, builder, newBuilder)
		assert.EqualError(t, err, "some file error")
	})
	t.Run("no file error", func(t *testing.T) {
		builder.reader = fs.MockReader{Files: map[string]fs.MockFile{"myfile.json": fs.MockFile{Data: []byte("hello")}}}
		newBuilder, err := builder.LoadComponentsFile("myfile.json")
		assert.Equal(t, builder, newBuilder)
		assert.EqualError(t, err, "invalid character 'h' looking for beginning of value")
	})
}

func TestSetBackgroundImageData(t *testing.T) {
	type testSet struct {
		name     string
		builder  ImageBuilder
		template Template
		result   ImageBuilder
		err      error
	}
	testFunc := func(test testSet, t *testing.T) {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.builder.setBackgroundImage(test.template)
			assert.Equal(t, result, test.result)
			if test.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.err.Error())
			}
		})
	}
	tests := []testSet{
		testSet{
			name:   "no base image properties",
			result: ImageBuilder{}.SetCanvas(ImageBuilder{}.GetCanvas()).(ImageBuilder),
		},
		testSet{
			name:     "invalid exclusive properties",
			template: Template{BaseImage: BaseImage{FileName: "something", Data: "something else"}},
			err:      fmt.Errorf("cannot load base image from file and load from data string and generate from base colour, specify only data or fileName or base colour"),
		},
		testSet{
			name:     "valid base colour",
			template: Template{BaseImage: BaseImage{BaseWidth: "1", BaseHeight: "1", BaseColour: BaseColour{Red: "250", Green: "12", Blue: "80", Alpha: "190"}}},
			result:   ImageBuilder{}.SetCanvas(render.ImageCanvas{}.SetUnderlyingImage(&image.NRGBA{Pix: []uint8{250, 12, 80, 190}, Rect: image.Rect(0, 0, 1, 1), Stride: 4})).(ImageBuilder),
		},
		testSet{
			name:     "base colour invalid width",
			template: Template{BaseImage: BaseImage{BaseWidth: "a", BaseHeight: "2", BaseColour: BaseColour{Red: "1"}}},
			err:      fmt.Errorf("strconv.ParseInt: parsing \"a\": invalid syntax"),
		},
		testSet{
			name:     "base colour invalid height",
			template: Template{BaseImage: BaseImage{BaseWidth: "2", BaseHeight: "a", BaseColour: BaseColour{Red: "1"}}},
			err:      fmt.Errorf("strconv.ParseInt: parsing \"a\": invalid syntax"),
		},
		testSet{
			name:     "base colour invalid red",
			template: Template{BaseImage: BaseImage{BaseWidth: "2", BaseHeight: "2", BaseColour: BaseColour{Red: "a"}}},
			err:      fmt.Errorf("strconv.ParseUint: parsing \"a\": invalid syntax"),
		},
		testSet{
			name:     "base colour invalid green",
			template: Template{BaseImage: BaseImage{BaseWidth: "2", BaseHeight: "2", BaseColour: BaseColour{Red: "1", Green: "a"}}},
			err:      fmt.Errorf("strconv.ParseUint: parsing \"a\": invalid syntax"),
		},
		testSet{
			name:     "base colour invalid blue",
			template: Template{BaseImage: BaseImage{BaseWidth: "2", BaseHeight: "2", BaseColour: BaseColour{Red: "1", Green: "1", Blue: "a"}}},
			err:      fmt.Errorf("strconv.ParseUint: parsing \"a\": invalid syntax"),
		},
		testSet{
			name:     "base colour invalid alpha",
			template: Template{BaseImage: BaseImage{BaseWidth: "2", BaseHeight: "2", BaseColour: BaseColour{Red: "1", Green: "1", Blue: "1", Alpha: "a"}}},
			err:      fmt.Errorf("strconv.ParseUint: parsing \"a\": invalid syntax"),
		},
	}
	for _, test := range tests {
		testFunc(test, t)
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
