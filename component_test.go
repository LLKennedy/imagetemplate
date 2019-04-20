package imagetemplate

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDataValue(t *testing.T) {
	type parseTest struct {
		name               string
		value              string
		hasNamedProperties bool
		cleanValues        []string
		propNames          []string
		err                error
	}
	testArray := []parseTest{
		parseTest{
			name:               "simple value",
			value:              "some data",
			hasNamedProperties: false,
			cleanValues:        []string{"some data"},
			propNames:          []string{},
			err:                nil,
		},
		parseTest{
			name:               "simple value 2",
			value:              "123",
			hasNamedProperties: false,
			cleanValues:        []string{"123"},
			propNames:          []string{},
			err:                nil,
		},
		parseTest{
			name:               "simple value 3",
			value:              "#4: This. Is. Just. Data!~",
			hasNamedProperties: false,
			cleanValues:        []string{"#4: This. Is. Just. Data!~"},
			propNames:          []string{},
			err:                nil,
		},
		parseTest{
			name:               "escaped dollars",
			value:              "some $$ data 2 be pr$$ocessed!",
			hasNamedProperties: false,
			cleanValues:        []string{"some $ data 2 be pr$ocessed!"},
			propNames:          []string{},
			err:                nil,
		},
		parseTest{
			name:               "simple property",
			value:              "$username$",
			hasNamedProperties: true,
			cleanValues:        []string{"", ""},
			propNames:          []string{"username"},
			err:                nil,
		},
		parseTest{
			name:               "multiple properties with surrounding text",
			value:              "Hello there, $title$. $username$!",
			hasNamedProperties: true,
			cleanValues:        []string{"Hello there, ", ". ", "!"},
			propNames:          []string{"title", "username"},
			err:                nil,
		},
		parseTest{
			name:               "unclosed property",
			value:              "Hello there, $title$. $username!",
			hasNamedProperties: true,
			cleanValues:        []string{"Hello there, "},
			propNames:          []string{"title"},
			err:                errors.New("Unclosed named property in Hello there, $title$. $username!"),
		},
		parseTest{
			name:               "empty value",
			value:              "",
			hasNamedProperties: false,
			cleanValues:        []string{},
			propNames:          []string{},
			err:                errors.New("Could not parse empty property"),
		},
	}
	for _, test := range testArray {
		t.Run(test.name, func(t *testing.T) {
			hasNamedProperties, cleanValues, propNames, err := ParseDataValue(test.value)
			assert.Equal(t, test.hasNamedProperties, hasNamedProperties)
			assert.Equal(t, test.cleanValues, cleanValues)
			assert.Equal(t, test.propNames, propNames)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestConditionals(t *testing.T) {
	type testProperty struct {
		name   string
		value  interface{}
		setErr error
	}
	type testSet struct {
		name            string
		conditional     ComponentConditional
		namedProperties []testProperty
		validateResult  bool
		validateError   error
	}
	testArray := []testSet{
		testSet{
			name: "basic single condition",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "equals",
				Value:    "john smith",
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "username",
					value:  "john smith",
					setErr: nil,
				},
			},
			validateResult: true,
			validateError:  nil,
		},
	}
	for _, test := range testArray {
		t.Run(test.name, func(t *testing.T) {
			for _, prop := range test.namedProperties {
				t.Run(prop.name, func(t *testing.T) {
					err := test.conditional.SetValue(prop.name, prop.value)
					assert.Equal(t, prop.setErr, err)
				})
			}
			success, err := test.conditional.Validate()
			assert.Equal(t, test.validateResult, success)
			assert.Equal(t, test.validateError, err)
		})
	}
}
