package cutils

import "github.com/LLKennedy/imagetemplate/v3/render"

// ExtractString extracts a string or variable(s) from the raw JSON data
func ExtractString(raw, name string, props map[string][]string) (string, map[string][]string, error) {
	newProps, newVal, err := render.ExtractSingleProp(raw, name, render.StringType, props)
	if err != nil {
		return "", props, err
	}
	foundString := ""
	if newVal != nil {
		foundString = newVal.(string)
	}
	return foundString, newProps, nil
}
