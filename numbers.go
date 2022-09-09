package gomme

import "strconv"

func Number[I Bytes]() Parser[I, float64] {
	return func(input I) Result[float64, I] {
		parser := Recognize(
			Sequence(
				Optional(Token[I]("-")),
				Digit1[I](),
				Optional(Recognize(Pair(Token[I]("."), Digit1[I]()))),
			),
		)

		result := parser(input)
		if result.Err != nil {
			return Failure[I, float64](result.Err, input)
		}

		number, err := strconv.ParseFloat(string(result.Output), 64)
		if err != nil {
			return Failure[I, float64](NewError(input, "number"), input)
		}

		return Success(number, result.Remaining)
	}
}
