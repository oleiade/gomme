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
			parser:        Count(Token("abc"), 2),
			input:         "abcabc",
			wantErr:       false,
			wantOutput:    []string{"abc", "abc"},
			wantRemaining: "",
		},
		{
			name:          "parsing more than count should succeed",
			parser:        Count(Token("abc"), 2),
			input:         "abcabcabc",
			wantErr:       false,
			wantOutput:    []string{"abc", "abc"},
			wantRemaining: "abc",
		},
		{
			name:          "parsing less than count should fail",
			parser:        Count(Token("abc"), 2),
			input:         "abc123",
			wantErr:       true,
			wantOutput:    nil,
			wantRemaining: "abc123",
		},
		{
			name:          "parsing no count should fail",
			parser:        Count(Token("abc"), 2),
			input:         "123123",
			wantErr:       true,
			wantOutput:    nil,
			wantRemaining: "123123",
		},
		{
			name:          "parsing empty input should fail",
			parser:        Count(Token("abc"), 2),
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

func TestMany(t *testing.T) {
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
				p: Many(Char('#')),
			},
			wantErr:       false,
			wantOutput:    []rune{'#', '#', '#'},
			wantRemaining: "",
		},
		{
			name:  "no match should succeed",
			input: "abc",
			args: args{
				p: Many(Char('#')),
			},
			wantErr:       false,
			wantOutput:    []rune{},
			wantRemaining: "abc",
		},
		{
			name:  "empty input should succeed",
			input: "",
			args: args{
				p: Many(Char('#')),
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
