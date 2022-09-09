package gomme

import "testing"

func TestNumber(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		parser        Parser[string, float64]
		input         string
		wantErr       bool
		wantOutput    float64
		wantRemaining string
	}{
		{
			name:          "parsing a floating point number with only an integral part should succeed",
			parser:        Number[string](),
			input:         "123",
			wantErr:       false,
			wantOutput:    123,
			wantRemaining: "",
		},
		{
			name:          "parsing a floating point number without sign should succeed",
			parser:        Number[string](),
			input:         "123.456",
			wantErr:       false,
			wantOutput:    123.456,
			wantRemaining: "",
		},
		{
			name:          "parsing a negative floating point number without sign should succeed",
			parser:        Number[string](),
			input:         "-123.456",
			wantErr:       false,
			wantOutput:    -123.456,
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
