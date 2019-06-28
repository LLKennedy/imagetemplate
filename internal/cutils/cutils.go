// Package cutils provides common parsing/conversion code for components to cut down on duplication
package cutils

import (
	"fmt"
)

// CombineErrors combines two errors, maintaining a history of errors separated by newlines
func CombineErrors(history, latest error) error {
	switch {
	case history == nil:
		return latest
	case latest == nil:
		return history
	}
	return fmt.Errorf("%v\n%v", history, latest)
}
