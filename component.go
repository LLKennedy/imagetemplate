package imagetemplate

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

// Component provides a generic interface for operations to perform on a canvas
type Component interface {
	Write(canvas Canvas) error
	SetNamedProperties(properties []NamedProperty) error
	GetJSONFormat() interface{}
}

// CircleComponent implements the Component interface for circles
type CircleComponent struct {
	NamedPropertiesMap map[string][]string
	Centre             image.Point
	Radius             int
	Colour             color.NRGBA
}

// Write draws a circle on the canvas
func (component *CircleComponent) Write(canvas Canvas) error {
	if len(component.NamedPropertiesMap) != 0 {
		return fmt.Errorf("Cannot draw circle, not all named properties are set: %v", component.NamedPropertiesMap)
	}
	return canvas.Circle(component.Centre, component.Radius, component.Colour)
}

// SetNamedProperties processes the named properties and sets them into the circle properties
func (component *CircleComponent) SetNamedProperties(properties []NamedProperty) error {
	setFunc := func(name string, value interface{}) error {
		if strings.Contains("RGBA", name) && len(name) == 1 {
			//Process colours
			colourVal, ok := value.(uint8)
			if !ok {
				return fmt.Errorf("Error converting %v to uint8", value)
			}
			switch name {
			case "R":
				component.Colour.R = colourVal
				return nil
			case "G":
				component.Colour.G = colourVal
				return nil
			case "B":
				component.Colour.B = colourVal
				return nil
			case "A":
				component.Colour.A = colourVal
				return nil
			default:
				//What? How did you get here?
				return fmt.Errorf("Name was a string inside RGBA and Value was a valid uint8, but Name wasn't R, G, B, or A. Name was: %v", name)
			}
		}
		numberVal, ok := value.(int)
		if !ok {
			return fmt.Errorf("Error converting %v to int", value)
		}
		switch name {
		case "centreX":
			component.Centre.X = numberVal
			return nil
		case "centreY":
			component.Centre.Y = numberVal
			return nil
		case "radius":
			component.Radius = numberVal
			return nil
		default:
			return fmt.Errorf("Invalid component property in named property map: %v", name)
		}
	}
	var err error
	component.NamedPropertiesMap, err = ProcessNamedProperties(properties, component.NamedPropertiesMap, setFunc)
	if err != nil {
		return err
	}
	return nil
}

func (component *CircleComponent) GetJSONFormat() interface{} {
	type format struct {
		CentreX string `json:"centreX"`
		CentreY string `json:"centreY"`
		Radius  string `json:"radius"`
		Colour  struct {
			Red   string `json:"R"`
			Green string `json:"G"`
			Blue  string `json:"B"`
			Alpha string `json:"A"`
		} `json:"colour"`
	}
	return &format{}
}

// PropertySetFunc maps property names and values to component inner properties
type PropertySetFunc func(string, interface{}) error

// ProcessNamedProperties iterates over all named properties, retrieves their value, and calls the provided function to map properties to inner component properties. Each implementation of Component should call this within its SetNamedProperties function.
func ProcessNamedProperties(properties []NamedProperty, propMap map[string][]string, setFunc PropertySetFunc) (leftovers map[string][]string, err error) {
	for _, prop := range properties {
		name := prop.GetName()
		innerPropNames := propMap[name]
		if len(innerPropNames) <= 0 {
			// Not matching props, keep going
			continue
		}
		value, err := prop.GetValue()
		if err != nil {
			return propMap, err
		}
		for _, innerName := range innerPropNames {
			err = setFunc(innerName, value)
			if err != nil {
				return propMap, err
			}
		}
		delete(propMap, name)
	}
	return propMap, nil
}
