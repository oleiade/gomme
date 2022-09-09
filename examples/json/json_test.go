package json

import (
	"testing"

	"github.com/oleiade/gomme"
	"github.com/stretchr/testify/assert"
)

func TestParseJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        gomme.Parser[string, JSONValue]
		input         string
		wantErr       bool
		wantOutput    JSONValue
		wantRemaining string
	}{
		// // Null
		// {
		// 	name:    "parsing null should succeed",
		// 	parser:  Value(),
		// 	input:   `null`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind: JSONNullKind,
		// 		Null: JSONNull{},
		// 	},
		// 	wantRemaining: "",
		// },

		// // BOOLEAN

		// {
		// 	name:    "parsing true should succeed",
		// 	parser:  Value(),
		// 	input:   `true`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind: JSONBoolKind,
		// 		Bool: JSONBool(true),
		// 	},
		// 	wantRemaining: "",
		// },
		// {
		// 	name:    "parsing false should succeed",
		// 	parser:  Value(),
		// 	input:   `false`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind: JSONBoolKind,
		// 		Bool: JSONBool(false),
		// 	},
		// 	wantRemaining: "",
		// },

		// // STRINGS
		// {
		// 	name:    "parsing empty string should succeed",
		// 	parser:  Value(),
		// 	input:   `""`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind:   JSONStringKind,
		// 		String: JSONString(""),
		// 	},
		// 	wantRemaining: "",
		// },
		// {
		// 	name:    "parsing continuous string should succeed",
		// 	parser:  Value(),
		// 	input:   `"abc"`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind:   JSONStringKind,
		// 		String: JSONString("abc"),
		// 	},
		// 	wantRemaining: "",
		// },
		// // TODO: strings with whitespace characters
		// // TODO: strings with escaped characters

		// // NUMBERS
		// {
		// 	name:    "parsing positive integer number should succeed",
		// 	parser:  Value(),
		// 	input:   `123`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind:   JSONNumberKind,
		// 		Number: JSONNumber(123),
		// 	},
		// 	wantRemaining: "",
		// },
		// {
		// 	name:    "parsing negative integer number should succeed",
		// 	parser:  Value(),
		// 	input:   `-123`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind:   JSONNumberKind,
		// 		Number: JSONNumber(-123),
		// 	},
		// 	wantRemaining: "",
		// },
		// {
		// 	name:    "parsing positive floating point number should succeed",
		// 	parser:  Value(),
		// 	input:   `123.456`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind:   JSONNumberKind,
		// 		Number: JSONNumber(123.456),
		// 	},
		// 	wantRemaining: "",
		// },
		// {
		// 	name:    "parsing negative floating point number should succeed",
		// 	parser:  Value(),
		// 	input:   `-123.456`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind:   JSONNumberKind,
		// 		Number: JSONNumber(-123.456),
		// 	},
		// 	wantRemaining: "",
		// },

		// // OBJECTS
		// {
		// 	name:    "parsing empty object should succeed",
		// 	parser:  Value(),
		// 	input:   `{}`,
		// 	wantErr: false,
		// 	wantOutput: JSONValue{
		// 		Kind:   JSONObjectKind,
		// 		Object: JSONObject{},
		// 	},
		// 	wantRemaining: "",
		// },
		{
			name:    "parsing simples object should succeed",
			parser:  Value(),
			input:   `{"abc":"123"}`,
			wantErr: false,
			wantOutput: JSONValue{
				Kind: JSONObjectKind,
				Object: JSONObject(map[JSONString]JSONValue{
					JSONString("abc"): {
						Kind:   JSONNumberKind,
						String: JSONString("123"),
					},
				}),
			},
			wantRemaining: "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotResult := tc.parser(tc.input)
			if (gotResult.Err != nil) != tc.wantErr {
				t.Errorf("got error %v, want error %v", gotResult.Err, tc.wantErr)
			}

			assert.Equal(t,
				tc.wantOutput,
				gotResult.Output,
				"got output %v, want output %v", gotResult.Output, tc.wantOutput,
			)

			if gotResult.Remaining != tc.wantRemaining {
				t.Errorf("got remaining %v, want remaining %v", gotResult.Remaining, tc.wantRemaining)
			}
		})
	}
}

func TestNull(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        gomme.Parser[string, JSONNull]
		input         string
		wantErr       bool
		wantOutput    JSONNull
		wantRemaining string
	}{
		{
			name:          "parsing null should succeed",
			parser:        Null(),
			input:         "null",
			wantErr:       false,
			wantOutput:    JSONNull(struct{}{}),
			wantRemaining: "",
		},
		{
			name:          "parsing non-matching should fail",
			parser:        Null(),
			input:         "abc",
			wantErr:       true,
			wantOutput:    JSONNull{},
			wantRemaining: "abc",
		},
		{
			name:          "parsing empty input should fail",
			parser:        Null(),
			input:         "",
			wantErr:       true,
			wantOutput:    JSONNull{},
			wantRemaining: "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotResult := tc.parser(tc.input)
			if (gotResult.Err != nil) != tc.wantErr {
				t.Errorf("got error %v, want error %v", gotResult.Err, tc.wantErr)
			}

			if gotResult.Output != tc.wantOutput {
				t.Errorf("got output %v, want output %v", gotResult.Output, tc.wantOutput)
			}

			if gotResult.Remaining != tc.wantRemaining {
				t.Errorf("got remaining %v, want remaining %v", gotResult.Remaining, tc.wantRemaining)
			}
		})
	}
}

func TestBoolean(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        gomme.Parser[string, JSONBool]
		input         string
		wantErr       bool
		wantOutput    JSONBool
		wantRemaining string
	}{
		{
			name:          "parsing true should succeed",
			parser:        Boolean(),
			input:         "true",
			wantErr:       false,
			wantOutput:    JSONBool(true),
			wantRemaining: "",
		},
		{
			name:          "parsing false should succeed",
			parser:        Boolean(),
			input:         "false",
			wantErr:       false,
			wantOutput:    JSONBool(false),
			wantRemaining: "",
		},
		{
			name:          "parsing invalid input should fail",
			parser:        Boolean(),
			input:         "invalid",
			wantErr:       true,
			wantOutput:    JSONBool(false),
			wantRemaining: "invalid",
		},
		{
			name:          "parsing empty input should fail",
			parser:        Boolean(),
			input:         "",
			wantErr:       true,
			wantOutput:    JSONBool(false),
			wantRemaining: "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotResult := tc.parser(tc.input)
			if (gotResult.Err != nil) != tc.wantErr {
				t.Errorf("got error %v, want error %v", gotResult.Err, tc.wantErr)
			}

			if gotResult.Output != tc.wantOutput {
				t.Errorf("got output %v, want output %v", gotResult.Output, tc.wantOutput)
			}

			if gotResult.Remaining != tc.wantRemaining {
				t.Errorf("got remaining %v, want remaining %v", gotResult.Remaining, tc.wantRemaining)
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        gomme.Parser[string, JSONString]
		input         string
		wantErr       bool
		wantOutput    JSONString
		wantRemaining string
	}{
		{
			name:          "parsing quoted alpha chars string should succeed",
			parser:        String(),
			input:         "\"bonjour\"",
			wantErr:       false,
			wantOutput:    JSONString("bonjour"),
			wantRemaining: "",
		},
		{
			name:          "parsing quoted empty string should succeed",
			parser:        String(),
			input:         "\"\"",
			wantErr:       false,
			wantOutput:    JSONString(""),
			wantRemaining: "",
		},
		{
			name:          "parsing unopened quotes string should fail",
			parser:        String(),
			input:         "bonjour\"",
			wantErr:       true,
			wantOutput:    JSONString(""),
			wantRemaining: "bonjour\"",
		},
		{
			name:          "parsing unclosed quotes string should fail",
			parser:        String(),
			input:         "\"bonjour",
			wantErr:       true,
			wantOutput:    JSONString(""),
			wantRemaining: "\"bonjour",
		},
		{
			name:          "parsing unquoted string should fail",
			parser:        String(),
			input:         "\"bonjour",
			wantErr:       true,
			wantOutput:    JSONString(""),
			wantRemaining: "\"bonjour",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotResult := tc.parser(tc.input)
			if (gotResult.Err != nil) != tc.wantErr {
				t.Errorf("got error %v, want error %v", gotResult.Err, tc.wantErr)
			}

			if gotResult.Output != tc.wantOutput {
				t.Errorf("got output %v, want output %v", gotResult.Output, tc.wantOutput)
			}

			if gotResult.Remaining != tc.wantRemaining {
				t.Errorf("got remaining %v, want remaining %v", gotResult.Remaining, tc.wantRemaining)
			}
		})
	}
}
