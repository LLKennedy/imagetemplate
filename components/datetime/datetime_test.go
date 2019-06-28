package datetime

import (
	"fmt"
	"image"
	"image/color"
	"runtime/debug"
	"testing"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/LLKennedy/gosysfonts"
	"github.com/LLKennedy/imagetemplate/v3/cutils"
	"github.com/LLKennedy/imagetemplate/v3/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
)

func TestDateTimeWrite(t *testing.T) {
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
	t.Run("datetime error", func(t *testing.T) {
		expectedFont := truetype.NewFace(goreg, &truetype.Options{Size: 14, Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: float64(72)})
		canvas := new(render.MockCanvas)
		canvas.On("GetPPI").Return(float64(72))
		canvas.On("TryText", "", image.Point{}, expectedFont, color.NRGBA{}, 100).Return(true, 10)
		canvas.On("Text", "", image.Point{}, expectedFont, color.NRGBA{}, 100).Return(canvas, fmt.Errorf("some error"))
		timeVal := time.Now()
		c := Component{Font: goreg, Size: 14, MaxWidth: 100, Time: &timeVal}
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
		timeVal := time.Now()
		c := Component{Font: goreg, Size: 24, MaxWidth: 100, Time: &timeVal}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.NoError(t, err)
		canvas.AssertExpectations(t)
	})
	t.Run("can't ever fit", func(t *testing.T) {
		expectedFont := truetype.NewFace(goreg, &truetype.Options{Size: float64(24), Hinting: font.HintingFull, SubPixelsX: 64, SubPixelsY: 64, DPI: float64(72)})
		canvas := new(render.MockCanvas)
		canvas.On("GetPPI").Return(float64(72))
		timeVal := time.Now()
		canvas.On("TryText", timeVal.Format(time.RFC822), image.Point{}, expectedFont, color.NRGBA{}, 100).Return(false, 100)
		c := Component{Font: goreg, Size: 24, MaxWidth: 100, Time: &timeVal, TimeFormat: time.RFC822}
		modifiedCanvas, err := c.Write(canvas)
		assert.Equal(t, canvas, modifiedCanvas)
		assert.EqualError(t, err, fmt.Sprintf("unable to fit datetime %s into maxWidth 100 after 10 tries", timeVal.Format(time.RFC822)))
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
		timeVal := time.Now()
		t.Run("left", func(t *testing.T) {
			c := Component{Font: goreg, Size: 24, MaxWidth: 100, Time: &timeVal, TextAlignment: cutils.TextAlignmentLeft}
			modifiedCanvas, err := c.Write(canvas)
			assert.Equal(t, canvas, modifiedCanvas)
			assert.NoError(t, err)
		})
		t.Run("right", func(t *testing.T) {
			c := Component{Font: goreg, Size: 24, MaxWidth: 100, Time: &timeVal, TextAlignment: cutils.TextAlignmentRight}
			modifiedCanvas, err := c.Write(canvas)
			assert.Equal(t, canvas, modifiedCanvas)
			assert.NoError(t, err)
		})
		t.Run("centre", func(t *testing.T) {
			c := Component{Font: goreg, Size: 24, MaxWidth: 100, Time: &timeVal, TextAlignment: cutils.TextAlignmentCentre}
			modifiedCanvas, err := c.Write(canvas)
			assert.Equal(t, canvas, modifiedCanvas)
			assert.NoError(t, err)
		})
		t.Run("default", func(t *testing.T) {
			c := Component{Font: goreg, Size: 24, MaxWidth: 100, Time: &timeVal, TextAlignment: cutils.TextAlignment(12)}
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

func TestDateTimeSetNamedProperties(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input render.NamedProperties
		res   Component
		err   string
	}
	timestamp, _ := time.Parse(time.RFC822, time.Now().Format(time.RFC822))
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
			name: "invalid time (not valid type)",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"time"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"time"},
				},
			},
			err: "error converting 12 to []string, *time.Time or time.Time",
		},
		{
			name: "invalid time (bad string slice)",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"time"},
				},
			},
			input: render.NamedProperties{
				"aProp": []string{"hello"},
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"time"},
				},
			},
			err: "error converting [hello] to []string, *time.Time or time.Time",
		},
		{
			name: "invalid time (bad string slice contents)",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"time"},
				},
			},
			input: render.NamedProperties{
				"aProp": []string{"hello", "there"},
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"time"},
				},
			},
			err: "cannot convert time string there to time format hello",
		},
		{
			name: "valid time (time.Time)",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"time"},
				},
			},
			input: render.NamedProperties{
				"aProp": timestamp,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Time:               &timestamp,
			},
		},
		{
			name: "valid time (*time.Time)",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"time"},
				},
			},
			input: render.NamedProperties{
				"aProp": &timestamp,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Time:               &timestamp,
			},
		},
		{
			name: "valid time ([]string)",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"time"},
				},
			},
			input: render.NamedProperties{
				"aProp": []string{time.RFC822, timestamp.Format(time.RFC822)},
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				Time:               &timestamp,
			},
		},
		{
			name: "invalid time format",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"timeFormat"},
				},
			},
			input: render.NamedProperties{
				"aProp": 12,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"timeFormat"},
				},
			},
			err: "error converting 12 to string",
		},
		{
			name: "valid time format",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"timeFormat"},
				},
			},
			input: render.NamedProperties{
				"aProp": time.RFC822,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				TimeFormat:         time.RFC822,
			},
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
				TextAlignment:      cutils.TextAlignmentLeft,
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
				TextAlignment: cutils.TextAlignmentLeft,
				fs:            ttfFS,
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
				TextAlignment: cutils.TextAlignmentLeft,
				fs:            ttfFS,
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
			err: "could not convert 12 to datetime alignment or string",
		},
		{
			name: "alignment constant valid",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"alignment"},
				},
			},
			input: render.NamedProperties{
				"aProp": cutils.TextAlignmentLeft,
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{},
				TextAlignment:      cutils.TextAlignmentLeft,
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
				TextAlignment:      cutils.TextAlignmentLeft,
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
				TextAlignment:      cutils.TextAlignmentRight,
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
				TextAlignment:      cutils.TextAlignmentCentre,
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
				TextAlignment:      cutils.TextAlignmentLeft,
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
			err: "invalid component property in named property map: not a prop",
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
			name: "invalid startX type",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"startX"},
				},
			},
			input: render.NamedProperties{
				"aProp": "15",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"startX"},
				},
			},
			err: "error converting 15 to int",
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
			name: "invalid startY type",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"startY"},
				},
			},
			input: render.NamedProperties{
				"aProp": "15",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"startY"},
				},
			},
			err: "error converting 15 to int",
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
			name: "invalid maxWidth type",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"maxWidth"},
				},
			},
			input: render.NamedProperties{
				"aProp": "15",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"maxWidth"},
				},
			},
			err: "error converting 15 to int",
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
	fileSystemsToCheck := []*filesystem.MockFileSystem{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
				}
			}()
			res, err := test.start.SetNamedProperties(test.input)
			assert.Equal(t, test.res, res)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.err)
			}
			if mockFs, isMock := test.start.fs.(*filesystem.MockFileSystem); isMock {
				fileSystemsToCheck = append(fileSystemsToCheck, mockFs)
			}
		})
	}
	for _, fs := range fileSystemsToCheck {
		fs.AssertExpectations(t)
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
	ttfFS := filesystem.NewMockFileSystem(
		filesystem.NewMockFile("myFont.ttf", goregular.TTF),
		filesystem.NewMockFile("badfont.TTF", []byte("hello")),
	)
	ttfFS.On("Open", "nilfont.TTF").Return(filesystem.NilFile, nil)
	tests := []testSet{
		{
			name:  "incorrect format data",
			start: Component{},
			input: "hello",
			res:   Component{},
			props: render.NamedProperties{},
			err:   "failed to convert returned data to component properties",
		},
		{
			name: "error extracting font",
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			props: render.NamedProperties{},
			err:   "exactly one of (fontName,fontFile,fontURL) must be set",
		},
		{
			name: "bad font name",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "bad",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err:   "bad font requested",
		},
		{
			name: "valid font name",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "good",
				},
				Time:          "$some time$",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fontPool:           fakeSysFonts{},
				Font:               func() *truetype.Font { f, _ := truetype.Parse(goregular.TTF); return f }(),
				TimeFormat:         time.RFC822,
				Start:              image.Pt(12, 12),
				MaxWidth:           67,
				Size:               89,
				TextAlignment:      cutils.TextAlignmentLeft,
				Colour:             color.NRGBA{R: 6, G: 53, B: 197, A: 244},
				NamedPropertiesMap: map[string][]string{"some time": {"time"}},
			},
			props: render.NamedProperties{"some time": struct{ Message string }{Message: "Please replace me with real data"}},
		},
		{
			name: "bad font reader returned from filesystem",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "nilfont.TTF",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err:   "cannot read from nil file",
		},
		{
			name: "bad font TTF data",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "badfont.TTF",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err:   "freetype: invalid TrueType format: TTF data is too short",
		},
		{
			name: "working font file",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "$some time$",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs:                 ttfFS,
				Font:               func() *truetype.Font { f, _ := truetype.Parse(goregular.TTF); return f }(),
				TimeFormat:         time.RFC822,
				Start:              image.Pt(12, 12),
				MaxWidth:           67,
				Size:               89,
				TextAlignment:      cutils.TextAlignmentLeft,
				Colour:             color.NRGBA{R: 6, G: 53, B: 197, A: 244},
				NamedPropertiesMap: map[string][]string{"some time": {"time"}},
			},
			props: render.NamedProperties{"some time": struct{ Message string }{Message: "Please replace me with real data"}},
		},
		{
			name:  "font URL not implemented",
			start: Component{},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontURL: "anything",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res:   Component{},
			props: render.NamedProperties{},
			err:   "fontURL not implemented",
		},
		{
			name: "empty time",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "good",
				},
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property time: could not parse empty property",
		},
		{
			name: "empty time format",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "good",
				},
				Time:          "3h",
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property timeFormat: could not parse empty property",
		},
		{
			name: "empty time and time format",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "good",
				},
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property time: could not parse empty property\nerror parsing data for property timeFormat: could not parse empty property",
		},
		{
			name: "empty startX",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "good",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartY:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property startX: could not parse empty property",
		},
		{
			name: "empty startY",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "good",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property startY: could not parse empty property",
		},
		{
			name: "empty startX and startY",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontName: "good",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property startX: could not parse empty property\nerror parsing data for property startY: could not parse empty property",
		},
		{
			name: "empty maxWidth",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				Size:          "12",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property maxWidth: could not parse empty property",
		},
		{
			name: "empty size",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "12",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property size: could not parse empty property",
		},
		{
			name: "empty alignment",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:       "3h",
				TimeFormat: time.RFC822,
				StartX:     "12",
				StartY:     "12",
				MaxWidth:   "67",
				Size:       "89",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property alignment: could not parse empty property",
		},
		{
			name: "valid alignment (left)",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "$some time$",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "12",
				Size:          "12",
				TextAlignment: "left",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs:                 ttfFS,
				Font:               func() *truetype.Font { f, _ := truetype.Parse(goregular.TTF); return f }(),
				TimeFormat:         time.RFC822,
				Start:              image.Pt(12, 12),
				MaxWidth:           12,
				Size:               12,
				TextAlignment:      cutils.TextAlignmentLeft,
				Colour:             color.NRGBA{R: 6, G: 53, B: 197, A: 244},
				NamedPropertiesMap: map[string][]string{"some time": {"time"}},
			},
			props: render.NamedProperties{"some time": struct{ Message string }{Message: "Please replace me with real data"}},
		},
		{
			name: "valid alignment (right)",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "$some time$",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "12",
				Size:          "12",
				TextAlignment: "right",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs:                 ttfFS,
				Font:               func() *truetype.Font { f, _ := truetype.Parse(goregular.TTF); return f }(),
				TimeFormat:         time.RFC822,
				Start:              image.Pt(12, 12),
				MaxWidth:           12,
				Size:               12,
				TextAlignment:      cutils.TextAlignmentRight,
				Colour:             color.NRGBA{R: 6, G: 53, B: 197, A: 244},
				NamedPropertiesMap: map[string][]string{"some time": {"time"}},
			},
			props: render.NamedProperties{"some time": struct{ Message string }{Message: "Please replace me with real data"}},
		},
		{
			name: "valid alignment (centre)",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "$some time$",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "12",
				Size:          "12",
				TextAlignment: "centre",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs:                 ttfFS,
				Font:               func() *truetype.Font { f, _ := truetype.Parse(goregular.TTF); return f }(),
				TimeFormat:         time.RFC822,
				Start:              image.Pt(12, 12),
				MaxWidth:           12,
				Size:               12,
				TextAlignment:      cutils.TextAlignmentCentre,
				Colour:             color.NRGBA{R: 6, G: 53, B: 197, A: 244},
				NamedPropertiesMap: map[string][]string{"some time": {"time"}},
			},
			props: render.NamedProperties{"some time": struct{ Message string }{Message: "Please replace me with real data"}},
		},
		{
			name: "valid alignment (default)",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "$some time$",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "12",
				Size:          "12",
				TextAlignment: "something else",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs:                 ttfFS,
				Font:               func() *truetype.Font { f, _ := truetype.Parse(goregular.TTF); return f }(),
				TimeFormat:         time.RFC822,
				Start:              image.Pt(12, 12),
				MaxWidth:           12,
				Size:               12,
				TextAlignment:      cutils.TextAlignmentLeft,
				Colour:             color.NRGBA{R: 6, G: 53, B: 197, A: 244},
				NamedPropertiesMap: map[string][]string{"some time": {"time"}},
			},
			props: render.NamedProperties{"some time": struct{ Message string }{Message: "Please replace me with real data"}},
		},
		{
			name: "empty red",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "12",
				Size:          "12",
				TextAlignment: "something else",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property R: could not parse empty property",
		},
		{
			name: "empty green",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "12",
				Size:          "12",
				TextAlignment: "something else",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property G: could not parse empty property",
		},
		{
			name: "empty blue",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "12",
				Size:          "12",
				TextAlignment: "something else",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Alpha: "244",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property B: could not parse empty property",
		},
		{
			name: "empty alpha",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "3h",
				TimeFormat:    time.RFC822,
				StartX:        "12",
				StartY:        "12",
				MaxWidth:      "12",
				Size:          "12",
				TextAlignment: "something else",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
				},
			},
			res: Component{
				fs: ttfFS,
			},
			props: render.NamedProperties{},
			err:   "error parsing data for property A: could not parse empty property",
		},
		{
			name: "valid everything",
			start: Component{
				fs: ttfFS,
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "myFont.ttf",
				},
				Time:          "$some time$",
				TimeFormat:    time.RFC822,
				StartX:        "123",
				StartY:        "45",
				MaxWidth:      "67",
				Size:          "89",
				TextAlignment: "something else",
				Colour: struct {
					Red   string `json:"R"`
					Green string `json:"G"`
					Blue  string `json:"B"`
					Alpha string `json:"A"`
				}{
					Red:   "6",
					Green: "53",
					Blue:  "197",
					Alpha: "244",
				},
			},
			res: Component{
				fs:                 ttfFS,
				Font:               func() *truetype.Font { f, _ := truetype.Parse(goregular.TTF); return f }(),
				TimeFormat:         time.RFC822,
				Start:              image.Pt(123, 45),
				MaxWidth:           67,
				Size:               89,
				TextAlignment:      cutils.TextAlignmentLeft,
				Colour:             color.NRGBA{R: 6, G: 53, B: 197, A: 244},
				NamedPropertiesMap: map[string][]string{"some time": {"time"}},
			},
			props: render.NamedProperties{"some time": struct{ Message string }{Message: "Please replace me with real data"}},
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
