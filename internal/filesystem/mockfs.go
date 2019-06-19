package filesystem

import (
	"bytes"
	"fmt"
	"os"

	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

// MockReader is a mock implementation of the Opener interface, for testing purposes
type MockReader struct {
	mock.Mock
}

// NewMockReader creates a new mock Reader
func NewMockReader(files ...mockFile) *MockReader {
	newReader := new(MockReader)
	_ = vfs.FileSystem(newReader) // this will fail to compile if the interface isn't met
	return newReader
}

// Open returns a pre-set data/error pair
func (m *MockReader) Open(filename string) (vfs.ReadSeekCloser, error) {
	args := m.Called(filename)
	return (args.Get(0)).(vfs.ReadSeekCloser), args.Error(1)
}

// Lstat does stat stuff
func (m *MockReader) Lstat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	return (args.Get(0)).(os.FileInfo), args.Error(1)
}

// Stat gets file stats
func (m *MockReader) Stat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	return (args.Get(0)).(os.FileInfo), args.Error(1)
}

// ReadDir walks the directory
func (m *MockReader) ReadDir(path string) ([]os.FileInfo, error) {
	args := m.Called(path)
	return (args.Get(0)).([]os.FileInfo), args.Error(1)
}

// RootType returns the Root Type
func (m *MockReader) RootType(path string) vfs.RootType {
	args := m.Called(path)
	return (args.Get(0)).(vfs.RootType)
}

func (m *MockReader) String() string {
	args := m.Called()
	return args.String(0)
}

// mockFile is the simple implementation of vfs.ReadSeekCloser
type mockFile struct {
	buf *bytes.Reader
}

// NewMockFile creates a new mock file for use with the mock file system
func NewMockFile(data []byte) vfs.ReadSeekCloser {
	file := &mockFile{
		buf: bytes.NewReader(data),
	}
	return file
}

// Close does nothing
func (m *mockFile) Close() error {
	return nil
}

// Read reads from the buffer
func (m *mockFile) Read(dst []byte) (int, error) {
	if m == nil || m.buf == nil {
		return 0, fmt.Errorf("cannot read from nil file")
	}
	return m.buf.Read(dst)
}

func (m *mockFile) Seek(offset int64, whence int) (int64, error) {
	if m == nil || m.buf == nil {
		return 0, fmt.Errorf("cannot seek in nil file")
	}
	return m.buf.Seek(offset, whence)
}
