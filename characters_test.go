package gomme

import (
	"testing"
)

func TestChar(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, rune]
		input         string
		wantErr       bool
		wantOutput    rune
		wantRemaining string
	}{
		{
			name:          "parsing char from single char input should succeed",
			parser:        Char('a'),
			input:         "a",
			wantErr:       false,
			wantOutput:    'a',
			wantRemaining: "",
		},
		{
			name:          "parsing valid char in longer input should succeed",
			parser:        Char('a'),
			input:         "abc",
			wantErr:       false,
			wantOutput:    'a',
			wantRemaining: "bc",
		},
		{
			name:          "parsing single non-char input should fail",
			parser:        Char('a'),
			input:         "123",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "123",
		},
		{
			name:          "parsing empty input should fail",
			parser:        Char('a'),
			input:         "",
			wantErr:       true,
			wantOutput:    rune(0),
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

func TestAnyChar(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, rune]
		input         string
		wantErr       bool
		wantOutput    rune
		wantRemaining string
	}{
		{
			name:          "parsing any char from single entry input should succeed",
			parser:        AnyChar(),
			input:         "a",
			wantErr:       false,
			wantOutput:    'a',
			wantRemaining: "",
		},
		{
			name:          "parsing valid any char from longer input should succeed",
			parser:        AnyChar(),
			input:         "abc",
			wantErr:       false,
			wantOutput:    'a',
			wantRemaining: "bc",
		},
		{
			name:          "parsing any char from empty input should fail",
			parser:        AnyChar(),
			input:         "",
			wantErr:       true,
			wantOutput:    rune(0),
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

func TestAlpha0(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing single alpha char from single alpha input should succeed",
			parser:        Alpha0(),
			input:         "a",
			wantErr:       false,
			wantOutput:    "a",
			wantRemaining: "",
		},
		{
			name:          "parsing alpha chars from multiple alpha input should succeed",
			parser:        Alpha0(),
			input:         "abc",
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "",
		},
		{
			name:          "parsing alpha chars until terminating char should succeed",
			parser:        Alpha0(),
			input:         "abc123",
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "123",
		},
		{
			name:          "parsing an empty input should succeed",
			parser:        Alpha0(),
			input:         "",
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:          "parsing non alpha chars should succeed",
			parser:        Alpha0(),
			input:         "123",
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "123",
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

func TestAlpha1(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing single alpha char from single alpha input should succeed",
			parser:        Alpha1(),
			input:         "a",
			wantErr:       false,
			wantOutput:    "a",
			wantRemaining: "",
		},
		{
			name:          "parsing alpha chars from multiple alpha input should succeed",
			parser:        Alpha1(),
			input:         "abc",
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "",
		},
		{
			name:          "parsing alpha chars until terminating char should succeed",
			parser:        Alpha1(),
			input:         "abc123",
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "123",
		},
		{
			name:          "parsing an empty input should fail",
			parser:        Alpha1(),
			input:         "",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:          "parsing input not starting with an alpha char should fail",
			parser:        Alpha1(),
			input:         "1c",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "1c",
		},
		{
			name:          "parsing non alpha chars should fail",
			parser:        Alpha1(),
			input:         "123",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "123",
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

func TestDigit0(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing single digit char from single digit input should succeed",
			parser:        Digit0(),
			input:         "1",
			wantErr:       false,
			wantOutput:    "1",
			wantRemaining: "",
		},
		{
			name:          "parsing digit chars from multiple digit input should succeed",
			parser:        Digit0(),
			input:         "123",
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "",
		},
		{
			name:          "parsing digit chars until terminating char should succeed",
			parser:        Digit0(),
			input:         "123abc",
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "abc",
		},
		{
			name:          "parsing an empty input should succeed",
			parser:        Digit0(),
			input:         "",
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:          "parsing non digit chars should succeed",
			parser:        Digit0(),
			input:         "abc",
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "abc",
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

func TestDigit1(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing single digit char from single digit input should succeed",
			parser:        Digit1(),
			input:         "1",
			wantErr:       false,
			wantOutput:    "1",
			wantRemaining: "",
		},
		{
			name:          "parsing digit chars from multiple digit input should succeed",
			parser:        Digit1(),
			input:         "123",
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "",
		},
		{
			name:          "parsing digit chars until terminating char should succeed",
			parser:        Digit1(),
			input:         "123abc",
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "abc",
		},
		{
			name:          "parsing an empty input should fail",
			parser:        Digit1(),
			input:         "",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:          "parsing input not starting with an digit char should fail",
			parser:        Digit1(),
			input:         "c1",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "c1",
		},
		{
			name:          "parsing non digit chars should fail",
			parser:        Digit1(),
			input:         "abc",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "abc",
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

func TestHexDigit0(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing single hex digit char from single hex digit input should succeed",
			parser:        HexDigit0(),
			input:         "1",
			wantErr:       false,
			wantOutput:    "1",
			wantRemaining: "",
		},
		{
			name:          "parsing hex digit chars from multiple hex digit input should succeed",
			parser:        HexDigit0(),
			input:         "1f3",
			wantErr:       false,
			wantOutput:    "1f3",
			wantRemaining: "",
		},
		{
			name:          "parsing hex digit chars until terminating char should succeed",
			parser:        HexDigit0(),
			input:         "1f3z",
			wantErr:       false,
			wantOutput:    "1f3",
			wantRemaining: "z",
		},
		{
			name:          "parsing an empty input should succeed",
			parser:        HexDigit0(),
			input:         "",
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:          "parsing non hex digit chars should succeed",
			parser:        HexDigit0(),
			input:         "ghi",
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "ghi",
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

func TestHexDigit1(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing single hex digit char from single hex digit input should succeed",
			parser:        HexDigit1(),
			input:         "1",
			wantErr:       false,
			wantOutput:    "1",
			wantRemaining: "",
		},
		{
			name:          "parsing hex digit chars from multiple hex digit input should succeed",
			parser:        HexDigit1(),
			input:         "1f3",
			wantErr:       false,
			wantOutput:    "1f3",
			wantRemaining: "",
		},
		{
			name:          "parsing hex digit chars until terminating char should succeed",
			parser:        HexDigit1(),
			input:         "1f3ghi",
			wantErr:       false,
			wantOutput:    "1f3",
			wantRemaining: "ghi",
		},
		{
			name:          "parsing an empty input should fail",
			parser:        HexDigit1(),
			input:         "",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:          "parsing input not starting with a hex digit char should fail",
			parser:        HexDigit1(),
			input:         "h1",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "h1",
		},
		{
			name:          "parsing non hex digit chars should fail",
			parser:        HexDigit1(),
			input:         "ghi",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "ghi",
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

func TestAlphanumeric0(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing single alpha char from single alphanumerical input should succeed",
			parser:        Alphanumeric0(),
			input:         "a",
			wantErr:       false,
			wantOutput:    "a",
			wantRemaining: "",
		},
		{
			name:          "parsing single digit char from single alphanumerical input should succeed",
			parser:        Alphanumeric0(),
			input:         "1",
			wantErr:       false,
			wantOutput:    "1",
			wantRemaining: "",
		},
		{
			name:          "parsing alpha chars from multiple alphanumerical input should succeed",
			parser:        Alphanumeric0(),
			input:         "abc",
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "",
		},
		{
			name:          "parsing digit chars from multiple alphanumerical input should succeed",
			parser:        Alphanumeric0(),
			input:         "123",
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "",
		},
		{
			name:          "parsing multiple alphanumerical input should succeed",
			parser:        Alphanumeric0(),
			input:         "a1b2c3",
			wantErr:       false,
			wantOutput:    "a1b2c3",
			wantRemaining: "",
		},
		{
			name:          "parsing alph chars until terminating char should succeed",
			parser:        Alphanumeric0(),
			input:         "abc$%^",
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "$%^",
		},
		{
			name:          "parsing digit chars until terminating char should succeed",
			parser:        Alphanumeric0(),
			input:         "123$%^",
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "$%^",
		},
		{
			name:          "parsing alphanumerical chars until terminating char should succeed",
			parser:        Alphanumeric0(),
			input:         "a1b2c3$%^",
			wantErr:       false,
			wantOutput:    "a1b2c3",
			wantRemaining: "$%^",
		},
		{
			name:          "parsing an empty input should succeed",
			parser:        Alphanumeric0(),
			input:         "",
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:          "parsing non alphanumerical chars should succeed",
			parser:        Alphanumeric0(),
			input:         "$%^",
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "$%^",
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

func TestAlphanumeric1(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing single alpha char from single alphanumerical input should succeed",
			parser:        Alphanumeric1(),
			input:         "a",
			wantErr:       false,
			wantOutput:    "a",
			wantRemaining: "",
		},
		{
			name:          "parsing single digit char from single alphanumerical input should succeed",
			parser:        Alphanumeric1(),
			input:         "1",
			wantErr:       false,
			wantOutput:    "1",
			wantRemaining: "",
		},
		{
			name:          "parsing alpha chars from multiple alphanumerical input should succeed",
			parser:        Alphanumeric1(),
			input:         "abc",
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "",
		},
		{
			name:          "parsing digit chars from multiple alphanumerical input should succeed",
			parser:        Alphanumeric1(),
			input:         "123",
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "",
		},
		{
			name:          "parsing alphanumerical chars from multiple alphanumerical input should succeed",
			parser:        Alphanumeric1(),
			input:         "a1b2c3",
			wantErr:       false,
			wantOutput:    "a1b2c3",
			wantRemaining: "",
		},
		{
			name:          "parsing alpha chars until terminating char should succeed",
			parser:        Alphanumeric1(),
			input:         "abc$%^",
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "$%^",
		},
		{
			name:          "parsing digit chars until terminating char should succeed",
			parser:        Alphanumeric1(),
			input:         "123$%^",
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "$%^",
		},
		{
			name:          "parsing alphanumerical chars until terminating char should succeed",
			parser:        Alphanumeric1(),
			input:         "a1b2c3$%^",
			wantErr:       false,
			wantOutput:    "a1b2c3",
			wantRemaining: "$%^",
		},
		{
			name:          "parsing an empty input should fail",
			parser:        Alphanumeric1(),
			input:         "",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:          "parsing input not starting with an alphanumeric char should fail",
			parser:        Alphanumeric1(),
			input:         "$1",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "$1",
		},
		{
			name:          "parsing non digit chars should fail",
			parser:        Alphanumeric1(),
			input:         "$%^",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "$%^",
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

func TestLF(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, rune]
		input         string
		wantErr       bool
		wantOutput    rune
		wantRemaining string
	}{
		{
			name:          "parsing single line-feed char from single line-feed input should succeed",
			parser:        LF(),
			input:         "\n",
			wantErr:       false,
			wantOutput:    rune('\n'),
			wantRemaining: "",
		},
		{
			name:          "parsing single line-feed char from multiple char input should succeed",
			parser:        LF(),
			input:         "\nabc",
			wantErr:       false,
			wantOutput:    rune('\n'),
			wantRemaining: "abc",
		},
		{
			name:          "parsing empty input should fail",
			parser:        LF(),
			input:         "",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "",
		},
		{
			name:          "parsing single line-feed char from single non-line-feed input should fail",
			parser:        LF(),
			input:         "1",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "1",
		},
		{
			name:          "parsing single line-feed from multiple non-line-feed input should fail",
			parser:        LF(),
			input:         "123",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "123",
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

func TestCR(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, rune]
		input         string
		wantErr       bool
		wantOutput    rune
		wantRemaining string
	}{
		{
			name:          "parsing single carriage-return char from single carriage-return input should succeed",
			parser:        CR(),
			input:         "\r",
			wantErr:       false,
			wantOutput:    rune('\r'),
			wantRemaining: "",
		},
		{
			name:          "parsing single carriage-return char from multiple char input should succeed",
			parser:        CR(),
			input:         "\rabc",
			wantErr:       false,
			wantOutput:    rune('\r'),
			wantRemaining: "abc",
		},
		{
			name:          "parsing empty input should fail",
			parser:        CR(),
			input:         "",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "",
		},
		{
			name:          "parsing single carriage-return char from single non-carriage-return input should fail",
			parser:        CR(),
			input:         "1",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "1",
		},
		{
			name:          "parsing single carriage-return from multiple non-carriage-return input should fail",
			parser:        CR(),
			input:         "123",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "123",
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

func TestCRLF(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing single CRLF char from single CRLF input should succeed",
			parser:        CRLF(),
			input:         "\r\n",
			wantErr:       false,
			wantOutput:    "\r\n",
			wantRemaining: "",
		},
		{
			name:          "parsing single CRLF char from multiple char input should succeed",
			parser:        CRLF(),
			input:         "\r\nabc",
			wantErr:       false,
			wantOutput:    "\r\n",
			wantRemaining: "abc",
		},
		{
			name:          "parsing empty input should fail",
			parser:        CRLF(),
			input:         "",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:          "parsing incomplete CRLF input should fail",
			parser:        CRLF(),
			input:         "\r",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "\r",
		},
		{
			name:          "parsing single CRLF char from single non-CRLF input should fail",
			parser:        CRLF(),
			input:         "1",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "1",
		},
		{
			name:          "parsing single CRLF from multiple non-CRLF input should fail",
			parser:        CRLF(),
			input:         "123",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "123",
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

func TestSatisfy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, rune]
		input         string
		wantErr       bool
		wantOutput    rune
		wantRemaining string
	}{
		{
			name:          "parsing single alpha char satisfying constraint should succeed",
			parser:        Satisfy(IsAlpha),
			input:         "a",
			wantErr:       false,
			wantOutput:    'a',
			wantRemaining: "",
		},
		{
			name:          "parsing alpha char satisfying constraint from mixed input should succeed",
			parser:        Satisfy(IsAlpha),
			input:         "a1",
			wantErr:       false,
			wantOutput:    'a',
			wantRemaining: "1",
		},
		{
			name:          "parsing char not satisfying constraint should succeed",
			parser:        Satisfy(IsAlpha),
			input:         "1",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "1",
		},
		{
			name:          "parsing empty input should succeed",
			parser:        Satisfy(IsAlpha),
			input:         "",
			wantErr:       true,
			wantOutput:    rune(0),
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

func TestSpace(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, rune]
		input         string
		wantErr       bool
		wantOutput    rune
		wantRemaining string
	}{
		{
			name:          "parsing single space char from single space input should succeed",
			parser:        Space(),
			input:         " ",
			wantErr:       false,
			wantOutput:    rune(' '),
			wantRemaining: "",
		},
		{
			name:          "parsing single space char from multiple char input should succeed",
			parser:        Space(),
			input:         " abc",
			wantErr:       false,
			wantOutput:    rune(' '),
			wantRemaining: "abc",
		},
		{
			name:          "parsing empty input should fail",
			parser:        Space(),
			input:         "",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "",
		},
		{
			name:          "parsing single space char from single non-space input should fail",
			parser:        Space(),
			input:         "1",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "1",
		},
		{
			name:          "parsing single space from multiple non-space input should fail",
			parser:        Space(),
			input:         "123",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "123",
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

func TestTab(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, rune]
		input         string
		wantErr       bool
		wantOutput    rune
		wantRemaining string
	}{
		{
			name:          "parsing single space char from single space input should succeed",
			parser:        Tab(),
			input:         "\t",
			wantErr:       false,
			wantOutput:    rune('\t'),
			wantRemaining: "",
		},
		{
			name:          "parsing single space char from multiple char input should succeed",
			parser:        Tab(),
			input:         "\tabc",
			wantErr:       false,
			wantOutput:    rune('\t'),
			wantRemaining: "abc",
		},
		{
			name:          "parsing empty input should fail",
			parser:        Tab(),
			input:         "",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "",
		},
		{
			name:          "parsing single space char from single non-space input should fail",
			parser:        Tab(),
			input:         "1",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "1",
		},
		{
			name:          "parsing single space from multiple non-space input should fail",
			parser:        Tab(),
			input:         "123",
			wantErr:       true,
			wantOutput:    rune(0),
			wantRemaining: "123",
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

func TestToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, string]
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:          "parsing a token from an input starting with it should succeed",
			parser:        Token("Bonjour"),
			input:         "Bonjour tout le monde",
			wantErr:       false,
			wantOutput:    "Bonjour",
			wantRemaining: " tout le monde",
		},
		{
			name:          "parsing a token from an non-matching input should fail",
			parser:        Token("Bonjour"),
			input:         "Hello tout le monde",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "Hello tout le monde",
		},
		{
			name:          "parsing a token from an empty input should fail",
			parser:        Token("Bonjour"),
			input:         "",
			wantErr:       true,
			wantOutput:    "",
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

func TestInt64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, int64]
		input         string
		wantErr       bool
		wantOutput    int64
		wantRemaining string
	}{
		{
			name:          "parsing positive integer should succeed",
			parser:        Int64(),
			input:         "123",
			wantErr:       false,
			wantOutput:    123,
			wantRemaining: "",
		},
		{
			name:          "parsing negative integer should succeed",
			parser:        Int64(),
			input:         "-123",
			wantErr:       false,
			wantOutput:    -123,
			wantRemaining: "",
		},
		{
			name:          "parsing positive integer prefix should succeed",
			parser:        Int64(),
			input:         "123abc",
			wantErr:       false,
			wantOutput:    123,
			wantRemaining: "abc",
		},
		{
			name:          "parsing negative integer should succeed",
			parser:        Int64(),
			input:         "-123abc",
			wantErr:       false,
			wantOutput:    -123,
			wantRemaining: "abc",
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

func TestInt8(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, int8]
		input         string
		wantErr       bool
		wantOutput    int8
		wantRemaining string
	}{
		{
			name:          "parsing positive integer should succeed",
			parser:        Int8(),
			input:         "123",
			wantErr:       false,
			wantOutput:    123,
			wantRemaining: "",
		},
		{
			name:          "parsing negative integer should succeed",
			parser:        Int8(),
			input:         "-123",
			wantErr:       false,
			wantOutput:    -123,
			wantRemaining: "",
		},
		{
			name:          "parsing positive integer prefix should succeed",
			parser:        Int8(),
			input:         "123abc",
			wantErr:       false,
			wantOutput:    123,
			wantRemaining: "abc",
		},
		{
			name:          "parsing negative integer should succeed",
			parser:        Int8(),
			input:         "-123abc",
			wantErr:       false,
			wantOutput:    -123,
			wantRemaining: "abc",
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

