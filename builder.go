package imagetemplate

import ()

// Builder manipulates Canvas objects and outputs to a bitmap
type Builder interface {
	NewCanvas(width, height int) (*Canvas, error)
}
