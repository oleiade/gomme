package gomme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlternative(t *testing.T) {
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
			name:  "head matching parser should succeed",
			input: "123",
			args: args{
				p: Alternative(Digit1(), Alpha0()),
			},
			wantErr:       false,
			wantOutput:    "123",
			wantRemaining: "",
		},
		{
			name:  "matching parser should succeed",
			input: "1",
			args: args{
				p: Alternative(Digit1(), Alpha0()),
			},
			wantErr:       false,
			wantOutput:    "1",
			wantRemaining: "",
		},
		{
			name:  "no matching parser should fail",
			input: "$%^*",
			args: args{
				p: Alternative(Digit1(), Alpha1()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "$%^*",
		},
		{
			name:  "empty input should fail",
			input: "",
			args: args{
				p: Alternative(Digit1(), Alpha1()),
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

func BenchmarkAlternative(b *testing.B) {
	p := Alternative(Digit1(), Alpha1())

	for i := 0; i < b.N; i++ {
		p("123")
	}
}
