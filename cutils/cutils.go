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

// TextAlignment is a text alignment.
type TextAlignment int

const (
	// TextAlignmentLeft aligns text left
	TextAlignmentLeft TextAlignment = iota
	// TextAlignmentRight aligns text right
	TextAlignmentRight
	// TextAlignmentCentre aligns text centrally
	TextAlignmentCentre
)

// StringToAlignment converts strings to TextAlignments, defaulting to Left
func StringToAlignment(alignment string) (converted TextAlignment) {
	switch alignment {
	case "left":
		converted = TextAlignmentLeft
	case "right":
		converted = TextAlignmentRight
	case "centre":
		converted = TextAlignmentCentre
	default:
		converted = TextAlignmentLeft
	}
	return
}
