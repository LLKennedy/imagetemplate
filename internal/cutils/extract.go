package cutils

import "github.com/LLKennedy/imagetemplate/v3/render"

// ExtractString extracts a string or variable(s) from the raw JSON data
func ExtractString(raw, name string, props map[string][]string) (string, map[string][]string, error) {
	newProps, newVal, err := render.ExtractSingleProp(raw, name, render.StringType, props)
	if err != nil {
		return "", props, err
	}
	var foundString string
	if newVal != nil {
		foundString = newVal.(string)
	}
	return foundString, newProps, nil
}

// ExtractInt extracts an integer or variable(s) from the raw JSON data
func ExtractInt(raw, name string, props map[string][]string) (int, map[string][]string, error) {
	newProps, newVal, err := render.ExtractSingleProp(raw, name, render.IntType, props)
	if err != nil {
		return 0, props, err
	}
	var foundInt int
	if newVal != nil {
		foundInt = newVal.(int)
	}
	return foundInt, newProps, nil
}

// ExtractFloat extracts a flot64 or variable(s) from the raw JSON data
func ExtractFloat(raw, name string, props map[string][]string) (float64, map[string][]string, error) {
	newProps, newVal, err := render.ExtractSingleProp(raw, name, render.Float64Type, props)
	if err != nil {
		return 0, props, err
	}
	var foundInt float64
	if newVal != nil {
		foundInt = newVal.(float64)
	}
	return foundInt, newProps, nil
}
