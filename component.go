package imagetemplate

import (
	"fmt"
)

// Component provides a generic interface for operations to perform on a canvas
type Component interface {
	Write(canvas Canvas) error
	SetNamedProperties(properties []NamedProperty) error
	GetJSONFormat() interface{}
	VerifyAndSetJSONData(interface{}) ([]NamedProperty, error)
}

// PropertySetFunc maps property names and values to component inner properties
type PropertySetFunc func(string, interface{}) error

// StandardSetNamedProperties iterates over all named properties, retrieves their value, and calls the provided function to map properties to inner component properties. Each implementation of Component should call this within its SetNamedProperties function.
func StandardSetNamedProperties(properties []NamedProperty, propMap map[string][]string, setFunc PropertySetFunc) (leftovers map[string][]string, err error) {
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
		if value[j] != '$' {
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

All properties will be assumed to be either strings or integers based on the operator.

String operators: "equals", "contains", "startswith", "endswith", "ci_equals", "ci_contains", "ci_startswith", "ci_endswith". Operators including "ci_" are case-insensitive variants.

Integer operators: ">", "<", "<=", ">=".

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
	validated bool
}
