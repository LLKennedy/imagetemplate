package render

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/tools/godoc/vfs"
)

var registry = map[string](func(vfs.FileSystem) Component){}

// RegisterComponent adds a new component to the registry, returning an error if duplicate names exist.
func RegisterComponent(name string, generator func(vfs.FileSystem) Component) error {
	if registry == nil {
		registry = map[string](func(vfs.FileSystem) Component){}
	}
	if registry[name] != nil {
		return fmt.Errorf("cannot register component, %v is already registered", name)
	}
	registry[name] = generator
	return nil
}

// Decode searches the registry for a component matching the provided name and returns a new blank component of that type.
func Decode(name string) (Component, error) {
	if registry == nil || registry[name] == nil {
		return nil, fmt.Errorf("component error: no component registered for name %v", name)
	}
	return registry[name](vfs.OS(".")), nil
}

// NamedProperties is a map of property names to property values - application variables to be set.
type NamedProperties map[string]interface{}

// Component provides a generic interface for operations to perform on a canvas.
type Component interface {
	Write(canvas Canvas) (Canvas, error)
	SetNamedProperties(properties NamedProperties) (Component, error)
	GetJSONFormat() interface{}
	VerifyAndSetJSONData(interface{}) (Component, NamedProperties, error)
}

// PropertySetFunc maps property names and values to component inner properties.
type PropertySetFunc func(string, interface{}) error

// StandardSetNamedProperties iterates over all named properties, retrieves their value, and calls the provided function to map properties to inner component properties. Each implementation of Component should call this within its SetNamedProperties function.
func StandardSetNamedProperties(properties NamedProperties, propMap map[string][]string, setFunc PropertySetFunc) (leftovers map[string][]string, err error) {
	for name, value := range properties {
		innerPropNames := propMap[name]
		if len(innerPropNames) <= 0 {
			// Not matching props, keep going
			continue
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

// DeconstructedDataValue is a string broken down into static values and property names. The reconstruction always starts with a static value, always has one more static value than props, and always alternates static, prop, static, prop... if any props exist.
type DeconstructedDataValue struct {
	// StaticValues are the true string components of the processed JSON value.
	StaticValues []string
	// PropNames are the extracted variable names from the processed JSON value.
	PropNames []string
}

func isSingleProp(d DeconstructedDataValue) bool {
	return len(d.PropNames) == 1 && len(d.StaticValues) == 2 && d.StaticValues[0] == "" && d.StaticValues[1] == ""
}

// PropType represents the types of properties which can be parsed.
type PropType string

const (
	// IntType is an int
	IntType PropType = "int"
	// StringType is a string
	StringType PropType = "string"
	// BoolType is a bool
	BoolType PropType = "bool"
	// Uint8Type is a uint8
	Uint8Type PropType = "uint8"
	// Float64Type is a float64
	Float64Type PropType = "float64"
	// TimeType is a *time.Time
	TimeType PropType = "time"
)

// ExtractSingleProp parses the loaded property configuration and application inputs and returns the desired property if it exists.
func ExtractSingleProp(inputVal, propName string, typeName PropType, namedPropsMap map[string][]string) (returnedPropsMap map[string][]string, ExtractedValue interface{}, err error) {
	var foundProps bool
	if returnedPropsMap, foundProps, err = parseToProps(inputVal, propName, namedPropsMap); err != nil || foundProps {
		return
	}
	switch typeName {
	case IntType:
		var int64Val int64
		int64Val, err = strconv.ParseInt(inputVal, 10, 64) //Use ParseInt instead of Atoi for compatibility with go 1.7
		if err != nil {
			err = fmt.Errorf("failed to convert property %v to integer: %v", propName, err)
		} else {
			ExtractedValue = int(int64Val)
		}
	case StringType:
		ExtractedValue = inputVal
	case BoolType:
		ExtractedValue, err = strconv.ParseBool(inputVal)
		if err != nil {
			err = fmt.Errorf("failed to convert property %v to bool: %v", propName, err)
		}
	case Uint8Type:
		var rawU uint64
		rawU, err = strconv.ParseUint(inputVal, 0, 8)
		if err != nil {
			err = fmt.Errorf("failed to convert property %v to uint8: %v", propName, err)
		}
		ExtractedValue = uint8(rawU)
	case Float64Type:
		ExtractedValue, err = strconv.ParseFloat(inputVal, 64)
		if err != nil {
			err = fmt.Errorf("failed to convert property %v to float64: %v", propName, err)
		}
	case TimeType:
		ExtractedValue, err = time.ParseDuration(inputVal)
		if err != nil {
			err = fmt.Errorf("failed to convert property %v to time.Duration: %v", propName, err)
		}
	default:
		err = fmt.Errorf("cannot convert property %v to unsupported type %v", propName, typeName)
	}
	if err != nil {
		returnedPropsMap = namedPropsMap
	}
	return returnedPropsMap, ExtractedValue, err
}

func parseToProps(inputVal, propName string, existingProps map[string][]string) (map[string][]string, bool, error) {
	npm := existingProps
	if npm == nil {
		npm = make(map[string][]string)
	}
	hasNamedProps, deconstructed, err := ParseDataValue(inputVal)
	if err != nil {
		return existingProps, false, fmt.Errorf("error parsing data for property %v: %v", propName, err)
	}
	if hasNamedProps {
		if !isSingleProp(deconstructed) {
			return existingProps, false, fmt.Errorf("composite properties are not yet supported: %v", inputVal)
		}
		customPropName := deconstructed.PropNames[0]
		npm[customPropName] = append(npm[propName], propName)
		return npm, true, nil
	}
	return npm, false, nil
}

// PropData is a matched triplet of input property data for use with extraction of exclusive properties.
type PropData struct {
	// InputValue is the raw JSON string data.
	InputValue string
	// PropName is the name of the property being sought, to associate with discovered variables.
	PropName string
	// Type is the conversion to attempt on the string.
	Type PropType
}

// ExtractExclusiveProp parses the loaded property configuration and application inputs and returns the desired property if it exists and if only one of the desired options exists.
func ExtractExclusiveProp(data []PropData, namedPropsMap map[string][]string) (returnedPropsMap map[string][]string, ExtractedValue interface{}, validIndex int, err error) {
	listSize := len(data)
	type result struct {
		props map[string][]string
		err   error
		value interface{}
	}
	resultArray := make([]*result, listSize)
	for i := range resultArray {
		resultArray[i] = &result{props: make(map[string][]string)}
	}
	setCount := 0
	validIndex = -1
	for i, datum := range data {
		aResult := resultArray[i]
		aResult.props, aResult.value, aResult.err = ExtractSingleProp(datum.InputValue, datum.PropName, datum.Type, aResult.props)
		if len(aResult.props) != 0 || aResult.err == nil { //This is an || because if a property has been added to the blank array, the function succeeded
			setCount++
			validIndex = i
		}
	}
	if setCount != 1 {
		concatString := ""
		for i, datum := range data {
			if i != 0 {
				concatString = concatString + ","
			}
			concatString = concatString + datum.PropName
		}
		return namedPropsMap, nil, -1, fmt.Errorf("exactly one of (%v) must be set", concatString)
	}
	returnedPropsMap = namedPropsMap
	for key, value := range resultArray[validIndex].props {
		if returnedPropsMap == nil {
			returnedPropsMap = make(map[string][]string)
		}
		returnedPropsMap[key] = append(returnedPropsMap[key], value...)
	}
	ExtractedValue = resultArray[validIndex].value
	err = nil
	return
}

// ParseDataValue determines whether a string represents raw data or a named variable and returns this information as well as the data cleaned of any variable definitions.
func ParseDataValue(value string) (hasNamedProperties bool, deconstructed DeconstructedDataValue, err error) {
	deconstructed = DeconstructedDataValue{}
	cleanString := ""
	if len(value) == 0 {
		err = fmt.Errorf("could not parse empty property")
		return
	}
	for i := 0; i < len(value); i++ {
		if value[i] != '$' {
			cleanString = cleanString + string(value[i])
			continue
		}
		var j int
		for j = i + 1; j < len(value); j++ {
			if value[j] == '$' {
				break
			}
		}
		if j >= len(value) || value[j] != '$' {
			err = fmt.Errorf("unclosed named property in '%v'", value)
			return
		}
		subString := value[i+1 : j]
		i = j
		if len(subString) == 0 {
			cleanString = cleanString + "$"
			continue
		}
		hasNamedProperties = true
		deconstructed.StaticValues = append(deconstructed.StaticValues, cleanString)
		cleanString = ""
		deconstructed.PropNames = append(deconstructed.PropNames, subString)
	}
	deconstructed.StaticValues = append(deconstructed.StaticValues, cleanString)
	return
}
