package gomme

import (
	"fmt"
	"strconv"
	"strings"
)

// Char parses a single character and matches it with
// a provided candidate.
func Char[I Bytes](character rune) Parser[I, rune] {
	return func(input I) Result[rune, I] {
		if len(input) == 0 || rune(input[0]) != character {
			return Failure[I, rune](NewGenericError(input, string(character)), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// AnyChar parses any single character.
func AnyChar[I Bytes]() Parser[I, rune] {
	return func(input I) Result[rune, I] {
		if len(input) == 0 {
			return Failure[I, rune](NewGenericError(input, "any character"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// Alpha0 parses a zero or more lowercase or uppercase alphabetic characters: a-z, A-Z.
// In the cases where the input is empty, or no terminating character is found, the parser
// returns the input as is.
func Alpha0[I Bytes]() Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Success(input, input)
		}

		lastAlphaPos := 0
		for idx, c := range input {
			if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
				return Success(input[:idx], input[idx:])
			}

			lastAlphaPos++
		}

		return Success(input[:lastAlphaPos], input[lastAlphaPos:])
	}
}

// Alpha1 parses one or more lowercase or uppercase alphabetic characters: a-z, A-Z.
// In the cases where the input doesn't hold enough data, or a terminating character
// is found before any matching ones were, the parser returns an error result.
func Alpha1[I Bytes]() Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewGenericError(input, "alpha1"), input)
		}

		if (input[0] < 'a' || input[0] > 'z') && (input[0] < 'A' || input[0] > 'Z') {
			return Failure[I, I](NewGenericError(input, "alpha1"), input)
		}

		lastAlphaPos := 1
		for idx, c := range input[1:] {
			if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
				return Success(input[:idx+1], input[idx+1:])
			}

			lastAlphaPos++
		}

		return Success(input[:lastAlphaPos], input[lastAlphaPos:])
	}
}

// Alphanumeric0 parses zero or more ASCII alphabetical or numerical characters: a-z, A-Z, 0-9.
// In the cases where the input is empty, or no terminating character is found, the parser
// returns the input as is.
func Alphanumeric0[I Bytes]() Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Success(input, input)
		}

		lastDigitPos := 0
		for idx, c := range input {
			if !isAlphanumeric(c) {
				return Success(input[:idx], input[idx:])
			}

			lastDigitPos++
		}

		return Success(input[:lastDigitPos], input[lastDigitPos:])
	}
}

// Alphanumeric1 parses one or more alphabetical or numerical characters: a-z, A-Z, 0-9.
// In the cases where the input doesn't hold enough data, or a terminating character
// is found before any matching ones were, the parser returns an error result.
func Alphanumeric1[I Bytes]() Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewGenericError(input, "digit1"), input)
		}

		if !isAlphanumeric(rune(input[0])) {
			return Failure[I, I](NewGenericError(input, "digit1"), input)
		}

		lastDigitPos := 1
		for idx, c := range input[1:] {
			if !isAlphanumeric(c) {
				return Success(input[:idx+1], input[idx+1:])
			}

			lastDigitPos++
		}

		return Success(input[:lastDigitPos], input[lastDigitPos:])
	}
}

// Digit0 parses zero or more ASCII numerical characters: 0-9.
// In the cases where the input is empty, or no terminating character is found, the parser
// returns the input as is.
func Digit0[I Bytes]() Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Success(input, input)
		}

		lastDigitPos := 0
		for idx, c := range input {
			if c < '0' || c > '9' {
				return Success(input[:idx], input[idx:])
			}

			lastDigitPos++
		}

		return Success(input[:lastDigitPos], input[lastDigitPos:])
	}
}

// Digit1 parses one or more numerical characters: 0-9.
// In the cases where the input doesn't hold enough data, or a terminating character
// is found before any matching ones were, the parser returns an error result.
func Digit1[I Bytes]() Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewGenericError(input, "digit1"), input)
		}

		if input[0] < '0' || input[0] > '9' {
			return Failure[I, I](NewGenericError(input, "digit1"), input)
		}

		lastDigitPos := 1
		for idx, c := range input[1:] {
			if c < '0' || c > '9' {
				return Success(input[:idx+1], input[idx+1:])
			}

			lastDigitPos++
		}

		return Success(input[:lastDigitPos], input[lastDigitPos:])
	}
}

// HexDigit0 parses zero or more ASCII hexadecimal characters: a-f, A-F, 0-9.
// In the cases where the input is empty, or no terminating character is found, the parser
// returns the input as is.
func HexDigit0[I Bytes]() Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Success(input, input)
		}

		lastDigitPos := 0
		for idx, c := range input {
			if !isHexDigit(c) {
				return Success(input[:idx], input[idx:])
			}

			lastDigitPos++
		}

		return Success(input[:lastDigitPos], input[lastDigitPos:])
	}
}

// HexDigit1 parses one or more ASCII hexadecimal characters: a-f, A-F, 0-9.
// In the cases where the input doesn't hold enough data, or a terminating character
// is found before any matching ones were, the parser returns an error result.
func HexDigit1[I Bytes]() Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewGenericError(input, "digit1"), input)
		}

		if !isHexDigit(rune(input[0])) {
			return Failure[I, I](NewGenericError(input, "digit1"), input)
		}

		lastDigitPos := 1
		for idx, c := range input[1:] {
			if !isHexDigit(c) {
				return Success(input[:idx+1], input[idx+1:])
			}

			lastDigitPos++
		}

		return Success(input[:lastDigitPos], input[lastDigitPos:])
	}
}

// LF parses a line feed `\n` character.
func LF[I Bytes]() Parser[I, rune] {
	return func(input I) Result[rune, I] {
		if len(input) == 0 || input[0] != '\n' {
			return Failure[I, rune](NewGenericError(input, "line feed ('\\n')"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// CR parses a carriage return `\r` character.
func CR[I Bytes]() Parser[I, rune] {
	return func(input I) Result[rune, I] {
		if len(input) == 0 || input[0] != '\r' {
			return Failure[I, rune](NewGenericError(input, "carriage return ('\\r')"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// CRLF parses the string `\r\n`.
func CRLF[I Bytes]() Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) < 2 || (input[0] != '\r' || input[1] != '\n') {
			return Failure[I, I](NewGenericError(input, "CRLF ('\\r\\n')"), input)
		}

		return Success(input[:2], input[2:])
	}
}

// Space parses a space character.
func Space[I Bytes]() Parser[I, rune] {
	return func(input I) Result[rune, I] {
		if len(input) == 0 || input[0] != ' ' {
			return Failure[I, rune](NewGenericError(input, "space"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// Tab parses a tab character.
func Tab[I Bytes]() Parser[I, rune] {
	return func(input I) Result[rune, I] {
		if len(input) == 0 || input[0] != '\t' {
			return Failure[I, rune](NewGenericError(input, "tab"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// Token parses a token from the input, and returns the part of the input that
// matched the token.
// If the token could not be found, the parser returns an error result.
func Token[I Bytes](token string) Parser[I, I] {
	return func(input I) Result[I, I] {
		if !strings.HasPrefix(string(input), token) {
			return Failure[I, I](NewGenericError(input, fmt.Sprintf("tag(%s)", token)), input)
		}

		return Success(input[:len(token)], input[len(token):])
	}
}

// Int64 parses an integer from the input, and returns the part of the input that
// matched the integer.
func Int64[I Bytes]() Parser[I, int64] {
	return func(input I) Result[int64, I] {
		parser := Recognize(Sequence(Optional(Token[I]("-")), Digit1[I]()))

		result := parser(input)
		if result.Err != nil {
			return Failure[I, int64](NewGenericError(input, "int64"), input)
		}

		n, err := strconv.ParseInt(string(result.Output), 10, 64)
		if err != nil {
			return Failure[I, int64](NewGenericError(input, "int64"), input)
		}

		return Success(n, result.Remaining)
	}
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlphanumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func isHexDigit(c rune) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
