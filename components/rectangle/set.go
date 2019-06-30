package rectangle

import (
	"fmt"

	"github.com/LLKennedy/imagetemplate/v3/cutils"
)

func (component *Component) delegatedSetProperties(name string, value interface{}) (err error) {
	switch name {
	case "R":
		component.Colour.R, err = cutils.SetUint8(value)
	case "G":
		component.Colour.G, err = cutils.SetUint8(value)
	case "B":
		component.Colour.B, err = cutils.SetUint8(value)
	case "A":
		component.Colour.A, err = cutils.SetUint8(value)
	case "topLeftX":
		component.TopLeft.X, err = cutils.SetInt(value)
	case "topLeftY":
		component.TopLeft.Y, err = cutils.SetInt(value)
	case "width":
		component.Width, err = cutils.SetInt(value)
	case "height":
		component.Height, err = cutils.SetInt(value)
	default:
		err = fmt.Errorf("invalid component property in named property map: %v", name)
	}
	return
}
