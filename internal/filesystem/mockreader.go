package filesystem

// MockFile is some mock file data
type MockFile struct {
	Data []byte
	Err  error
}

// MockReader is a mock implementation of the Reader interface, for testing purposes
type MockReader struct {
	Files map[string]MockFile
}

// ReadFile returns a pre-set data/error pair
func (m MockReader) ReadFile(filename string) ([]byte, error) {
	return m.Files[filename].Data, m.Files[filename].Err
}
