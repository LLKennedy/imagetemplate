package filesystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadFile(t *testing.T) {
	r := IoutilFileReader{}
	data, err := r.ReadFile("///////\\some file you can't possibly find")
	assert.Equal(t, []byte(nil), data)
	assert.EqualError(t, err, "open ///////\\some file you can't possibly find: no such file or directory")
}
