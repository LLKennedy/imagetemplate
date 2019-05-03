package imagetemplate

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
)

func TestWrite(t *testing.T) {
	newCircle := CircleComponent{
		Centre: image.Pt(4,4),
		Radius: 3,
		Colour: color.NRGBA{R:255,A:255},		
	}
	newCanvas := mockCanvas{FixedCircleError: nil}
	c, err := newCircle.Write(newCanvas)
	assert.NoError(t, err)
	assert.Equal(t, newCanvas, c)
}