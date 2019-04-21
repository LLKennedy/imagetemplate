package imagetemplate

import (
	"fmt"
	"strconv"
	"strings"
)

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

// ParseDataValue determines whether a string represents raw data or a named variable and returns this information as well as the data cleaned of any variable definitions
func ParseDataValue(value string) (hasNamedProperties bool, cleanValues []string, propNames []string, err error) {
	cleanValues = []string{}
	propNames = []string{}
	cleanString := ""
	if len(value) == 0 {
		err = fmt.Errorf("Could not parse empty property")
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
			err = fmt.Errorf("Unclosed named property in %v", value)
			return
		}
		subString := value[i+1 : j]
		i = j
		if len(subString) == 0 {
			cleanString = cleanString + "$"
			continue
		}
		hasNamedProperties = true
		cleanValues = append(cleanValues, cleanString)
		cleanString = ""
		propNames = append(propNames, subString)
	}
	cleanValues = append(cleanValues, cleanString)
	return
}

type conditionalOperator string

const (
	equals         conditionalOperator = "equals"
	contains       conditionalOperator = "contains"
	startswith     conditionalOperator = "startswith"
	endswith       conditionalOperator = "endswith"
	ci_equals      conditionalOperator = "ci_equals"
	ci_contains    conditionalOperator = "ci_contains"
	ci_startswith  conditionalOperator = "ci_startswith"
	ci_endswith    conditionalOperator = "ci_endswith"
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

/* ComponentConditional enables or disables a component based on named properties.

All properties will be assumed to be either strings or floats based on the operator.

String operators: "equals", "contains", "startswith", "endswith", "ci_equals", "ci_contains", "ci_startswith", "ci_endswith". Operators including "ci_" are case-insensitive variants.

Float operators: "=", ">", "<", "<=", ">=".

Group operators can be "and", "or", "nand", "nor", "xor".*/
type ComponentConditional struct {
	Name     string              `json:"name"`
	Not      bool                `json:"boolNot"`
	Operator conditionalOperator `json:"operator"`
	Value    string              `json:"value"`
	Group    struct {
		Operator     groupOperator          `json:"groupOperator"`
		Conditionals []ComponentConditional `json:"conditionals"`
	} `json:"group"`
	valueSet  bool // Represents whether this individual component has had its value set and its condition evaluated at least once
	validated bool // Represents whether this individual component at this level is validated. Use ComponentConditional.Validate() to evaluate the logic of entire groups.
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
		case equals, contains, startswith, endswith, ci_equals, ci_contains, ci_startswith, ci_endswith:
			// Handle string operators
			stringVal, ok := value.(string)
			if !ok {
				return c, fmt.Errorf("Invalid value for string operator: %v", value)
			}
			conVal := conditional.Value
			switch conditional.Operator {
			case ci_equals:
				conVal = strings.ToLower(conVal)
				stringVal = strings.ToLower(stringVal)
				fallthrough
			case equals:
				conditional.validated = conVal == stringVal
			case ci_contains:
				conVal = strings.ToLower(conVal)
				stringVal = strings.ToLower(stringVal)
				fallthrough
			case contains:
				conditional.validated = strings.Contains(stringVal, conVal)
			case ci_startswith:
				conVal = strings.ToLower(conVal)
				stringVal = strings.ToLower(stringVal)
				fallthrough
			case startswith:
				if len(conVal) > len(stringVal) {
					conditional.validated = false
					break
				}
				conditional.validated = stringVal[:len(conVal)] == conVal
			case ci_endswith:
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
					return c, fmt.Errorf("Invalid value for float operator: %v", value)
				}
				floatVal = float64(intVal)
			}
			conVal, err := strconv.ParseFloat(conditional.Value, 64)
			if err != nil {
				return c, fmt.Errorf("Failed to convert conditional value to float: %v", conditional.Value)
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
			return c, fmt.Errorf("Invalid conditional operator %v", conditional.Operator)
		}
		if conditional.Not {
			conditional.validated = !conditional.validated
		}
		conditional.valueSet = true
	}
	return conditional, nil
}

// Validate validates this conditional chain, erroring if a value down the line has not been set and evaluated
func (conditional ComponentConditional) Validate() (bool, error) {
	if !conditional.valueSet {
		return false, fmt.Errorf("Attempted to validate conditional %v %v %v without setting %v", conditional.Name, conditional.Operator, conditional.Value, conditional.Name)
	}
	group := conditional.Group.Conditionals
	if len(group) == 0 {
		return conditional.validated, nil
	}
	op := conditional.Group.Operator
	if op == xor {
		//Evaluate XOR on a group as meaning only one of all results in the list can be true, and one must be true.
		trueCount := 0
		if conditional.validated {
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
	result = conditional.validated
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
	return false, fmt.Errorf("Invalid group operator %v", op)
}

// GetNamedProps returns a list of all named props found in the conditional
func (conditional ComponentConditional) GetNamedPropertiesList() NamedProperties {
	var results NamedProperties
	type invalidData struct {
		Message string
	}
	results[conditional.Name] = invalidData{Message: "Please replace this struct with real data"}
	for _, subConditional := range conditional.Group.Conditionals {
		subResults := subConditional.GetNamedPropertiesList()
		for key, value := range subResults {
			results[key] = value
		}
	}
	return results
}
