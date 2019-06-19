package filesystem

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs"
)

func TestMockReadFile(t *testing.T) {
	m := new(MockReader)
	buffer := vfs.ReadSeekCloser(NewMockFile([]byte("hello!")))
	m.On("Open", "a file").Return(buffer, errors.New("an error"))
	res, err := m.Open("a file")
	assert.EqualError(t, err, "an error")
	readRes, err := ioutil.ReadAll(res)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello!"), readRes)
}
