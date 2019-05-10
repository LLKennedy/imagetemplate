package imagetemplate

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"image/color"
	"image"
	"errors"
)

func TestWrite(t *testing.T) {
	t.Run("not all props set", func(t *testing.T) {
		canvas := mockCanvas{}
		c := CircleComponent{NamedPropertiesMap: map[string][]string{"not set":[]string{"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "cannot draw circle, not all named properties are set: map[not set:[something]]")
	})
	t.Run("circle error", func(t *testing.T) {
		canvas := mockCanvas{FixedCircleError: errors.New("some error")}
		c := CircleComponent{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "some error")
	})
	t.Run("passing", func(t *testing.T) {
		canvas := mockCanvas{FixedCircleError: nil}
		c := CircleComponent{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
	})
}

func TestSetNamedProperties(t *testing.T) {
	type testSet struct{
		name string
		start CircleComponent
		input NamedProperties
		res CircleComponent
		err string
	}
	tests := []testSet{
		testSet{
			name: "no props",
			start: CircleComponent{},
			input: NamedProperties{},
			res: CircleComponent{},
			err: "",
		},
		testSet{
			name: "RGBA invalid",
			start: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"R"},
				},
			},
			input: NamedProperties{
				"aProp": "not a number",
			},
			res: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"R"},
				},
			},
			err: "error converting not a number to uint8",
		},
		testSet{
			name: "RGBA valid",
			start: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"R","G","B","A"},
				},
			},
			input: NamedProperties{
				"aProp": uint8(1),
			},
			res: CircleComponent{
				NamedPropertiesMap: map[string][]string{},
				Colour: color.NRGBA{R:uint8(1),G:uint8(1),B:uint8(1),A:uint8(1)},
			},
			err: "",
		},
		testSet{
			name: "non-RGBA invalid type",
			start: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			input: NamedProperties{
				"aProp": "not a number",
			},
			res: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			err: "error converting not a number to int",
		},
		testSet{
			name: "non-RGBA invalid name",
			start: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			input: NamedProperties{
				"aProp": 12,
			},
			res: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"not a prop"},
				},
			},
			err: "invalid component property in named property map: not a prop",
		},
		testSet{
			name: "centreX",
			start: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"centreX"},
				},
			},
			input: NamedProperties{
				"aProp": 15,
			},
			res: CircleComponent{
				NamedPropertiesMap: map[string][]string{},
				Centre: image.Pt(15, 0),
			},
			err: "",
		},
		testSet{
			name: "centreY",
			start: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"centreY"},
				},
			},
			input: NamedProperties{
				"aProp": 15,
			},
			res: CircleComponent{
				NamedPropertiesMap: map[string][]string{},
				Centre: image.Pt(0, 15),
			},
			err: "",
		},
		testSet{
			name: "radius",
			start: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"aProp":[]string{"radius"},
				},
			},
			input: NamedProperties{
				"aProp": 15,
			},
			res: CircleComponent{
				NamedPropertiesMap: map[string][]string{},
				Radius: 15,
			},
			err: "",
		},
		testSet{
			name: "full prop set, multiple sources, unused props",
			start: CircleComponent{
				NamedPropertiesMap: map[string][]string{
					"col1": []string{"R","G","A"},
					"left": []string{"centreX"},
					"wide": []string{"radius", "centreY"},
					"col3": []string{"B"},
					"what": []string{"R", "G", "B", "A", "centreX"},
				},
			},
			input: NamedProperties{
				"col1": uint8(15),
				"col2": uint8(6),
				"up": 50,
				"some other thing": "doesn't matter",
				"col3": uint8(150),
				"wide": 80,
				"left": 3,
			},
			res: CircleComponent{
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

func TestGetJSONFormat(t *testing.T) {
	c := CircleComponent{}
	expectedFormat := &circleFormat{}
	format := c.GetJSONFormat()
	assert.Equal(t, expectedFormat, format)
}

func TestVerifyAndTestJSONData(t *testing.T) {
	type testSet struct{
		name string
		start CircleComponent
		input interface{}
		res CircleComponent
		props NamedProperties
		err string
	}
	tests := []testSet{
		testSet{
			name: "incorrect format data",
			start: CircleComponent{},
			input: "hello",
			res: CircleComponent{},
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