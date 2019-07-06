package render

import (
	"errors"
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Text draws text on the canvas.
func (canvas ImageCanvas) Text(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (Canvas, error) {
	if maxWidth <= 0 {
		return canvas, errors.New("invalid maxWidth")
	}
	if canvas.Image == nil {
		return canvas, errors.New("no image set for canvas to draw on")
	}
	c := canvas
	drawer := &font.Drawer{
		Dot:  fixed.Point26_6{X: fixed.I(start.X), Y: fixed.I(start.Y)},
		Dst:  c.Image,
		Face: typeFace,
		Src:  image.NewUniform(colour),
	}
	width := drawer.MeasureString(text).Ceil()
	if width > maxWidth {
		return canvas, errors.New("resultant drawn text was longer than maxWidth")
	}
	drawer.DrawString(text)
	return c, nil
}

// TryText returns whether the text would fit on the canvas, and the width the text would currently use up.
func (canvas ImageCanvas) TryText(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (bool, int) {
	if maxWidth <= 0 {
		return false, -1
	}
	if canvas.Image == nil {
		return false, -2
	}
	drawer := &font.Drawer{
		Dot:  fixed.Point26_6{X: fixed.I(start.X), Y: fixed.I(start.Y)},
		Dst:  canvas.Image,
		Face: typeFace,
		Src:  image.NewUniform(colour),
	}
	width := drawer.MeasureString(text).Ceil()
	return width <= maxWidth, width
}
