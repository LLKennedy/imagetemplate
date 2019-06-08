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
	newBuilder := NewBuilder()
	assert.Equal(t, ImageBuilder{reader: fs.IoutilFileReader{}}, newBuilder)
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
		canvas := new(render.MockCanvas)
		canvas.On("GetUnderlyingImage").Return(&img)
		newBuilder = newBuilder.SetCanvas(canvas)
		data, err := newBuilder.WriteToBMP()
		assert.EqualError(t, err, "bmp: negative bounds")
		assert.Nil(t, data)
		canvas.AssertExpectations(t)
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
			assert.Equal(t, test.result, result)
			if test.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.err.Error())
			}
		})
	}
	canvas := new(render.MockCanvas)
	canvas.On("GetWidth").Return(10)
	canvas.On("GetHeight").Return(10)
	canvas.On("DrawImage", image.Point{}, &image.NRGBA{Pix:[]uint8{0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe, 0xfa, 0xc, 0x50, 0xbe}, Stride:40, Rect:image.Rectangle{Min:image.Point{X:0, Y:0}, Max:image.Point{X:10, Y:10}}}).Return(canvas, fmt.Errorf("error drawing on canvas"))
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
			name:     "scaled base colour (wide)",
			builder:  ImageBuilder{}.SetCanvas(render.ImageCanvas{}.SetUnderlyingImage(&image.NRGBA{Pix: []uint8{31, 63, 127, 255}, Rect: image.Rect(0, 0, 1, 1), Stride: 4})).(ImageBuilder),
			template: Template{BaseImage: BaseImage{BaseWidth: "2", BaseHeight: "1", BaseColour: BaseColour{Red: "250", Green: "12", Blue: "80", Alpha: "190"}}},
			result:   ImageBuilder{}.SetCanvas(render.ImageCanvas{}.SetUnderlyingImage(&image.NRGBA{Pix: []uint8{0xc2, 0x19, 0x5c, 0xff}, Rect: image.Rect(0, 0, 1, 1), Stride: 4})).(ImageBuilder),
		},
		testSet{
			name:     "scaled base colour (tall)",
			builder:  ImageBuilder{}.SetCanvas(render.ImageCanvas{}.SetUnderlyingImage(&image.NRGBA{Pix: []uint8{31, 63, 127, 255}, Rect: image.Rect(0, 0, 1, 1), Stride: 4})).(ImageBuilder),
			template: Template{BaseImage: BaseImage{BaseWidth: "1", BaseHeight: "2", BaseColour: BaseColour{Red: "250", Green: "12", Blue: "80", Alpha: "190"}}},
			result:   ImageBuilder{}.SetCanvas(render.ImageCanvas{}.SetUnderlyingImage(&image.NRGBA{Pix: []uint8{0xc2, 0x19, 0x5c, 0xff}, Rect: image.Rect(0, 0, 1, 1), Stride: 4})).(ImageBuilder),
		},
		testSet{
			name:     "scaled base colour (too big)",
			builder:  ImageBuilder{}.SetCanvas(render.ImageCanvas{}.SetUnderlyingImage(&image.NRGBA{Pix: []uint8{31, 63, 127, 255}, Rect: image.Rect(0, 0, 1, 1), Stride: 4})).(ImageBuilder),
			template: Template{BaseImage: BaseImage{BaseWidth: "2", BaseHeight: "2", BaseColour: BaseColour{Red: "250", Green: "12", Blue: "80", Alpha: "190"}}},
			result:   ImageBuilder{}.SetCanvas(render.ImageCanvas{}.SetUnderlyingImage(&image.NRGBA{Pix: []uint8{0xc2, 0x19, 0x5c, 0xff}, Rect: image.Rect(0, 0, 1, 1), Stride: 4})).(ImageBuilder),
		},
		testSet{
			name:     "error drawing resized image",
			builder:  ImageBuilder{}.SetCanvas(canvas).(ImageBuilder),
			template: Template{BaseImage: BaseImage{BaseWidth: "2", BaseHeight: "2", BaseColour: BaseColour{Red: "250", Green: "12", Blue: "80", Alpha: "190"}}},
			result:   ImageBuilder{}.SetCanvas(canvas).(ImageBuilder),
			err:      fmt.Errorf("error drawing on canvas"),
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
		testSet{
			name:     "invalid base64",
			template: Template{BaseImage: BaseImage{Data: "a"}},
			err:      fmt.Errorf("unexpected EOF"),
		},
		testSet{
			name:     "base64 non-image",
			template: Template{BaseImage: BaseImage{Data: "TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlzIHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2YgdGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGludWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRoZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4="}},
			err:      fmt.Errorf("image: unknown format"),
		},
		testSet{
			name:     "valid jpeg",
			builder:  ImageBuilder{reader: fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Data: []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01, 0x01, 0x00, 0x60, 0x00, 0x60, 0x00, 0x00, 0xff, 0xdb, 0x00, 0x43, 0x00, 0x02, 0x01, 0x01, 0x02, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x05, 0x03, 0x03, 0x03, 0x03, 0x03, 0x06, 0x04, 0x04, 0x03, 0x05, 0x07, 0x06, 0x07, 0x07, 0x07, 0x06, 0x07, 0x07, 0x08, 0x09, 0x0b, 0x09, 0x08, 0x08, 0x0a, 0x08, 0x07, 0x07, 0x0a, 0x0d, 0x0a, 0x0a, 0x0b, 0x0c, 0x0c, 0x0c, 0x0c, 0x07, 0x09, 0x0e, 0x0f, 0x0d, 0x0c, 0x0e, 0x0b, 0x0c, 0x0c, 0x0c, 0xff, 0xdb, 0x00, 0x43, 0x01, 0x02, 0x02, 0x02, 0x03, 0x03, 0x03, 0x06, 0x03, 0x03, 0x06, 0x0c, 0x08, 0x07, 0x08, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0xff, 0xc0, 0x00, 0x11, 0x08, 0x00, 0x02, 0x00, 0x02, 0x03, 0x01, 0x22, 0x00, 0x02, 0x11, 0x01, 0x03, 0x11, 0x01, 0xff, 0xc4, 0x00, 0x1f, 0x00, 0x00, 0x01, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0xff, 0xc4, 0x00, 0xb5, 0x10, 0x00, 0x02, 0x01, 0x03, 0x03, 0x02, 0x04, 0x03, 0x05, 0x05, 0x04, 0x04, 0x00, 0x00, 0x01, 0x7d, 0x01, 0x02, 0x03, 0x00, 0x04, 0x11, 0x05, 0x12, 0x21, 0x31, 0x41, 0x06, 0x13, 0x51, 0x61, 0x07, 0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xa1, 0x08, 0x23, 0x42, 0xb1, 0xc1, 0x15, 0x52, 0xd1, 0xf0, 0x24, 0x33, 0x62, 0x72, 0x82, 0x09, 0x0a, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xc4, 0x00, 0x1f, 0x01, 0x00, 0x03, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0xff, 0xc4, 0x00, 0xb5, 0x11, 0x00, 0x02, 0x01, 0x02, 0x04, 0x04, 0x03, 0x04, 0x07, 0x05, 0x04, 0x04, 0x00, 0x01, 0x02, 0x77, 0x00, 0x01, 0x02, 0x03, 0x11, 0x04, 0x05, 0x21, 0x31, 0x06, 0x12, 0x41, 0x51, 0x07, 0x61, 0x71, 0x13, 0x22, 0x32, 0x81, 0x08, 0x14, 0x42, 0x91, 0xa1, 0xb1, 0xc1, 0x09, 0x23, 0x33, 0x52, 0xf0, 0x15, 0x62, 0x72, 0xd1, 0x0a, 0x16, 0x24, 0x34, 0xe1, 0x25, 0xf1, 0x17, 0x18, 0x19, 0x1a, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xda, 0x00, 0x0c, 0x03, 0x01, 0x00, 0x02, 0x11, 0x03, 0x11, 0x00, 0x3f, 0x00, 0xfd, 0x7a, 0xfd, 0x9e, 0xbf, 0x64, 0xff, 0x00, 0x85, 0x9f, 0x11, 0xfe, 0x01, 0x78, 0x1f, 0xc4, 0x3e, 0x21, 0xf8, 0x6b, 0xe0, 0x0d, 0x7b, 0x5f, 0xd7, 0xbc, 0x3f, 0x61, 0xa8, 0xea, 0x7a, 0x9e, 0xa3, 0xe1, 0xeb, 0x4b, 0xab, 0xcd, 0x46, 0xea, 0x6b, 0x68, 0xe4, 0x96, 0x79, 0xa5, 0x78, 0xcb, 0xc9, 0x2b, 0xbb, 0x33, 0x33, 0xb1, 0x2c, 0xcc, 0xc4, 0x92, 0x49, 0xa2, 0x8a, 0x2b, 0xe3, 0x73, 0x0f, 0xf7, 0xaa, 0x9f, 0xe2, 0x7f, 0x99, 0xfd, 0x71, 0xc3, 0xdf, 0xf2, 0x2a, 0xc3, 0x7f, 0xd7, 0xb8, 0x7f, 0xe9, 0x28, 0xff, 0xd9}}}}},
			template: Template{BaseImage: BaseImage{FileName: "somefile.jpg"}},
			result:   ImageBuilder{Canvas: render.ImageCanvas{Image: &image.NRGBA{Pix: []uint8{0xff, 0xf7, 0xff, 0xff, 0x63, 0x56, 0x66, 0xff, 0x82, 0x75, 0x86, 0xff, 0xff, 0xfa, 0xff, 0xff}, Stride: 8, Rect: image.Rect(0, 0, 2, 2)}}, reader: fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Data: []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01, 0x01, 0x00, 0x60, 0x00, 0x60, 0x00, 0x00, 0xff, 0xdb, 0x00, 0x43, 0x00, 0x02, 0x01, 0x01, 0x02, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x05, 0x03, 0x03, 0x03, 0x03, 0x03, 0x06, 0x04, 0x04, 0x03, 0x05, 0x07, 0x06, 0x07, 0x07, 0x07, 0x06, 0x07, 0x07, 0x08, 0x09, 0x0b, 0x09, 0x08, 0x08, 0x0a, 0x08, 0x07, 0x07, 0x0a, 0x0d, 0x0a, 0x0a, 0x0b, 0x0c, 0x0c, 0x0c, 0x0c, 0x07, 0x09, 0x0e, 0x0f, 0x0d, 0x0c, 0x0e, 0x0b, 0x0c, 0x0c, 0x0c, 0xff, 0xdb, 0x00, 0x43, 0x01, 0x02, 0x02, 0x02, 0x03, 0x03, 0x03, 0x06, 0x03, 0x03, 0x06, 0x0c, 0x08, 0x07, 0x08, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0x0c, 0xff, 0xc0, 0x00, 0x11, 0x08, 0x00, 0x02, 0x00, 0x02, 0x03, 0x01, 0x22, 0x00, 0x02, 0x11, 0x01, 0x03, 0x11, 0x01, 0xff, 0xc4, 0x00, 0x1f, 0x00, 0x00, 0x01, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0xff, 0xc4, 0x00, 0xb5, 0x10, 0x00, 0x02, 0x01, 0x03, 0x03, 0x02, 0x04, 0x03, 0x05, 0x05, 0x04, 0x04, 0x00, 0x00, 0x01, 0x7d, 0x01, 0x02, 0x03, 0x00, 0x04, 0x11, 0x05, 0x12, 0x21, 0x31, 0x41, 0x06, 0x13, 0x51, 0x61, 0x07, 0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xa1, 0x08, 0x23, 0x42, 0xb1, 0xc1, 0x15, 0x52, 0xd1, 0xf0, 0x24, 0x33, 0x62, 0x72, 0x82, 0x09, 0x0a, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xc4, 0x00, 0x1f, 0x01, 0x00, 0x03, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0xff, 0xc4, 0x00, 0xb5, 0x11, 0x00, 0x02, 0x01, 0x02, 0x04, 0x04, 0x03, 0x04, 0x07, 0x05, 0x04, 0x04, 0x00, 0x01, 0x02, 0x77, 0x00, 0x01, 0x02, 0x03, 0x11, 0x04, 0x05, 0x21, 0x31, 0x06, 0x12, 0x41, 0x51, 0x07, 0x61, 0x71, 0x13, 0x22, 0x32, 0x81, 0x08, 0x14, 0x42, 0x91, 0xa1, 0xb1, 0xc1, 0x09, 0x23, 0x33, 0x52, 0xf0, 0x15, 0x62, 0x72, 0xd1, 0x0a, 0x16, 0x24, 0x34, 0xe1, 0x25, 0xf1, 0x17, 0x18, 0x19, 0x1a, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xda, 0x00, 0x0c, 0x03, 0x01, 0x00, 0x02, 0x11, 0x03, 0x11, 0x00, 0x3f, 0x00, 0xfd, 0x7a, 0xfd, 0x9e, 0xbf, 0x64, 0xff, 0x00, 0x85, 0x9f, 0x11, 0xfe, 0x01, 0x78, 0x1f, 0xc4, 0x3e, 0x21, 0xf8, 0x6b, 0xe0, 0x0d, 0x7b, 0x5f, 0xd7, 0xbc, 0x3f, 0x61, 0xa8, 0xea, 0x7a, 0x9e, 0xa3, 0xe1, 0xeb, 0x4b, 0xab, 0xcd, 0x46, 0xea, 0x6b, 0x68, 0xe4, 0x96, 0x79, 0xa5, 0x78, 0xcb, 0xc9, 0x2b, 0xbb, 0x33, 0x33, 0xb1, 0x2c, 0xcc, 0xc4, 0x92, 0x49, 0xa2, 0x8a, 0x2b, 0xe3, 0x73, 0x0f, 0xf7, 0xaa, 0x9f, 0xe2, 0x7f, 0x99, 0xfd, 0x71, 0xc3, 0xdf, 0xf2, 0x2a, 0xc3, 0x7f, 0xd7, 0xb8, 0x7f, 0xe9, 0x28, 0xff, 0xd9}}}}},
		},
		testSet{
			name:     "file error",
			builder:  ImageBuilder{reader: fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Err: fmt.Errorf("some file error")}}}},
			template: Template{BaseImage: BaseImage{FileName: "somefile.jpg"}},
			result:   ImageBuilder{reader: fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Err: fmt.Errorf("some file error")}}}},
			err:      fmt.Errorf("some file error"),
		},
	}
	for _, test := range tests {
		testFunc(test, t)
	}
}

type mockComponent struct{ data int }

func (c mockComponent) Write(canvas render.Canvas) (render.Canvas, error) {
	return nil, nil
}

func (c mockComponent) SetNamedProperties(properties render.NamedProperties) (render.Component, error) {
	return c, nil
}

func (c mockComponent) GetJSONFormat() interface{} {
	return &struct {
		SomeProp string `json:"someProp"`
	}{}
}

func (c mockComponent) VerifyAndSetJSONData(data interface{}) (render.Component, render.NamedProperties, error) {
	realData := data.(*struct {
		SomeProp string `json:"someProp"`
	})
	if realData.SomeProp == "fail" {
		return c, nil, fmt.Errorf("failed to set props from JSON")
	} else if realData.SomeProp == "giveProp" {
		return c, render.NamedProperties{"aprop": struct{ Message string }{Message: "Please replace this struct with real data"}}, nil
	}
	return c, nil, nil
}

func newMock() render.Component {
	return &mockComponent{}
}

func TestParseComponents(t *testing.T) {
	render.RegisterComponent("mock", newMock)
	type testSet struct {
		name        string
		templates   []ComponentTemplate
		toggleables []ToggleableComponent
		props       render.NamedProperties
		err         error
	}
	testFunc := func(test testSet, t *testing.T) {
		toggleables, props, err := parseComponents(test.templates)
		assert.Equal(t, test.toggleables, toggleables)
		assert.Equal(t, test.props, props)
		if test.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.err.Error())
		}
	}
	tests := []testSet{
		testSet{
			name:  "empty templates",
			props: render.NamedProperties{},
		},
		testSet{
			name:      "error decoding type",
			templates: []ComponentTemplate{ComponentTemplate{Type: "something wrong", Conditional: render.ComponentConditional{Name: "$myVar$"}}},
			props:     render.NamedProperties{"$myVar$": struct{ Message string }{Message: "Please replace this struct with real data"}},
			err:       fmt.Errorf("component error: no component registered for name something wrong"),
		},
		testSet{
			name:        "mock component with named props",
			templates:   []ComponentTemplate{ComponentTemplate{Type: "mock", Properties: []byte(`{"someProp":"giveProp"}`)}},
			props:       render.NamedProperties{"aprop": struct{ Message string }{Message: "Please replace this struct with real data"}},
			toggleables: []ToggleableComponent{ToggleableComponent{Component: mockComponent{}}},
		},
		testSet{
			name:      "mock component fails to set data",
			templates: []ComponentTemplate{ComponentTemplate{Type: "mock", Properties: []byte(`{"someProp":"fail"}`)}},
			props:     render.NamedProperties{},
			err:       fmt.Errorf("failed to set props from JSON"),
		},
		testSet{
			name:      "mock component with invalid JSON",
			templates: []ComponentTemplate{ComponentTemplate{Type: "mock", Properties: []byte(`:"fail"}`)}},
			props:     render.NamedProperties{},
			err:       fmt.Errorf("invalid character ':' looking for beginning of value"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { testFunc(test, t) })
	}
}

func TestLoadComponentsData(t *testing.T) {
	render.RegisterComponent("mock", newMock)
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
	t.Run("failing on setbackgroundimage", func(t *testing.T) {

	})
}

func TestGetComponents(t *testing.T) {
	goodCondition := render.ComponentConditional{Name: "myVar", Operator: "ci_equals", Value: "something"}
	goodCondition, _ = goodCondition.SetValue("myVar", "something")
	badCondition := render.ComponentConditional{Name: "myVar", Operator: "ci_equals", Value: "something"}
	badCondition, _ = badCondition.SetValue("myVar", "nothing")
	builder := ImageBuilder{Components: []ToggleableComponent{
		ToggleableComponent{
			Conditional: goodCondition,
			Component:   mockComponent{data: 1},
		},
		ToggleableComponent{
			Conditional: badCondition,
			Component:   mockComponent{data: 2},
		},
	}}
	components := builder.GetComponents()
	expectedComponents := []render.Component{
		mockComponent{data: 1},
	}
	assert.Equal(t, expectedComponents, components)
}

func TestSetComponents(t *testing.T) {
	builder := ImageBuilder{}
	assert.Equal(t, []ToggleableComponent(nil), builder.Components)
	builder = builder.SetComponents([]ToggleableComponent{}).(ImageBuilder)
	assert.Equal(t, []ToggleableComponent{}, builder.Components)
}

func TestGetNamedPropertiesList(t *testing.T) {
	builder := ImageBuilder{}
	assert.Equal(t, render.NamedProperties(nil), builder.GetNamedPropertiesList())
	builder.NamedProperties = render.NamedProperties{"something": "something else"}
	assert.Equal(t, render.NamedProperties{"something": "something else"}, builder.GetNamedPropertiesList())
}
