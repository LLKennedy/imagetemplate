package filesystem

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestMockReadFile(t *testing.T) {
	m := MockReader{
		Files: map[string]MockFile{
			"a file": MockFile{
				Data: []byte("hello!"), 
				Err: errors.New("an error"),
			},
		},
	}
	res, err := m.ReadFile("a file")
	assert.Equal(t, []byte("hello!"), res)
	assert.EqualError(t, err, "an error")
}