package gomme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelimited(t *testing.T) {
	t.Parallel()

	type args struct {
		p Parser[string, string]
	}
	testCases := []struct {
		name          string
		args          args
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:  "matching parser should succeed",
			input: "+1\r\n",
			args: args{
				p: Delimited(Char('+'), Digit1(), CRLF()),
			},
			wantErr:       false,
			wantOutput:    "1",
			wantRemaining: "",
		},
		{
			name:  "no prefix match should fail",
			input: "1\r\n",
			args: args{
				p: Delimited(Char('+'), Digit1(), CRLF()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "1\r\n",
		},
		{
			name:  "no parser match should fail",
			input: "+\r\n",
			args: args{
				p: Delimited(Char('+'), Digit1(), CRLF()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "+\r\n",
		},
		{
			name:  "no suffix match should fail",
			input: "+1",
			args: args{
				p: Delimited(Char('+'), Digit1(), CRLF()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "+1",
		},
		{
			name:  "empty input should fail",
			input: "",
			args: args{
				p: Delimited(Char('+'), Digit1(), CRLF()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "",
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotResult := tc.args.p(tc.input)
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

func BenchmarkDelimited(b *testing.B) {
	parser := Delimited(Char('+'), Digit1(), CRLF())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser("+1\r\n")
	}
}

func TestPair(t *testing.T) {
	t.Parallel()

	type args struct {
		leftParser  Parser[string, string]
		rightParser Parser[string, string]
	}
	testCases := []struct {
		name          string
		args          args
		input         string
		wantErr       bool
		wantOutput    PairContainer[string, string]
		wantRemaining string
	}{
		{
			name:  "matching parsers should succeed",
			input: "1abc\r\n",
			args: args{
				leftParser:  Digit1(),
				rightParser: TakeUntil(CRLF()),
			},
			wantErr:       false,
			wantOutput:    PairContainer[string, string]{"1", "abc"},
			wantRemaining: "\r\n",
		},
		{
			name:  "matching left parser, failing right parser, should fail",
			input: "1abc",
			args: args{
				leftParser:  Digit1(),
				rightParser: TakeWhileOneOf('d', 'e', 'f'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "1abc",
		},
		{
			name:  "failing left parser, matching right parser, should fail",
			input: "adef",
			args: args{
				leftParser:  Digit1(),
				rightParser: TakeWhileOneOf('d', 'e', 'f'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "adef",
		},
		{
			name:  "failing left parser, failing right parser, should fail",
			input: "123",
			args: args{
				leftParser:  Digit1(),
				rightParser: TakeWhileOneOf('d', 'e', 'f'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "123",
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			parser := Pair(tc.args.leftParser, tc.args.rightParser)

			gotResult := parser(tc.input)
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

func BenchmarkPair(b *testing.B) {
	parser := Pair(Digit1(), TakeUntil(CRLF()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser("1abc\r\n")
	}
}

func TestPreceded(t *testing.T) {
	t.Parallel()

	type args struct {
		p Parser[string, string]
	}
	testCases := []struct {
		name          string
		args          args
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:  "matching parser should succeed",
			input: "+123",
			args: args{
				p: Preceded(Char('+'), Digit1()),
			},
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "",
		},
		{
			name:  "no prefix match should fail",
			input: "+123",
			args: args{
				p: Preceded(Char('-'), Digit1()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "+123",
		},
		{
			name:  "no parser match should succeed",
			input: "+",
			args: args{
				p: Preceded(Char('+'), Digit1()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "+",
		},
		{
			name:  "empty input should fail",
			input: "",
			args: args{
				p: Preceded(Char('+'), Digit1()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "",
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotResult := tc.args.p(tc.input)
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

func BenchmarkPreceded(b *testing.B) {
	parser := Preceded(Char('+'), Digit1())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser("+123")
	}
}

func TestSeparatedPair(t *testing.T) {
	t.Parallel()

	type args struct {
		leftParser      Parser[string, string]
		separatorParser Parser[string, rune]
		rightParser     Parser[string, string]
	}
	testCases := []struct {
		name          string
		args          args
		input         string
		wantErr       bool
		wantOutput    PairContainer[string, string]
		wantRemaining string
	}{
		// { true, true, true }
		{
			name:  "matching parsers should succeed",
			input: "1|abc\r\n",
			args: args{
				leftParser:      Digit1(),
				separatorParser: Char('|'),
				rightParser:     TakeUntil(CRLF()),
			},
			wantErr:       false,
			wantOutput:    PairContainer[string, string]{"1", "abc"},
			wantRemaining: "\r\n",
		},
		// { true, true, false }
		{
			name:  "matching left parser, matching separator, failing right parser, should fail",
			input: "1|abc",
			args: args{
				leftParser:      Digit1(),
				separatorParser: Char('|'),
				rightParser:     TakeWhileOneOf('d', 'e', 'f'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "1|abc",
		},
		// { true, false, true }
		{
			name:  "matching left parser, failing separator, matching right parser, should fail",
			input: "1^abc",
			args: args{
				leftParser:      Digit1(),
				separatorParser: Char('|'),
				rightParser:     TakeWhileOneOf('a', 'b', 'c'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "1^abc",
		},
		// { true, false, false }
		{
			name:  "matching left parser, failing separator, failing right parser, should fail",
			input: "1^abc",
			args: args{
				leftParser:      Digit1(),
				separatorParser: Char('|'),
				rightParser:     TakeWhileOneOf('d', 'e', 'f'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "1^abc",
		},
		// { false, true, true }
		{
			name:  "failing left parser, matching separator, matching right parser, should fail",
			input: "a|def",
			args: args{
				leftParser:      Digit1(),
				separatorParser: Char('|'),
				rightParser:     TakeWhileOneOf('d', 'e', 'f'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "a|def",
		},
		// { false, true, false }
		{
			name:  "failing left parser, matching separator, failing right parser, should fail",
			input: "a|123",
			args: args{
				leftParser:      Digit1(),
				separatorParser: Char('|'),
				rightParser:     TakeWhileOneOf('d', 'e', 'f'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "a|123",
		},
		// { false, false, true }
		{
			name:  "failing left parser, failing separator, matching right parser, should fail",
			input: "a^def",
			args: args{
				leftParser:      Digit1(),
				separatorParser: Char('|'),
				rightParser:     TakeWhileOneOf('d', 'e', 'f'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "a^def",
		},
		// { false, false, false }
		{
			name:  "failing left parser, failing separator, failing right parser, should fail",
			input: "a^123",
			args: args{
				leftParser:      Digit1(),
				separatorParser: Char('|'),
				rightParser:     TakeWhileOneOf('d', 'e', 'f'),
			},
			wantErr:       true,
			wantOutput:    PairContainer[string, string]{},
			wantRemaining: "a^123",
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			parser := SeparatedPair(tc.args.leftParser, tc.args.separatorParser, tc.args.rightParser)

			gotResult := parser(tc.input)
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

func BenchmarkSeparatedPair(b *testing.B) {
	parser := SeparatedPair(Digit1(), Char('|'), TakeUntil(CRLF()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser("1|abc\r\n")
	}
}

func TestSequence(t *testing.T) {
	t.Parallel()

	type args struct {
		p Parser[string, []string]
	}
	testCases := []struct {
		name          string
		args          args
		input         string
		wantErr       bool
		wantOutput    []string
		wantRemaining string
	}{
		{
			name:  "matching parsers should succeed",
			input: "1a3",
			args: args{
				p: Sequence(Digit1(), Alpha0(), Digit1()),
			},
			wantErr:       false,
			wantOutput:    []string{"1", "a", "3"},
			wantRemaining: "",
		},
		{
			name:  "matching parsers in longer input should succeed",
			input: "1a3bcd",
			args: args{
				p: Sequence(Digit1(), Alpha0(), Digit1()),
			},
			wantErr:       false,
			wantOutput:    []string{"1", "a", "3"},
			wantRemaining: "bcd",
		},
		{
			name:  "partially matching parsers should fail",
			input: "1a3",
			args: args{
				p: Sequence(Digit1(), Digit1(), Digit1()),
			},
			wantErr:       true,
			wantOutput:    nil,
			wantRemaining: "1a3",
		},
		{
			name:  "too short input should fail",
			input: "12",
			args: args{
				p: Sequence(Digit1(), Digit1(), Digit1()),
			},
			wantErr:       true,
			wantOutput:    nil,
			wantRemaining: "12",
		},
		{
			name:  "empty input should succeed",
			input: "",
			args: args{
				p: Sequence(Digit1(), Digit1(), Digit1()),
			},
			wantErr:       true,
			wantOutput:    nil,
			wantRemaining: "",
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotResult := tc.args.p(tc.input)
			if (gotResult.Err != nil) != tc.wantErr {
				t.Errorf("got error %v, want error %v", gotResult.Err, tc.wantErr)
			}

			// testify makes it easier comparing slices
			assert.Equal(t,
				tc.wantOutput, gotResult.Output,
				"got output %v, want output %v", gotResult.Output, tc.wantOutput,
			)

			if gotResult.Remaining != tc.wantRemaining {
				t.Errorf("got remaining %v, want remaining %v", gotResult.Remaining, tc.wantRemaining)
			}
		})
	}
}

func BenchmarkSequence(b *testing.B) {
	parser := Sequence(Digit1(), Alpha0(), Digit1())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser("123")
	}
}

func TestTerminated(t *testing.T) {
	t.Parallel()

	type args struct {
		p Parser[string, string]
	}
	testCases := []struct {
		name          string
		args          args
		input         string
		wantErr       bool
		wantOutput    string
		wantRemaining string
	}{
		{
			name:  "matching parser should succeed",
			input: "1+23",
			args: args{
				p: Terminated(Digit1(), Char('+')),
			},
			wantErr:       false,
			wantOutput:    "1",
			wantRemaining: "23",
		},
		{
			name:  "no suffix match should fail",
			input: "1-23",
			args: args{
				p: Terminated(Digit1(), Char('+')),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "1-23",
		},
		{
			name:  "no parser match should succeed",
			input: "+",
			args: args{
				p: Terminated(Digit1(), Char('+')),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "+",
		},
		{
			name:  "empty input should fail",
			input: "",
			args: args{
				p: Terminated(Digit1(), Char('+')),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "",
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotResult := tc.args.p(tc.input)
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

func BenchmarkTerminated(b *testing.B) {
	parser := Terminated(Digit1(), Char('+'))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser("123+")
	}
}
