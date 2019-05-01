package imagetemplate

import (
	"fmt"
	"strconv"
	"strings"
)

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

type propType string

const (
	intType    propType = "int"
	stringType propType = "string"
	boolType   propType = "bool"
	uint8Type  propType = "uint8"
	float64Type propType = "float64"
)

func extractSingleProp(inputVal, propName string, typeName propType, namedPropsMap map[string][]string) (returnedPropsMap map[string][]string, extractedValue interface{}, err error) {
	npm := namedPropsMap
	hasNamedProps, deconstructed, err := ParseDataValue(inputVal)
	if err != nil {
		return namedPropsMap, nil, err
	}
	if hasNamedProps {
		if !isSingleProp(deconstructed) {
			return namedPropsMap, nil, fmt.Errorf("composite properties are not yet supported: %v", inputVal)
		}
		propName := deconstructed.PropNames[0]
		npm[propName] = append(npm[propName], propName)
		return npm, nil, nil
	}
	switch typeName {
	case intType:
		intVal, err := strconv.Atoi(inputVal)
		if err != nil {
			return namedPropsMap, nil, fmt.Errorf("failed to convert property %v to integer: %v", propName, err)
		}
		return npm, intVal, nil
	case stringType:
		return npm, inputVal, nil
	case boolType:
		boolVal, err := strconv.ParseBool(inputVal)
		if err != nil {
			return namedPropsMap, nil, fmt.Errorf("failed to convert property %v to bool: %v", propName, err)
		}
		return npm, boolVal, nil
	case uint8Type:
		uintVal, err := strconv.ParseUint(inputVal, 0, 8)
		if err != nil {
			return namedPropsMap, nil, err
		}
		uint8Val := uint8(uintVal)
		return npm, uint8Val, nil
	case float64Type:
		float64Val, err := strconv.ParseFloat(inputVal, 64)
		if err != nil {
			return namedPropsMap, nil, err
		}
		return npm, float64Val, nil
	}
	return namedPropsMap, nil, fmt.Errorf("cannot convert property %v to unsupported type %v", propName, typeName)
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
	if !c.valueSet {
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
	type invalidData struct {
		Message string
	}
	results[c.Name] = invalidData{Message: "Please replace this struct with real data"}
	for _, subConditional := range c.Group.Conditionals {
		subResults := subConditional.GetNamedPropertiesList()
		for key, value := range subResults {
			results[key] = value
		}
	}
	return results
}
