package imagetemplate

import (
	"testing"

	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/stretchr/testify/assert"
)

func TestLoadTemplate(t *testing.T) {
	_, props, err := New(nil).Load().FromFile("//\\/")
	assert.Equal(t, render.NamedProperties(nil), props)
	assert.EqualError(t, err, "Open: //\\/ is a directory")
}
