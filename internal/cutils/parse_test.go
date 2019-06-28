package cutils

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
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
