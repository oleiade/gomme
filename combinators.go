package gomme

import "fmt"

// FIXME: Ideally, I would want the combinators working with sequences
// to produce somewhat detailed errors, and tell me which of the combinators failed

// Bytes is a generic type alias for string
type Bytes interface {
	string
}

// Separator is a generic type alias for separator characters
type Separator interface {
	rune | byte | string
}

// Result is a generic type alias for Result
type Result[Output any, Remaining Bytes] struct {
	Output    Output
	Err       *GenericError[Remaining]
	Remaining Remaining
}

// Parser is a generic type alias for Parser
type Parser[Input Bytes, Output any] func(input Input) Result[Output, Input]

// Success creates a Result with a output set from
// the result of a successful parsing.
func Success[O any, Remaining Bytes](output O, r Remaining) Result[O, Remaining] {
	return Result[O, Remaining]{output, nil, r}
}

// Failure creates a Result with an error set from
// the result of a failed parsing.
// TODO: The Error type could be generic too
func Failure[I Bytes, O any](err *GenericError[I], input I) Result[O, I] {
	var output O
	return Result[O, I]{output, err, input}
}

// TakeWhileOneOf parses any number of characters present in the
// provided collection of runes.
func TakeWhileOneOf[I Bytes](collection ...rune) Parser[I, I] {
	index := make(map[rune]struct{}, len(collection))

	for _, r := range collection {
		index[r] = struct{}{}
	}

	expected := fmt.Sprintf("chars(%v)", string(collection))

	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewGenericError(input, expected), input)
		}

		pos := 0
		for ; pos < len(input); pos++ {
			_, exists := index[rune(input[pos])]
			if !exists {
				if pos == 0 {
					return Failure[I, I](NewGenericError(input, expected), input)
				}

				break
			}
		}

		return Success(input[:pos], input[pos:])
	}
}

// TakeUntil parses any number of characters until the provided parser is successful.
// If the provided parser is not successful, the parser fails, and the entire input is
// returned as the Result's Remaining.
func TakeUntil[I Bytes, O any](p Parser[I, O]) Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewGenericError(input, "take until"), input)
		}

		pos := 0
		for ; pos < len(input); pos++ {
			current := input[pos:]
			res := p(current)
			if res.Err == nil {
				return Success(input[:pos], input[pos:])
			}

			continue
		}

		return Failure[I, I](NewGenericError(input, "take until"), input)
	}
}

// Take returns a subset of the input of size `count`.
func Take[I Bytes](count uint) Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 && count > 0 {
			return Failure[I, I](NewGenericError(input, "take until"), input)
		}

		if uint(len(input)) < count {
			return Failure[I, I](NewGenericError(input, "take"), input)
		}

		return Success(input[:count], input[count:])
	}
}

func TakeWhileMN[I Bytes](atLeast, atMost uint, predicate func(rune) bool) Parser[I, I] {
	return func(input I) Result[I, I] {
		if len(input) == 0 {
			return Failure[I, I](NewGenericError(input, "TakeWhileMN"), input)
		}

		// Input is shorter than the minimum expected matching length,
		// it is thus not possible to match it within the established
		// constraints.
		if uint(len(input)) < atLeast {
			return Failure[I, I](NewGenericError(input, "TakeWhileMN"), input)
		}

		lastValidPos := 0
		for idx, c := range input {
			if uint(idx) == atMost {
				break
			}

			matched := predicate(c)
			if !matched {
				if uint(idx) < atLeast {
					return Failure[I, I](NewGenericError(input, "TakeWhileMN"), input)
				}

				return Success(input[:idx], input[idx:])
			}

			lastValidPos++
		}

		return Success(input[:lastValidPos], input[lastValidPos:])
	}
}

// Map applies a function to the result of a parser.
func Map[I Bytes, PO any, MO any](p Parser[I, PO], fn func(PO) (MO, error)) Parser[I, MO] {
	return func(input I) Result[MO, I] {
		res := p(input)
		if res.Err != nil {
			return Failure[I, MO](NewGenericError(input, "map"), input)
		}

		output, err := fn(res.Output)
		if err != nil {
			return Failure[I, MO](NewGenericError(input, err.Error()), input)
		}

		return Success(output, res.Remaining)
	}
}

// Optional applies a an optional child parser. Will return nil
// if not successful.
//
// N.B: unless a FatalError is encountered, Optional will ignore
// any parsing failures and errors.
func Optional[I Bytes, O any](p Parser[I, O]) Parser[I, O] {
	return func(input I) Result[O, I] {
		result := p(input)
		if result.Err != nil && !result.Err.IsFatal() {
			result.Err = nil
		}

		return Success(result.Output, result.Remaining)
	}
}

// Peek tries to apply the provided parser without consuming any input.
// It effectively allows to look ahead in the input.
func Peek[I Bytes, O any](p Parser[I, O]) Parser[I, O] {
	return func(input I) Result[O, I] {
		result := p(input)
		if result.Err != nil {
			return Failure[I, O](result.Err, input)
		}

		return Success(result.Output, input)
	}
}

// Recognize returns the consumed input as the produced value when
// the provided parser succeeds.
func Recognize[I Bytes, O any](p Parser[I, O]) Parser[I, I] {
	return func(input I) Result[I, I] {
		result := p(input)
		if result.Err != nil {
			return Failure[I, I](result.Err, input)
		}

		return Success(input[:len(input)-len(result.Remaining)], result.Remaining)
	}
}

// Assign returns the provided value if the parser succeeds, otherwise
// it returns an error result.
func Assign[I Bytes, O1, O2 any](value O1, p Parser[I, O2]) Parser[I, O1] {
	return func(input I) Result[O1, I] {
		result := p(input)
		if result.Err != nil {
			return Failure[I, O1](result.Err, input)
		}

		return Success(value, result.Remaining)

	}
}
