package image

import (
	"github.com/LLKennedy/imagetemplate/v2/render"
	"github.com/stretchr/testify/assert"
	"testing"
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
	expectedErr := "open /!!!\\!!!!!!/\\/\\/\\/\\/!!!!/\\!!!\\\\\\: no such file or directory"
	res, err := c.SetNamedProperties(input)
	assert.Equal(t, expected, res)
	assert.EqualError(t, err, expectedErr)
}
