package render

import (
	"errors"
	"image"
	"image/draw"
)

// DrawImage draws another image on the canvas.
func (canvas ImageCanvas) DrawImage(start image.Point, subImage image.Image) (Canvas, error) {
	c := canvas
	if c.Image == nil {
		return canvas, errors.New("no image set for canvas to draw on")
	}
	subBounds := subImage.Bounds()
	width := subBounds.Max.X - subBounds.Min.X
	height := subBounds.Max.Y - subBounds.Min.Y
	draw.Draw(c.Image, image.Rect(start.X, start.Y, start.X+width, start.Y+height), subImage, image.ZP, draw.Over)
	return c, nil
}
