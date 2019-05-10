package imagetemplate

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"image"
	"errors"
	"encoding/base64"
	"bytes"
	_ "golang.org/x/image/bmp"  // bmp imported for image decoding
)

func TestImageWrite(t *testing.T) {
	t.Run("not all props set", func(t *testing.T) {
		canvas := mockCanvas{}
		c := ImageComponent{NamedPropertiesMap: map[string][]string{"not set":[]string{"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "cannot draw image, not all named properties are set: map[not set:[something]]")
	})
	t.Run("image error", func(t *testing.T) {
		canvas := mockCanvas{FixedDrawImageError: errors.New("some error")}
		c := ImageComponent{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "some error")
	})
	t.Run("passing", func(t *testing.T) {
		canvas := mockCanvas{FixedCircleError: nil}
		c := ImageComponent{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
	})
}

func TestImageSetNamedProperties(t *testing.T) {
	type testSet struct{
		name string
		start ImageComponent
		input NamedProperties
		res ImageComponent
		err string
	}
	//Pure white 2x2 image
	sampleTinyImageData := []byte{0x42,0x4d,0x46,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x36,0x00,0x00,0x00,0x28,0x00,0x00,0x00,0x02,0x00,0x00,0x00,0x02,0x00,0x00,0x00,0x01,0x00,0x18,0x00,0x00,0x00,0x00,0x00,0x10,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0xff,0xff,0xff,0xff,0xff,0xff,0x00,0x00,0xff,0xff,0xff,0xff,0xff,0xff,0x00,0x00}
	sampleTinyImageString := base64.RawStdEncoding.EncodeToString(sampleTinyImageData)
	sampleTinyImageBuffer := bytes.NewBuffer(sampleTinyImageData)
	sampleTinyImage, _, err := image.Decode(bytes.NewBuffer(sampleTinyImageData))
	assert.NoError(t, err, "failed to import sample image")
	tests := []testSet{
		testSet{
			name: "no props",
			start: ImageComponent{},
			input: NamedProperties{},
			res: ImageComponent{},
			err: "",
		},
		testSet{
			name: "data invalid",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"data"},
				},
			},
			input: NamedProperties{
				"aProp": 3,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"data"},
				},
			},
			err: "error converting 3 to []byte, string or io.Reader",
		},
		testSet{
			name: "data bytes",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"data"},
				},
			},
			input: NamedProperties{
				"aProp": sampleTinyImageData,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{},
				Image: sampleTinyImage,
			},
			err: "",
		},
		testSet{
			name: "data string",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"data"},
				},
			},
			input: NamedProperties{
				"aProp": sampleTinyImageString,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{},
				Image: sampleTinyImage,
			},
			err: "",
		},
		testSet{
			name: "data reader",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"data"},
				},
			},
			input: NamedProperties{
				"aProp": sampleTinyImageBuffer,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{},
				Image: sampleTinyImage,
			},
			err: "",
		},
		testSet{
			name: "data image error",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"data"},
				},
			},
			input: NamedProperties{
				"aProp": []byte{0x00,0x00,0x00,0x00},
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"data"},
				},
			},
			err: "image: unknown format",
		},
		testSet{
			name: "filename invalid",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"fileName"},
				},
			},
			input: NamedProperties{
				"aProp": 3,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"fileName"},
				},
			},
			err: "error converting 3 to string",
		},
		testSet{
			name: "file load error",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"fileName"},
				},
				reader: mockReader{files: map[string]mockFile{"somefile.jpg":mockFile{data: sampleTinyImageData, err:errors.New("file not found")}}},
			},
			input: NamedProperties{
				"aProp": "somefile.jpg",
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"fileName"},
				},
				reader: mockReader{files: map[string]mockFile{"somefile.jpg":mockFile{data: sampleTinyImageData, err:errors.New("file not found")}}},
			},
			err: "file not found",
		},
		testSet{
			name: "image file data invalid",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"fileName"},
				},
				reader: mockReader{files: map[string]mockFile{"somefile.jpg":mockFile{data: []byte{0x00,0x00,0x00,0x00}, err:nil}}},
			},
			input: NamedProperties{
				"aProp": "somefile.jpg",
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"fileName"},
				},
				reader: mockReader{files: map[string]mockFile{"somefile.jpg":mockFile{data: []byte{0x00,0x00,0x00,0x00}, err:nil}}},
			},
			err: "image: unknown format",
		},
		testSet{
			name: "filename valid",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"fileName"},
				},
				reader: mockReader{files: map[string]mockFile{"somefile.jpg":mockFile{data: sampleTinyImageData, err:nil}}},
			},
			input: NamedProperties{
				"aProp": "somefile.jpg",
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{},
				Image: sampleTinyImage,
				reader: mockReader{files: map[string]mockFile{"somefile.jpg":mockFile{data: sampleTinyImageData, err:nil}}},
			},
			err: "",
		},
		testSet{
			name: "other invalid type",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			input: NamedProperties{
				"aProp": "not a number",
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			err: "error converting not a number to int",
		},
		testSet{
			name: "other invalid name",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			input: NamedProperties{
				"aProp": 12,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			err: "invalid component property in named property map: not a prop",
		},
		testSet{
			name: "topLeftX",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"topLeftX"},
				},
			},
			input: NamedProperties{
				"aProp": 15,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{},
				TopLeft: image.Pt(15, 0),
			},
			err: "",
		},
		testSet{
			name: "topLeftY",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"topLeftY"},
				},
			},
			input: NamedProperties{
				"aProp": 15,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{},
				TopLeft: image.Pt(0, 15),
			},
			err: "",
		},
		testSet{
			name: "width",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"width"},
				},
			},
			input: NamedProperties{
				"aProp": 15,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{},
				Width: 15,
			},
			err: "",
		},
		testSet{
			name: "height",
			start: ImageComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"height"},
				},
			},
			input: NamedProperties{
				"aProp": 15,
			},
			res: ImageComponent{
				NamedPropertiesMap: map[string][]string{},
				Height: 15,
			},
			err: "",
		},
		// testSet{
		// 	name: "full prop set, multiple sources, unused props",
		// 	start: ImageComponent{
		// 		NamedPropertiesMap: map[string][]string{
		// 			"col1": []string{"R","G","A"},
		// 			"left": []string{"centreX"},
		// 			"wide": []string{"radius", "centreY"},
		// 			"col3": []string{"B"},
		// 			"what": []string{"R", "G", "B", "A", "centreX"},
		// 		},
		// 	},
		// 	input: NamedProperties{
		// 		"col1": uint8(15),
		// 		"col2": uint8(6),
		// 		"up": 50,
		// 		"some other thing": "doesn't matter",
		// 		"col3": uint8(150),
		// 		"wide": 80,
		// 		"left": 3,
		// 	},
		// 	res: ImageComponent{
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
			assert.Equal(t, test.res, res)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.err)
			}
		})
	}
}

func TestImageGetJSONFormat(t *testing.T) {
	c := ImageComponent{}
	expectedFormat := &imageFormat{}
	format := c.GetJSONFormat()
	assert.Equal(t, expectedFormat, format)
}

func TestImageVerifyAndTestImageJSONData(t *testing.T) {
	type testSet struct{
		name string
		start ImageComponent
		input interface{}
		res ImageComponent
		props NamedProperties
		err string
	}
	tests := []testSet{
		testSet{
			name: "incorrect format data",
			start: ImageComponent{},
			input: "hello",
			res: ImageComponent{},
			props: NamedProperties{},
			err: "failed to convert returned data to component properties",
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