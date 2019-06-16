package filesystem

import "github.com/stretchr/testify/mock"

// MockReader is a mock implementation of the Reader interface, for testing purposes
type MockReader struct {
	mock.Mock
}

// ReadFile returns a pre-set data/error pair
func (m MockReader) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	return (args.Get(0)).([]byte), args.Error(1)
}
