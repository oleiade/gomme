package gomme

// Count runs the provided parser `count` times.
func Count[I Bytes, O any](p Parser[I, O], count uint) Parser[I, []O] {
	return func(input I) Result[[]O, I] {
		if len(input) == 0 || count == 0 {
			return Failure[I, []O](NewGenericError(input, "count"), input)
		}

		outputs := make([]O, 0, int(count))
		remaining := input
		for i := 0; uint(i) < count; i++ {
			result := p(remaining)
			if result.Err != nil {
				return Failure[I, []O](result.Err, input)
			}

			remaining = result.Remaining
			outputs = append(outputs, result.Output)
		}

		return Success(outputs, remaining)
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
