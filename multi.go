package gomme

// Count runs the provided parser `count` times.
//
// If the provided parser cannot be successfully applied `count` times, the operation
// fails and the Result will contain an error.
func Count[Input Bytes, Output any](parse Parser[Input, Output], count uint) Parser[Input, []Output] {
	return func(input Input) Result[[]Output, Input] {
		if len(input) == 0 || count == 0 {
			return Failure[Input, []Output](NewError(input, "Count"), input)
		}

		outputs := make([]Output, 0, int(count))
		remaining := input
		for i := 0; uint(i) < count; i++ {
			result := parse(remaining)
			if result.Err != nil {
				return Failure[Input, []Output](result.Err, input)
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
func Many0[Input Bytes, Output any](parse Parser[Input, Output]) Parser[Input, []Output] {
	return func(input Input) Result[[]Output, Input] {
		results := []Output{}

		remaining := input
		for {
			res := parse(remaining)
			if res.Err != nil {
				return Success(results, remaining)
			}

			// Checking for infinite loops, if nothing was consumed,
			// the provided parser would make us go around in circles.
			if len(res.Remaining) == len(remaining) {
				return Failure[Input, []Output](NewError(input, "Many0"), input)
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
func Many1[Input Bytes, Output any](parse Parser[Input, Output]) Parser[Input, []Output] {
	return func(input Input) Result[[]Output, Input] {
		first := parse(input)
		if first.Err != nil {
			return Failure[Input, []Output](first.Err, input)
		}

		// Checking for infinite loops, if nothing was consumed,
		// the provided parser would make us go around in circles.
		if len(first.Remaining) == len(input) {
			return Failure[Input, []Output](NewError(input, "Many1"), input)
		}

		results := []Output{first.Output}
		remaining := first.Remaining

		for {
			res := parse(remaining)
			if res.Err != nil {
				return Success(results, remaining)
			}

			// Checking for infinite loops, if nothing was consumed,
			// the provided parser would make us go around in circles.
			if len(res.Remaining) == len(remaining) {
				return Failure[Input, []Output](NewError(input, "Many1"), input)
			}

			results = append(results, res.Output)
			remaining = res.Remaining
		}
	}
}

// SeparatedList0 applies an element parser and a separator parser repeatedly in order
// to produce a list of elements.
//
// Note that SeparatedList0 will succeed even if the element parser fails to match at all.
// It will however fail if the provided element parser accepts empty inputs (such as
// `Digit0`, or `Alpha0`) in order to prevent infinite loops.
//
// Because the `SeparatedList0` is really looking to produce a list of elements resulting
// from the provided main parser, it will succeed even if the separator parser fails to
// match at all. It will however fail if the provided separator parser accepts empty
// inputs in order to prevent infinite loops.
func SeparatedList0[Input Bytes, Output any, S Separator](
	parse Parser[Input, Output],
	separator Parser[Input, S],
) Parser[Input, []Output] {
	return func(input Input) Result[[]Output, Input] {
		results := []Output{}

		res := parse(input)
		if res.Err != nil {
			return Success(results, input)
		}

		// Checking for infinite loops, if nothing was consumed,
		// the provided parser would make us go around in circles.
		if len(res.Remaining) == len(input) {
			return Failure[Input, []Output](NewError(input, "SeparatedList0"), input)
		}

		results = append(results, res.Output)
		remaining := res.Remaining

		for {
			separatorResult := separator(remaining)
			if separatorResult.Err != nil {
				return Success(results, remaining)
			}

			// Checking for infinite loops, if nothing was consumed,
			// the provided parser would make us go around in circles.
			if len(separatorResult.Remaining) == len(remaining) {
				return Failure[Input, []Output](NewError(input, "SeparatedList0"), input)
			}

			parserResult := parse(separatorResult.Remaining)
			if parserResult.Err != nil {
				return Success(results, remaining)
			}

			results = append(results, parserResult.Output)

			remaining = parserResult.Remaining
		}
	}
}

// SeparatedList1 applies an element parser and a separator parser repeatedly in order
// to produce a list of elements.
//
// Note that SeparatedList1 will fail if the element parser fails to match at all.
//
// Because the `SeparatedList1` is really looking to produce a list of elements resulting
// from the provided main parser, it will succeed even if the separator parser fails to
// match at all.
func SeparatedList1[Input Bytes, Output any, S Separator](
	parse Parser[Input, Output],
	separator Parser[Input, S],
) Parser[Input, []Output] {
	return func(input Input) Result[[]Output, Input] {
		results := []Output{}

		res := parse(input)
		if res.Err != nil {
			return Failure[Input, []Output](res.Err, input)
		}

		// Checking for infinite loops, if nothing was consumed,
		// the provided parser would make us go around in circles.
		if len(res.Remaining) == len(input) {
			return Failure[Input, []Output](NewError(input, "SeparatedList0"), input)
		}

		results = append(results, res.Output)
		remaining := res.Remaining

		for {
			separatorResult := separator(remaining)
			if separatorResult.Err != nil {
				return Success(results, remaining)
			}

			// Checking for infinite loops, if nothing was consumed,
			// the provided parser would make us go around in circles.
			if len(separatorResult.Remaining) == len(remaining) {
				return Failure[Input, []Output](NewError(input, "SeparatedList0"), input)
			}

			parserResult := parse(separatorResult.Remaining)
			if parserResult.Err != nil {
				return Success(results, remaining)
			}

			results = append(results, parserResult.Output)

			remaining = parserResult.Remaining
		}
	}
}
