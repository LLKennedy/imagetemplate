package cutils

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/LLKennedy/imagetemplate/v3/internal/filesystem"
	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font/gofont/goregular"
)

func TestParseColourStrings(t *testing.T) {
	t.Run("all succeeding", func(t *testing.T) {
		colour, props, err := ParseColourStrings("1", "2", "3", "4", map[string][]string{})
		assert.Equal(t, color.NRGBA{R: 1, G: 2, B: 3, A: 4}, colour)
		assert.Equal(t, map[string][]string{}, props)
		assert.NoError(t, err)
	})
	t.Run("all failing", func(t *testing.T) {
		colour, props, err := ParseColourStrings("a", "b", "c", "d", map[string][]string{})
		assert.Equal(t, color.NRGBA{}, colour)
		assert.Equal(t, map[string][]string{}, props)
		assert.EqualError(t, err, "failed to convert property R to uint8: strconv.ParseUint: parsing \"a\": invalid syntax\nfailed to convert property G to uint8: strconv.ParseUint: parsing \"b\": invalid syntax\nfailed to convert property B to uint8: strconv.ParseUint: parsing \"c\": invalid syntax\nfailed to convert property A to uint8: strconv.ParseUint: parsing \"d\": invalid syntax")
	})
}

type fakeSysFonts struct{}

func (f fakeSysFonts) GetFont(req string) (*truetype.Font, error) {
	if req == "good" {
		return truetype.Parse(goregular.TTF)
	}
	return nil, fmt.Errorf("bad font requested")
}

func TestParseFont(t *testing.T) {
	t.Run("no options, empty props", func(t *testing.T) {
		font, props, err := ParseFont("", "", "", ParseFontOptions{})
		assert.Nil(t, font)
		assert.Equal(t, map[string][]string{}, props)
		assert.EqualError(t, err, "exactly one of (fontName,fontFile,fontURL) must be set")
	})
	t.Run("valid font name", func(t *testing.T) {
		validFont, err := truetype.Parse(goregular.TTF)
		assert.NoError(t, err)
		font, props, err := ParseFont("good", "", "", ParseFontOptions{FontPool: fakeSysFonts{}})
		assert.Equal(t, validFont, font)
		assert.Equal(t, map[string][]string{}, props)
		assert.NoError(t, err)
	})
	t.Run("valid font file", func(t *testing.T) {
		fs := filesystem.NewMockFileSystem(filesystem.NewMockFile("font.ttf", goregular.TTF))
		validFont, err := truetype.Parse(goregular.TTF)
		font, props, err := ParseFont("", "font.ttf", "", ParseFontOptions{FileSystem: fs})
		assert.Equal(t, validFont, font)
		assert.Equal(t, map[string][]string{}, props)
		assert.NoError(t, err)
		fs.AssertExpectations(t)
	})
	t.Run("URL not implemented", func(t *testing.T) {
		font, props, err := ParseFont("", "", "anything", ParseFontOptions{})
		assert.Nil(t, font)
		assert.Equal(t, map[string][]string{}, props)
		assert.EqualError(t, err, "fontURL not implemented")
	})
}

func TestParsePoint(t *testing.T) {
	t.Run("error in x", func(t *testing.T) {
		point, props, err := ParsePoint("", "12", "x", "y", map[string][]string{})
		assert.Equal(t, image.Pt(0, 12), point)
		assert.Equal(t, map[string][]string{}, props)
		assert.EqualError(t, err, "error parsing data for property x: could not parse empty property")
	})
	t.Run("error in y", func(t *testing.T) {
		point, props, err := ParsePoint("6", "", "x", "y", map[string][]string{})
		assert.Equal(t, image.Pt(6, 0), point)
		assert.Equal(t, map[string][]string{}, props)
		assert.EqualError(t, err, "error parsing data for property y: could not parse empty property")
	})
	t.Run("error in x and y", func(t *testing.T) {
		point, props, err := ParsePoint("", "", "x", "y", map[string][]string{})
		assert.Equal(t, image.Pt(0, 0), point)
		assert.Equal(t, map[string][]string{}, props)
		assert.EqualError(t, err, "error parsing data for property x: could not parse empty property\nerror parsing data for property y: could not parse empty property")
	})
	t.Run("valid point", func(t *testing.T) {
		point, props, err := ParsePoint("6", "12", "x", "y", map[string][]string{})
		assert.Equal(t, image.Pt(6, 12), point)
		assert.Equal(t, map[string][]string{}, props)
		assert.NoError(t, err)
	})
}
