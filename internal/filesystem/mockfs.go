package filesystem

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

// NilFile is a nil mockfile for easy reference
var NilFile *MockFile

// MockFileSystem is a mock implementation of the Opener interface, for testing purposes
type MockFileSystem struct {
	mock.Mock
}

// NewMockFileSystem creates a new mock Reader
func NewMockFileSystem(files ...*MockFile) *MockFileSystem {
	newReader := new(MockFileSystem)
	for _, file := range files {
		newReader.On("Open", file.name).Return(file, nil)
	}
	_ = vfs.FileSystem(newReader) // this will fail to compile if the interface isn't met
	return newReader
}

// Open returns a pre-set data/error pair
func (m *MockFileSystem) Open(filename string) (vfs.ReadSeekCloser, error) {
	args := m.Called(filename)
	return (args.Get(0)).(vfs.ReadSeekCloser), args.Error(1)
}

// Lstat does stat stuff
func (m *MockFileSystem) Lstat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	return (args.Get(0)).(os.FileInfo), args.Error(1)
}

// Stat gets file stats
func (m *MockFileSystem) Stat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	return (args.Get(0)).(os.FileInfo), args.Error(1)
}

// ReadDir walks the directory
func (m *MockFileSystem) ReadDir(path string) ([]os.FileInfo, error) {
	args := m.Called(path)
	return (args.Get(0)).([]os.FileInfo), args.Error(1)
}

// RootType returns the Root Type
func (m *MockFileSystem) RootType(path string) vfs.RootType {
	args := m.Called(path)
	return (args.Get(0)).(vfs.RootType)
}

func (m *MockFileSystem) String() string {
	args := m.Called()
	return args.String(0)
}

// MockFile is the simple implementation of vfs.ReadSeekCloser
type MockFile struct {
	name string
	mock.Mock
	buf *bytes.Reader
}

// NewMockFile creates a new mock file for use with the mock file system
func NewMockFile(name string, data []byte) *MockFile {
	file := &MockFile{
		name: name,
		buf:  bytes.NewReader(data),
	}
	_ = vfs.ReadSeekCloser(file)
	_ = os.FileInfo(file) // check interface conformation
	return file
}

// Close does nothing
func (m *MockFile) Close() error {
	return nil
}

// Read reads from the buffer
func (m *MockFile) Read(dst []byte) (int, error) {
	if m == nil || m.buf == nil {
		return 0, fmt.Errorf("cannot read from nil file")
	}
	return m.buf.Read(dst)
}

// Seek sets the read/write header
func (m *MockFile) Seek(offset int64, whence int) (int64, error) {
	if m == nil || m.buf == nil {
		return 0, fmt.Errorf("cannot seek in nil file")
	}
	return m.buf.Seek(offset, whence)
}

// Name returns the base name of the file
func (m *MockFile) Name() string {
	args := m.Called()
	return args.String(0)
}

// Size returns the length in bytes for regular files; system-dependent for others
func (m *MockFile) Size() int64 {
	args := m.Called()
	return (args.Get(0)).(int64)
}

// Mode returns the file mode bits
func (m *MockFile) Mode() os.FileMode {
	args := m.Called()
	return (args.Get(0)).(os.FileMode)
}

// ModTime returns the modification time
func (m *MockFile) ModTime() time.Time {
	args := m.Called()
	return (args.Get(0)).(time.Time)
}

// IsDir is an abbreviation for Mode().IsDir()
func (m *MockFile) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

// Sys returns the underlying data source (can return nil)
func (m *MockFile) Sys() interface{} {
	args := m.Called()
	return args.Get(0)
}
