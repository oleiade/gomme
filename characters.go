package gomme

import (
	"strconv"
)

// Char parses a single character and matches it with
// a provided candidate.
func Char[Input Bytes](character rune) Parser[Input, rune] {
	return func(input Input) Result[rune, Input] {
		if len(input) == 0 || rune(input[0]) != character {
			return Failure[Input, rune](NewError(input, string(character)), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// AnyChar parses any single character.
func AnyChar[Input Bytes]() Parser[Input, rune] {
	return func(input Input) Result[rune, Input] {
		if len(input) == 0 {
			return Failure[Input, rune](NewError(input, "AnyChar"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// Alpha0 parses a zero or more lowercase or uppercase alphabetic characters: a-z, A-Z.
// In the cases where the input is empty, or no terminating character is found, the parser
// returns the input as is.
func Alpha0[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Success(input, input)
		}

		lastAlphaPos := 0
		for idx := 0; idx < len(input); idx++ {
			if !IsAlpha(rune(input[idx])) {
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
func Alpha1[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Failure[Input, Input](NewError(input, "Alpha1"), input)
		}

		if !IsAlpha(rune(input[0])) {
			return Failure[Input, Input](NewError(input, "Alpha1"), input)
		}

		lastAlphaPos := 1
		for idx := 1; idx < len(input); idx++ {
			if !IsAlpha(rune(input[idx])) {
				return Success(input[:idx], input[idx:])
			}

			lastAlphaPos++
		}

		return Success(input[:lastAlphaPos], input[lastAlphaPos:])
	}
}

// Alphanumeric0 parses zero or more ASCII alphabetical or numerical characters: a-z, A-Z, 0-9.
// In the cases where the input is empty, or no terminating character is found, the parser
// returns the input as is.
func Alphanumeric0[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Success(input, input)
		}

		lastDigitPos := 0
		for idx := 0; idx < len(input); idx++ {
			if !IsAlphanumeric(rune(input[idx])) {
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
func Alphanumeric1[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Failure[Input, Input](NewError(input, "Digit1"), input)
		}

		if !IsAlphanumeric(rune(input[0])) {
			return Failure[Input, Input](NewError(input, "Digit1"), input)
		}

		lastDigitPos := 1
		for idx := 1; idx < len(input); idx++ {
			if !IsAlphanumeric(rune(input[idx])) {
				return Success(input[:idx], input[idx:])
			}

			lastDigitPos++
		}

		return Success(input[:lastDigitPos], input[lastDigitPos:])
	}
}

// Digit0 parses zero or more ASCII numerical characters: 0-9.
// In the cases where the input is empty, or no terminating character is found, the parser
// returns the input as is.
func Digit0[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Success(input, input)
		}

		lastDigitPos := 0
		for idx := 0; idx < len(input); idx++ {
			if !IsDigit(rune(input[idx])) {
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
func Digit1[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Failure[Input, Input](NewError(input, "Digit1"), input)
		}

		if !IsDigit(rune(input[0])) {
			return Failure[Input, Input](NewError(input, "Digit1"), input)
		}

		lastDigitPos := 1
		for idx := 1; idx < len(input); idx++ {
			if !IsDigit(rune(input[idx])) {
				return Success(input[:idx], input[idx:])
			}

			lastDigitPos++
		}

		return Success(input[:lastDigitPos], input[lastDigitPos:])
	}
}

// HexDigit0 parses zero or more ASCII hexadecimal characters: a-f, A-F, 0-9.
// In the cases where the input is empty, or no terminating character is found, the parser
// returns the input as is.
func HexDigit0[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Success(input, input)
		}

		lastDigitPos := 0
		for idx := 0; idx < len(input); idx++ {
			if !IsHexDigit(rune(input[idx])) {
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
func HexDigit1[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Failure[Input, Input](NewError(input, "HexDigit1"), input)
		}

		if !IsHexDigit(rune(input[0])) {
			return Failure[Input, Input](NewError(input, "HexDigit1"), input)
		}

		lastDigitPos := 1
		for idx := 1; idx < len(input); idx++ {
			if !IsHexDigit(rune(input[idx])) {
				return Success(input[:idx], input[idx:])
			}

			lastDigitPos++
		}

		return Success(input[:lastDigitPos], input[lastDigitPos:])
	}
}

// Whitespace0 parses zero or more whitespace characters: ' ', '\t', '\n', '\r'.
// In the cases where the input is empty, or no terminating character is found, the parser
// returns the input as is.
func Whitespace0[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Success(input, input)
		}

		lastPos := 0
		for idx := 0; idx < len(input); idx++ {
			if !IsWhitespace(rune(input[idx])) {
				return Success(input[:idx], input[idx:])
			}

			lastPos++
		}

		return Success(input[:lastPos], input[lastPos:])
	}
}

// Whitespace1 parses one or more whitespace characters: ' ', '\t', '\n', '\r'.
// In the cases where the input doesn't hold enough data, or a terminating character
// is found before any matching ones were, the parser returns an error result.
func Whitespace1[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Failure[Input, Input](NewError(input, "WhiteSpace1"), input)
		}

		if !IsWhitespace(rune(input[0])) {
			return Failure[Input, Input](NewError(input, "WhiteSpace1"), input)
		}

		lastPos := 1
		for idx := 1; idx < len(input); idx++ {
			if !IsWhitespace(rune(input[idx])) {
				return Success(input[:idx], input[idx:])
			}

			lastPos++
		}

		return Success(input[:lastPos], input[lastPos:])
	}
}

// LF parses a line feed `\n` character.
func LF[Input Bytes]() Parser[Input, rune] {
	return func(input Input) Result[rune, Input] {
		if len(input) == 0 || input[0] != '\n' {
			return Failure[Input, rune](NewError(input, "LF"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// CR parses a carriage return `\r` character.
func CR[Input Bytes]() Parser[Input, rune] {
	return func(input Input) Result[rune, Input] {
		if len(input) == 0 || input[0] != '\r' {
			return Failure[Input, rune](NewError(input, "CR"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// CRLF parses the string `\r\n`.
func CRLF[Input Bytes]() Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) < 2 || (input[0] != '\r' || input[1] != '\n') {
			return Failure[Input, Input](NewError(input, "CRLF"), input)
		}

		return Success(input[:2], input[2:])
	}
}

// OneOf parses a single character from the given set of characters.
func OneOf[Input Bytes](collection ...rune) Parser[Input, rune] {
	return func(input Input) Result[rune, Input] {
		if len(input) == 0 {
			return Failure[Input, rune](NewError(input, "OneOf"), input)
		}

		for _, c := range collection {
			if rune(input[0]) == c {
				return Success(rune(input[0]), input[1:])
			}
		}

		return Failure[Input, rune](NewError(input, "OneOf"), input)
	}
}

// Satisfy parses a single character, and ensures that it satisfies the given predicate.
func Satisfy[Input Bytes](predicate func(rune) bool) Parser[Input, rune] {
	return func(input Input) Result[rune, Input] {
		if len(input) == 0 {
			return Failure[Input, rune](NewError(input, "Satisfy"), input)
		}

		if !predicate(rune(input[0])) {
			return Failure[Input, rune](NewError(input, "Satisfy"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// Space parses a space character.
func Space[Input Bytes]() Parser[Input, rune] {
	return func(input Input) Result[rune, Input] {
		if len(input) == 0 || input[0] != ' ' {
			return Failure[Input, rune](NewError(input, "Space"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// Tab parses a tab character.
func Tab[Input Bytes]() Parser[Input, rune] {
	return func(input Input) Result[rune, Input] {
		if len(input) == 0 || input[0] != '\t' {
			return Failure[Input, rune](NewError(input, "Tab"), input)
		}

		return Success(rune(input[0]), input[1:])
	}
}

// Int64 parses an integer from the input, and returns the part of the input that
// matched the integer.
func Int64[Input Bytes]() Parser[Input, int64] {
	return func(input Input) Result[int64, Input] {
		parser := Recognize(Sequence(Optional(Token[Input]("-")), Digit1[Input]()))

		result := parser(input)
		if result.Err != nil {
			return Failure[Input, int64](NewError(input, "Int64"), input)
		}

		n, err := strconv.ParseInt(string(result.Output), 10, 64)
		if err != nil {
			return Failure[Input, int64](NewError(input, "Int64"), input)
		}

		return Success(n, result.Remaining)
	}
}

// Int8 parses an 8-bit integer from the input,
// and returns the part of the input that matched the integer.
func Int8[Input Bytes]() Parser[Input, int8] {
	return func(input Input) Result[int8, Input] {
		parser := Recognize(Sequence(Optional(Token[Input]("-")), Digit1[Input]()))

		result := parser(input)
		if result.Err != nil {
			return Failure[Input, int8](NewError(input, "Int8"), input)
		}

		n, err := strconv.ParseInt(string(result.Output), 10, 8)
		if err != nil {
			return Failure[Input, int8](NewError(input, "Int8"), input)
		}

		return Success(int8(n), result.Remaining)
	}
}

// UInt8 parses an 8-bit integer from the input,
// and returns the part of the input that matched the integer.
func UInt8[Input Bytes]() Parser[Input, uint8] {
	return func(input Input) Result[uint8, Input] {
		result := Digit1[Input]()(input)
		if result.Err != nil {
			return Failure[Input, uint8](NewError(input, "UInt8"), input)
		}

		n, err := strconv.ParseUint(string(result.Output), 10, 8)
		if err != nil {
			return Failure[Input, uint8](NewError(input, "UInt8"), input)
		}

		return Success(uint8(n), result.Remaining)
	}
}

// IsAlpha returns true if the rune is an alphabetic character.
func IsAlpha(c rune) bool {
	return IsLowAlpha(c) || IsUpAlpha(c)
}

// IsLowAlpha returns true if the rune is a lowercase alphabetic character.
func IsLowAlpha(c rune) bool {
	return c >= 'a' && c <= 'z'
}

// IsUpAlpha returns true if the rune is an uppercase alphabetic character.
func IsUpAlpha(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

// IsDigit returns true if the rune is a digit.
func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// IsAlphanumeric returns true if the rune is an alphanumeric character.
func IsAlphanumeric(c rune) bool {
	return IsAlpha(c) || IsDigit(c)
}

// IsHexDigit returns true if the rune is a hexadecimal digit.
func IsHexDigit(c rune) bool {
	return IsDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// IsControl returns true if the rune is a control character.
func IsControl(c rune) bool {
	return (c >= 0 && c < 32) || c == 127
}

func IsWhitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
