package cutils

import (
	"image"
	"image/color"
	"runtime/debug"
	"testing"
	"time"

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

func TestSetString(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		converted, err := SetString("a")
		assert.Equal(t, "a", converted)
		assert.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		converted, err := SetString(1)
		assert.Equal(t, "", converted)
		assert.EqualError(t, err, "error converting 1 to string")
	})
}

func TestSetInt(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		converted, err := SetInt(1)
		assert.Equal(t, 1, converted)
		assert.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		converted, err := SetInt("a")
		assert.Equal(t, 0, converted)
		assert.EqualError(t, err, "error converting a to int")
	})
}

func TestSetUint8(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		converted, err := SetUint8(uint8(1))
		assert.Equal(t, uint8(1), converted)
		assert.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		converted, err := SetUint8("a")
		assert.Equal(t, uint8(0), converted)
		assert.EqualError(t, err, "error converting a to uint8")
	})
}

func TestSetFloat64(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		converted, err := SetFloat64(float64(1))
		assert.Equal(t, float64(1), converted)
		assert.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		converted, err := SetFloat64("a")
		assert.Equal(t, float64(0), converted)
		assert.EqualError(t, err, "error converting a to float64")
	})
}

func TestSetBool(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		converted, err := SetBool(true)
		assert.Equal(t, true, converted)
		assert.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		converted, err := SetBool("a")
		assert.Equal(t, false, converted)
		assert.EqualError(t, err, "error converting a to bool")
	})
}

func TestSetTime(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ts := time.Now()
		converted, err := SetTime(ts)
		assert.Equal(t, ts, converted)
		assert.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		converted, err := SetTime("a")
		assert.Equal(t, time.Time{}, converted)
		assert.EqualError(t, err, "error converting a to time")
	})
}

func TestSetTimePointer(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ts := time.Now()
		converted, err := SetTimePointer(&ts)
		assert.Equal(t, &ts, converted)
		assert.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		converted, err := SetTimePointer("a")
		assert.Equal(t, (*time.Time)(nil), converted)
		assert.EqualError(t, err, "error converting a to time pointer")
	})
}

func TestSetImage(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		img := image.NRGBA{}
		converted, err := SetImage(&img)
		assert.Equal(t, &img, converted)
		assert.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
			}
		}()
		converted, err := SetImage("a")
		assert.Equal(t, image.Image(nil), converted)
		assert.EqualError(t, err, "error converting a to image")
	})
}

func TestSetColour(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		clr := color.NRGBA{}
		converted, err := SetColour(&clr)
		assert.Equal(t, &clr, converted)
		assert.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
			}
		}()
		converted, err := SetColour("a")
		assert.Equal(t, color.Color(nil), converted)
		assert.EqualError(t, err, "error converting a to colour")
	})
}
