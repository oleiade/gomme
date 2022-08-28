package gomme

import (
	"strconv"
	"testing"
)

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

func TestMap(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Foo int
		Bar string
	}

	type args struct {
		parser Parser[string, TestStruct]
	}
	testCases := []struct {
		name          string
		input         string
		args          args
		wantErr       bool
		wantOutput    TestStruct
		wantRemaining string
	}{
		{
			name:  "matching parser should succeed",
			input: "1abc\r\n",
			args: args{
				Map(Pair(Digit1[string](), TakeUntil(CRLF[string]())), func(p PairContainer[string, string]) (TestStruct, error) {
					left, _ := strconv.Atoi(p.Left)
					return TestStruct{
						Foo: left,
						Bar: p.Right,
					}, nil
				}),
			},
			wantErr: false,
			wantOutput: TestStruct{
				Foo: 1,
				Bar: "abc",
			},
			wantRemaining: "\r\n",
		},
		{
			name:  "failing parser should fail",
			input: "abc\r\n",
			args: args{
				Map(Pair(Digit1[string](), TakeUntil(CRLF[string]())), func(p PairContainer[string, string]) (TestStruct, error) {
					left, _ := strconv.Atoi(p.Left)

					return TestStruct{
						Foo: left,
						Bar: p.Right,
					}, nil
				}),
			},
			wantErr:       true,
			wantOutput:    TestStruct{},
			wantRemaining: "abc\r\n",
		},
		{
			name:  "empty input should fail",
			input: "",
			args: args{
				Map(Pair(Digit1[string](), TakeUntil(CRLF[string]())), func(p PairContainer[string, string]) (TestStruct, error) {
					left, _ := strconv.Atoi(p.Left)

					return TestStruct{
						Foo: left,
						Bar: p.Right,
					}, nil
				}),
			},
			wantErr:       true,
			wantOutput:    TestStruct{},
			wantRemaining: "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotResult := tc.args.parser(tc.input)
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

func BenchmarkMap(b *testing.B) {
	type TestStruct struct {
		Foo int
		Bar string
	}

	p := Map(Pair(Digit1[string](), TakeUntil(CRLF[string]())), func(p PairContainer[string, string]) (TestStruct, error) {
		left, _ := strconv.Atoi(p.Left)

		return TestStruct{
			Foo: left,
			Bar: p.Right,
		}, nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p("1abc\r\n")
	}
}

func TestOptional(t *testing.T) {
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
			input: "\r\n123",
			args: args{
				p: Optional(CRLF[string]()),
			},
			wantErr:       false,
			wantOutput:    "\r\n",
			wantRemaining: "123",
		},
		{
			name:  "no match should succeed",
			input: "123",
			args: args{
				p: Optional(CRLF[string]()),
			},
			wantErr:       false,
			wantOutput:    "",
			wantRemaining: "123",
		},
		{
			name:  "empty input should succeed",
			input: "",
			args: args{
				p: Optional(CRLF[string]()),
			},
			wantErr:       false,
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

func BenchmarkOptional(b *testing.B) {
	p := Optional(CRLF[string]())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p("\r\n123")
	}
}

func TestPeek(t *testing.T) {
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
			input: "abcd;",
			args: args{
				p: Peek(Alpha1[string]()),
			},
			wantErr:       false,
			wantOutput:    "abcd",
			wantRemaining: "abcd;",
		},
		{
			name:  "non matching parser should fail",
			input: "123;",
			args: args{
				p: Peek(Alpha1[string]()),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "123;",
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

func BenchmarkPeek(b *testing.B) {
	p := Peek(Alpha1[string]())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p("abcd;")
	}
}

func TestRecognize(t *testing.T) {
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
			input: "123abc",
			args: args{
				p: Recognize(Pair(Digit1[string](), Alpha1[string]())),
			},
			wantErr:       false,
			wantOutput:    "123abc",
			wantRemaining: "",
		},
		{
			name:  "no prefix match should fail",
			input: "abc",
			args: args{
				p: Recognize(Pair(Digit1[string](), Alpha1[string]())),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "abc",
		},
		{
			name:  "no parser match should fail",
			input: "123",
			args: args{
				p: Recognize(Pair(Digit1[string](), Alpha1[string]())),
			},
			wantErr:       true,
			wantOutput:    "",
			wantRemaining: "123",
		},
		{
			name:  "empty input should fail",
			input: "",
			args: args{
				p: Recognize(Pair(Digit1[string](), Alpha1[string]())),
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

func BenchmarkRecognize(b *testing.B) {
	p := Recognize(Pair(Digit1[string](), Alpha1[string]()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p("123abc")
	}
}

func TestAssign(t *testing.T) {
	t.Parallel()

	type args struct {
		p Parser[string, int]
	}
	testCases := []struct {
		name          string
		args          args
		input         string
		wantErr       bool
		wantOutput    int
		wantRemaining string
	}{
		{
			name:  "matching parser should succeed",
			input: "abcd",
			args: args{
				p: Assign(1234, Alpha1[string]()),
			},
			wantErr:       false,
			wantOutput:    1234,
			wantRemaining: "",
		},
		{
			name:  "non matching parser should fail",
			input: "123abcd;",
			args: args{
				p: Assign(1234, Alpha1[string]()),
			},
			wantErr:       true,
			wantOutput:    0,
			wantRemaining: "123abcd;",
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

func BenchmarkAssign(b *testing.B) {
	p := Assign(1234, Alpha1[string]())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p("abcd")
	}
}
