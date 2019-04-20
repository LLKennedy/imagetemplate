package imagetemplate

import (
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
