package image

import (
	"testing"

	"github.com/LLKennedy/imagetemplate/v2/render"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs"
)

func TestImageSetNamedPropertiesLoadRealFile(t *testing.T) {
	c := Component{
		NamedPropertiesMap: map[string][]string{
			"aProp": {"fileName"},
		},
	}
	input := render.NamedProperties{
		"aProp": "!!!\\!!!!!!//\\//\\//\\/\\/!!!!//\\!!!\\\\\\////",
	}
	expected := Component{
		NamedPropertiesMap: map[string][]string{
			"aProp": {"fileName"},
		},
		fs: vfs.OS(""),
	}
	expectedErr := "open \\!!!\\!!!!!!\\!!!!\\!!!: The system cannot find the path specified."
	res, err := c.SetNamedProperties(input)
	assert.Equal(t, expected, res)
	assert.EqualError(t, err, expectedErr)
}
