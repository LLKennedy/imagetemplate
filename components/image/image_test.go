package image

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"runtime/debug"
	"testing"

	fs "github.com/LLKennedy/imagetemplate/v3/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/stretchr/testify/assert"
	_ "golang.org/x/image/bmp" // bmp imported for image decoding
	"golang.org/x/tools/godoc/vfs"
)

func TestImageWrite(t *testing.T) {
	t.Run("not all props set", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		c := Component{NamedPropertiesMap: map[string][]string{"not set": {"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "cannot draw image, not all named properties are set: map[not set:[something]]")
		canvas.AssertExpectations(t)
	})
	t.Run("image error", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		canvas.On("DrawImage", image.Point{}, &image.NRGBA{}).Return(canvas, fmt.Errorf("some error"))
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "some error")
		canvas.AssertExpectations(t)
	})
	t.Run("passing", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		canvas.On("DrawImage", image.Point{}, &image.NRGBA{}).Return(canvas, nil)
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
		canvas.AssertExpectations(t)
	})
}

func TestImageSetNamedProperties(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input render.NamedProperties
		res   Component
		err   string
	}
	//Pure white 2x2 image
	sampleTinyImageData := []byte{0x42, 0x4d, 0x46, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x36, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00}
	sampleTinyImageString := base64.StdEncoding.EncodeToString(sampleTinyImageData)
	sampleTinyImageBuffer := bytes.NewBuffer(sampleTinyImageData)
	sampleTinyImage, _, err := image.Decode(bytes.NewBuffer(sampleTinyImageData))
	assert.NoError(t, err, "failed to import sample image")
	tests := []testSet{
		{
			name:  "no props",
			start: Component{},
			input: render.NamedProperties{},
			res:   Component{},
			err:   "",
		},
		{
			name: "data invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"data"},
				},
			},
			input: render.NamedProperties{
				"aProp": 3,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"data"},
				},
			},
			err: "error converting 3 to []byte, string or io.Reader",
		},
		{
			name: "data bytes",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"data"},
				},
			},
			input: render.NamedProperties{
				"aProp": sampleTinyImageData,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Image:              sampleTinyImage,
			},
			err: "",
		},
		{
			name: "data string",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"data"},
				},
			},
			input: render.NamedProperties{
				"aProp": sampleTinyImageString,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Image:              sampleTinyImage,
			},
			err: "",
		},
		{
			name: "data reader",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"data"},
				},
			},
			input: render.NamedProperties{
				"aProp": sampleTinyImageBuffer,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Image:              sampleTinyImage,
			},
			err: "",
		},
		{
			name: "data image error",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"data"},
				},
			},
			input: render.NamedProperties{
				"aProp": []byte{0x00, 0x00, 0x00, 0x00},
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"data"},
				},
			},
			err: "image: unknown format",
		},
		{
			name: "filename invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fileName"},
				},
			},
			input: render.NamedProperties{
				"aProp": 3,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fileName"},
				},
			},
			err: "error converting 3 to string",
		},
		{
			name: "file load error",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fileName"},
				},
				fs: func() vfs.FileSystem {
					reader := fs.NewMockFileSystem()
					reader.On("Open", "somefile.jpg").Return(fs.NewMockFile("", sampleTinyImageData), errors.New("file not found"))
					return reader
				}(),
			},
			input: render.NamedProperties{
				"aProp": "somefile.jpg",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fileName"},
				},
				fs: func() vfs.FileSystem {
					reader := fs.NewMockFileSystem()
					reader.On("Open", "somefile.jpg").Return(fs.NewMockFile("", sampleTinyImageData), errors.New("file not found"))
					return reader
				}(),
			},
			err: "file not found",
		},
		{
			name: "image file data invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fileName"},
				},
				fs: func() vfs.FileSystem {
					reader := fs.NewMockFileSystem()
					reader.On("Open", "somefile.jpg").Return(fs.NewMockFile("", []byte{0x00, 0x00, 0x00, 0x00}), nil)
					return reader
				}(),
			},
			input: render.NamedProperties{
				"aProp": "somefile.jpg",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fileName"},
				},
				fs: func() vfs.FileSystem {
					reader := fs.NewMockFileSystem()
					reader.On("Open", "somefile.jpg").Return(fs.NewMockFile("", []byte{0x00, 0x00, 0x00, 0x00}), nil)
					return reader
				}(),
			},
			err: "image: unknown format",
		},
		{
			name: "filename valid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fileName"},
				},
				fs: func() vfs.FileSystem {
					reader := fs.NewMockFileSystem()
					reader.On("Open", "somefile.jpg").Return(fs.NewMockFile("", sampleTinyImageData), nil)
					return reader
				}(),
			},
			input: render.NamedProperties{
				"aProp": "somefile.jpg",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Image:              sampleTinyImage,
				fs: func() vfs.FileSystem {
					reader := fs.NewMockFileSystem()
					reader.On("Open", "somefile.jpg").Return(fs.NewMockFile("", sampleTinyImageData), nil)
					return reader
				}(),
			},
			err: "",
		},
		{
			name: "other invalid type",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"not a prop"},
				},
			},
			input: render.NamedProperties{
				"aProp": "not a number",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"not a prop"},
				},
			},
			err: "error converting not a number to int",
		},
		{
			name: "other invalid name",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"not a prop"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"not a prop"},
				},
			},
			err: "invalid component property in named property map: not a prop",
		},
		{
			name: "topLeftX",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"topLeftX"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				TopLeft:            image.Pt(15, 0),
			},
			err: "",
		},
		{
			name: "topLeftY",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"topLeftY"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				TopLeft:            image.Pt(0, 15),
			},
			err: "",
		},
		{
			name: "width",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"width"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Width:              15,
			},
			err: "",
		},
		{
			name: "height",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"height"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Height:             15,
			},
			err: "",
		},
		// testSet{
		// 	name: "full prop set, multiple sources, unused props",
		// 	start: Component{
		// 		NamedPropertiesMap: map[string][]string{
		// 			"col1": []string{"R","G","A"},
		// 			"left": []string{"centreX"},
		// 			"wide": []string{"radius", "centreY"},
		// 			"col3": []string{"B"},
		// 			"what": []string{"R", "G", "B", "A", "centreX"},
		// 		},
		// 	},
		// 	input: render.NamedProperties{
		// 		"col1": uint8(15),
		// 		"col2": uint8(6),
		// 		"up": 50,
		// 		"some other thing": "doesn't matter",
		// 		"col3": uint8(150),
		// 		"wide": 80,
		// 		"left": 3,
		// 	},
		// 	res: Component{
		// 		NamedPropertiesMap: map[string][]string{"what": []string{"R", "G", "B", "A", "centreX"}},
		// 		Centre: image.Pt(3, 80),
		// 		Radius: 80,
		// 		Colour: color.NRGBA{R: uint8(15), G: uint8(15), B: uint8(150), A: uint8(15)},
		// 	},
		// 	err: "",
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					assert.NoError(t, r.(error))
				}
			}()
			res, err := test.start.SetNamedProperties(test.input)
			test.res.fs = nil
			ICres := res.(Component)
			ICres.fs = nil
			assert.Equal(t, test.res, ICres)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.err)
			}
		})
	}
}

func TestImageGetJSONFormat(t *testing.T) {
	c := Component{}
	expectedFormat := &imageFormat{}
	format := c.GetJSONFormat()
	assert.Equal(t, expectedFormat, format)
}

func TestImageVerifyAndTestImageJSONData(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input interface{}
		res   Component
		props render.NamedProperties
		err   string
	}
	imageFS := fs.NewMockFileSystem(
		fs.NewMockFile("myImage.png", []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xde, 0x00, 0x00, 0x00, 0x01, 0x73, 0x52, 0x47, 0x42, 0x00, 0xae, 0xce, 0x1c, 0xe9, 0x00, 0x00, 0x00, 0x04, 0x67, 0x41, 0x4d, 0x41, 0x00, 0x00, 0xb1, 0x8f, 0x0b, 0xfc, 0x61, 0x05, 0x00, 0x00, 0x00, 0x09, 0x70, 0x48, 0x59, 0x73, 0x00, 0x00, 0x0e, 0xc3, 0x00, 0x00, 0x0e, 0xc3, 0x01, 0xc7, 0x6f, 0xa8, 0x64, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44, 0x41, 0x54, 0x18, 0x57, 0x63, 0xf8, 0xff, 0xff, 0x3f, 0x00, 0x05, 0xfe, 0x02, 0xfe, 0xa7, 0x35, 0x81, 0x84, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}),
		fs.NewMockFile("badImage.bmp", []byte{}),
	)
	imageFS.On("Open", "nilImage.bmp").Return(fs.NilFile, nil)
	tests := []testSet{
		{
			name:  "incorrect format data",
			input: "hello",
			props: render.NamedProperties{},
			err:   "failed to convert returned data to component properties",
		},
		{
			name:  "invalid image type",
			input: &imageFormat{},
			props: render.NamedProperties{},
			err:   "exactly one of (fileName,data) must be set",
		},
		{
			name:  "nil image file",
			start: Component{
				fs: imageFS,
			},
			input: &imageFormat{
				FileName: "nilImage.bmp",
			},
			res: Component{
				fs: imageFS,
			},
			props: render.NamedProperties{},
			err:   "image: unknown format",
		},
		{
			name:  "bad image file",
			start: Component{
				fs: imageFS,
			},
			input: &imageFormat{
				FileName: "badImage.bmp",
			},
			res: Component{
				fs: imageFS,
			},
			props: render.NamedProperties{},
			err:   "image: unknown format",
		},
		{
			name:  "valid image file",
			start: Component{
				fs: imageFS,
			},
			input: &imageFormat{
				FileName: "myImage.png",
			},
			res: Component{
				fs: imageFS,
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property topLeftX: could not parse empty property",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
				}
			}()
			res, props, err := test.start.VerifyAndSetJSONData(test.input)
			assert.Equal(t, test.res, res)
			assert.Equal(t, test.props, props)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.err)
			}
		})
	}
	imageFS.AssertExpectations(t)
}

func TestInit(t *testing.T) {
	c, err := render.Decode("image")
	assert.NoError(t, err)
	assert.Equal(t, Component{fs: vfs.OS(".")}, c)
}
