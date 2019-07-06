package render

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
)

// Rectangle draws a rectangle of a specific colour on the canvas.
func (canvas ImageCanvas) Rectangle(topLeft image.Point, width, height int, colour color.Color) (Canvas, error) {
	c := canvas
	if width <= 0 && height <= 0 {
		return canvas, errors.New("invalid width and height")
	} else if width <= 0 {
		return canvas, errors.New("invalid width")
	} else if height <= 0 {
		return canvas, errors.New("invalid height")
	}
	if c.Image == nil {
		return canvas, errors.New("no image set for canvas to draw on")
	}
	draw.Draw(c.Image, image.Rect(topLeft.X, topLeft.Y, topLeft.X+width, topLeft.Y+height), image.NewUniform(colour), topLeft, draw.Over)
	return c, nil
}

// Circle draws a circle of a specific colour on the canvas.
func (canvas ImageCanvas) Circle(centre image.Point, radius int, colour color.Color) (Canvas, error) {
	c := canvas
	if radius <= 0 {
		return canvas, errors.New("invalid radius")
	}
	if c.Image == nil {
		return canvas, errors.New("no image set for canvas to draw on")
	}
	mask := &circle{p: centre, r: radius}
	draw.DrawMask(c.Image, mask.Bounds(), image.NewUniform(colour), image.ZP, mask, mask.Bounds().Min, draw.Over)
	return c, nil
}

// Steal the circle example code from https://blog.golang.org/go-imagedraw-package.
type circle struct {
	p image.Point
	r int
}
