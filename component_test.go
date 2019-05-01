package imagetemplate

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		setPropTest{
			name: "single property success",
			properties: NamedProperties{
				"username": "john smith",
			},
			propMap: map[string][]string{
				"username": []string{"innerPropUsername"},
			},
			setFunc:         successSetFunc,
			resultLeftovers: map[string][]string{},
			resultErr:       nil,
		},
		setPropTest{
			name: "single property success",
			properties: NamedProperties{
				"username": "john smith",
			},
			propMap: map[string][]string{
				"username": []string{"innerPropUsername"},
			},
			setFunc: errorSetFunc,
			resultLeftovers: map[string][]string{
				"username": []string{"innerPropUsername"},
			},
			resultErr: errors.New("failed to set property"),
		},
		setPropTest{
			name: "many properties success",
			properties: NamedProperties{
				"username": "john smith",
				"age":      57,
				"title":    "Mr.",
			},
			propMap: map[string][]string{
				"username": []string{"innerPropUsername", "innerPropEmail"},
				"age":      []string{"innerPropAge"},
				"title":    []string{"innerPropTitle", "TITLE", "innerPropRank", "some random field"},
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
			assert.Equal(t, test.resultErr, err)
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
				"username": []string{"innerPropUsername", "innerPropEmail"},
				"age":      []string{"innerPropAge"},
				"title":    []string{"innerPropTitle", "TITLE", "innerPropRank", "some random field"},
			},
			setFunc:         testSetFunc,
			resultLeftovers: map[string][]string{},
			resultErr:       nil,
		}
		leftovers, err := StandardSetNamedProperties(test.properties, test.propMap, test.setFunc)
		assert.Equal(t, test.resultLeftovers, leftovers)
		assert.Equal(t, test.resultErr, err)
		assert.Equal(t, len(expectedValues), len(valuesSet), "valuesSet length changes, something was added or deleted improperly")
		finalResult := true
		for _, result := range valuesSet {
			finalResult = finalResult && result
		}
		assert.True(t, finalResult)
	})
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
		parseTest{
			name:               "simple value",
			value:              "some data",
			hasNamedProperties: false,
			cleanValues:        []string{"some data"},
			propNames:          []string(nil),
			err:                nil,
		},
		parseTest{
			name:               "simple value 2",
			value:              "123",
			hasNamedProperties: false,
			cleanValues:        []string{"123"},
			propNames:          []string(nil),
			err:                nil,
		},
		parseTest{
			name:               "simple value 3",
			value:              "#4: This. Is. Just. Data!~",
			hasNamedProperties: false,
			cleanValues:        []string{"#4: This. Is. Just. Data!~"},
			propNames:          []string(nil),
			err:                nil,
		},
		parseTest{
			name:               "escaped dollars",
			value:              "some $$ data 2 be pr$$ocessed!",
			hasNamedProperties: false,
			cleanValues:        []string{"some $ data 2 be pr$ocessed!"},
			propNames:          []string(nil),
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
			err:                errors.New("unclosed named property in 'Hello there, $title$. $username!'"),
		},
		parseTest{
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
	var test testSet
	testFunc := func(t *testing.T) {
		for _, prop := range test.namedProperties {
			t.Run(prop.name, func(t *testing.T) {
				var err error
				test.conditional, err = test.conditional.SetValue(prop.name, prop.value)
				assert.Equal(t, prop.setErr, err)
			})
		}
		success, err := test.conditional.Validate()
		assert.Equal(t, test.validateResult, success)
		assert.Equal(t, test.validateError, err)
	}

	testArray := []testSet{
		testSet{
			name: "single string condition",
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
		testSet{
			name: "string condition, int set value",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "equals",
				Value:    "john smith",
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "username",
					value:  18,
					setErr: errors.New("invalid value for string operator: 18"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional username equals john smith without setting username"),
		},
		testSet{
			name: "string condition, mismatched value name",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "equals",
				Value:    "john smith",
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "UserName",
					value:  "john smith",
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional username equals john smith without setting username"),
		},
		testSet{
			name: "string equals, fails on case sensitivity",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "equals",
				Value:    "john smith",
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "username",
					value:  "John Smith",
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  nil,
		},
		testSet{
			name: "string ci_equals fixed cs problems",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "ci_equals",
				Value:    "john smith",
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "username",
					value:  "John Smith",
					setErr: nil,
				},
			},
			validateResult: true,
			validateError:  nil,
		},
		testSet{
			name: "single int condition",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "age",
					value:  18,
					setErr: nil,
				},
			},
			validateResult: true,
			validateError:  nil,
		},
		testSet{
			name: "int condition, string value",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "age",
					value:  "john smith",
					setErr: errors.New("invalid value for float operator: john smith"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional age >= 18 without setting age"),
		},
		testSet{
			name: "int condition, mismatched value name",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "age",
					value:  "john smith",
					setErr: errors.New("invalid value for float operator: john smith"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional age >= 18 without setting age"),
		},
		testSet{
			name: "failing xor",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: xor,
					Conditionals: []ComponentConditional{
						ComponentConditional{
							Name:     "height",
							Not:      false,
							Operator: "==",
							Value:    "180",
						},
					},
				},
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "age",
					value:  18,
					setErr: nil,
				},
				testProperty{
					name:   "height",
					value:  180,
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  nil,
		},
		testSet{
			name: "overflowing endswith",
			conditional: ComponentConditional{
				Name:     "username",
				Not:      false,
				Operator: "endswith",
				Value:    "aklsdijghyaos;idjghasldkf",
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "username",
					value:  "john smith",
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  nil,
		},
		testSet{
			name: "invalid float condition",
			conditional: ComponentConditional{
				Name:     "age1",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: and,
					Conditionals: []ComponentConditional{
						ComponentConditional{
							Name:     "age2",
							Not:      false,
							Operator: ">",
							Value:    "smith",
						},
					},
				},
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "age1",
					value:  18,
					setErr: nil,
				},
				testProperty{
					name:   "age2",
					value:  20,
					setErr: errors.New("failed to convert conditional value to float: smith"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional age2 > smith without setting age2"),
		},
		testSet{
			name: "invalid operator",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: and,
					Conditionals: []ComponentConditional{
						ComponentConditional{
							Name:     "username",
							Not:      false,
							Operator: "is exactly",
							Value:    "john smith",
						},
					},
				},
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "age",
					value:  18,
					setErr: nil,
				},
				testProperty{
					name:   "username",
					value:  "john smith",
					setErr: errors.New("invalid conditional operator is exactly"),
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional username is exactly john smith without setting username"),
		},
		testSet{
			name: "unset inner conditional",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: xor,
					Conditionals: []ComponentConditional{
						ComponentConditional{
							Name:     "username",
							Not:      false,
							Operator: "equals",
							Value:    "john smith",
						},
					},
				},
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "age",
					value:  18,
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  errors.New("attempted to validate conditional username equals john smith without setting username"),
		},
		testSet{
			name: "unset inner conditional",
			conditional: ComponentConditional{
				Name:     "age",
				Not:      false,
				Operator: ">=",
				Value:    "18",
				Group: conditionalGroup{
					Operator: "some other operator",
					Conditionals: []ComponentConditional{
						ComponentConditional{
							Name:     "username",
							Not:      false,
							Operator: "equals",
							Value:    "john smith",
						},
					},
				},
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "age",
					value:  18,
					setErr: nil,
				},
				testProperty{
					name:   "username",
					value:  "john smith",
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  errors.New("invalid group operator some other operator"),
		},
		testSet{
			name: "one of everything, all passing",
			conditional: ComponentConditional{
				Name:     "prop1",
				Not:      false,
				Operator: "ci_equals",
				Value:    "vAlUe1!",
				Group: conditionalGroup{
					Operator: and,
					Conditionals: []ComponentConditional{
						ComponentConditional{
							Name:     "prop2",
							Not:      false,
							Operator: "equals",
							Value:    "vAlUe2!",
						},
						ComponentConditional{
							Name:     "prop3",
							Not:      false,
							Operator: "ci_contains",
							Value:    "al",
							Group: conditionalGroup{
								Operator: or,
								Conditionals: []ComponentConditional{
									ComponentConditional{
										Name:     "prop4",
										Not:      false,
										Operator: "contains",
										Value:    "al",
									},
									ComponentConditional{
										Name:     "prop5",
										Not:      false,
										Operator: "ci_startswith",
										Value:    "va",
									},
									ComponentConditional{
										Name:     "prop6",
										Not:      false,
										Operator: "startswith",
										Value:    "va",
									},
									ComponentConditional{
										Name:     "prop7",
										Not:      false,
										Operator: "startswith",
										Value:    "vaasoidfgha;sodigkfhasldkfhjas",
									},
								},
							},
						},
						ComponentConditional{
							Name:     "prop8",
							Not:      false,
							Operator: "ci_endswith",
							Value:    "E8!",
							Group: conditionalGroup{
								Operator: nand,
								Conditionals: []ComponentConditional{
									ComponentConditional{
										Name:     "prop9",
										Not:      true,
										Operator: "endswith",
										Value:    "E9!",
										Group: conditionalGroup{
											Operator: nor,
											Conditionals: []ComponentConditional{
												ComponentConditional{
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
						ComponentConditional{
							Name:     "prop11",
							Not:      false,
							Operator: ">",
							Value:    "50",
							Group: conditionalGroup{
								Operator: xor,
								Conditionals: []ComponentConditional{
									ComponentConditional{
										Name:     "prop12",
										Not:      false,
										Operator: "<=",
										Value:    "6",
									},
									ComponentConditional{
										Name:     "prop13",
										Not:      false,
										Operator: ">=",
										Value:    "9",
									},
								},
							},
						},
						ComponentConditional{
							Name:     "prop14",
							Not:      false,
							Operator: "<=",
							Value:    "6",
						},
						ComponentConditional{
							Name:     "prop15",
							Not:      false,
							Operator: ">=",
							Value:    "9",
						},
						ComponentConditional{
							Name:     "prop16",
							Not:      false,
							Operator: ">=",
							Value:    "9",
						},
						ComponentConditional{
							Name:     "prop17",
							Not:      false,
							Operator: "==",
							Value:    "52",
						},
					},
				},
			},
			namedProperties: []testProperty{
				testProperty{
					name:   "prop1",
					value:  "value1!",
					setErr: nil,
				},
				testProperty{
					name:   "prop2",
					value:  "vAlUe2!",
					setErr: nil,
				},
				testProperty{
					name:   "prop3",
					value:  "vAlUe3!",
					setErr: nil,
				},
				testProperty{
					name:   "prop4",
					value:  "vAlUe4!",
					setErr: nil,
				},
				testProperty{
					name:   "prop5",
					value:  "vAlUe5!",
					setErr: nil,
				},
				testProperty{
					name:   "prop6",
					value:  "vAlUe6!",
					setErr: nil,
				},
				testProperty{
					name:   "prop7",
					value:  "vAlUe7!",
					setErr: nil,
				},
				testProperty{
					name:   "prop8",
					value:  "vAlUe8!",
					setErr: nil,
				},
				testProperty{
					name:   "prop9",
					value:  "vAlUe9!",
					setErr: nil,
				},
				testProperty{
					name:   "prop10",
					value:  10,
					setErr: nil,
				},
				testProperty{
					name:   "prop11",
					value:  10,
					setErr: nil,
				},
				testProperty{
					name:   "prop12",
					value:  4,
					setErr: nil,
				},
				testProperty{
					name:   "prop13",
					value:  2,
					setErr: nil,
				},
				testProperty{
					name:   "prop14",
					value:  6,
					setErr: nil,
				},
				testProperty{
					name:   "prop15",
					value:  9,
					setErr: nil,
				},
				testProperty{
					name:   "prop16",
					value:  100,
					setErr: nil,
				},
				testProperty{
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
