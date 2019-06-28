package cutils

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/LLKennedy/imagetemplate/v3/internal/filesystem"
	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font/gofont/goregular"
)

func TestCombineErrors(t *testing.T) {
	t.Run("nil history", func(t *testing.T) {
		latest := fmt.Errorf("error b")
		combined := CombineErrors(nil, latest)
		assert.EqualError(t, combined, "error b")
	})
	t.Run("nil latest", func(t *testing.T) {
		history := fmt.Errorf("error a")
		combined := CombineErrors(history, nil)
		assert.EqualError(t, combined, "error a")
	})
	t.Run("nil latest", func(t *testing.T) {
		history := fmt.Errorf("error a")
		latest := fmt.Errorf("error b")
		combined := CombineErrors(history, latest)
		assert.EqualError(t, combined, "error a\nerror b")
	})
}

func TestLoadFontFile(t *testing.T) {
	t.Run("failed conversion to string", func(t *testing.T) {
		font, err := LoadFontFile(nil, 12)
		assert.Nil(t, font)
		assert.EqualError(t, err, "error converting 12 to string")
	})
	t.Run("failed to open file", func(t *testing.T) {
		var fs *filesystem.MockFileSystem
		font, err := LoadFontFile(fs, "some file")
		assert.Nil(t, font)
		assert.EqualError(t, err, "cannot open on nil file system")
	})
	t.Run("non-nil file system", func(t *testing.T) {
		fs := filesystem.NewMockFileSystem(filesystem.NewMockFile("goodFile.ttf", goregular.TTF))
		fs.On("Open", "myFile.ttf").Return(filesystem.NilFile, nil)
		t.Run("failed to read file contents", func(t *testing.T) {
			font, err := LoadFontFile(fs, "myFile.ttf")
			assert.Nil(t, font)
			assert.EqualError(t, err, "cannot read from nil file")
		})
		t.Run("valid file", func(t *testing.T) {
			validFont, err := truetype.Parse(goregular.TTF)
			assert.NoError(t, err)
			font, err := LoadFontFile(fs, "goodFile.ttf")
			assert.NoError(t, err)
			assert.Equal(t, validFont, font)
		})
	})
}

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
