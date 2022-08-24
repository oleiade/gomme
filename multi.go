package gomme

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
