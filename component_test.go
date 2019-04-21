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
	type testGroup struct {
		Operator     groupOperator
		Conditionals []ComponentConditional
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
					setErr: errors.New("Invalid value for string operator: 18"),
				},
			},
			validateResult: false,
			validateError:  errors.New("Attempted to validate conditional username equals john smith without setting username"),
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
			validateError:  errors.New("Attempted to validate conditional username equals john smith without setting username"),
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
					setErr: errors.New("Invalid value for integer operator: john smith"),
				},
			},
			validateResult: false,
			validateError:  errors.New("Attempted to validate conditional age >= 18 without setting age"),
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
					name:   "Age",
					value:  18,
					setErr: nil,
				},
			},
			validateResult: false,
			validateError:  errors.New("Attempted to validate conditional age >= 18 without setting age"),
		},
		testSet{
			name: "one of everything, all passing",
			conditional: ComponentConditional{
				Name:     "prop1",
				Not:      false,
				Operator: "ci_equals",
				Value:    "vAlUe1!",
				Group: struct {
					Operator     groupOperator          `json:"groupOperator"`
					Conditionals []ComponentConditional `json:"conditionals"`
				}(testGroup{
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
							Group: struct {
								Operator     groupOperator          `json:"groupOperator"`
								Conditionals []ComponentConditional `json:"conditionals"`
							}(testGroup{
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
							}),
						},
						ComponentConditional{
							Name:     "prop8",
							Not:      false,
							Operator: "ci_endswith",
							Value:    "E8!",
							Group: struct {
								Operator     groupOperator          `json:"groupOperator"`
								Conditionals []ComponentConditional `json:"conditionals"`
							}(testGroup{
								Operator: nand,
								Conditionals: []ComponentConditional{
									ComponentConditional{
										Name:     "prop9",
										Not:      true,
										Operator: "endswith",
										Value:    "E9!",
										Group: struct {
											Operator     groupOperator          `json:"groupOperator"`
											Conditionals []ComponentConditional `json:"conditionals"`
										}(testGroup{
											Operator: nor,
											Conditionals: []ComponentConditional{
												ComponentConditional{
													Name:     "prop10",
													Not:      false,
													Operator: "<",
													Value:    "100",
												},
											},
										}),
									},
								},
							}),
						},
						ComponentConditional{
							Name:     "prop11",
							Not:      false,
							Operator: ">",
							Value:    "50",
							Group: struct {
								Operator     groupOperator          `json:"groupOperator"`
								Conditionals []ComponentConditional `json:"conditionals"`
							}(testGroup{
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
							}),
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
							Name:     "prop15",
							Not:      false,
							Operator: ">=",
							Value:    "9",
						},
					},
				}),
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
			},
			validateResult: true,
			validateError:  nil,
		},
	}
	for _, test := range testArray {
		t.Run(test.name, func(t *testing.T) {
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
		})
	}
}
