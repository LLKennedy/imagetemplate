// Package render renders images onto a canvas.
package render

import (
	"errors"
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/tools/godoc/vfs"
)

// Canvas holds the image struct and associated properties.
type Canvas interface {
	SetUnderlyingImage(newImage image.Image) Canvas
	GetUnderlyingImage() image.Image
	GetWidth() int
	GetHeight() int
	GetPPI() float64
	SetPPI(float64) Canvas
	Rectangle(topLeft image.Point, width, height int, colour color.Color) (Canvas, error)
	Circle(centre image.Point, radius int, colour color.Color) (Canvas, error)
	Text(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (Canvas, error)
	TryText(text string, start image.Point, typeFace font.Face, colour color.Color, maxWidth int) (bool, int)
	DrawImage(start image.Point, subImage image.Image) (Canvas, error)
	Barcode(codeType BarcodeType, content []byte, extra BarcodeExtraData, start image.Point, width, height int, dataColour color.Color, bgColour color.Color) (Canvas, error)
}

// ImageCanvas uses golang's native Image package to implement the Canvas interface.
type ImageCanvas struct {
	// Image is the underlying drawable image used for rendering.
	Image draw.Image
	// fs is the file system
	fs vfs.FileSystem
	// pixelsPerInch is the PPI of the canvas
	pixelsPerInch float64
}

// NewCanvas generates a new canvas of the given width and height.
func NewCanvas(width, height int) (ImageCanvas, error) {
	if width <= 0 && height <= 0 {
		return ImageCanvas{}, errors.New("invalid width and height")
	} else if width <= 0 {
		return ImageCanvas{}, errors.New("invalid width")
	} else if height <= 0 {
		return ImageCanvas{}, errors.New("invalid height")
	}
	return ImageCanvas{
		Image: image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: width, Y: height},
		}),
		fs: vfs.OS("."),
	}, nil
}

// SetUnderlyingImage sets the internal Image property to the given object.
func (canvas ImageCanvas) SetUnderlyingImage(newImage image.Image) Canvas {
	if newImage == nil {
		canvas.Image = nil
		return canvas
	}
	drawImage, ok := newImage.(draw.Image)
	if !ok {
		bounds := newImage.Bounds()
		drawImage = image.NewNRGBA(bounds)
		draw.Draw(drawImage, bounds, newImage, bounds.Min, draw.Src)
	}
	canvas.Image = drawImage
	return canvas
}

// GetUnderlyingImage gets the internal Image property.
func (canvas ImageCanvas) GetUnderlyingImage() image.Image {
	if canvas.Image == nil {
		return image.NewNRGBA(image.Rect(0, 0, 0, 0))
	}
	return canvas.Image
}

// GetWidth returns the width of the underlying Image. Returns 0 if no canvas is set.
func (canvas ImageCanvas) GetWidth() int {
	if canvas.Image == nil {
		return 0
	}
	return canvas.Image.Bounds().Size().X
}

// GetHeight returns the height of the underlying Image. Returns 0 if no canvas is set.
func (canvas ImageCanvas) GetHeight() int {
	if canvas.Image == nil {
		return 0
	}
	return canvas.Image.Bounds().Size().Y
}

// SetPPI sets the pixels per inch of the canvas.
func (canvas ImageCanvas) SetPPI(ppi float64) Canvas {
	canvas.pixelsPerInch = ppi
	return canvas
}

// GetPPI returns the pixels per inch of the canvas.
func (canvas ImageCanvas) GetPPI() float64 {
	return canvas.pixelsPerInch
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
