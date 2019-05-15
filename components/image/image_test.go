package image

import (
	"bytes"
	"encoding/base64"
	"errors"
	fs "github.com/LLKennedy/imagetemplate/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/render"
	"github.com/stretchr/testify/assert"
	_ "golang.org/x/image/bmp" // bmp imported for image decoding
	"image"
	"testing"
)

func TestImageWrite(t *testing.T) {
	t.Run("not all props set", func(t *testing.T) {
		canvas := render.MockCanvas{}
		c := Component{NamedPropertiesMap: map[string][]string{"not set": []string{"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "cannot draw image, not all named properties are set: map[not set:[something]]")
	})
	t.Run("image error", func(t *testing.T) {
		canvas := render.MockCanvas{FixedDrawImageError: errors.New("some error")}
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "some error")
	})
	t.Run("passing", func(t *testing.T) {
		canvas := render.MockCanvas{FixedCircleError: nil}
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
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
	sampleTinyImageString := base64.RawStdEncoding.EncodeToString(sampleTinyImageData)
	sampleTinyImageBuffer := bytes.NewBuffer(sampleTinyImageData)
	sampleTinyImage, _, err := image.Decode(bytes.NewBuffer(sampleTinyImageData))
	assert.NoError(t, err, "failed to import sample image")
	tests := []testSet{
		testSet{
			name:  "no props",
			start: Component{},
			input: render.NamedProperties{},
			res:   Component{},
			err:   "",
		},
		testSet{
			name: "data invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"data"},
				},
			},
			input: render.NamedProperties{
				"aProp": 3,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"data"},
				},
			},
			err: "error converting 3 to []byte, string or io.Reader",
		},
		testSet{
			name: "data bytes",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"data"},
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
		testSet{
			name: "data string",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"data"},
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
		testSet{
			name: "data reader",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"data"},
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
		testSet{
			name: "data image error",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"data"},
				},
			},
			input: render.NamedProperties{
				"aProp": []byte{0x00, 0x00, 0x00, 0x00},
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"data"},
				},
			},
			err: "image: unknown format",
		},
		testSet{
			name: "filename invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"fileName"},
				},
			},
			input: render.NamedProperties{
				"aProp": 3,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"fileName"},
				},
			},
			err: "error converting 3 to string",
		},
		testSet{
			name: "file load error",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"fileName"},
				},
				reader: fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Data: sampleTinyImageData, Err: errors.New("file not found")}}},
			},
			input: render.NamedProperties{
				"aProp": "somefile.jpg",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"fileName"},
				},
				reader: fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Data: sampleTinyImageData, Err: errors.New("file not found")}}},
			},
			err: "file not found",
		},
		testSet{
			name: "image file data invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"fileName"},
				},
				reader: fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Data: []byte{0x00, 0x00, 0x00, 0x00}, Err: nil}}},
			},
			input: render.NamedProperties{
				"aProp": "somefile.jpg",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"fileName"},
				},
				reader: fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Data: []byte{0x00, 0x00, 0x00, 0x00}, Err: nil}}},
			},
			err: "image: unknown format",
		},
		testSet{
			name: "filename valid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"fileName"},
				},
				reader: fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Data: sampleTinyImageData, Err: nil}}},
			},
			input: render.NamedProperties{
				"aProp": "somefile.jpg",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Image:              sampleTinyImage,
				reader:             fs.MockReader{Files: map[string]fs.MockFile{"somefile.jpg": fs.MockFile{Data: sampleTinyImageData, Err: nil}}},
			},
			err: "",
		},
		testSet{
			name: "other invalid type",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"not a prop"},
				},
			},
			input: render.NamedProperties{
				"aProp": "not a number",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"not a prop"},
				},
			},
			err: "error converting not a number to int",
		},
		testSet{
			name: "other invalid name",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"not a prop"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"not a prop"},
				},
			},
			err: "invalid component property in named property map: not a prop",
		},
		testSet{
			name: "topLeftX",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"topLeftX"},
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
		testSet{
			name: "topLeftY",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"topLeftY"},
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
		testSet{
			name: "width",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"width"},
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
		testSet{
			name: "height",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"height"},
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
			res, err := test.start.SetNamedProperties(test.input)
			test.res.reader = nil
			ICres := res.(Component)
			ICres.reader = nil
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
	tests := []testSet{
		testSet{
			name:  "incorrect format data",
			start: Component{},
			input: "hello",
			res:   Component{},
			props: render.NamedProperties{},
			err:   "failed to convert returned data to component properties",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
}
