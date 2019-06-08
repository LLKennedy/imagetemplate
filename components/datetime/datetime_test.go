package datetime

import (
	"fmt"
	"image"
	"image/color"
	"testing"
	"time"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/LLKennedy/imagetemplate/render"
	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
)

func TestDateTimeWrite(t *testing.T) {
	goreg, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("not all props set", func(t *testing.T) {
		canvas := render.MockCanvas{}
		c := Component{NamedPropertiesMap: map[string][]string{"not set": []string{"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "failed to write to canvas: runtime error: invalid memory address or nil pointer dereference")
		canvas.AssertExpectations(t)
	})
	t.Run("datetime error", func(t *testing.T) {
		canvas := render.MockCanvas{FixedTextError: fmt.Errorf("some error"), FixedTryTextBool: true, FixedTryTextInt: 10, FixedPixelsPerInch: 72}
		timeVal := time.Now()
		c := Component{Font: goreg, Size: 14, MaxWidth: 100, Time: &timeVal}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "some error")
		canvas.AssertExpectations(t)
	})
	t.Run("passing", func(t *testing.T) {
		canvas := render.MockCanvas{FixedTextError: nil, FixedTryTextBool: true, FixedTryTextInt: 10, FixedPixelsPerInch: 72}
		timeVal := time.Now()
		c := Component{Font: goreg, Size: 14, MaxWidth: 100, Time: &timeVal}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
		canvas.AssertExpectations(t)
	})
}

func TestDateTimeSetNamedProperties(t *testing.T) {
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
					"aProp": []string{"R"},
				},
			},
			input: render.NamedProperties{
				"aProp": "not a number",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"R"},
				},
			},
			err: "error converting not a number to uint8",
		},
		testSet{
			name: "RGBA valid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"R", "G", "B", "A"},
				},
			},
			input: render.NamedProperties{
				"aProp": uint8(1),
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Colour:             color.NRGBA{R: uint8(1), G: uint8(1), B: uint8(1), A: uint8(1)},
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
			name: "startX",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"startX"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Start:              image.Pt(15, 0),
			},
			err: "",
		},
		testSet{
			name: "startY",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"startY"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Start:              image.Pt(0, 15),
			},
			err: "",
		},
		testSet{
			name: "size",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"size"},
				},
			},
			input: render.NamedProperties{
				"aProp": float64(15),
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Size:               15,
			},
			err: "",
		},
		testSet{
			name: "full prop set, multiple sources, unused props",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"col1": []string{"R", "G", "A"},
					"left": []string{"startX", "startY"},
					"wide": []string{"size"},
					"col3": []string{"B"},
					"what": []string{"R", "G", "B", "A", "startX"},
				},
			},
			input: render.NamedProperties{
				"col1":             uint8(15),
				"col2":             uint8(6),
				"up":               50,
				"some other thing": "doesn't matter",
				"col3":             uint8(150),
				"wide":             float64(80),
				"left":             3,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{"what": []string{"R", "G", "B", "A", "startX"}},
				Start:              image.Pt(3, 3),
				Size:               80,
				Colour:             color.NRGBA{R: uint8(15), G: uint8(15), B: uint8(150), A: uint8(15)},
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

func TestDateTimeGetJSONFormat(t *testing.T) {
	c := Component{}
	expectedFormat := &datetimeFormat{}
	format := c.GetJSONFormat()
	assert.Equal(t, expectedFormat, format)
}

func TestDateTimeVerifyAndTestDateTimeJSONData(t *testing.T) {
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

func TestInit(t *testing.T) {
	c, err := render.Decode("date")
	assert.NoError(t, err)
	assert.Equal(t, Component{}, c)
}
