package filesystem

import (
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockFS(t *testing.T) {
	m := NewMockReader()
	buffer := NewMockFile([]byte("hello!"))
	m.On("Open", "a file").Return(buffer, errors.New("an error"))
	res, err := m.Open("a file")
	assert.EqualError(t, err, "an error")
	t.Run("valid read", func(t *testing.T) {
		readRes, err := ioutil.ReadAll(res)
		assert.NoError(t, err)
		assert.Equal(t, []byte("hello!"), readRes)
	})
	t.Run("valid seek", func(t *testing.T) {
		seekRes, err := res.Seek(3, io.SeekStart)
		assert.NoError(t, err)
		assert.Equal(t, seekRes, int64(3))
	})
	t.Run("invalid read", func(t *testing.T) {
		mFile, isMock := res.(*mockFile)
		assert.True(t, isMock)
		mFile.buf = nil
		readRes, err := ioutil.ReadAll(res)
		assert.EqualError(t, err, "cannot read from nil file")
		assert.Empty(t, readRes)
	})
	t.Run("invalid seek", func(t *testing.T) {
		mFile, isMock := res.(*mockFile)
		assert.True(t, isMock)
		mFile.buf = nil
		seekRes, err := res.Seek(3, io.SeekStart)
		assert.EqualError(t, err, "cannot seek in nil file")
		assert.Equal(t, seekRes, int64(0))
	})
	t.Run("close", func(t *testing.T) {
		assert.NoError(t, res.Close())
	})
}
