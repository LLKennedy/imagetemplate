package image

import (
	fs "github.com/LLKennedy/imagetemplate/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/render"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImageSetNamedPropertiesLoadRealFile(t *testing.T) {
	c := Component{
		NamedPropertiesMap: map[string][]string{
			"aProp": []string{"fileName"},
		},
	}
	input := render.NamedProperties{
		"aProp": "!!!\\!!!!!!//\\//\\//\\/\\/!!!!//\\!!!\\\\\\////",
	}
	expected := Component{
		NamedPropertiesMap: map[string][]string{
			"aProp": []string{"fileName"},
		},
		reader: fs.IoutilFileReader{},
	}
	expectedErr := "open !!!\\!!!!!!//\\//\\//\\/\\/!!!!//\\!!!\\\\\\////: The system cannot find the path specified."
	res, err := c.SetNamedProperties(input)
	assert.Equal(t, expected, res)
	assert.EqualError(t, err, expectedErr)
}
