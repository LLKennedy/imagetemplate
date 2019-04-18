package imagetemplate

import (
	"image"
	"image/color"
)

// Canvas holds the image struct and associated properties
type Canvas interface {
	Rectangle(topLeftX, topLeftY, width, height int, colour color.Color)
}

type TemplateCanvas struct {
	Image image.Image
}
