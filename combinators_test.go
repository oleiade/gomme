package gomme

import (
	"errors"
	"strconv"
	"testing"
)

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
			name:  "failing mapper should fail",
			input: "1abc\r\n",
			args: args{
				Map(Pair(Digit1[string](), TakeUntil(CRLF[string]())), func(p PairContainer[string, string]) (TestStruct, error) {
					return TestStruct{}, errors.New("unexpected error")
				}),
			},
			wantErr:       true,
			wantOutput:    TestStruct{},
			wantRemaining: "1abc\r\n",
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
