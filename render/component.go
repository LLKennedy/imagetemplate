package render

import (
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"strconv"
	"strings"
	"time"
)

var registry = map[string](func(vfs.FileSystem) Component){}

// RegisterComponent adds a new component to the registry, returning an error if duplicate names exist
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

// Decode searches the registry for a component matching the provided name and returns a new blank component of that type
func Decode(name string) (Component, error) {
	if registry == nil || registry[name] == nil {
		return nil, fmt.Errorf("component error: no component registered for name %v", name)
	}
	return registry[name](vfs.OS(".")), nil
}

// NamedProperties is a map of property names to property values - application variables to be set
type NamedProperties map[string]interface{}

// Component provides a generic interface for operations to perform on a canvas
type Component interface {
	Write(canvas Canvas) (Canvas, error)
	SetNamedProperties(properties NamedProperties) (Component, error)
	GetJSONFormat() interface{}
	VerifyAndSetJSONData(interface{}) (Component, NamedProperties, error)
}

// PropertySetFunc maps property names and values to component inner properties
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
	StaticValues []string
	PropNames    []string
}

func isSingleProp(d DeconstructedDataValue) bool {
	return len(d.PropNames) == 1 && len(d.StaticValues) == 2 && d.StaticValues[0] == "" && d.StaticValues[1] == ""
}

// PropType represents the types of properties which can be parsed
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

// ExtractSingleProp parses the loaded property configuration and application inputs and returns the desired property if it exists
func ExtractSingleProp(inputVal, propName string, typeName PropType, namedPropsMap map[string][]string) (returnedPropsMap map[string][]string, ExtractedValue interface{}, err error) {
	npm := namedPropsMap
	if npm == nil {
		npm = make(map[string][]string)
	}
	hasNamedProps, deconstructed, err := ParseDataValue(inputVal)
	if err != nil {
		return namedPropsMap, nil, fmt.Errorf("error parsing data for property %v: %v", propName, err)
	}
	if hasNamedProps {
		if !isSingleProp(deconstructed) {
			return namedPropsMap, nil, fmt.Errorf("composite properties are not yet supported: %v", inputVal)
		}
		customPropName := deconstructed.PropNames[0]
		npm[customPropName] = append(npm[propName], propName)
		return npm, nil, nil
	}
	switch typeName {
	case IntType:
		int64Val, err := strconv.ParseInt(inputVal, 10, 64) //Use ParseInt instead of Atoi for compatibility with go 1.7
		if err != nil {
			return namedPropsMap, nil, fmt.Errorf("failed to convert property %v to integer: %v", propName, err)
		}
		intVal := int(int64Val)
		return npm, intVal, nil
	case StringType:
		return npm, inputVal, nil
	case BoolType:
		boolVal, err := strconv.ParseBool(inputVal)
		if err != nil {
			return namedPropsMap, nil, fmt.Errorf("failed to convert property %v to bool: %v", propName, err)
		}
		return npm, boolVal, nil
	case Uint8Type:
		uintVal, err := strconv.ParseUint(inputVal, 0, 8)
		if err != nil {
			return namedPropsMap, nil, fmt.Errorf("failed to convert property %v to uint8: %v", propName, err)
		}
		uint8Val := uint8(uintVal)
		return npm, uint8Val, nil
	case Float64Type:
		float64Val, err := strconv.ParseFloat(inputVal, 64)
		if err != nil {
			return namedPropsMap, nil, fmt.Errorf("failed to convert property %v to float64: %v", propName, err)
		}
		return npm, float64Val, nil
	case TimeType:
		durationVal, err := time.ParseDuration(inputVal)
		if err != nil {
			return namedPropsMap, nil, fmt.Errorf("failed to convert property %v to time.Duration: %v", propName, err)
		}
		return npm, durationVal, nil
	}
	return namedPropsMap, nil, fmt.Errorf("cannot convert property %v to unsupported type %v", propName, typeName)
}

// PropData is a matched triplet of input property data for use with extraction of exclusive properties
type PropData struct {
	InputValue string
	PropName   string
	Type       PropType
}

// ExtractExclusiveProp parses the loaded property configuration and application inputs and returns the desired property if it exists and if only one of the desired options exists
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

// ParseDataValue determines whether a string represents raw data or a named variable and returns this information as well as the data cleaned of any variable definitions
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

type conditionalOperator string

const (
	equals         conditionalOperator = "equals"
	contains       conditionalOperator = "contains"
	startswith     conditionalOperator = "startswith"
	endswith       conditionalOperator = "endswith"
	ciEquals       conditionalOperator = "ci_equals"
	ciContains     conditionalOperator = "ci_contains"
	ciStartswith   conditionalOperator = "ci_startswith"
	ciEndswith     conditionalOperator = "ci_endswith"
	numequals      conditionalOperator = "=="
	lessthan       conditionalOperator = "<"
	greaterthan    conditionalOperator = ">"
	lessorequal    conditionalOperator = "<="
	greaterorequal conditionalOperator = ">="
)

type groupOperator string

const (
	or   groupOperator = "or"
	and  groupOperator = "and"
	nor  groupOperator = "nor"
	nand groupOperator = "nand"
	xor  groupOperator = "xor"
)

/*ComponentConditional enables or disables a component based on named properties.

All properties will be assumed to be either strings or floats based on the operator.

String operators: "equals", "contains", "startswith", "endswith", "ci_equals", "ci_contains", "ci_startswith", "ci_endswith". Operators including "ci_" are case-insensitive variants.

Float operators: "=", ">", "<", "<=", ">=".

Group operators can be "and", "or", "nand", "nor", "xor".*/
type ComponentConditional struct {
	Name      string              `json:"name"`
	Not       bool                `json:"boolNot"`
	Operator  conditionalOperator `json:"operator"`
	Value     string              `json:"value"`
	Group     conditionalGroup    `json:"group"`
	valueSet  bool                // Represents whether this individual component has had its value set and its condition evaluated at least once
	validated bool                // Represents whether this individual component at this level is validated. Use ComponentConditional.Validate() to evaluate the logic of entire groups.
}

type conditionalGroup struct {
	Operator     groupOperator          `json:"groupOperator"`
	Conditionals []ComponentConditional `json:"conditionals"`
}

// SetValue sets the value of a specific named property through this conditional chain, evaluating any conditions along the way
func (c ComponentConditional) SetValue(name string, value interface{}) (ComponentConditional, error) {
	conditional := c
	for conIndex, con := range conditional.Group.Conditionals {
		var err error
		conditional.Group.Conditionals[conIndex], err = con.SetValue(name, value)
		if err != nil {
			return c, err
		}
	}
	if conditional.Name == "" && !conditional.valueSet {
		conditional.validated = true
		conditional.valueSet = true
		return conditional, nil
	}
	if conditional.Name == name {
		switch conditional.Operator {
		case equals, contains, startswith, endswith, ciEquals, ciContains, ciStartswith, ciEndswith:
			// Handle string operators
			stringVal, ok := value.(string)
			if !ok {
				return c, fmt.Errorf("invalid value for string operator: %v", value)
			}
			conVal := conditional.Value
			switch conditional.Operator {
			case ciEquals:
				conVal = strings.ToLower(conVal)
				stringVal = strings.ToLower(stringVal)
				fallthrough
			case equals:
				conditional.validated = conVal == stringVal
			case ciContains:
				conVal = strings.ToLower(conVal)
				stringVal = strings.ToLower(stringVal)
				fallthrough
			case contains:
				conditional.validated = strings.Contains(stringVal, conVal)
			case ciStartswith:
				conVal = strings.ToLower(conVal)
				stringVal = strings.ToLower(stringVal)
				fallthrough
			case startswith:
				if len(conVal) > len(stringVal) {
					conditional.validated = false
					break
				}
				conditional.validated = stringVal[:len(conVal)] == conVal
			case ciEndswith:
				conVal = strings.ToLower(conVal)
				stringVal = strings.ToLower(stringVal)
				fallthrough
			case endswith:
				if len(conVal) > len(stringVal) {
					conditional.validated = false
					break
				}
				conditional.validated = stringVal[len(stringVal)-len(conVal):] == conVal
			}
		case numequals, lessthan, greaterthan, lessorequal, greaterorequal:
			// Handle float operators
			floatVal, ok := value.(float64)
			if !ok {
				intVal, ok := value.(int)
				if !ok {
					return c, fmt.Errorf("invalid value for float operator: %v", value)
				}
				floatVal = float64(intVal)
			}
			conVal, err := strconv.ParseFloat(conditional.Value, 64)
			if err != nil {
				return c, fmt.Errorf("failed to convert conditional value to float: %v", conditional.Value)
			}
			switch conditional.Operator {
			case numequals:
				conditional.validated = floatVal == conVal
			case lessthan:
				conditional.validated = floatVal < conVal
			case greaterthan:
				conditional.validated = floatVal > conVal
			case lessorequal:
				conditional.validated = floatVal <= conVal
			case greaterorequal:
				conditional.validated = floatVal >= conVal
			}
		default:
			return c, fmt.Errorf("invalid conditional operator %v", conditional.Operator)
		}
		if conditional.Not {
			conditional.validated = !conditional.validated
		}
		conditional.valueSet = true
	}
	return conditional, nil
}

// Validate validates this conditional chain, erroring if a value down the line has not been set and evaluated
func (c ComponentConditional) Validate() (bool, error) {
	if !c.valueSet && c.Name != "" {
		return false, fmt.Errorf("attempted to validate conditional %v %v %v without setting %v", c.Name, c.Operator, c.Value, c.Name)
	}
	group := c.Group.Conditionals
	if len(group) == 0 {
		return c.validated, nil
	}
	op := c.Group.Operator
	if op == xor {
		//Evaluate XOR on a group as meaning only one of all results in the list can be true, and one must be true.
		trueCount := 0
		if c.validated {
			trueCount++
		}
		for _, subConditional := range group {
			result, err := subConditional.Validate()
			if err != nil {
				return false, err
			}
			if result {
				trueCount++
			}
		}
		return trueCount == 1, nil
	}
	var result bool
	var negate bool
	if op == nand || op == nor {
		negate = true
	}
	result = c.validated
	if op == and || op == nand || op == or || op == nor {
		for _, subConditional := range group {
			subResult, err := subConditional.Validate()
			if err != nil {
				return false, err
			}
			if op == and || op == nand {
				result = result && subResult
			} else {
				result = result || subResult
			}
		}
		if negate {
			result = !result
		}
		return result, nil
	}
	return false, fmt.Errorf("invalid group operator %v", op)
}

// GetNamedPropertiesList returns a list of all named props found in the conditional
func (c ComponentConditional) GetNamedPropertiesList() NamedProperties {
	results := NamedProperties{}
	if c.Name == "" && len(c.Group.Conditionals) == 0 {
		return results
	}
	results[c.Name] = struct {
		Message string
	}{Message: "Please replace this struct with real data"}
	for _, subConditional := range c.Group.Conditionals {
		subResults := subConditional.GetNamedPropertiesList()
		for key, value := range subResults {
			results[key] = value
		}
	}
	return results
}
