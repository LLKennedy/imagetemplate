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

// ExclusiveOr returns true if one and only one of the passed in booleans is true
func ExclusiveOr(args ...bool) bool {
	trueCount := 0
	for _, arg := range args {
		if arg {
			trueCount++
		}
	}
	return trueCount == 1
}

// ExclusiveNor returns false if one and only one of the passed in booleans is true
func ExclusiveNor(args ...bool) bool {
	return !ExclusiveOr(args...)
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

// ScaleFontsToWidth scales the input float to match the font size and alignment parameters
func ScaleFontsToWidth(currentSize float64, currentWidth, maxWidth int, alignment TextAlignment) (newSize float64, alignmentOffset int) {
	newSize = currentSize
	if currentWidth > maxWidth {
		ratio := float64(maxWidth) / float64(currentWidth)
		newSize = ratio * currentSize
	} else if currentWidth < maxWidth {
		remainingWidth := float64(maxWidth) - float64(currentWidth)
		switch alignment {
		case TextAlignmentLeft:
			alignmentOffset = 0
		case TextAlignmentRight:
			alignmentOffset = int(remainingWidth)
		case TextAlignmentCentre:
			alignmentOffset = int(remainingWidth / 2)
		default:
			alignmentOffset = 0
		}
	}
	return
}
