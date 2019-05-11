package filesystem

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	r := IoutilFileReader{}
	data, err := r.ReadFile("///////\\some file you can't possibly find")
	assert.Equal(t, []byte(nil), data)
	assert.EqualError(t, err, "open ///////\\some file you can't possibly find: The specified path is invalid.")
}