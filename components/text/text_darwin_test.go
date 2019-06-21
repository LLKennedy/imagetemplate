package text

import (
	"testing"

	"golang.org/x/tools/godoc/vfs"

	"github.com/LLKennedy/gosysfonts"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	c, err := render.Decode("text")
	assert.NoError(t, err)
	assert.Equal(t, Component{fs: vfs.OS("."), fontPool: gosysfonts.OSXPool{}}, c)
}
