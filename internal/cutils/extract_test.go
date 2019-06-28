package cutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractString(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		str, props, err := ExtractString("", "myProp", map[string][]string{})
		assert.Equal(t, "", str)
		assert.Equal(t, map[string][]string{}, props)
		assert.EqualError(t, err, "error parsing data for property myProp: could not parse empty property")
	})
	t.Run("valid value", func(t *testing.T) {
		str, props, err := ExtractString("hello", "myProp", map[string][]string{})
		assert.Equal(t, "hello", str)
		assert.Equal(t, map[string][]string{}, props)
		assert.NoError(t, err)
	})
	t.Run("extracted props", func(t *testing.T) {
		str, props, err := ExtractString("$hello$", "myProp", map[string][]string{"preExisting": {"something"}})
		assert.Equal(t, "", str)
		assert.Equal(t, map[string][]string{"preExisting": {"something"}, "hello": {"myProp"}}, props)
		assert.NoError(t, err)
	})
}

func TestExtractInt(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		i, props, err := ExtractInt("", "myProp", map[string][]string{})
		assert.Equal(t, 0, i)
		assert.Equal(t, map[string][]string{}, props)
		assert.EqualError(t, err, "error parsing data for property myProp: could not parse empty property")
	})
	t.Run("valid value", func(t *testing.T) {
		i, props, err := ExtractInt("72", "myProp", map[string][]string{})
		assert.Equal(t, 72, i)
		assert.Equal(t, map[string][]string{}, props)
		assert.NoError(t, err)
	})
	t.Run("extracted props", func(t *testing.T) {
		i, props, err := ExtractInt("$hello$", "myProp", map[string][]string{"preExisting": {"something"}})
		assert.Equal(t, 0, i)
		assert.Equal(t, map[string][]string{"preExisting": {"something"}, "hello": {"myProp"}}, props)
		assert.NoError(t, err)
	})
}
