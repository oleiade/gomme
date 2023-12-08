// Package gomme implements a parser combinator library.
// It provides a toolkit for developers to build reliable, fast, flexible, and easy-to-develop and maintain parsers
// for both textual and binary formats. It extensively uses the recent introduction of Generics in the Go programming
// language to offer flexibility in how combinators can be mixed and matched to produce the desired output while
// providing as much compile-time type safety as possible.
package gomme

// FIXME: Ideally, I would want the combinators working with sequences
// to produce somewhat detailed errors, and tell me which of the combinators failed

// Bytes is a generic type alias for string
type Bytes interface {
	string | []byte
}

// Separator is a generic type alias for separator characters
type Separator interface {
	rune | byte | string
}

// Result is a generic type alias for Result
type Result[Output any, Remaining Bytes] struct {
	Output    Output
	Err       *Error[Remaining]
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
func Failure[I Bytes, O any](err *Error[I], input I) Result[O, I] {
	var output O
	return Result[O, I]{output, err, input}
}

// Map applies a function to the result of a parser.
func Map[I Bytes, PO any, MO any](parse Parser[I, PO], fn func(PO) (MO, error)) Parser[I, MO] {
	return func(input I) Result[MO, I] {
		res := parse(input)
		if res.Err != nil {
			return Failure[I, MO](NewError(input, "Map"), input)
		}

		output, err := fn(res.Output)
		if err != nil {
			return Failure[I, MO](NewError(input, err.Error()), input)
		}

		return Success(output, res.Remaining)
	}
}

// Optional applies a an optional child parser. Will return nil
// if not successful.
//
// N.B: unless a FatalError is encountered, Optional will ignore
// any parsing failures and errors.
func Optional[I Bytes, O any](parse Parser[I, O]) Parser[I, O] {
	return func(input I) Result[O, I] {
		result := parse(input)
		if result.Err != nil && !result.Err.IsFatal() {
			result.Err = nil
		}

		return Success(result.Output, result.Remaining)
	}
}

// Peek tries to apply the provided parser without consuming any input.
// It effectively allows to look ahead in the input.
func Peek[I Bytes, O any](parse Parser[I, O]) Parser[I, O] {
	return func(input I) Result[O, I] {
		result := parse(input)
		if result.Err != nil {
			return Failure[I, O](result.Err, input)
		}

		return Success(result.Output, input)
	}
}

// Recognize returns the consumed input as the produced value when
// the provided parser succeeds.
func Recognize[I Bytes, O any](parse Parser[I, O]) Parser[I, I] {
	return func(input I) Result[I, I] {
		result := parse(input)
		if result.Err != nil {
			return Failure[I, I](result.Err, input)
		}

		return Success(input[:len(input)-len(result.Remaining)], result.Remaining)
	}
}

// Assign returns the provided value if the parser succeeds, otherwise
// it returns an error result.
func Assign[I Bytes, O1, O2 any](value O1, parse Parser[I, O2]) Parser[I, O1] {
	return func(input I) Result[O1, I] {
		result := parse(input)
		if result.Err != nil {
			return Failure[I, O1](result.Err, input)
		}

		return Success(value, result.Remaining)
	}
}
