package imagetemplate

import (
	"testing"

	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/stretchr/testify/assert"
)

func TestLoadTemplate(t *testing.T) {
	props, _, err := LoadTemplate("//\\/")
	assert.Equal(t, render.NamedProperties(nil), props)
	assert.EqualError(t, err, "open /\\: no such file or directory")
}
