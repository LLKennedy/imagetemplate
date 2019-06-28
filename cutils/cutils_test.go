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
