package redis

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParseRESPMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}
	testCases := []struct {
		name    string
		args    args
		want    RESPMessage
		wantErr bool
	}{
		//
		// General
		//
		{
			name: "empty message should fail",
			args: args{
				input: "",
			},
			want:    RESPMessage{},
			wantErr: true,
		},
		{
			name: "message with only a prefix should fail",
			args: args{
				input: "+",
			},
			want:    RESPMessage{},
			wantErr: true,
		},
		{
			name: "message with only a CRLF should fail",
			args: args{
				input: "\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},
		{
			name: "message with an invalid prefix should fail",
			args: args{
				input: "?\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},

		//
		// SIMPLE STRINGS
		//

		{
			name: "proper simple string should succeed",
			args: args{
				"+OK\r\n",
			},
			want: RESPMessage{
				Kind:         SimpleStringKind,
				SimpleString: &SimpleStringMessage{Content: "OK"},
			},
			wantErr: false,
		},
		{
			name: "empty simple string should succeed",
			args: args{
				"+\r\n",
			},
			want: RESPMessage{
				Kind:         SimpleStringKind,
				SimpleString: &SimpleStringMessage{Content: ""},
			},
			wantErr: false,
		},
		{
			name: "malformed simple string containing a \\r should fail",
			args: args{
				"+Hello\rWorld\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},
		{
			name: "malformed simple string containing a \\n should fail",
			args: args{
				"+Hello\nWorld\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},
		{
			name: "malformed simple string containing a \\n\\r should fail",
			args: args{
				"+Hello\n\rWorld\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},

		// //
		// // ERROR STRINGS
		// //

		{
			name: "proper error string should succeed",
			args: args{
				"-Error message\r\n",
			},
			want: RESPMessage{
				Kind: ErrorKind,
				Error: &ErrorStringMessage{
					Kind:    "ERR",
					Message: "Error message",
				},
			},
			wantErr: false,
		},
		{
			name: "malformed error string containing a \\r should fail",
			args: args{
				"-Error\r message\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},
		{
			name: "malformed error string containing a \\n should fail",
			args: args{
				"-Error\n message\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},
		{
			name: "malformed error string containing a \\n\\r should fail",
			args: args{
				"-Error\n\r message\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},

		// //
		// // INTEGER
		// //

		{
			name: "proper integer should succeed",
			args: args{
				":123\r\n",
			},
			want: RESPMessage{
				Kind: IntegerKind,
				Integer: &IntegerMessage{
					Value: 123,
				},
			},
			wantErr: false,
		},

		//
		// Bulk Strings
		//

		{
			name: "proper bulk string should succeed",
			args: args{
				"$5\r\nhello\r\n",
			},
			want: RESPMessage{
				Kind: BulkStringKind,
				BulkString: &BulkStringMessage{
					Data: []byte("hello"),
				},
			},
			wantErr: false,
		},
		{
			name: "nil bulk string should succeed",
			args: args{
				"$-1\r\n",
			},
			want: RESPMessage{
				Kind: BulkStringKind,
				BulkString: &BulkStringMessage{
					Data: []byte(""),
				},
			},
			wantErr: false,
		},
		{
			name: "bulk string with negative size != -1 should fail",
			args: args{
				"$-2\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},
		{
			name: "malformed bulk string with actual length different from declared length should fail",
			args: args{
				"$5\r\nhello world\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},

		//
		// ARRAYS
		//

		{
			name: "proper array of simple strings should succeed",
			args: args{
				"*2\r\n+hello\r\n+world\r\n",
			},
			want: RESPMessage{
				Kind: ArrayKind,
				Array: &ArrayMessage{
					Elements: []RESPMessage{
						{
							Kind: SimpleStringKind,
							SimpleString: &SimpleStringMessage{
								Content: "hello",
							},
						},
						{
							Kind: SimpleStringKind,
							SimpleString: &SimpleStringMessage{
								Content: "world",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "proper array of errors should succeed",
			args: args{
				"*2\r\n-Error Message\r\n-Other error\r\n",
			},
			want: RESPMessage{
				Kind: ArrayKind,
				Array: &ArrayMessage{
					Elements: []RESPMessage{
						{
							Kind: ErrorKind,
							Error: &ErrorStringMessage{
								Kind:    "ERR",
								Message: "Error Message",
							},
						},
						{
							Kind: ErrorKind,
							Error: &ErrorStringMessage{
								Kind:    "ERR",
								Message: "Other error",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "proper array of integers should succeed",
			args: args{
				"*2\r\n:0\r\n:1000\r\n",
			},
			want: RESPMessage{
				Kind: ArrayKind,
				Array: &ArrayMessage{
					Elements: []RESPMessage{
						{
							Kind: IntegerKind,
							Integer: &IntegerMessage{
								Value: 0,
							},
						},
						{
							Kind: IntegerKind,
							Integer: &IntegerMessage{
								Value: 1000,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "proper array of bulk strings should succeed",
			args: args{
				"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			},
			want: RESPMessage{
				Kind: ArrayKind,
				Array: &ArrayMessage{
					Elements: []RESPMessage{
						{
							Kind: BulkStringKind,
							BulkString: &BulkStringMessage{
								Data: []byte("hello"),
							},
						},
						{
							Kind: BulkStringKind,
							BulkString: &BulkStringMessage{
								Data: []byte("world"),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "proper array of mixed types should succeed",
			args: args{
				"*4\r\n$5\r\nhello\r\n:123\r\n+OK\r\n-Error Message\r\n",
			},
			want: RESPMessage{
				Kind: ArrayKind,
				Array: &ArrayMessage{
					Elements: []RESPMessage{
						{
							Kind: BulkStringKind,
							BulkString: &BulkStringMessage{
								Data: []byte("hello"),
							},
						},
						{
							Kind: IntegerKind,
							Integer: &IntegerMessage{
								Value: 123,
							},
						},
						{
							Kind: SimpleStringKind,
							SimpleString: &SimpleStringMessage{
								Content: "OK",
							},
						},
						{
							Kind: ErrorKind,
							Error: &ErrorStringMessage{
								Kind:    "ERR",
								Message: "Error Message",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty array should succeed",
			args: args{
				"*-1\r\n",
			},
			want: RESPMessage{
				Kind: ArrayKind,
				Array: &ArrayMessage{
					Elements: []RESPMessage{},
				},
			},
			wantErr: false,
		},
		{
			name: "array with non matching size prefix should fail",
			args: args{
				"*2\r\n+OK\r\n",
			},
			want:    RESPMessage{},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseRESPMessage(tc.args.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseRESPMessage() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseRESPMessage() = %v, want %v", got, tc.want)
			}
		})
	}
}

func BenchmarkParseMessage(b *testing.B) {
	var benchmarks = []struct {
		kind string
		data string
		size string
	}{
		{"simple_string", "+OK\r\n", "2"},
		{"simple_string", simpleStringProducer(128 * Byte), "128b"},
		{"simple_string", simpleStringProducer(1 * KiloBytes), "1kb"},
		{"simple_string", simpleStringProducer(1 * MegaBytes), "1mb"},
		{"error_string", "-Error\r\n", "5"},
		{"error_string", errorStringProducer(128 * Byte), "128b"},
		{"error_string", errorStringProducer(1 * KiloBytes), "1kb"},
		{"integer", ":1\r\n", "1"},
		{"integer", ":9,223,372,036,854,775,807\r\n", "biggest integer"},
		{"integer", ":-9223372036854775808\r\n", "smallest integer"},
		{"bulk_string", bulkStringProducer(128 * Byte), "128b"},
		{"bulk_string", bulkStringProducer(1 * KiloBytes), "1kb"},
		{"bulk_string", bulkStringProducer(1 * MegaBytes), "1mb"},
		{"bulk_string", bulkStringProducer(512 * MegaBytes), "512mb"},
		{"array", arrayProducer(10000, 128*Byte), "10000 * 128b"},
		{"array", arrayProducer(1000, 1*KiloBytes), "1000 * 1kb"},
		{"array", arrayProducer(100, 1*MegaBytes), "100 * 1mb"},
	}

	for _, tt := range benchmarks {
		b.Run(fmt.Sprintf("%s_with_size_%s", tt.kind, tt.size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				//nolint:errcheck,gosec
				ParseRESPMessage(tt.data)
			}
		})
	}
}

const (
	Byte      = 1
	KiloBytes = Byte * 1024
	MegaBytes = KiloBytes * 1024
	GigaBytes = MegaBytes * 1024
	TeraBytes = GigaBytes * 1024
)

// TODO: add fuzz tests input for other kind of messages,
// and handled their expected format too.
func FuzzTestParseMessage(f *testing.F) {
	testCases := []string{
		"+OK\r\n",
		"+Hello world\r\n",
		"+This is a string\r\n",
	}

	for _, testCase := range testCases {
		f.Add(testCase)
	}

	f.Fuzz(func(t *testing.T, message string) {
		_, err := ParseRESPMessage(message)
		if err != nil {
			if errors.Is(err, ErrMessageTooShort) || errors.Is(err, ErrInvalidPrefix) || errors.Is(err, ErrInvalidSuffix) {
				t.Skip("skipping expected error")
			}

			if strings.Count(message, "\r") > 1 || strings.Count(message, "\n") > 1 {
				t.Skip("skipping simple string message with multiple \\r or \\n")
			}

			t.Errorf("ParseRESPMessage() error = %v", err)
		}
	})
}

func simpleStringProducer(messageSize int) string {
	return strings.Join(
		[]string{
			"+",
			stringWithinCharset(messageSize, alnumCharset),
			"\r\n",
		},
		"",
	)
}

func errorStringProducer(messageSize int) string {
	return strings.Join(
		[]string{
			"-",
			stringWithinCharset(messageSize, alnumCharset),
			"\r\n",
		},
		"",
	)
}

func bulkStringProducer(messageSize int) string {
	return strings.Join(
		[]string{
			"$",
			strconv.Itoa(messageSize),
			"\r\n",
			stringWithinCharset(messageSize, alnumCharset),
			"\r\n",
		},
		"",
	)
}

func arrayProducer(arraySize, messageSize int) string {
	messages := make([]string, 0, arraySize)

	for i := 0; i < arraySize; i++ {
		messageKind := i % 4

		switch messageKind {
		case 0:
			messages = append(messages, simpleStringProducer(messageSize))
		case 1:
			messages = append(messages, errorStringProducer(messageSize))
		case 2:
			messages = append(messages, bulkStringProducer(messageSize))
		case 3:
			messages = append(messages, strconv.Itoa(rand.Int()))
		}
	}

	return strings.Join(
		[]string{
			"*",
			strings.Join(messages, ""),
			"\r\n",
		},
		"",
	)
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const alnumCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func stringWithinCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
