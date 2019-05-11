// Package filesystem provides file system access
package filesystem

import "io/ioutil"

// FileReader is a wrapper for file system access
type FileReader interface {
	ReadFile(string) ([]byte, error)
}

// IoutilFileReader uses io/ioutil to implement FileReader
type IoutilFileReader struct{}

// ReadFile returns the file data found at the file path
func (r IoutilFileReader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}