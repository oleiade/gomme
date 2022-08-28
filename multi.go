package gomme

// Count runs the provided parser `count` times.
func Count[I Bytes, O any](parse Parser[I, O], count uint) Parser[I, []O] {
	return func(input I) Result[[]O, I] {
		if len(input) == 0 || count == 0 {
			return Failure[I, []O](NewError(input, "Count"), input)
		}

		outputs := make([]O, 0, int(count))
		remaining := input
		for i := 0; uint(i) < count; i++ {
			result := parse(remaining)
			if result.Err != nil {
				return Failure[I, []O](result.Err, input)
			}

			remaining = result.Remaining
			outputs = append(outputs, result.Output)
		}

		return Success(outputs, remaining)
	}
}

// Many0 applies a parser repeatedly until it fails, and returns a slice of all
// the results as the Result's Output.
//
// Note that Many0 will succeed even if the parser fails to match at all. It will
// however fail if the provided parser accepts empty inputs (such as `Digit0`, or
// `Alpha0`) in order to prevent infinite loops.
func Many0[I Bytes, O any](parse Parser[I, O]) Parser[I, []O] {
	return func(input I) Result[[]O, I] {
		results := []O{}

		remaining := input
		for {
			res := parse(remaining)
			if res.Err != nil {
				return Success(results, remaining)
			}

			// Checking for infinite loops, if nothing was consumed,
			// the provided parser would make us go around in circles.
			if len(res.Remaining) == len(remaining) {
				return Failure[I, []O](NewError(input, "Many0"), input)
			}

			results = append(results, res.Output)
			remaining = res.Remaining
		}
	}
}

// Many1 applies a parser repeatedly until it fails, and returns a slice of all
// the results as the Result's Output. Many1 will fail if the parser fails to
// match at least once.
//
// Note that Many1 will fail if the provided parser accepts empty
// inputs (such as `Digit0`, or `Alpha0`) in order to prevent infinite loops.
func Many1[I Bytes, O any](parse Parser[I, O]) Parser[I, []O] {
	return func(input I) Result[[]O, I] {
		first := parse(input)
		if first.Err != nil {
			return Failure[I, []O](first.Err, input)
		}

		// Checking for infinite loops, if nothing was consumed,
		// the provided parser would make us go around in circles.
		if len(first.Remaining) == len(input) {
			return Failure[I, []O](NewError(input, "Many1"), input)
		}

		results := []O{first.Output}
		remaining := first.Remaining

		for {
			res := parse(remaining)
			if res.Err != nil {
				return Success(results, remaining)
			}

			// Checking for infinite loops, if nothing was consumed,
			// the provided parser would make us go around in circles.
			if len(res.Remaining) == len(remaining) {
				return Failure[I, []O](NewError(input, "Many1"), input)
			}

			results = append(results, res.Output)
			remaining = res.Remaining
		}
	}
}
