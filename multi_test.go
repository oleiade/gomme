package gomme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, []string]
		input         string
		wantErr       bool
		wantOutput    []string
		wantRemaining string
	}{
		{
			name:          "parsing exact count should succeed",
			parser:        Count(Token[string]("abc"), 2),
			input:         "abcabc",
			wantErr:       false,
			wantOutput:    []string{"abc", "abc"},
			wantRemaining: "",
		},
		{
			name:          "parsing more than count should succeed",
			parser:        Count(Token[string]("abc"), 2),
			input:         "abcabcabc",
			wantErr:       false,
			wantOutput:    []string{"abc", "abc"},
			wantRemaining: "abc",
		},
		{
			name:          "parsing less than count should fail",
			parser:        Count(Token[string]("abc"), 2),
			input:         "abc123",
			wantErr:       true,
			wantOutput:    nil,
			wantRemaining: "abc123",
		},
		{
			name:          "parsing no count should fail",
			parser:        Count(Token[string]("abc"), 2),
			input:         "123123",
			wantErr:       true,
			wantOutput:    nil,
			wantRemaining: "123123",
		},
		{
			name:          "parsing empty input should fail",
			parser:        Count(Token[string]("abc"), 2),
			input:         "",
			wantErr:       true,
			wantOutput:    nil,
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

func BenchmarkCount(b *testing.B) {
	parser := Count(Char[string]('#'), 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser("###")
	}
}

func TestMany0(t *testing.T) {
	t.Parallel()

	type args struct {
		p Parser[string, []rune]
	}
	testCases := []struct {
		name          string
		args          args
		input         string
		wantErr       bool
		wantOutput    []rune
		wantRemaining string
	}{
		{
			name:  "matching parser should succeed",
			input: "###",
			args: args{
				p: Many0(Char[string]('#')),
			},
			wantErr:       false,
			wantOutput:    []rune{'#', '#', '#'},
			wantRemaining: "",
		},
		{
			name:  "no match should succeed",
			input: "abc",
			args: args{
				p: Many0(Char[string]('#')),
			},
			wantErr:       false,
			wantOutput:    []rune{},
			wantRemaining: "abc",
		},
		{
			name:  "empty input should succeed",
			input: "",
			args: args{
				p: Many0(Char[string]('#')),
			},
			wantErr:       false,
			wantOutput:    []rune{},
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

func TestMany0DetectsInfiniteLoops(t *testing.T) {
	t.Parallel()

	// Digit0 accepts empty input, and would cause an infinite loop if not detected
	input := "abcdef"
	parser := Many0(Digit0[string]())

	result := parser(input)

	assert.Error(t, result.Err)
	assert.Nil(t, result.Output)
	assert.Equal(t, input, result.Remaining)
}

func BenchmarkMany0(b *testing.B) {
	parser := Many0(Char[string]('#'))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser("###")
	}
}

func TestMany1(t *testing.T) {
	t.Parallel()

	type args struct {
		p Parser[string, []rune]
	}
	testCases := []struct {
		name          string
		args          args
		input         string
		wantErr       bool
		wantOutput    []rune
		wantRemaining string
	}{
		{
			name:  "matching parser should succeed",
			input: "###",
			args: args{
				p: Many1(Char[string]('#')),
			},
			wantErr:       false,
			wantOutput:    []rune{'#', '#', '#'},
			wantRemaining: "",
		},
		{
			name:  "matching at least once should succeed",
			input: "#abc",
			args: args{
				p: Many1(Char[string]('#')),
			},
			wantErr:       false,
			wantOutput:    []rune{'#'},
			wantRemaining: "abc",
		},
		{
			name:  "not matching at least once should fail",
			input: "a##",
			args: args{
				p: Many1(Char[string]('#')),
			},
			wantErr:       true,
			wantOutput:    nil,
			wantRemaining: "a##",
		},
		{
			name:  "no match should fail",
			input: "abc",
			args: args{
				p: Many1(Char[string]('#')),
			},
			wantErr:       true,
			wantOutput:    nil,
			wantRemaining: "abc",
		},
		{
			name:  "empty input should fail",
			input: "",
			args: args{
				p: Many1(Char[string]('#')),
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

func TestMany1DetectsInfiniteLoops(t *testing.T) {
	t.Parallel()

	// Digit0 accepts empty input, and would cause an infinite loop if not detected
	input := "abcdef"
	parser := Many1(Digit0[string]())

	result := parser(input)

	assert.Error(t, result.Err)
	assert.Nil(t, result.Output)
	assert.Equal(t, input, result.Remaining)
}

func BenchmarkMany1(b *testing.B) {
	parser := Many1(Char[string]('#'))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser("###")
	}
}
