package filesystem

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs"
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
		mFile, isMock := res.(*MockFile)
		assert.True(t, isMock)
		mFile.buf = nil
		readRes, err := ioutil.ReadAll(res)
		assert.EqualError(t, err, "cannot read from nil file")
		assert.Empty(t, readRes)
	})
	t.Run("invalid seek", func(t *testing.T) {
		mFile, isMock := res.(*MockFile)
		assert.True(t, isMock)
		mFile.buf = nil
		seekRes, err := res.Seek(3, io.SeekStart)
		assert.EqualError(t, err, "cannot seek in nil file")
		assert.Equal(t, seekRes, int64(0))
	})
	t.Run("close", func(t *testing.T) {
		assert.NoError(t, res.Close())
	})
	t.Run("lstat", func(t *testing.T) {
		m.On("Lstat", "a").Return(NewMockFile(nil), nil)
		stats, err := m.Lstat("a")
		assert.Equal(t, NewMockFile(nil), stats)
		assert.NoError(t, err)
	})
	t.Run("stat", func(t *testing.T) {
		m.On("Stat", "b").Return(NewMockFile([]byte{}), errors.New("some error"))
		stats, err := m.Stat("b")
		assert.Equal(t, NewMockFile([]byte{}), stats)
		assert.EqualError(t, err, "some error")
	})
	t.Run("readdir", func(t *testing.T) {
		m.On("ReadDir", "c").Return([]os.FileInfo{NewMockFile(nil)}, nil)
		files, err := m.ReadDir("c")
		assert.NoError(t, err)
		assert.Equal(t, []os.FileInfo{NewMockFile(nil)}, files)
	})
	t.Run("roottype", func(t *testing.T) {
		m.On("RootType", "d").Return(vfs.RootTypeGoPath)
		rtype := m.RootType("d")
		assert.Equal(t, vfs.RootTypeGoPath, rtype)
	})
	t.Run("string", func(t *testing.T) {
		m.On("String").Return("hello")
		str := m.String()
		assert.Equal(t, "hello", str)
	})
	m.AssertExpectations(t)
}

func TestMockFile(t *testing.T) {
	file := NewMockFile(nil)
	t.Run("name", func(t *testing.T) {
		file.On("Name").Return("a file")
		name := file.Name()
		assert.Equal(t, "a file", name)
	})
	t.Run("size", func(t *testing.T) {
		file.On("Size").Return(int64(100))
		size := file.Size()
		assert.Equal(t, int64(100), size)
	})
	t.Run("mode", func(t *testing.T) {
		file.On("Mode").Return(os.ModeCharDevice)
		mode := file.Mode()
		assert.Equal(t, os.ModeCharDevice, mode)
	})
	t.Run("modtime", func(t *testing.T) {
		orig := time.Now()
		file.On("ModTime").Return(orig)
		modTime := file.ModTime()
		assert.Equal(t, orig, modTime)
	})
	t.Run("isdir", func(t *testing.T) {
		file.On("IsDir").Return(true)
		assert.True(t, file.IsDir())
	})
	t.Run("sys", func(t *testing.T) {
		file.On("Sys").Return(struct{ a int }{a: 12})
		sys := file.Sys()
		assert.Equal(t, struct{ a int }{a: 12}, sys)
	})
	file.AssertExpectations(t)
}
