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

func TestExclusiveOr(t *testing.T) {
	t.Run("no parameters", func(t *testing.T) {
		assert.False(t, ExclusiveOr())
	})
	t.Run("one true", func(t *testing.T) {
		assert.True(t, ExclusiveOr(true))
	})
	t.Run("one false", func(t *testing.T) {
		assert.False(t, ExclusiveOr(false))
	})
	t.Run("one true, one false", func(t *testing.T) {
		assert.True(t, ExclusiveOr(true, false))
	})
	t.Run("one true, many false", func(t *testing.T) {
		assert.True(t, ExclusiveOr(true, false, false, false, false, false, false, false, false))
	})
	t.Run("two true, one false", func(t *testing.T) {
		assert.False(t, ExclusiveOr(true, true, false))
	})
	t.Run("two true, many false", func(t *testing.T) {
		assert.False(t, ExclusiveOr(true, true, false, false, false, false, false, false, false))
	})
}

func TestExclusiveNor(t *testing.T) {
	t.Run("no parameters", func(t *testing.T) {
		assert.True(t, ExclusiveNor())
	})
	t.Run("one true", func(t *testing.T) {
		assert.False(t, ExclusiveNor(true))
	})
	t.Run("one false", func(t *testing.T) {
		assert.True(t, ExclusiveNor(false))
	})
	t.Run("one true, one false", func(t *testing.T) {
		assert.False(t, ExclusiveNor(true, false))
	})
	t.Run("one true, many false", func(t *testing.T) {
		assert.False(t, ExclusiveNor(true, false, false, false, false, false, false, false, false))
	})
	t.Run("two true, one false", func(t *testing.T) {
		assert.True(t, ExclusiveNor(true, true, false))
	})
	t.Run("two true, many false", func(t *testing.T) {
		assert.True(t, ExclusiveNor(true, true, false, false, false, false, false, false, false))
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
