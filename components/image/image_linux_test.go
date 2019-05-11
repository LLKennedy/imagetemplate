package image

import (
	"testing"
	"github.com/stretchr/testify/assert"
	fs "github.com/LLKennedy/imagetemplate/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/render"
)

func TestImageSetNamedPropertiesLoadRealFile(t *testing.T) {
	c := Component{
		NamedPropertiesMap: map[string][]string{
			"aProp":[]string{"fileName"},
		},
	}
	input := render.NamedProperties{
		"aProp": "!!!\\!!!!!!//\\//\\//\\/\\/!!!!//\\!!!\\\\\\////",
	}
	expected := Component{
		NamedPropertiesMap: map[string][]string{
			"aProp":[]string{"fileName"},
		},
		reader: fs.IoutilFileReader{},
	}
	expectedErr := "open !!!\\!!!!!!//\\//\\//\\/\\/!!!!//\\!!!\\\\\\////: no such file or directory"
	res, err := c.SetNamedProperties(input)
	assert.Equal(t, expected, res)
	assert.EqualError(t, err, expectedErr)
}