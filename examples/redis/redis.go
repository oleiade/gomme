// Package redis demonstrates the usage of the gomme package to parse Redis'
// [RESP protocol] messages.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
package redis

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/oleiade/gomme"
)

// ParseRESPMESSAGE parses a Redis' [RESP protocol] message.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
func ParseRESPMessage(input string) (RESPMessage, error) {
	if len(input) < 3 {
		return RESPMessage{}, fmt.Errorf("malformed message %s; reason: %w", input, ErrMessageTooShort)
	}

	if !isValidMessageKind(MessageKind(input[0])) {
		return RESPMessage{}, fmt.Errorf("malformed message %s; reason: %w %c", input, ErrInvalidPrefix, input[0])
	}

	if input[len(input)-2] != '\r' || input[len(input)-1] != '\n' {
		return RESPMessage{}, fmt.Errorf("malformed message %s; reason: %w", input, ErrInvalidSuffix)
	}

	parser := gomme.Alternative(
		SimpleString(),
		Error(),
		Integer(),
		BulkString(),
		Array(),
	)

	result := parser(input)
	if result.Err != nil {
		return RESPMessage{}, result.Err
	}

	return result.Output, nil
}

// ErrMessageTooShort is returned when a message is too short to be valid.
// A [RESP protocol] message is at least 3 characters long: the message kind
// prefix, the message content (which can be empty), and the gomme.CRLF suffix.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
var ErrMessageTooShort = errors.New("message too short")

// ErrInvalidPrefix is returned when a message kind prefix is not recognized.
// Valid [RESP Protocol] message kind prefixes are "+", "-", ":", and "$".
//
// [RESP Protocol]: https://redis.io/docs/reference/protocol-spec/
var ErrInvalidPrefix = errors.New("invalid message prefix")

// ErrInvalidSuffix is returned when a message suffix is not recognized.
// Every [RESP protocol] message ends with a gomme.CRLF.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
var ErrInvalidSuffix = errors.New("invalid message suffix")

// RESPMessage is a parsed Redis' [RESP protocol] message.
//
// It can hold either a simple string, an error, an integer, a bulk string,
// or an array. The kind of the message is available in the Kind field.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
type RESPMessage struct {
	Kind         MessageKind
	SimpleString *SimpleStringMessage
	Error        *ErrorStringMessage
	Integer      *IntegerMessage
	BulkString   *BulkStringMessage
	Array        *ArrayMessage
}

// MessageKind is the kind of a Redis' [RESP protocol] message.
type MessageKind string

// The many different kinds of Redis' [RESP protocol] messages map
// to their respective protocol message's prefixes.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
const (
	SimpleStringKind MessageKind = "+"
	ErrorKind        MessageKind = "-"
	IntegerKind      MessageKind = ":"
	BulkStringKind   MessageKind = "$"
	ArrayKind        MessageKind = "*"
	InvalidKind      MessageKind = "?"
)

// SimpleStringMessage is a simple string message parsed from a Redis'
// [RESP protocol] message.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
type SimpleStringMessage struct {
	Content string
}

// SimpleString is a parser for Redis' RESP protocol simple strings.
//
// Simple strings are strings that are not expected to contain newlines.
// Simple strings start with a "+" character, and end with a gomme.CRLF.
//
// Once parsed, the content of the simple string is available in the
// simpleString field of the result's RESPMessage.
func SimpleString() gomme.Parser[string, RESPMessage] {
	mapFn := func(message string) (RESPMessage, error) {
		if strings.ContainsAny(message, "\r\n") {
			return RESPMessage{}, fmt.Errorf("malformed simple string: %s", message)
		}

		return RESPMessage{
			Kind: SimpleStringKind,
			SimpleString: &SimpleStringMessage{
				Content: message,
			},
		}, nil
	}

	return gomme.Delimited(
		gomme.Token(string(SimpleStringKind)),
		gomme.Map(gomme.TakeUntil(gomme.CRLF()), mapFn),
		gomme.CRLF(),
	)
}

// ErrorStringMessage is a parsed error string message from a Redis'
// [RESP protocol] message.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
type ErrorStringMessage struct {
	Kind    string
	Message string
}

// Error is a parser for Redis' RESP protocol errors.
//
// Errors are strings that start with a "-" character, and end with a gomme.CRLF.
//
// The error message is available in the Error field of the result's
// RESPMessage.
func Error() gomme.Parser[string, RESPMessage] {
	mapFn := func(message string) (RESPMessage, error) {
		if strings.ContainsAny(message, "\r\n") {
			return RESPMessage{}, fmt.Errorf("malformed error string: %s", message)
		}

		return RESPMessage{
			Kind: ErrorKind,
			Error: &ErrorStringMessage{
				Kind:    "ERR",
				Message: message,
			},
		}, nil
	}

	return gomme.Delimited(
		gomme.Token(string(ErrorKind)),
		gomme.Map(gomme.TakeUntil(gomme.CRLF()), mapFn),
		gomme.CRLF(),
	)
}

// IntegerMessage is a parsed integer message from a Redis' [RESP protocol]
// message.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
type IntegerMessage struct {
	Value int
}

// Integer is a parser for Redis' RESP protocol integers.
//
// Integers are signed nummerical values represented as string messages
// that start with a ":" character, and end with a gomme.CRLF.
//
// The integer value is available in the IntegerMessage field of the result's
// RESPMessage.
func Integer() gomme.Parser[string, RESPMessage] {
	mapFn := func(message string) (RESPMessage, error) {
		value, err := strconv.Atoi(message)
		if err != nil {
			return RESPMessage{}, err
		}

		return RESPMessage{
			Kind: IntegerKind,
			Integer: &IntegerMessage{
				Value: value,
			},
		}, nil
	}

	return gomme.Delimited(
		gomme.Token(string(IntegerKind)),
		gomme.Map(gomme.TakeUntil(gomme.CRLF()), mapFn),
		gomme.CRLF(),
	)
}

// BulkStringMessage is a parsed bulk string message from a Redis' [RESP protocol]
// message.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
type BulkStringMessage struct {
	Data []byte
}

// BulkString is a parser for Redis' RESP protocol bulk strings.
//
// Bulk strings are binary-safe strings up to 512MB in size.
// Bulk strings start with a "$" character, and end with a gomme.CRLF.
//
// The bulk string's data is available in the BulkString field of the result's
// RESPMessage.
func BulkString() gomme.Parser[string, RESPMessage] {
	mapFn := func(message gomme.PairContainer[int64, string]) (RESPMessage, error) {
		if message.Left < 0 {
			if message.Left < -1 {
				return RESPMessage{}, fmt.Errorf(
					"unable to parse bulk string; "+
						"reason: negative length %d",
					message.Left,
				)
			}

			if message.Left == -1 && len(message.Right) != 0 {
				return RESPMessage{}, fmt.Errorf(
					"malformed array: declared message size -1, and actual size differ %d",
					len(message.Right),
				)
			}
		} else if len(message.Right) != int(message.Left) {
			return RESPMessage{}, fmt.Errorf(
				"malformed array: declared message size %d, and actual size differ %d",
				message.Left,
				len(message.Right),
			)
		}

		return RESPMessage{
			Kind: BulkStringKind,
			BulkString: &BulkStringMessage{
				Data: []byte(message.Right),
			},
		}, nil
	}

	return gomme.Map(
		gomme.Pair(
			sizePrefix(gomme.Token(string(BulkStringKind))),
			gomme.Optional(
				gomme.Terminated(gomme.TakeUntil(gomme.CRLF()), gomme.CRLF()),
			),
		),
		mapFn,
	)
}

// ArrayMessage is a parsed array message from a Redis' [RESP protocol] message.
//
// [RESP protocol]: https://redis.io/docs/reference/protocol-spec/
type ArrayMessage struct {
	Elements []RESPMessage
}

// Array is a parser for Redis' RESP protocol arrays.
//
// Arrays are sequences of RESP messages.
// Arrays start with a "*" character, and end with a gomme.CRLF.
//
// The array's messages are available in the Array field of the result's
// RESPMessage.
func Array() gomme.Parser[string, RESPMessage] {
	mapFn := func(message gomme.PairContainer[int64, []RESPMessage]) (RESPMessage, error) {
		if int(message.Left) == -1 {
			if len(message.Right) != 0 {
				return RESPMessage{}, fmt.Errorf(
					"malformed array: declared message size -1, and actual size differ %d",
					len(message.Right),
				)
			}
		} else {
			if len(message.Right) != int(message.Left) {
				return RESPMessage{}, fmt.Errorf(
					"malformed array: declared message size %d, and actual size differ %d",
					message.Left,
					len(message.Right),
				)
			}
		}

		messages := make([]RESPMessage, 0, len(message.Right))
		messages = append(messages, message.Right...)

		return RESPMessage{
			Kind: ArrayKind,
			Array: &ArrayMessage{
				Elements: messages,
			},
		}, nil
	}

	return gomme.Map(
		gomme.Pair(
			sizePrefix(gomme.Token(string(ArrayKind))),
			gomme.Many(
				gomme.Alternative(
					SimpleString(),
					Error(),
					Integer(),
					BulkString(),
				),
			),
		),
		mapFn,
	)
}

func sizePrefix(prefix gomme.Parser[string, string]) gomme.Parser[string, int64] {
	return gomme.Delimited(
		prefix,
		gomme.Int64(),
		gomme.CRLF(),
	)
}

func isValidMessageKind(kind MessageKind) bool {
	return kind == SimpleStringKind ||
		kind == ErrorKind ||
		kind == IntegerKind ||
		kind == BulkStringKind ||
		kind == ArrayKind
}
