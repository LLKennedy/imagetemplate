package barcode

import (
	"errors"
	"image"
	"image/color"
	"testing"

	"github.com/LLKennedy/imagetemplate/render"
	"github.com/stretchr/testify/assert"
)

func TestBarcodeWrite(t *testing.T) {
	t.Run("not all props set", func(t *testing.T) {
		canvas := render.MockCanvas{}
		c := Component{NamedPropertiesMap: map[string][]string{"not set": []string{"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "cannot draw barcode, not all named properties are set: map[not set:[something]]")
	})
	t.Run("barcode error", func(t *testing.T) {
		canvas := render.MockCanvas{FixedBarcodeError: errors.New("some error")}
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "some error")
	})
	t.Run("passing", func(t *testing.T) {
		canvas := render.MockCanvas{FixedBarcodeError: nil}
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
	})
}

func TestBarcodeSetNamedProperties(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input render.NamedProperties
		res   Component
		err   string
	}
	tests := []testSet{
		testSet{
			name:  "no props",
			start: Component{},
			input: render.NamedProperties{},
			res:   Component{},
			err:   "",
		},
		testSet{
			name: "RGBA invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"dR"},
				},
			},
			input: render.NamedProperties{
				"aProp": "not a number",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"dR"},
				},
			},
			err: "error converting not a number to uint8",
		},
		testSet{
			name: "dRGBA valid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"dR", "dG", "dB", "dA"},
				},
			},
			input: render.NamedProperties{
				"aProp": uint8(1),
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				DataColour:         color.NRGBA{R: uint8(1), G: uint8(1), B: uint8(1), A: uint8(1)},
			},
			err: "",
		},
		testSet{
			name: "non-RGBA invalid type",
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
			name: "non-RGBA invalid name",
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
			name: "full prop set, multiple sources, unused props",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"col1": []string{"dR", "dG", "dA"},
					"left": []string{"topLeftX"},
					"wide": []string{"width", "height", "topLeftY"},
					"col3": []string{"dB"},
					"what": []string{"bR", "bG", "bB", "bA", "topLeftX"},
				},
			},
			input: render.NamedProperties{
				"col1":             uint8(15),
				"col2":             uint8(6),
				"up":               50,
				"some other thing": "doesn't matter",
				"col3":             uint8(150),
				"wide":             80,
				"left":             3,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{"what": []string{"bR", "bG", "bB", "bA", "topLeftX"}},
				TopLeft:            image.Pt(3, 80),
				Width:              80,
				Height:             80,
				DataColour:         color.NRGBA{R: uint8(15), G: uint8(15), B: uint8(150), A: uint8(15)},
			},
			err: "",
		},
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

func TestBarcodeGetJSONFormat(t *testing.T) {
	c := Component{}
	expectedFormat := &barcodeFormat{}
	format := c.GetJSONFormat()
	assert.Equal(t, expectedFormat, format)
}

func TestBarcodeVerifyAndTestBarcodeJSONData(t *testing.T) {
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
