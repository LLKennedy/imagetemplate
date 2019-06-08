package barcode

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/LLKennedy/imagetemplate/v2/render"
	"github.com/boombuler/barcode/qr"
	"github.com/stretchr/testify/assert"
)

func TestBarcodeWrite(t *testing.T) {
	t.Run("not all props set", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		c := Component{NamedPropertiesMap: map[string][]string{"not set": []string{"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "cannot draw barcode, not all named properties are set: map[not set:[something]]")
		canvas.AssertExpectations(t)
	})
	t.Run("barcode error", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		canvas.On("Barcode", render.BarcodeType(""), []byte{}, render.BarcodeExtraData{}, image.Point{}, 0, 0, color.NRGBA{}, color.NRGBA{}).Return(canvas, fmt.Errorf("some error"))
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "some error")
		canvas.AssertExpectations(t)
	})
	t.Run("passing", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		canvas.On("Barcode", render.BarcodeType(""), []byte{}, render.BarcodeExtraData{}, image.Point{}, 0, 0, color.NRGBA{}, color.NRGBA{}).Return(canvas, nil)
		c := Component{}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
		canvas.AssertExpectations(t)
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
			name: "content invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"content"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"content"},
				},
			},
			err: "error converting 15 to string",
		},
		testSet{
			name: "type invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"barcodeType"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"barcodeType"},
				},
			},
			err: "error converting 15 to barcode type",
		},
		testSet{
			name: "unsupported type",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"barcodeType"},
				},
			},
			input: render.NamedProperties{
				"aProp": "15",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"barcodeType"},
				},
			},
			err: "error converting 15 to barcode type",
		},
		testSet{
			name: "invalid colour code",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"Rd"},
				},
			},
			input: render.NamedProperties{
				"aProp": uint8(12),
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": []string{"Rd"},
				},
			},
			err: "name was a string inside RGBA and Value was a valid uint8, but Name wasn't R, G, B, or A. Name was: Rd",
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
			name: "aztec",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": []string{"barcodeType"},
				},
			},
			input: render.NamedProperties{
				"type": render.BarcodeTypeAztec,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Type:               render.BarcodeTypeAztec,
				Extra:              render.BarcodeExtraData{AztecMinECCPercent: 50, AztecUserSpecifiedLayers: 4},
			},
			err: "",
		},
		testSet{
			name: "39",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": []string{"barcodeType"},
				},
			},
			input: render.NamedProperties{
				"type": render.BarcodeTypeCode39,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Type:               render.BarcodeTypeCode39,
				Extra:              render.BarcodeExtraData{Code39IncludeChecksum: true, Code39FullASCIIMode: true},
			},
			err: "",
		},
		testSet{
			name: "93",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": []string{"barcodeType"},
				},
			},
			input: render.NamedProperties{
				"type": render.BarcodeTypeCode93,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Type:               render.BarcodeTypeCode93,
				Extra:              render.BarcodeExtraData{Code93IncludeChecksum: true, Code93FullASCIIMode: true},
			},
			err: "",
		},
		testSet{
			name: "pdf",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": []string{"barcodeType"},
				},
			},
			input: render.NamedProperties{
				"type": render.BarcodeTypePDF,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Type:               render.BarcodeTypePDF,
				Extra:              render.BarcodeExtraData{PDFSecurityLevel: 4},
			},
			err: "",
		},
		testSet{
			name: "qr",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": []string{"barcodeType"},
				},
			},
			input: render.NamedProperties{
				"type": render.BarcodeTypeQR,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Type:               render.BarcodeTypeQR,
				Extra:              render.BarcodeExtraData{QRLevel: qr.Q, QRMode: qr.Unicode},
			},
			err: "",
		},
		testSet{
			name: "full prop set, multiple sources, unused props",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"col1":             []string{"dR", "dG", "dA"},
					"left":             []string{"topLeftX"},
					"wide":             []string{"width", "height", "topLeftY"},
					"col3":             []string{"dB", "bR", "bG", "bB", "bA"},
					"what":             []string{"bR", "bG", "bB", "bA", "topLeftX"},
					"some other thing": []string{"content"},
					"type":             []string{"barcodeType"},
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
				"type":             render.BarcodeTypeDataMatrix,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{"what": []string{"bR", "bG", "bB", "bA", "topLeftX"}},
				TopLeft:            image.Pt(3, 80),
				Width:              80,
				Height:             80,
				DataColour:         color.NRGBA{R: uint8(15), G: uint8(15), B: uint8(150), A: uint8(15)},
				BackgroundColour:   color.NRGBA{R: uint8(150), G: uint8(150), B: uint8(150), A: uint8(150)},
				Type:               render.BarcodeTypeDataMatrix,
				Content:            "doesn't matter",
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

func TestInit(t *testing.T) {
	c, err := render.Decode("barcode")
	assert.NoError(t, err)
	assert.Equal(t, Component{}, c)
}
