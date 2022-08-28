package gomme

// Alternative tests a list of parsers in order, one by one, until one
// succeeds.
//
// If none of the parsers succeed, this combinator produces an error Result.
func Alternative[I Bytes, O any](parsers ...Parser[I, O]) Parser[I, O] {
	return func(input I) Result[O, I] {
		for _, parse := range parsers {
			result := parse(input)
			if result.Err == nil {
				return result
			}
		}

		return Failure[I, O](NewError(input, "Alternative"), input)
	}
}
