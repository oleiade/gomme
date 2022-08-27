package gomme

// Delimited parses and discards the result from the prefix parser, then
// parses the result of the main parser, and finally parses and discards
// the result of the suffix parser.
func Delimited[I Bytes, OP, O, OS any](prefix Parser[I, OP], parser Parser[I, O], suffix Parser[I, OS]) Parser[I, O] {
	return func(input I) Result[O, I] {
		return Terminated(Preceded(prefix, parser), suffix)(input)
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
