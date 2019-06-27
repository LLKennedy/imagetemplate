package cutils

import (
	"testing"

	"github.com/LLKennedy/imagetemplate/v3/internal/filesystem"
	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font/gofont/goregular"
)

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
