package hexcolor

import "testing"

func TestParseRGBColor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		input     string
		wantErr   bool
		wantColor RGBColor
	}{
		{
			name:      "parsing minimum hexadecimal color should succeed",
			input:     "#000000",
			wantErr:   false,
			wantColor: RGBColor{0, 0, 0},
		},
		{
			name:      "parsing maximum hexadecimal color should succeed",
			input:     "#ffffff",
			wantErr:   false,
			wantColor: RGBColor{255, 255, 255},
		},
		{
			name:      "parsing out of bound color component should fail",
			input:     "#fffffg",
			wantErr:   true,
			wantColor: RGBColor{},
		},
		{
			name:      "omitting leading # character should fail",
			input:     "ffffff",
			wantErr:   true,
			wantColor: RGBColor{},
		},
		{
			name:      "parsing insufficient number of characters should fail",
			input:     "#fffff",
			wantErr:   true,
			wantColor: RGBColor{},
		},
		{
			name:      "empty input should fail",
			input:     "",
			wantErr:   true,
			wantColor: RGBColor{},
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotColor, gotErr := ParseRGBColor(tc.input)
			if (gotErr != nil) != tc.wantErr {
				t.Errorf("got error %v, want error %v", gotErr, tc.wantErr)
			}

			if gotColor != tc.wantColor {
				t.Errorf("got color %v, want color %v", gotColor, tc.wantColor)
			}
		})
	}
}
