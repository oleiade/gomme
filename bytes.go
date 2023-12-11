package gomme

import (
	"fmt"
	"strings"
)

// Take returns a subset of the input of size `count`.
func Take[Input Bytes](count uint) Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 && count > 0 {
			return Failure[Input, Input](NewError(input, "TakeUntil"), input)
		}

		if uint(len(input)) < count {
			return Failure[Input, Input](NewError(input, "Take"), input)
		}

		return Success(input[:count], input[count:])
	}
}

// TakeUntil parses any number of characters until the provided parser is successful.
// If the provided parser is not successful, the parser fails, and the entire input is
// returned as the Result's Remaining.
func TakeUntil[Input Bytes, Output any](parse Parser[Input, Output]) Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Failure[Input, Input](NewError(input, "TakeUntil"), input)
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

		return Failure[Input, Input](NewError(input, "TakeUntil"), input)
	}
}

// TakeWhileMN returns the longest input subset that matches the predicates, within
// the boundaries of `atLeast` <= len(input) <= `atMost`.
//
// If the provided parser is not successful or the pattern is out of the
// `atLeast` <= len(input) <= `atMost` range, the parser fails, and the entire
// input is returned as the Result's Remaining.
func TakeWhileMN[Input Bytes](atLeast, atMost uint, predicate func(rune) bool) Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if len(input) == 0 {
			return Failure[Input, Input](NewError(input, "TakeWhileMN"), input)
		}

		// Input is shorter than the minimum expected matching length,
		// it is thus not possible to match it within the established
		// constraints.
		if uint(len(input)) < atLeast {
			return Failure[Input, Input](NewError(input, "TakeWhileMN"), input)
		}

		lastValidPos := 0
		for idx := 0; idx < len(input); idx++ {
			if uint(idx) == atMost {
				break
			}

			matched := predicate(rune(input[idx]))
			if !matched {
				if uint(idx) < atLeast {
					return Failure[Input, Input](NewError(input, "TakeWhileMN"), input)
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
func Token[Input Bytes](token string) Parser[Input, Input] {
	return func(input Input) Result[Input, Input] {
		if !strings.HasPrefix(string(input), token) {
			return Failure[Input, Input](NewError(input, fmt.Sprintf("Token(%s)", token)), input)
		}

		return Success(input[:len(token)], input[len(token):])
	}
}
