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
		c := Component{NamedPropertiesMap: map[string][]string{"not set": {"something"}}}
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
					"aProp": {"dR"},
				},
			},
			input: render.NamedProperties{
				"aProp": "not a number",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"dR"},
				},
			},
			err: "error converting not a number to uint8",
		},
		{
			name: "dRGBA valid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"dR", "dG", "dB", "dA"},
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
		{
			name: "non-RGBA invalid type",
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
			name: "content invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"content"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"content"},
				},
			},
			err: "error converting 15 to string",
		},
		{
			name: "type invalid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"barcodeType"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"barcodeType"},
				},
			},
			err: "error converting 15 to barcode type",
		},
		{
			name: "unsupported type",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"barcodeType"},
				},
			},
			input: render.NamedProperties{
				"aProp": "15",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"barcodeType"},
				},
			},
			err: "error converting 15 to barcode type",
		},
		{
			name: "invalid colour code",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"Rd"},
				},
			},
			input: render.NamedProperties{
				"aProp": uint8(12),
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"Rd"},
				},
			},
			err: "name was a string inside RGBA and Value was a valid uint8, but Name wasn't R, G, B, or A. Name was: Rd",
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
			name: "aztec",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": {"barcodeType"},
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
		{
			name: "39",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": {"barcodeType"},
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
		{
			name: "93",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": {"barcodeType"},
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
		{
			name: "pdf",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": {"barcodeType"},
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
		{
			name: "qr",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"type": {"barcodeType"},
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
		{
			name: "full prop set, multiple sources, unused props",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"col1":             {"dR", "dG", "dA"},
					"left":             {"topLeftX"},
					"wide":             {"width", "height", "topLeftY"},
					"col3":             {"dB", "bR", "bG", "bB", "bA"},
					"what":             {"bR", "bG", "bB", "bA", "topLeftX"},
					"some other thing": {"content"},
					"type":             {"barcodeType"},
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
				NamedPropertiesMap: map[string][]string{"what": {"bR", "bG", "bB", "bA", "topLeftX"}},
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
		{
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
