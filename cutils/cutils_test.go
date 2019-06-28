package cutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	t.Run("valid both", func(t *testing.T) {
		history := fmt.Errorf("error a")
		latest := fmt.Errorf("error b")
		combined := CombineErrors(history, latest)
		assert.EqualError(t, combined, "error a\nerror b")
	})
	t.Run("nil both", func(t *testing.T) {
		combined := CombineErrors(nil, nil)
		assert.NoError(t, combined)
	})
}

func TestStringToAlignment(t *testing.T) {
	assert.Equal(t, TextAlignmentLeft, StringToAlignment("left"))
	assert.Equal(t, TextAlignmentCentre, StringToAlignment("centre"))
	assert.Equal(t, TextAlignmentRight, StringToAlignment("right"))
	assert.Equal(t, TextAlignmentLeft, StringToAlignment("gibberish"))
}

func TestScaleFontsToWidth(t *testing.T) {
	t.Run("perfect size", func(t *testing.T) {
		size, offset := ScaleFontsToWidth(10, 150, 150, TextAlignmentRight)
		assert.Equal(t, float64(10), size)
		assert.Equal(t, 0, offset)
	})
	t.Run("too large", func(t *testing.T) {
		size, offset := ScaleFontsToWidth(10, 300, 150, TextAlignmentRight)
		assert.Equal(t, float64(5), size)
		assert.Equal(t, 0, offset)
	})
	t.Run("too small", func(t *testing.T) {
		t.Run("left", func(t *testing.T) {
			size, offset := ScaleFontsToWidth(10, 100, 200, TextAlignmentLeft)
			assert.Equal(t, float64(10), size)
			assert.Equal(t, 0, offset)
		})
		t.Run("right", func(t *testing.T) {
			size, offset := ScaleFontsToWidth(10, 100, 200, TextAlignmentRight)
			assert.Equal(t, float64(10), size)
			assert.Equal(t, 100, offset)
		})
		t.Run("centre", func(t *testing.T) {
			size, offset := ScaleFontsToWidth(10, 100, 200, TextAlignmentCentre)
			assert.Equal(t, float64(10), size)
			assert.Equal(t, 50, offset)
		})
		t.Run("default", func(t *testing.T) {
			size, offset := ScaleFontsToWidth(10, 100, 200, TextAlignment(-1))
			assert.Equal(t, float64(10), size)
			assert.Equal(t, 0, offset)
		})
	})
}
