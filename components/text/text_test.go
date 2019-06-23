package text

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/LLKennedy/gosysfonts"
	"github.com/LLKennedy/imagetemplate/v3/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
)

func TestTextWrite(t *testing.T) {
	goreg, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("not all props set", func(t *testing.T) {
		canvas := new(render.MockCanvas)
		c := Component{NamedPropertiesMap: map[string][]string{"not set": {"something"}}}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.Error(t, err)
		canvas.AssertExpectations(t)
	})
	t.Run("text error", func(t *testing.T) {
		expectedFont := truetype.NewFace(goreg, &truetype.Options{Size: 14, Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: float64(72)})
		canvas := new(render.MockCanvas)
		canvas.On("GetPPI").Return(float64(72))
		canvas.On("TryText", "", image.Point{}, expectedFont, color.NRGBA{}, 100).Return(true, 10)
		canvas.On("Text", "", image.Point{}, expectedFont, color.NRGBA{}, 100).Return(canvas, fmt.Errorf("some error"))
		c := Component{Font: goreg, Size: 14, MaxWidth: 100}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "some error")
		canvas.AssertExpectations(t)
	})
	t.Run("multiple passes required", func(t *testing.T) {
		expectedFont := truetype.NewFace(goreg, &truetype.Options{Size: float64(24), Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: float64(72)})
		expectedFont2 := truetype.NewFace(goreg, &truetype.Options{Size: float64(12), Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: float64(72)})
		expectedFont3 := truetype.NewFace(goreg, &truetype.Options{Size: float64(8), Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: float64(72)})
		canvas := new(render.MockCanvas)
		canvas.On("GetPPI").Return(float64(72))
		canvas.On("TryText", "", image.Point{}, expectedFont, color.NRGBA{}, 100).Return(false, 200)
		canvas.On("TryText", "", image.Point{}, expectedFont2, color.NRGBA{}, 100).Return(false, 150)
		canvas.On("TryText", "", image.Point{}, expectedFont3, color.NRGBA{}, 100).Return(true, 100)
		canvas.On("Text", "", image.Point{}, expectedFont3, color.NRGBA{}, 100).Return(canvas, nil)
		c := Component{Font: goreg, Size: 24, MaxWidth: 100}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
		canvas.AssertExpectations(t)
	})
	t.Run("can't ever fit", func(t *testing.T) {
		expectedFont := truetype.NewFace(goreg, &truetype.Options{Size: float64(24), Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: float64(72)})
		canvas := new(render.MockCanvas)
		canvas.On("GetPPI").Return(float64(72))
		canvas.On("TryText", "", image.Point{}, expectedFont, color.NRGBA{}, 100).Return(false, 100)
		c := Component{Font: goreg, Size: 24, MaxWidth: 100}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, "unable to fit text  into maxWidth 100 after 10 tries")
		canvas.AssertExpectations(t)
	})
	t.Run("different alignments", func(t *testing.T) {
		expectedFont := truetype.NewFace(goreg, &truetype.Options{Size: float64(24), Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: float64(72)})
		canvas := new(render.MockCanvas)
		canvas.On("GetPPI").Return(float64(72))
		canvas.On("TryText", "", image.Point{}, expectedFont, color.NRGBA{}, 100).Return(true, 50)
		canvas.On("Text", "", image.Point{}, expectedFont, color.NRGBA{}, 100).Return(canvas, nil)
		canvas.On("Text", "", image.Pt(25, 0), expectedFont, color.NRGBA{}, 100).Return(canvas, nil)
		canvas.On("Text", "", image.Pt(50, 0), expectedFont, color.NRGBA{}, 100).Return(canvas, nil)
		t.Run("left", func(t *testing.T) {
			c := Component{Font: goreg, Size: 24, MaxWidth: 100, Alignment: AlignmentLeft}
			modifiedCanvas, err := c.Write(canvas)
			assert.Equal(t, canvas, modifiedCanvas)
			assert.NoError(t, err)
		})
		t.Run("right", func(t *testing.T) {
			c := Component{Font: goreg, Size: 24, MaxWidth: 100, Alignment: AlignmentRight}
			modifiedCanvas, err := c.Write(canvas)
			assert.Equal(t, canvas, modifiedCanvas)
			assert.NoError(t, err)
		})
		t.Run("centre", func(t *testing.T) {
			c := Component{Font: goreg, Size: 24, MaxWidth: 100, Alignment: AlignmentCentre}
			modifiedCanvas, err := c.Write(canvas)
			assert.Equal(t, canvas, modifiedCanvas)
			assert.NoError(t, err)
		})
		t.Run("default", func(t *testing.T) {
			c := Component{Font: goreg, Size: 24, MaxWidth: 100, Alignment: Alignment(12)}
			modifiedCanvas, err := c.Write(canvas)
			assert.Equal(t, canvas, modifiedCanvas)
			assert.NoError(t, err)
		})
		canvas.AssertExpectations(t)
	})
}

type fakeSysFonts struct{}

func (f fakeSysFonts) GetFont(req string) (*truetype.Font, error) {
	if req == "good" {
		return truetype.Parse(goregular.TTF)
	}
	return nil, fmt.Errorf("bad font requested")
}

func TestTextSetNamedProperties(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input render.NamedProperties
		res   Component
		err   string
	}
	ttfFS := filesystem.NewMockFileSystem(
		filesystem.NewMockFile("myFont.ttf", goregular.TTF),
		filesystem.NewMockFile("badfont.TTF", []byte("hello")),
	)
	ttfFS.On("Open", "nilfont.TTF").Return(filesystem.NilFile, nil)
	tests := []testSet{
		{
			name:  "no props",
			start: Component{},
			input: render.NamedProperties{},
			res:   Component{},
			err:   "",
		},
		{
			name: "valid content",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"content"},
				},
			},
			input: render.NamedProperties{
				"aProp": "good",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Content: "good",
			},
		},
		{
			name: "invalid content",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"content"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"content"},
				},
			},
			err: "error converting 12 to string",
		},
		{
			name: "invalid font name",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontName"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontName"},
				},
			},
			err: "error converting 12 to string",
		},
		{
			name: "error requesting font",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontName"},
				},
				fontPool: fakeSysFonts{},
			},
			input: render.NamedProperties{
				"aProp": "bad",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontName"},
				},
				fontPool: fakeSysFonts{},
			},
			err: "bad font requested",
		},
		{
			name: "valid font name",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontName"},
				},
				fontPool: fakeSysFonts{},
			},
			input: render.NamedProperties{
				"aProp": "good",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				fontPool:           fakeSysFonts{},
				Font:               func() *truetype.Font { f, _ := truetype.Parse(goregular.TTF); return f }(),
			},
		},
		{
			name: "invalid font file",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontFile"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontFile"},
				},
			},
			err: "error converting 12 to string",
		},
		{
			name: "valid font file",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontFile"},
				},
				fs: ttfFS,
			},
			input: render.NamedProperties{
				"aProp": "myFont.ttf",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Alignment:          AlignmentLeft,
				fs:                 ttfFS,
				Font:               func() *truetype.Font { f, _ := truetype.Parse(goregular.TTF); return f }(),
			},
			err: "",
		},
		{
			name: "error parsing font file",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontFile"},
				},
				fs: ttfFS,
			},
			input: render.NamedProperties{
				"aProp": "badfont.TTF",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontFile"},
				},
				Alignment: AlignmentLeft,
				fs:        ttfFS,
			},
			err: "freetype: invalid TrueType format: TTF data is too short",
		},
		{
			name: "error reading font data",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontFile"},
				},
				fs: ttfFS,
			},
			input: render.NamedProperties{
				"aProp": "nilfont.TTF",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontFile"},
				},
				Alignment: AlignmentLeft,
				fs:        ttfFS,
			},
			err: "cannot read from nil file",
		},
		{
			name: "invalid font url",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontURL"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontURL"},
				},
			},
			err: "fontURL not implemented",
		},
		{
			name: "invalid alignment type",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"alignment"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"alignment"},
				},
			},
			err: "could not convert 12 to text alignment or string",
		},
		{
			name: "alignment constant valid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"alignment"},
				},
			},
			input: render.NamedProperties{
				"aProp": AlignmentLeft,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Alignment:          AlignmentLeft,
			},
			err: "",
		},
		{
			name: "alignment string left",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"alignment"},
				},
			},
			input: render.NamedProperties{
				"aProp": "left",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Alignment:          AlignmentLeft,
			},
			err: "",
		},
		{
			name: "alignment string right",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"alignment"},
				},
			},
			input: render.NamedProperties{
				"aProp": "right",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Alignment:          AlignmentRight,
			},
			err: "",
		},
		{
			name: "alignment string centre",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"alignment"},
				},
			},
			input: render.NamedProperties{
				"aProp": "centre",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Alignment:          AlignmentCentre,
			},
			err: "",
		},
		{
			name: "alignment string default",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"alignment"},
				},
			},
			input: render.NamedProperties{
				"aProp": "gibberish",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Alignment:          AlignmentLeft,
			},
			err: "",
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
			name: "startX",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"startX"},
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
		{
			name: "startY",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"startY"},
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
		{
			name: "maxWidth",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"maxWidth"},
				},
			},
			input: render.NamedProperties{
				"aProp": 15,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				MaxWidth:           15,
			},
			err: "",
		},
		{
			name: "size",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"size"},
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
		{
			name: "invalid size",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"size"},
				},
			},
			input: render.NamedProperties{
				"aProp": "a",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"size"},
				},
			},
			err: "error converting a to float64",
		},
		{
			name: "full prop set, multiple sources, unused props",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"col1": {"R", "G", "A"},
					"left": {"startX", "startY"},
					"wide": {"size"},
					"col3": {"B"},
					"what": {"R", "G", "B", "A", "startX"},
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
				NamedPropertiesMap: map[string][]string{"what": {"R", "G", "B", "A", "startX"}},
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
	ttfFS.AssertExpectations(t)
}

func TestTextGetJSONFormat(t *testing.T) {
	c := Component{}
	expectedFormat := &textFormat{}
	format := c.GetJSONFormat()
	assert.Equal(t, expectedFormat, format)
}

func TestTextVerifyAndTestTextJSONData(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input interface{}
		res   Component
		props render.NamedProperties
		err   string
	}
	ttfFS := filesystem.NewMockFileSystem(
		filesystem.NewMockFile("myFont.ttf", goregular.TTF),
		filesystem.NewMockFile("badfont.TTF", []byte("hello")),
	)
	ttfFS.On("Open", "nilfont.TTF").Return(filesystem.NilFile, nil)
	tests := []testSet{
		{
			name:  "incorrect format data",
			input: "hello",
			props: render.NamedProperties{},
			err:   "failed to convert returned data to component properties",
		},
		{
			name: "error extracting font",
			input: &textFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					
				},
			},
			props: render.NamedProperties{},
			err: "exactly one of (fontName,fontFile,fontURL) must be set",
		},
		{
			name: "bad font name",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &textFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "bad",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err: "bad font requested",
		},
		{
			name: "valid font name",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &textFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "good",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err: "error parsing data for property content: could not parse empty property",
		},
		{
			name: "bad font reader returned from filesystem",
			start: Component{
				fs: ttfFS,
			},
			input: &textFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "nilfont.TTF",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err: "cannot read from nil file",
		},
		{
			name: "bad font TTF data",
			start: Component{
				fs: ttfFS,
			},
			input: &textFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "badfont.TTF",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err: "freetype: invalid TrueType format: TTF data is too short",
		},
		{
			name: "working font file",
			start: Component{
				fs: ttfFS,
			},
			input: &textFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err: "error parsing data for property content: could not parse empty property",
		},
		{
			name: "font URL not implemented",
			start: Component{
			},
			input: &textFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontURL: "anything",
				},
			},
			res: Component{
			},
			props: render.NamedProperties{},
			err: "fontURL not implemented",
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
	ttfFS.AssertExpectations(t)
}

func TestGetFontPool(t *testing.T) {
	assert.Equal(t, gosysfonts.New(), Component{}.getFontPool())
}
