package rectangle

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/stretchr/testify/assert"
)

func TestRectangleWrite(t *testing.T) {
	t.Run("not all props set", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		c := Component{NamedPropertiesMap: map[string][]string{"not set": {"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "cannot draw rectangle, not all named properties are set: map[not set:[something]]")
		canvas.AssertExpectations(t)
	})
	t.Run("rectangle error", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		canvas.On("Rectangle", image.Pt(0, 0), 0, 0, color.NRGBA{}).Return(canvas, fmt.Errorf("some error"))
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "some error")
		canvas.AssertExpectations(t)
	})
	t.Run("passing", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		canvas.On("Rectangle", image.Pt(0, 0), 0, 0, color.NRGBA{}).Return(canvas, nil)
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
		canvas.AssertExpectations(t)
	})
}

func TestRectangleSetNamedProperties(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input render.NamedProperties
		res   Component
		err   string
	}
	tests := []testSet{
		{
			name:  "no props",
			start: Component{},
			input: render.NamedProperties{},
			res:   Component{},
			err:   "",
		},
		{
			name: "RGBA invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"R"},
				},
			},
			input: render.NamedProperties{
				"aProp": "not a number",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"R"},
				},
			},
			err: "error converting not a number to uint8",
		},
		{
			name: "RGBA valid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"R", "G", "B", "A"},
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
		{
			name: "non-RGBA invalid name",
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
		{
			name: "full prop set, multiple sources, unused props",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"col1": {"R", "G", "A"},
					"left": {"topLeftX"},
					"wide": {"width", "topLeftY"},
					"col3": {"B"},
					"what": {"R", "G", "B", "A", "topLeftX"},
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
				NamedPropertiesMap: map[string][]string{"what": {"R", "G", "B", "A", "topLeftX"}},
				TopLeft:            image.Pt(3, 80),
				Width:              80,
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

func TestRectangleGetJSONFormat(t *testing.T) {
	c := Component{}
	expectedFormat := &rectangleFormat{}
	format := c.GetJSONFormat()
	assert.Equal(t, expectedFormat, format)
}

func TestRectangleVerifyAndTestRectangleJSONData(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input interface{}
		res   Component
		props render.NamedProperties
		err   string
	}
	tests := []testSet{
		{
			name:  "incorrect format data",
			input: "hello",
			props: render.NamedProperties{},
			err:   "failed to convert returned data to component properties",
		},
		{
			name:  "invalid topLeftX",
			start: Component{},
			input: &rectangleFormat{},
			res:   Component{},
			props: render.NamedProperties{},
			err:   "error parsing data for property topLeftX: could not parse empty property",
		},
		{
			name:  "valid topLeftX",
			start: Component{},
			input: &rectangleFormat{
				TopLeftX: "180",
			},
			res:   Component{},
			props: render.NamedProperties{},
			err:   "error parsing data for property topLeftY: could not parse empty property",
		},
		{
			name:  "valid topLeftY",
			start: Component{},
			input: &rectangleFormat{
				TopLeftX: "180",
				TopLeftY: "37",
			},
			res:   Component{},
			props: render.NamedProperties{},
			err:   "error parsing data for property width: could not parse empty property",
		},
		{
			name:  "valid width",
			start: Component{},
			input: &rectangleFormat{
				TopLeftX: "180",
				TopLeftY: "37",
				Width:    "300",
			},
			res:   Component{},
			props: render.NamedProperties{},
			err:   "error parsing data for property height: could not parse empty property",
		},
		{
			name:  "valid height",
			start: Component{},
			input: &rectangleFormat{
				TopLeftX: "180",
				TopLeftY: "37",
				Width:    "300",
				Height:   "150",
			},
			res:   Component{},
			props: render.NamedProperties{},
			err:   "error parsing data for property R: could not parse empty property",
		},
		{
			name:  "valid R",
			start: Component{},
			input: &rectangleFormat{
				TopLeftX: "180",
				TopLeftY: "37",
				Width:    "300",
				Height:   "150",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red: "192",
				},
			},
			res:   Component{},
			props: render.NamedProperties{},
			err:   "error parsing data for property G: could not parse empty property",
		},
		{
			name:  "valid G",
			start: Component{},
			input: &rectangleFormat{
				TopLeftX: "180",
				TopLeftY: "37",
				Width:    "300",
				Height:   "150",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "192",
					Green: "1",
				},
			},
			res:   Component{},
			props: render.NamedProperties{},
			err:   "error parsing data for property B: could not parse empty property",
		},
		{
			name:  "valid B",
			start: Component{},
			input: &rectangleFormat{
				TopLeftX: "180",
				TopLeftY: "37",
				Width:    "300",
				Height:   "150",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "192",
					Green: "1",
					Blue:  "66",
				},
			},
			res:   Component{},
			props: render.NamedProperties{},
			err:   "error parsing data for property A: could not parse empty property",
		},
		{
			name:  "valid everything",
			start: Component{},
			input: &rectangleFormat{
				TopLeftX: "180",
				TopLeftY: "37",
				Width:    "$myWidth$",
				Height:   "150",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "192",
					Green: "1",
					Blue:  "66",
					Alpha: "201",
				},
			},
			res: Component{
				TopLeft:            image.Pt(180, 37),
				Height:             150,
				Colour:             color.NRGBA{R: 192, G: 1, B: 66, A: 201},
				NamedPropertiesMap: map[string][]string{"myWidth": {"width"}},
			},
			props: render.NamedProperties{"myWidth": struct{ Message string }{Message: "Please replace me with real data"}},
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
	c, err := render.Decode("rectangle")
	assert.NoError(t, err)
	assert.Equal(t, Component{}, c)
}
