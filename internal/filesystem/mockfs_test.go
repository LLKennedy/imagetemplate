package filesystem

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockReadFile(t *testing.T) {
	m := MockReader{
		Files: map[string]MockFile{
			"a file": {
				Data: []byte("hello!"),
				Err:  errors.New("an error"),
			},
		},
	}
	res, err := m.ReadFile("a file")
	assert.Equal(t, []byte("hello!"), res)
	assert.EqualError(t, err, "an error")
}
