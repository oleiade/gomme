package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRGBColor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		input      string
		wantErr    bool
		wantOutput [][]string
	}{
		{
			name:       "parsing a single csv line should succeed",
			input:      "abc,def,ghi\r\n",
			wantErr:    false,
			wantOutput: [][]string{{"abc", "def", "ghi"}},
		},
		{
			name:    "parsing multie csv lines should succeed",
			input:   "abc,def,ghi\r\njkl,mno,pqr\r\n",
			wantErr: false,
			wantOutput: [][]string{
				{"abc", "def", "ghi"},
				{"jkl", "mno", "pqr"},
			},
		},
		{
			name:       "parsing a single csv line of escaped strings should succeed",
			input:      "\"abc\",\"def\",\"ghi\"\r\n",
			wantErr:    false,
			wantOutput: [][]string{{"abc", "def", "ghi"}},
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotOutput, gotErr := ParseCSV(tc.input)
			if (gotErr != nil) != tc.wantErr {
				t.Errorf("got error %v, want error %v", gotErr, tc.wantErr)
			}

			assert.Equal(t,
				tc.wantOutput,
				gotOutput,
				"got output %v, want output %v", gotOutput, tc.wantOutput,
			)
		})
	}
}
