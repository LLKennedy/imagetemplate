package circle

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"image/color"
	"image"
	"errors"
	"github.com/LLKennedy/imagetemplate/render"
)

func TestCircleWrite(t *testing.T) {
	t.Run("not all props set", func(t *testing.T) {
		canvas := render.MockCanvas{}
		c := Component{NamedPropertiesMap: map[string][]string{"not set":[]string{"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "cannot draw circle, not all named properties are set: map[not set:[something]]")
	})
	t.Run("circle error", func(t *testing.T) {
		canvas := render.MockCanvas{FixedCircleError: errors.New("some error")}
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

func TestCircleSetNamedProperties(t *testing.T) {
	type testSet struct{
		name string
		start Component
		input render.NamedProperties
		res Component
		err string
	}
	tests := []testSet{
		testSet{
			name: "no props",
			start: Component{},
			input: render.NamedProperties{},
			res: Component{},
			err: "",
		},
		testSet{
			name: "RGBA invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"R"},
				},
			},
			input: render.NamedProperties{
				"aProp": "not a number",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"R"},
				},
			},
			err: "error converting not a number to uint8",
		},
		testSet{
			name: "RGBA valid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"R","G","B","A"},
				},
			},
			input: render.NamedProperties{
				"aProp": uint8(1),
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Colour: color.NRGBA{R:uint8(1),G:uint8(1),B:uint8(1),A:uint8(1)},
			},
			err: "",
		},
		testSet{
			name: "non-RGBA invalid type",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			input: render.NamedProperties{
				"aProp": "not a number",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			err: "error converting not a number to int",
		},
		testSet{
			name: "non-RGBA invalid name",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			err: "invalid component property in named property map: not a prop",
		},
		testSet{
			name: "centreX",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"centreX"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Centre: image.Pt(15, 0),
			},
			err: "",
		},
		testSet{
			name: "centreY",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"centreY"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Centre: image.Pt(0, 15),
			},
			err: "",
		},
		testSet{
			name: "radius",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"radius"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Radius: 15,
			},
			err: "",
		},
		testSet{
			name: "full prop set, multiple sources, unused props",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"col1": []string{"R","G","A"},
					"left": []string{"centreX"},
					"wide": []string{"radius", "centreY"},
					"col3": []string{"B"},
					"what": []string{"R", "G", "B", "A", "centreX"},
				},
			},
			input: render.NamedProperties{
				"col1": uint8(15),
				"col2": uint8(6),
				"up": 50,
				"some other thing": "doesn't matter",
				"col3": uint8(150),
				"wide": 80,
				"left": 3,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{"what": []string{"R", "G", "B", "A", "centreX"}},
				Centre: image.Pt(3, 80),
				Radius: 80,
				Colour: color.NRGBA{R: uint8(15), G: uint8(15), B: uint8(150), A: uint8(15)},
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

func TestCircleGetJSONFormat(t *testing.T) {
	c := Component{}
	expectedFormat := &circleFormat{}
	format := c.GetJSONFormat()
	assert.Equal(t, expectedFormat, format)
}

func TestCircleVerifyAndTestCircleJSONData(t *testing.T) {
	type testSet struct{
		name string
		start Component
		input interface{}
		res Component
		props render.NamedProperties
		err string
	}
	tests := []testSet{
		testSet{
			name: "incorrect format data",
			start: Component{},
			input: "hello",
			res: Component{},
			props: render.NamedProperties{},
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