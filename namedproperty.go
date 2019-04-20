package imagetemplate

import (
	"fmt"
)

// NamedProperty represents a variable or other application-specific property, the value of which needs to be mapped to properties in Components. The application using imagetemplate should retrieve the list of named properties and set all values before attempting to write from the template.
type NamedProperty interface {
	GetName() string
	SetValue(interface{}) error
	GetValue() (interface{}, error)
}

type StringProperty struct {
	Name  string
	value string
	isSet bool
}

func (p *StringProperty) GetName() string {
	return p.Name
}

// SetValue sets the value of a string property
func (p *StringProperty) SetValue(val interface{}) error {
	stringVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("Invalid property assignment, property %v expected type string", p.Name)
	}
	p.value = stringVal
	p.isSet = true
	return nil
}

// GetValue returns the value of a string property
func (p *StringProperty) GetValue() (interface{}, error) {
	if !p.isSet {
		return nil, fmt.Errorf("Attempted to access unset property %v, please assign a value first", p.Name)
	}
	return p.value, nil
}
