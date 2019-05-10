package imagetemplate

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestImageSetNamedPropertiesLoadRealFile(t *testing.T) {
	c := ImageComponent{
		NamedPropertiesMap: map[string][]string{
			"aProp":[]string{"fileName"},
		},
	}
	input := NamedProperties{
		"aProp": "!!!\\!!!!!!//\\//\\//\\/\\/!!!!//\\!!!\\\\\\////",
	}
	expected := ImageComponent{
		NamedPropertiesMap: map[string][]string{
			"aProp":[]string{"fileName"},
		},
		reader: ioutilFileReader{},
	}
	expectedErr := "open !!!\\!!!!!!//\\//\\//\\/\\/!!!!//\\!!!\\\\\\////: The system cannot find the path specified."
	res, err := c.SetNamedProperties(input)
	assert.Equal(t, expected, res)
	assert.EqualError(t, err, expectedErr)
}