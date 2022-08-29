package gomme

import (
	"fmt"
	"strings"
)

// Take returns a subset of the input of size `count`.
func Take[I Bytes](count uint) Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 && count > 0 {
			return Failure[I, I](NewError(input, "TakeUntil"), input)
		}

		if uint(len(input)) < count {
			return Failure[I, I](NewError(input, "Take"), input)
		}

		return Success(input[:count], input[count:])
	}
}

// TakeUntil parses any number of characters until the provided parser is successful.
// If the provided parser is not successful, the parser fails, and the entire input is
// returned as the Result's Remaining.
func TakeUntil[I Bytes, O any](parse Parser[I, O]) Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewError(input, "TakeUntil"), input)
		}

		pos := 0
		for ; pos < len(input); pos++ {
			current := input[pos:]
			res := parse(current)
			if res.Err == nil {
				return Success(input[:pos], input[pos:])
			}

			continue
		}

		return Failure[I, I](NewError(input, "TakeUntil"), input)
	}
}

// TakeWhileMN returns the longest input subset that matches the predicates, within
// the boundaries of `atLeast` <= len(input) <= `atMost`.
//
// If the provided parser is not successful or the pattern is out of the
// `atLeast` <= len(input) <= `atMost` range, the parser fails, and the entire
// input is returned as the Result's Remaining.
func TakeWhileMN[I Bytes](atLeast, atMost uint, predicate func(rune) bool) Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewError(input, "TakeWhileMN"), input)
		}

		// Input is shorter than the minimum expected matching length,
		// it is thus not possible to match it within the established
		// constraints.
		if uint(len(input)) < atLeast {
			return Failure[I, I](NewError(input, "TakeWhileMN"), input)
		}

		lastValidPos := 0
		for idx := 0; idx < len(input); idx++ {
			if uint(idx) == atMost {
				break
			}

			matched := predicate(rune(input[idx]))
			if !matched {
				if uint(idx) < atLeast {
					return Failure[I, I](NewError(input, "TakeWhileMN"), input)
				}

				return Success(input[:idx], input[idx:])
			}

			lastValidPos++
		}

		return Success(input[:lastValidPos], input[lastValidPos:])
	}
}

// Token parses a token from the input, and returns the part of the input that
// matched the token.
// If the token could not be found, the parser returns an error result.
func Token[I Bytes](token string) Parser[I, I] {
	return func(input I) Result[I, I] {
		if !strings.HasPrefix(string(input), token) {
			return Failure[I, I](NewError(input, fmt.Sprintf("Token(%s)", token)), input)
		}

		return Success(input[:len(token)], input[len(token):])
	}
}
