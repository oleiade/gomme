package gomme

import (
	"fmt"
	"testing"
)

func TestTake(t *testing.T) {
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
			name:  "taking less than input size should succeed",
			input: "1234567",
			args: args{
				p: Take[string](6),
			},
			wantErr:       false,
			wantOutput:    "123456",
			wantRemaining: "7",
		},
		{
			name:  "taking exact input size should succeed",
			input: "123456",
			args: args{
				p: Take[string](6),
			},
			wantErr:       false,
			wantOutput:    "123456",
			wantRemaining: "",
		},
		{
			name:  "taking more than input size should fail",
			input: "123",
			args: args{
				p: Take[string](6),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "123",
		},
		{
			name:  "taking from empty input should fail",
			input: "",
			args: args{
				p: Take[string](6),
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

func BenchmarkTake(b *testing.B) {
	p := Take[string](6)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p("123456")
	}
}

func TestTakeUntil(t *testing.T) {
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
			input: "abc123",
			args: args{
				p: TakeUntil(Digit1[string]()),
			},
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "123",
		},
		{
			name:  "immediately matching parser should succeed",
			input: "123",
			args: args{
				p: TakeUntil(Digit1[string]()),
			},
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "123",
		},
		{
			name:  "no match should fail",
			input: "abcdef",
			args: args{
				p: TakeUntil(Digit1[string]()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "abcdef",
		},
		{
			name:  "empty input should fail",
			input: "",
			args: args{
				p: TakeUntil(Digit1[string]()),
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

func BenchmarkTakeUntil(b *testing.B) {
	p := TakeUntil(Digit1[string]())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p("abc123")
	}
}

func TestTakeWhileMN(t *testing.T) {
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
			name:  "parsing input with enough characters and partially matching predicated should succeed",
			input: "latin123",
			args: args{
				p: TakeWhileMN[string](3, 6, IsAlpha),
			},
			wantErr:       false,
			wantOutput:    "latin",
			wantRemaining: "123",
		},
		{
			name:  "parsing input longer than atLeast and atMost should succeed",
			input: "lengthy",
			args: args{
				p: TakeWhileMN[string](3, 6, IsAlpha),
			},
			wantErr:       false,
			wantOutput:    "length",
			wantRemaining: "y",
		},
		{
			name:  "parsing input longer than atLeast and shorter than atMost should succeed",
			input: "latin",
			args: args{
				p: TakeWhileMN[string](3, 6, IsAlpha),
			},
			wantErr:       false,
			wantOutput:    "latin",
			wantRemaining: "",
		},
		{
			name:  "parsing empty input should fail",
			input: "",
			args: args{
				p: TakeWhileMN[string](3, 6, IsAlpha),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "",
		},
		{
			name:  "parsing too short input should fail",
			input: "ed",
			args: args{
				p: TakeWhileMN[string](3, 6, IsAlpha),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "ed",
		},
		{
			name:  "parsing with non-matching predicate should fail",
			input: "12345",
			args: args{
				p: TakeWhileMN[string](3, 6, IsAlpha),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "12345",
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

func BenchmarkTakeWhileMN(b *testing.B) {
	p := TakeWhileMN[string](3, 6, IsAlpha)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p("latin")
	}
}

// TakeWhileOneOf parses any number of characters present in the
// provided collection of runes.
func TakeWhileOneOf[I Bytes](collection ...rune) Parser[I, I] {
	index := make(map[rune]struct{}, len(collection))

	for _, r := range collection {
		index[r] = struct{}{}
	}

	expected := fmt.Sprintf("chars(%v)", string(collection))

	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewError(input, expected), input)
		}

		pos := 0
		for ; pos < len(input); pos++ {
			_, exists := index[rune(input[pos])]
			if !exists {
				if pos == 0 {
					return Failure[I, I](NewError(input, expected), input)
				}

				break
			}
		}

		return Success(input[:pos], input[pos:])
	}
}

func TestTakeWhileOneOf(t *testing.T) {
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
			input: "abc123",
			args: args{
				p: TakeWhileOneOf[string]('a', 'b', 'c'),
			},
			wantErr:       false,
			wantOutput:    "abc",
			wantRemaining: "123",
		},
		{
			name:  "no match should fail",
			input: "123",
			args: args{
				p: TakeWhileOneOf[string]('a', 'b', 'c'),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "123",
		},
		{
			name:  "empty input should fail",
			input: "",
			args: args{
				p: TakeWhileOneOf[string]('a', 'b', 'c'),
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

func BenchmarkTakeWhileOneOf(b *testing.B) {
	p := TakeWhileOneOf[string]('a', 'b', 'c')

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p("abc123")
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
			parser:        Token[string]("Bonjour"),
			input:         "Bonjour tout le monde",
			wantErr:       false,
			wantOutput:    "Bonjour",
			wantRemaining: " tout le monde",
		},
		{
			name:          "parsing a token from an non-matching input should fail",
			parser:        Token[string]("Bonjour"),
			input:         "Hello tout le monde",
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "Hello tout le monde",
		},
		{
			name:          "parsing a token from an empty input should fail",
			parser:        Token[string]("Bonjour"),
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

func BenchmarkToken(b *testing.B) {
	parser := Token[string]("Bonjour")

	for i := 0; i < b.N; i++ {
		parser("Bonjour tout le monde")
	}
}
