package gomme

// Alternative tests a list of parsers in order, one by one, until one
// succeeds.
//
// If none of the parsers succeed, this combinator produces an error Result.
func Alternative[Input Bytes, Output any](parsers ...Parser[Input, Output]) Parser[Input, Output] {
	return func(input Input) Result[Output, Input] {
		for _, parse := range parsers {
			result := parse(input)
			if result.Err == nil {
				return result
			}
		}

		return Failure[Input, Output](NewError(input, "Alternative"), input)
	}
}
