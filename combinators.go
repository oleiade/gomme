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

// Pair applies two parsers and returns a Result containing a slice of
// size 2 as its output.
func Pair[I Bytes, LO, RO any, LP Parser[I, LO], RP Parser[I, RO]](leftParser LP, rightParser RP) Parser[I, PairContainer[LO, RO]] {
	return func(input I) Result[PairContainer[LO, RO], I] {
		leftResult := leftParser(input)
		if leftResult.Err != nil {
			return Failure[I, PairContainer[LO, RO]](NewGenericError(input, "pair"), input)
		}

		rightResult := rightParser(leftResult.Remaining)
		if rightResult.Err != nil {
			return Failure[I, PairContainer[LO, RO]](NewGenericError(input, "pair"), input)
		}

		return Success(PairContainer[LO, RO]{leftResult.Output, rightResult.Output}, rightResult.Remaining)
	}
}

// SeparatedPair applies two separated parsers and returns a Result containing a slice of
// size 2 as its output. The first element of the slice is the result of the left parser,
// and the second element is the result of the right parser. The result of the separator parser
// is discarded.
func SeparatedPair[I Bytes, LO, RO any, S Separator, LP Parser[I, LO], SP Parser[I, S], RP Parser[I, RO]](leftParser LP, separator SP, rightParser RP) Parser[I, PairContainer[LO, RO]] {
	return func(input I) Result[PairContainer[LO, RO], I] {
		leftResult := leftParser(input)
		if leftResult.Err != nil {
			return Failure[I, PairContainer[LO, RO]](NewGenericError(input, "separated pair"), input)
		}

		sepResult := separator(leftResult.Remaining)
		if sepResult.Err != nil {
			return Failure[I, PairContainer[LO, RO]](NewGenericError(input, "separated pair"), input)
		}

		rightResult := rightParser(sepResult.Remaining)
		if rightResult.Err != nil {
			return Failure[I, PairContainer[LO, RO]](NewGenericError(input, "pair"), input)
		}

		return Success(PairContainer[LO, RO]{leftResult.Output, rightResult.Output}, rightResult.Remaining)
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

// Alternative tests a list of parsers in order, one by one, until one
// succeeds.
//
// If none of the parsers succeed, this combinator produces an error Result.
func Alternative[I Bytes, O any](parsers ...Parser[I, O]) Parser[I, O] {
	return func(input I) Result[O, I] {
		for _, p := range parsers {
			res := p(input)
			if res.Err == nil {
				return res
			}
		}

		return Failure[I, O](NewGenericError(input, "alternative"), input)
	}
}

// Many applies a parser repeatedly until it fails, and returns a slice of all
// the results as the Result's Output.
func Many[I Bytes, O any](p Parser[I, O]) Parser[I, []O] {
	return func(input I) Result[[]O, I] {
		results := []O{}

		remaining := input
		for {
			res := p(remaining)
			if res.Err != nil {
				return Success(results, remaining)
			}

			results = append(results, res.Output)
			remaining = res.Remaining
		}
	}
}

// Sequence applies a sequence of parsers and returns either a
// slice of results or an error if any parser fails.
func Sequence[I Bytes, O any](parsers ...Parser[I, O]) Parser[I, []O] {
	return func(input I) Result[[]O, I] {
		remaining := input
		outputs := make([]O, 0, len(parsers))

		for _, parser := range parsers {
			res := parser(remaining)
			if res.Err != nil {
				return Failure[I, []O](res.Err, input)
			}

			outputs = append(outputs, res.Output)
			remaining = res.Remaining
		}

		return Success(outputs, remaining)
	}
}

// Preceded parses and discards a result from the prefix parser. It
// then parses a result from the main parser and returns its result.
//
// Preceded is effectively equivalent to applying DiscardAll(prefix),
// and then applying the main parser.
func Preceded[I Bytes, OP, O any](prefix Parser[I, OP], parser Parser[I, O]) Parser[I, O] {
	return func(input I) Result[O, I] {
		prefixResult := prefix(input)
		if prefixResult.Err != nil {
			return Failure[I, O](prefixResult.Err, input)
		}

		result := parser(prefixResult.Remaining)
		if result.Err != nil {
			return Failure[I, O](result.Err, input)
		}

		return Success(result.Output, result.Remaining)
	}
}

// Terminated parses a result from the main parser, it then
// parses the result from the suffix parser and discards it; only
// returning the result of the main parser.
func Terminated[I Bytes, O, OS any](parser Parser[I, O], suffix Parser[I, OS]) Parser[I, O] {
	return func(input I) Result[O, I] {
		result := parser(input)
		if result.Err != nil {
			return Failure[I, O](result.Err, input)
		}

		suffixResult := suffix(result.Remaining)
		if suffixResult.Err != nil {
			return Failure[I, O](suffixResult.Err, input)
		}

		return Success(result.Output, suffixResult.Remaining)
	}
}

// Delimited parses and discards the result from the prefix parser, then
// parses the result of the main parser, and finally parses and discards
// the result of the suffix parser.
func Delimited[I Bytes, OP, O, OS any](prefix Parser[I, OP], parser Parser[I, O], suffix Parser[I, OS]) Parser[I, O] {
	return func(input I) Result[O, I] {
		return Terminated(Preceded(prefix, parser), suffix)(input)
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
