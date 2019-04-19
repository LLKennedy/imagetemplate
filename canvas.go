package imagetemplate

import (
	"errors"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Canvas holds the image struct and associated properties
type Canvas interface {
	SetUnderlyingImage(newImage draw.Image)
	GetUnderlyingImage() image.Image
	GetWidth() int
	GetHeight() int
	Rectangle(topLeft image.Point, width, height int, colour color.Color) error
	Circle(centre image.Point, radius int, alignment CircleAlignment, colour color.Color) error
	Text(start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) error
	SubImage(start image.Point, subImage image.Image) error
}

// ImageCanvas uses golang's native Image package to implement the Canvas interface
type ImageCanvas struct {
	Image draw.Image
}

// NewCanvas generates a new canvas of the given width and height
func NewCanvas(width, height int) (*ImageCanvas, error) {
	if width <= 0 && height <= 0 {
		return nil, errors.New("Invalid width and height")
	} else if width <= 0 {
		return nil, errors.New("Invalid width")
	} else if height <= 0 {
		return nil, errors.New("Invalid height")
	}
	return &ImageCanvas{
		Image: image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: width, Y: height},
		}),
	}, nil
}

// SetUnderlyingImage sets the internal Image property to the given object
func (canvas *ImageCanvas) SetUnderlyingImage(newImage draw.Image) {
	canvas.Image = newImage
}

// GetUnderlyingImage gets the internal Image property
func (canvas *ImageCanvas) GetUnderlyingImage() image.Image {
	return canvas.Image
}

// GetWidth returns the width of the underlying Image
func (canvas *ImageCanvas) GetWidth() int {
	return canvas.Image.Bounds().Size().X
}

// GetWidth returns the width of the underlying Image
func (canvas *ImageCanvas) GetHeight() int {
	return canvas.Image.Bounds().Size().Y
}

// Rectangle draws a rectangle of a specific colour on the canvas
func (canvas *ImageCanvas) Rectangle(topLeft image.Point, width, height int, colour color.Color) error {
	colourPlane := image.Uniform{C: colour}
	if width <= 0 && height <= 0 {
		return errors.New("Invalid width and height")
	} else if width <= 0 {
		return errors.New("Invalid width")
	} else if height <= 0 {
		return errors.New("Invalid height")
	}
	rect := image.Rectangle{
		Min: image.Point{
			X: topLeft.X,
			Y: topLeft.Y,
		},
		Max: image.Point{
			X: topLeft.X + width,
			Y: topLeft.Y + height,
		},
	}
	draw.Draw(canvas.Image, rect, &colourPlane, topLeft, draw.Over)
	return nil
}

// CircleAlignment dicates where precisely a circle is centred
type CircleAlignment int

const (
	// CAlignPixelTopLeft puts the circle centre in the top left of the pizel, so the diameter is exactly 2r and the circle can be placed anywhere except the bottom and right edges
	CAlignPixelTopLeft = iota
	// CAlignPIxelTopRight puts the circle centre in the top right of the pixel, so the diameter is exactly 2r and the circle can be placed anywhere except the bottom and left edges
	CAlignPixelTopRight
	// CAlignPIxelBottomLeft puts the circle centre in the bottom left of the pixel, so the diameter is exactly 2r and the circle can be placed anywhere except the top and right edges
	CAlignPixelBottomLeft
	// CAlignPIxelBottomRight puts the circle centre in the bottom right of the pixel, so the diameter is exactly 2r and the circle can be placed anywhere except the top and left edges
	CAlignPixelBottomRight
	// CAlignPixelCentre puts the circle centre in the middle of the pixel, so the diameter is exactly 2r-1
	CAlignPixelCentre
	// CAlignPixelAuto attempts to determine the most mathematically accurate option between Centre and Top Left
	CAlignPixelAuto
)

// Circle draws a circle of a specific colour on the canvas
func (canvas *ImageCanvas) Circle(centre image.Point, radius int, alignment CircleAlignment, colour color.Color) error {
	colourPlane := image.Uniform{C: colour}
	var startX, startY, enclosingWidth int
	//TODO: Handle the rest of these
	switch alignment {
	case CAlignPixelTopLeft:
		startX = centre.X - radius
		startY = centre.Y - radius
		enclosingWidth = 2 * radius
		break
	case CAlignPixelTopRight:
		startX = centre.X - (radius - 1)
		startY = centre.Y - radius
		enclosingWidth = 2 * radius
		break
	case CAlignPixelBottomLeft:
		startX = centre.X - radius
		startY = centre.Y - (radius - 1)
		enclosingWidth = 2 * radius
		break
	case CAlignPixelBottomRight:
		startX = centre.X - (radius - 1)
		startY = centre.Y - (radius - 1)
		enclosingWidth = 2 * radius
		break
	case CAlignPixelCentre:
		startX = centre.X - (radius - 1)
		startY = centre.Y - (radius - 1)
		enclosingWidth = (2 * radius) - 1
		break
	case CAlignPixelAuto:
		return errors.New("Not implemented yet")
		break
	default:
		return errors.New("Invalid circle alignment setting: " + string(alignment))
	}
	enclosingRectangle := image.Rectangle{
		Min: image.Point{X: startX, Y: startY},
		Max: image.Point{X: startX + enclosingWidth - 1, Y: startY + enclosingWidth - 1},
	}
	mask := image.NewNRGBA(enclosingRectangle)
	for x := startX; x < startX+enclosingWidth; x++ {
		for y := startY; y < startY+enclosingWidth; y++ {
			// Set base distance
			centreDelta := image.Point{
				X: x - startX,
				Y: y - startY,
			}
			cDeltaAbs := image.Point{
				// Thanks for the implicit conversions Go
				X: int(math.Abs(float64(centreDelta.X))),
				Y: int(math.Abs(float64(centreDelta.Y))),
			}
			// Correct based on circle centre position
			switch alignment {
			case CAlignPixelTopLeft:
				if centreDelta.X >= 0 {
					centreDelta.X += 1
				}
				if centreDelta.Y >= 0 {
					centreDelta.Y += 1
				}
				break
			case CAlignPixelTopRight:
				if centreDelta.X <= 0 {
					centreDelta.X -= 1
				}
				if centreDelta.Y >= 0 {
					centreDelta.Y += 1
				}
				break
			case CAlignPixelBottomLeft:
				if centreDelta.X >= 0 {
					centreDelta.X += 1
				}
				if centreDelta.Y <= 0 {
					centreDelta.Y -= 1
				}
				break
			case CAlignPixelBottomRight:
				if centreDelta.X <= 0 {
					centreDelta.X -= 1
				}
				if centreDelta.Y <= 0 {
					centreDelta.Y -= 1
				}
				break
			case CAlignPixelCentre:
				// Reduce magnitude of each axis by 1
				if centreDelta.X != 0 {
					centreDelta.X += centreDelta.X / cDeltaAbs.X
				}
				if centreDelta.Y != 0 {
					centreDelta.Y += centreDelta.Y / cDeltaAbs.Y
				}
				break
			case CAlignPixelAuto:
				return errors.New("Not implemented yet")
				break
			default:
				return errors.New("Invalid circle alignment setting: " + string(alignment))
			}
			// Calculate current polar distance from the centre of the circle
			currentRadius := int(math.Sqrt(math.Pow(float64(centreDelta.X), 2) + math.Pow(float64(centreDelta.Y), 2)))
			if currentRadius <= radius {
				mask.Set(x, y, color.Opaque)
			} else {
				mask.Set(x, y, color.Transparent)
			}
		}
	}
	draw.DrawMask(canvas.Image, enclosingRectangle, &colourPlane, image.Point{X: startX, Y: startY}, mask, image.Point{}, draw.Over)
	return nil
}

// Text draws text on the canvas
func (canvas *ImageCanvas) Text(start image.Point, typeFace font.Face, colour color.Color, fontSize, maxWidth, maxLines int) error {
	return errors.New("Not implemented yet")
}

// SubImage draws another image on the canvas
func (canvas *ImageCanvas) SubImage(start image.Point, subImage image.Image) error {
	return errors.New("Not implemented yet")
}
