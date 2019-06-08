package imagetemplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTemplate(t *testing.T) {
	b, err := LoadTemplate("")
	assert.Nil(t, b)
	assert.NoError(t, err)
}
