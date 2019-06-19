package render

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs"
)

type mockComponent struct{}

func (c mockComponent) Write(canvas Canvas) (Canvas, error) {
	return nil, nil
}

func (c mockComponent) SetNamedProperties(properties NamedProperties) (Component, error) {
	return c, nil
}

func (c mockComponent) GetJSONFormat() interface{} {
	return nil
}

func (c mockComponent) VerifyAndSetJSONData(interface{}) (Component, NamedProperties, error) {
	return nil, nil, nil
}

func newMock(vfs.FileSystem) Component {
	return &mockComponent{}
}

func TestRegisterComponentAndDecode(t *testing.T) {
	registry = nil
	err := RegisterComponent("newComponent", newMock)
	assert.NoError(t, err)
	err = RegisterComponent("newComponent", newMock)
	assert.EqualError(t, err, "cannot register component, newComponent is already registered")
	registry = nil
	c, err := Decode("newComponent")
	assert.Nil(t, c)
	assert.EqualError(t, err, "component error: no component registered for name newComponent")
	err = RegisterComponent("newComponent", newMock)
	assert.NoError(t, err)
	c, err = Decode("wrong")
	assert.Nil(t, c)
	assert.EqualError(t, err, "component error: no component registered for name wrong")
	c, err = Decode("newComponent")
	assert.Equal(t, newMock(vfs.OS("")), c)
	assert.NoError(t, err)
}

func TestDecode(t *testing.T) {
}

func TestStandardSetNamedProperties(t *testing.T) {
	type setPropTest struct {
		name            string
		properties      NamedProperties
		propMap         map[string][]string
		setFunc         PropertySetFunc
		resultLeftovers map[string][]string
		resultErr       error
	}
	successSetFunc := func(string, interface{}) error {
		return nil
	}
	errorSetFunc := func(string, interface{}) error {
		return errors.New("failed to set property")
	}
	testArray := []setPropTest{
		{
			name: "single property success",
			properties: NamedProperties{
				"username": "john smith",
			},
			propMap: map[string][]string{
				"username": {"innerPropUsername"},
			},
			setFunc:         successSetFunc,
			resultLeftovers: map[string][]string{},
			resultErr:       nil,
		},
		{
			name: "single property success",
			properties: NamedProperties{
				"username": "john smith",
			},
			propMap: map[string][]string{
				"username": {"innerPropUsername"},
			},
			setFunc: errorSetFunc,
			resultLeftovers: map[string][]string{
				"username": {"innerPropUsername"},
			},
			resultErr: errors.New("failed to set property"),
		},
		{
			name: "many properties success",
			properties: NamedProperties{
				"username": "john smith",
				"age":      57,
				"title":    "Mr.",
			},
			propMap: map[string][]string{
				"username": {"innerPropUsername", "innerPropEmail"},
				"age":      {"innerPropAge"},
				"title":    {"innerPropTitle", "TITLE", "innerPropRank", "some random field"},
			},
			setFunc:         successSetFunc,
			resultLeftovers: map[string][]string{},
			resultErr:       nil,
		},
	}
	for _, test := range testArray {
		t.Run(test.name, func(t *testing.T) {
			leftovers, err := StandardSetNamedProperties(test.properties, test.propMap, test.setFunc)
			assert.Equal(t, test.resultLeftovers, leftovers)
			if test.resultErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.resultErr.Error())
			}
		})
	}
	t.Run("check all internal properties are passed through", func(t *testing.T) {
		valuesSet := map[string]bool{
			"innerPropUsername": false,
			"innerPropEmail":    false,
			"innerPropAge":      false,
			"innerPropTitle":    false,
			"TITLE":             false,
			"innerPropRank":     false,
			"some random field": false,
		}
		expectedValues := NamedProperties{
			"innerPropUsername": "john smith",
			"innerPropEmail":    "john smith",
			"innerPropAge":      57,
			"innerPropTitle":    "Mr.",
			"TITLE":             "Mr.",
			"innerPropRank":     "Mr.",
			"some random field": "Mr.",
		}
		testSetFunc := func(name string, value interface{}) error {
			assert.Equal(t, expectedValues[name], value)
			valuesSet[name] = true
			return nil
		}
		test := setPropTest{
			name: "all properties get passed through",
			properties: NamedProperties{
				"username": "john smith",
				"age":      57,
				"title":    "Mr.",
				"unused":   "whatever",
			},
			propMap: map[string][]string{
				"username": {"innerPropUsername", "innerPropEmail"},
				"age":      {"innerPropAge"},
				"title":    {"innerPropTitle", "TITLE", "innerPropRank", "some random field"},
			},
			setFunc:         testSetFunc,
			resultLeftovers: map[string][]string{},
			resultErr:       nil,
		}
		leftovers, err := StandardSetNamedProperties(test.properties, test.propMap, test.setFunc)
		assert.Equal(t, test.resultLeftovers, leftovers)
		if test.resultErr == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.resultErr.Error())
		}
		assert.Equal(t, len(expectedValues), len(valuesSet), "valuesSet length changes, something was added or deleted improperly")
		finalResult := true
		for _, result := range valuesSet {
			finalResult = finalResult && result
		}
		assert.True(t, finalResult)
	})
}

func TestIsSingleProp(t *testing.T) {
	t.Run("failing on nil props", func(t *testing.T) {
		val := DeconstructedDataValue{PropNames: nil}
		assert.False(t, isSingleProp(val))
	})
	t.Run("failing on prop length", func(t *testing.T) {
		val := DeconstructedDataValue{PropNames: []string{}}
		assert.False(t, isSingleProp(val))
	})
	t.Run("failing on nil static values", func(t *testing.T) {
		val := DeconstructedDataValue{PropNames: []string{"someProp"}, StaticValues: nil}
		assert.False(t, isSingleProp(val))
	})
	t.Run("failing on static values length (0)", func(t *testing.T) {
		val := DeconstructedDataValue{PropNames: []string{"someProp"}, StaticValues: []string{}}
		assert.False(t, isSingleProp(val))
	})
	t.Run("failing on static values length (1)", func(t *testing.T) {
		val := DeconstructedDataValue{PropNames: []string{"someProp"}, StaticValues: []string{"a"}}
		assert.False(t, isSingleProp(val))
	})
	t.Run("failing on static values length (3)", func(t *testing.T) {
		val := DeconstructedDataValue{PropNames: []string{"someProp"}, StaticValues: []string{"a", "b", "c"}}
		assert.False(t, isSingleProp(val))
	})
	t.Run("failing on first static value", func(t *testing.T) {
		val := DeconstructedDataValue{PropNames: []string{"someProp"}, StaticValues: []string{"not empty", ""}}
		assert.False(t, isSingleProp(val))
	})
	t.Run("failing on second static value", func(t *testing.T) {
		val := DeconstructedDataValue{PropNames: []string{"someProp"}, StaticValues: []string{"", "not empty"}}
		assert.False(t, isSingleProp(val))
	})
	t.Run("passing", func(t *testing.T) {
		val := DeconstructedDataValue{PropNames: []string{"someProp"}, StaticValues: []string{"", ""}}
		assert.True(t, isSingleProp(val))
	})
}

func TestExtractSingleProp(t *testing.T) {
	type testSet struct {
		name             string
		inputVal         string
		propName         string
		typeName         PropType
		namedPropsMap    map[string][]string
		returnedPropsMap map[string][]string
		extractedValue   interface{}
		err              error
	}
	testFunc := func(t *testing.T, test testSet) {
		returnedPropsMap, extractedValue, err := ExtractSingleProp(test.inputVal, test.propName, test.typeName, test.namedPropsMap)
		assert.Equal(t, test.returnedPropsMap, returnedPropsMap)
		assert.Equal(t, test.extractedValue, extractedValue)
		if test.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.err.Error())
		}
	}
	tests := []testSet{
		{
			name:     "error in input",
			inputVal: "$",
			propName: "x",
			err:      errors.New("error parsing data for property x: unclosed named property in '$'"),
		},
		{
			name:     "multiple properties",
			inputVal: "$a$$b$",
			err:      errors.New("composite properties are not yet supported: $a$$b$"),
		},
		{
			name:             "valid single prop",
			inputVal:         "$a$",
			propName:         "myInternalProp",
			returnedPropsMap: map[string][]string{"a": {"myInternalProp"}},
			err:              nil,
		},
		{
			name:     "invalid type",
			inputVal: "laskdjf;alsdf",
			propName: "bad data",
			typeName: PropType("something else"),
			err:      errors.New("cannot convert property bad data to unsupported type something else"),
		},
		{
			name:     "invalid int",
			inputVal: "laskdjf;alsdf",
			propName: "aProp",
			typeName: IntType,
			err:      errors.New("failed to convert property aProp to integer: strconv.ParseInt: parsing \"laskdjf;alsdf\": invalid syntax"),
		},
		{
			name:     "invalid bool",
			inputVal: "laskdjf;alsdf",
			propName: "aProp",
			typeName: BoolType,
			err:      errors.New("failed to convert property aProp to bool: strconv.ParseBool: parsing \"laskdjf;alsdf\": invalid syntax"),
		},
		{
			name:     "invalid uint8",
			inputVal: "laskdjf;alsdf",
			propName: "aProp",
			typeName: Uint8Type,
			err:      errors.New("failed to convert property aProp to uint8: strconv.ParseUint: parsing \"laskdjf;alsdf\": invalid syntax"),
		},
		{
			name:     "invalid float64",
			inputVal: "laskdjf;alsdf",
			propName: "aProp",
			typeName: Float64Type,
			err:      errors.New("failed to convert property aProp to float64: strconv.ParseFloat: parsing \"laskdjf;alsdf\": invalid syntax"),
		},
		{
			name:     "invalid float64",
			inputVal: "laskdjf;alsdf",
			propName: "aProp",
			typeName: TimeType,
			err:      errors.New("failed to convert property aProp to time.Duration: time: invalid duration laskdjf;alsdf"),
		},
		{
			name:             "valid int",
			inputVal:         "53",
			propName:         "aProp",
			typeName:         IntType,
			err:              nil,
			returnedPropsMap: map[string][]string{},
			extractedValue:   53,
		},
		{
			name:             "valid string",
			inputVal:         "53",
			propName:         "aProp",
			typeName:         StringType,
			err:              nil,
			returnedPropsMap: map[string][]string{},
			extractedValue:   "53",
		},
		{
			name:             "valid bool",
			inputVal:         "true",
			propName:         "aProp",
			typeName:         BoolType,
			err:              nil,
			returnedPropsMap: map[string][]string{},
			extractedValue:   true,
		},
		{
			name:             "valid uint8",
			inputVal:         "53",
			propName:         "aProp",
			typeName:         Uint8Type,
			err:              nil,
			returnedPropsMap: map[string][]string{},
			extractedValue:   uint8(53),
		},
		{
			name:             "valid float64",
			inputVal:         "53",
			propName:         "aProp",
			typeName:         Float64Type,
			err:              nil,
			returnedPropsMap: map[string][]string{},
			extractedValue:   float64(53),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFunc(t, test)
		})
	}
	t.Run("valid *time.Time", func(t *testing.T) {
		// 	inputVal:         "2000h",
		// 	propName:         "aProp",
		// 	typeName:         TimeType,
		// 	err:              nil,
		// 	returnedPropsMap: map[string][]string{},
		// 	extractedValue:   time.Now().Add(time.Duration(2000)),
		dur, err := time.ParseDuration("2000h")
		assert.NoError(t, err)
		returnedPropsMap, extractedValue, err := ExtractSingleProp("2000h", "aProp", TimeType, nil)
		assert.NoError(t, err)
		expectedTime := time.Now().Add(dur)
		assert.NoError(t, err)
		expectedString := expectedTime.Format("2006/01/02")
		assert.Equal(t, map[string][]string{}, returnedPropsMap)
		convertedTime, ok := extractedValue.(*time.Time)
		if !ok || convertedTime == nil {
			t.Fatal("failed to convert extracted value to *time.Time")
		}
		assert.Equal(t, expectedString, convertedTime.Format("2006/01/02"))
	})
}

func TestExtractExclusiveProp(t *testing.T) {
	type testSet struct {
		name               string
		propData           []PropData
		namedPropsMap      map[string][]string
		returnedPropsMap   map[string][]string
		extractedValue     interface{}
		returnedValidIndex int
		err                error
	}
	testFunc := func(t *testing.T, test testSet) {
		returnedPropsMap, extractedValue, validIndex, err := ExtractExclusiveProp(test.propData, test.namedPropsMap)
		assert.Equal(t, test.returnedPropsMap, returnedPropsMap)
		assert.Equal(t, test.extractedValue, extractedValue)
		assert.Equal(t, test.returnedValidIndex, validIndex)
		if test.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.err.Error())
		}
	}
	tests := []testSet{
		{
			name:               "no data",
			returnedValidIndex: -1,
			err:                errors.New("exactly one of () must be set"),
		},
		{
			name: "single invalid prop option",
			propData: []PropData{
				{
					InputValue: "a",
					PropName:   "myProp",
					Type:       IntType,
				},
			},
			returnedValidIndex: -1,
			err:                errors.New("exactly one of (myProp) must be set"),
		},
		{
			name: "single valid prop option",
			propData: []PropData{
				{
					InputValue: "6",
					PropName:   "myProp",
					Type:       IntType,
				},
			},
			returnedValidIndex: 0,
			extractedValue:     6,
			err:                nil,
		},
		{
			name: "multiple valid prop options",
			propData: []PropData{
				{
					InputValue: "6",
					PropName:   "myProp",
					Type:       IntType,
				},
				{
					InputValue: "something",
					PropName:   "anotherProp",
					Type:       StringType,
				},
			},
			returnedValidIndex: -1,
			err:                errors.New("exactly one of (myProp,anotherProp) must be set"),
		},
		{
			name: "single valid named prop option",
			propData: []PropData{
				{
					InputValue: "$setMeLater$",
					PropName:   "myProp",
					Type:       IntType,
				},
			},
			returnedValidIndex: 0,
			returnedPropsMap:   map[string][]string{"setMeLater": {"myProp"}},
			err:                nil,
		},
		{
			name: "multiple props, only one valid",
			propData: []PropData{
				{
					InputValue: "$setMeLater$",
					PropName:   "myProp",
					Type:       IntType,
				},
				{
					InputValue: "a",
					PropName:   "someProp",
					Type:       IntType,
				},
				{
					InputValue: "-67",
					PropName:   "nothing",
					Type:       Uint8Type,
				},
			},
			returnedValidIndex: 0,
			returnedPropsMap:   map[string][]string{"setMeLater": {"myProp"}},
			err:                nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFunc(t, test)
		})
	}
}

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
		{
			name:               "simple value",
			value:              "some data",
			hasNamedProperties: false,
			cleanValues:        []string{"some data"},
			propNames:          []string(nil),
			err:                nil,
		},
		{
			name:               "simple value 2",
			value:              "123",
			hasNamedProperties: false,
			cleanValues:        []string{"123"},
			propNames:          []string(nil),
			err:                nil,
		},
		{
			name:               "simple value 3",
			value:              "#4: This. Is. Just. Data!~",
			hasNamedProperties: false,
			cleanValues:        []string{"#4: This. Is. Just. Data!~"},
			propNames:          []string(nil),
			err:                nil,
		},
		{
			name:               "escaped dollars",
			value:              "some $$ data 2 be pr$$ocessed!",
			hasNamedProperties: false,
			cleanValues:        []string{"some $ data 2 be pr$ocessed!"},
			propNames:          []string(nil),
			err:                nil,
		},
		{
			name:               "simple property",
			value:              "$username$",
			hasNamedProperties: true,
			cleanValues:        []string{"", ""},
			propNames:          []string{"username"},
			err:                nil,
		},
		{
			name:               "multiple properties with surrounding text",
			value:              "Hello there, $title$. $username$!",
			hasNamedProperties: true,
			cleanValues:        []string{"Hello there, ", ". ", "!"},
			propNames:          []string{"title", "username"},
			err:                nil,
		},
		{
			name:               "unclosed property",
			value:              "Hello there, $title$. $username!",
			hasNamedProperties: true,
			cleanValues:        []string{"Hello there, "},
			propNames:          []string{"title"},
			err:                errors.New("unclosed named property in 'Hello there, $title$. $username!'"),
		},
		{
			name:               "empty value",
			value:              "",
			hasNamedProperties: false,
			cleanValues:        []string(nil),
			propNames:          []string(nil),
			err:                errors.New("could not parse empty property"),
		},
	}
	for _, test := range testArray {
		t.Run(test.name, func(t *testing.T) {
			hasNamedProperties, deconstructed, err := ParseDataValue(test.value)
			assert.Equal(t, test.hasNamedProperties, hasNamedProperties)
			assert.Equal(t, test.cleanValues, deconstructed.StaticValues)
			assert.Equal(t, test.propNames, deconstructed.PropNames)
			if test.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.err.Error())
			}
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
		foundNamedProps []string
		namedProperties []testProperty
		validateResult  bool
		validateError   error
	}
	var test testSet
	testFunc := func(t *testing.T) {
		namedProperties := test.conditional.GetNamedPropertiesList()
		assert.Equal(t, len(test.foundNamedProps), len(namedProperties))
		for key := range namedProperties {
			found := false
			for _, prop := range test.foundNamedProps {
				if prop == key {
					found = true
					break
				}
			}
			assert.True(t, found)
		}
		for _, prop := range test.namedProperties {
			t.Run(prop.name, func(t *testing.T) {
				var err error
				test.conditional, err = test.conditional.SetValue(prop.name, prop.value)
				if prop.setErr == nil {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, prop.setErr.Error())
				}
			})
		}
		success, err := test.conditional.Validate()
		assert.Equal(t, test.validateResult, success)
		if test.validateError == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.validateError.Error())
		}
	}

	testArray := []testSet{
		{
			name: "single string condition",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "equals",
				Value:    "john smith",
			},
			foundNamedProps: []string{"username"},
			namedProperties: []testProperty{
				{
					name:   "username",
					value:  "john smith",
					setErr: nil,
				},
			},
			validateResult: true,
			validateError:  nil,
		},
		{
			name: "string condition, int set value",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "equals",
				Value:    "john smith",
			},
			foundNamedProps: []string{"username"},
			namedProperties: []testProperty{
				{
					name:   "username",
					value:  18,
					setErr: errors.New("invalid value for string operator: 18"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional username equals john smith without setting username"),
		},
		{
			name: "string condition, mismatched value name",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "equals",
				Value:    "john smith",
			},
			foundNamedProps: []string{"username"},
			namedProperties: []testProperty{
				{
					name:   "UserName",
					value:  "john smith",
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional username equals john smith without setting username"),
		},
		{
			name: "string equals, fails on case sensitivity",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "equals",
				Value:    "john smith",
			},
			foundNamedProps: []string{"username"},
			namedProperties: []testProperty{
				{
					name:   "username",
					value:  "John Smith",
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  nil,
		},
		{
			name: "string ci_equals fixed cs problems",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "ci_equals",
				Value:    "john smith",
			},
			foundNamedProps: []string{"username"},
			namedProperties: []testProperty{
				{
					name:   "username",
					value:  "John Smith",
					setErr: nil,
				},
			},
			validateResult: true,
			validateError:  nil,
		},
		{
			name: "single int condition",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
			},
			foundNamedProps: []string{"age"},
			namedProperties: []testProperty{
				{
					name:   "age",
					value:  18,
					setErr: nil,
				},
			},
			validateResult: true,
			validateError:  nil,
		},
		{
			name: "int condition, string value",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
			},
			foundNamedProps: []string{"age"},
			namedProperties: []testProperty{
				{
					name:   "age",
					value:  "john smith",
					setErr: errors.New("invalid value for float operator: john smith"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional age >= 18 without setting age"),
		},
		{
			name: "int condition, mismatched value name",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
			},
			foundNamedProps: []string{"age"},
			namedProperties: []testProperty{
				{
					name:   "age",
					value:  "john smith",
					setErr: errors.New("invalid value for float operator: john smith"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional age >= 18 without setting age"),
		},
		{
			name: "failing xor",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: xor,
					Conditionals: []ComponentConditional{
						{
							Name:     "height",
							Not:      false,
							Operator: "==",
							Value:    "180",
						},
					},
				},
			},
			foundNamedProps: []string{"age", "height"},
			namedProperties: []testProperty{
				{
					name:   "age",
					value:  18,
					setErr: nil,
				},
				{
					name:   "height",
					value:  180,
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  nil,
		},
		{
			name: "overflowing endswith",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "endswith",
				Value:    "aklsdijghyaos;idjghasldkf",
			},
			foundNamedProps: []string{"username"},
			namedProperties: []testProperty{
				{
					name:   "username",
					value:  "john smith",
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  nil,
		},
		{
			name: "invalid float condition",
			conditional: ComponentConditional{
				Name:     "age1",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: and,
					Conditionals: []ComponentConditional{
						{
							Name:     "age2",
							Not:      false,
							Operator: ">",
							Value:    "smith",
						},
					},
				},
			},
			foundNamedProps: []string{"age1", "age2"},
			namedProperties: []testProperty{
				{
					name:   "age1",
					value:  18,
					setErr: nil,
				},
				{
					name:   "age2",
					value:  20,
					setErr: errors.New("failed to convert conditional value to float: smith"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional age2 > smith without setting age2"),
		},
		{
			name: "invalid operator",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: and,
					Conditionals: []ComponentConditional{
						{
							Name:     "username",
							Not:      false,
							Operator: "is exactly",
							Value:    "john smith",
						},
					},
				},
			},
			foundNamedProps: []string{"age", "username"},
			namedProperties: []testProperty{
				{
					name:   "age",
					value:  18,
					setErr: nil,
				},
				{
					name:   "username",
					value:  "john smith",
					setErr: errors.New("invalid conditional operator is exactly"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional username is exactly john smith without setting username"),
		},
		{
			name: "unset inner conditional",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: xor,
					Conditionals: []ComponentConditional{
						{
							Name:     "username",
							Not:      false,
							Operator: "equals",
							Value:    "john smith",
						},
					},
				},
			},
			foundNamedProps: []string{"age", "username"},
			namedProperties: []testProperty{
				{
					name:   "age",
					value:  18,
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional username equals john smith without setting username"),
		},
		{
			name: "unset inner conditional",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: "some other operator",
					Conditionals: []ComponentConditional{
						{
							Name:     "username",
							Not:      false,
							Operator: "equals",
							Value:    "john smith",
						},
					},
				},
			},
			foundNamedProps: []string{"age", "username"},
			namedProperties: []testProperty{
				{
					name:   "age",
					value:  18,
					setErr: nil,
				},
				{
					name:   "username",
					value:  "john smith",
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  errors.New("invalid group operator some other operator"),
		},
		{
			name: "one of everything, all passing",
			conditional: ComponentConditional{
				Name:     "prop1",
				Not:      false,
				Operator: "ci_equals",
				Value:    "vAlUe1!",
				Group: conditionalGroup{
					Operator: and,
					Conditionals: []ComponentConditional{
						{
							Name:     "prop2",
							Not:      false,
							Operator: "equals",
							Value:    "vAlUe2!",
						},
						{
							Name:     "prop3",
							Not:      false,
							Operator: "ci_contains",
							Value:    "al",
							Group: conditionalGroup{
								Operator: or,
								Conditionals: []ComponentConditional{
									{
										Name:     "prop4",
										Not:      false,
										Operator: "contains",
										Value:    "al",
									},
									{
										Name:     "prop5",
										Not:      false,
										Operator: "ci_startswith",
										Value:    "va",
									},
									{
										Name:     "prop6",
										Not:      false,
										Operator: "startswith",
										Value:    "va",
									},
									{
										Name:     "prop7",
										Not:      false,
										Operator: "startswith",
										Value:    "vaasoidfgha;sodigkfhasldkfhjas",
									},
								},
							},
						},
						{
							Name:     "prop8",
							Not:      false,
							Operator: "ci_endswith",
							Value:    "E8!",
							Group: conditionalGroup{
								Operator: nand,
								Conditionals: []ComponentConditional{
									{
										Name:     "prop9",
										Not:      true,
										Operator: "endswith",
										Value:    "E9!",
										Group: conditionalGroup{
											Operator: nor,
											Conditionals: []ComponentConditional{
												{
													Name:     "prop10",
													Not:      false,
													Operator: "<",
													Value:    "100",
												},
											},
										},
									},
								},
							},
						},
						{
							Name:     "prop11",
							Not:      false,
							Operator: ">",
							Value:    "50",
							Group: conditionalGroup{
								Operator: xor,
								Conditionals: []ComponentConditional{
									{
										Name:     "prop12",
										Not:      false,
										Operator: "<=",
										Value:    "6",
									},
									{
										Name:     "prop13",
										Not:      false,
										Operator: ">=",
										Value:    "9",
									},
								},
							},
						},
						{
							Name:     "prop14",
							Not:      false,
							Operator: "<=",
							Value:    "6",
						},
						{
							Name:     "prop15",
							Not:      false,
							Operator: ">=",
							Value:    "9",
						},
						{
							Name:     "prop16",
							Not:      false,
							Operator: ">=",
							Value:    "9",
						},
						{
							Name:     "prop17",
							Not:      false,
							Operator: "==",
							Value:    "52",
						},
						{
							Name:     "",
							Not:      true,
							Operator: "ci_equals",
							Value:    "some non-empty string",
						},
					},
				},
			},
			foundNamedProps: []string{"prop1", "prop2", "prop3", "prop4", "prop5", "prop6", "prop7", "prop8", "prop9", "prop10", "prop11", "prop12", "prop13", "prop14", "prop15", "prop16", "prop17"},
			namedProperties: []testProperty{
				{
					name:   "prop1",
					value:  "value1!",
					setErr: nil,
				},
				{
					name:   "prop2",
					value:  "vAlUe2!",
					setErr: nil,
				},
				{
					name:   "prop3",
					value:  "vAlUe3!",
					setErr: nil,
				},
				{
					name:   "prop4",
					value:  "vAlUe4!",
					setErr: nil,
				},
				{
					name:   "prop5",
					value:  "vAlUe5!",
					setErr: nil,
				},
				{
					name:   "prop6",
					value:  "vAlUe6!",
					setErr: nil,
				},
				{
					name:   "prop7",
					value:  "vAlUe7!",
					setErr: nil,
				},
				{
					name:   "prop8",
					value:  "vAlUe8!",
					setErr: nil,
				},
				{
					name:   "prop9",
					value:  "vAlUe9!",
					setErr: nil,
				},
				{
					name:   "prop10",
					value:  10,
					setErr: nil,
				},
				{
					name:   "prop11",
					value:  10,
					setErr: nil,
				},
				{
					name:   "prop12",
					value:  4,
					setErr: nil,
				},
				{
					name:   "prop13",
					value:  2,
					setErr: nil,
				},
				{
					name:   "prop14",
					value:  6,
					setErr: nil,
				},
				{
					name:   "prop15",
					value:  9,
					setErr: nil,
				},
				{
					name:   "prop16",
					value:  100,
					setErr: nil,
				},
				{
					name:   "prop17",
					value:  52,
					setErr: nil,
				},
			},
			validateResult: true,
			validateError:  nil,
		},
	}
	for _, test = range testArray {
		t.Run(test.name, testFunc)
	}
}
