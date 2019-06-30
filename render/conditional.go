package render

import (
	"fmt"
	"strconv"
	"strings"
)

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
	// Name is the variable to check against the specified value.
	Name string `json:"name"`
	// Not determines whether to negate the final result of the boolean operation.
	Not bool `json:"boolNot"`
	// Operator specifies which comparison operation to perform.
	Operator conditionalOperator `json:"operator"`
	// Value is the condition to operate against with the variable specified by Name.
	Value string `json:"value"`
	// Group is an optional set of other conditionals to check along with this one.
	Group conditionalGroup `json:"group"`
	/*
		valueSet represents whether this individual component has had its value set
		and its condition evaluated at least once.
	*/
	valueSet bool
	/*
		validated represents whether this individual component at this level is
		validated. Use ComponentConditional.Validate() to evaluate the logic of
		entire groups.
	*/
	validated bool
}

type conditionalGroup struct {
	Operator     groupOperator          `json:"groupOperator"`
	Conditionals []ComponentConditional `json:"conditionals"`
}

// SetValue sets the value of a specific named property through this conditional chain, evaluating any conditions along the way.
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

// Validate validates this conditional chain, erroring if a value down the line has not been set and evaluated.
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

// GetNamedPropertiesList returns a list of all named props found in the conditional.
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
