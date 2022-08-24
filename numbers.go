package gomme

// import "math"

// Float parses a sequence of numerical characters into a float64.
// The '.' character is used as the optional decimal delimiter. Any
// number without a decimal part will still be parsed as a float64.
//
// N.B: it is not the parser's role to make sure the floating point
// number you're attempting to parse fits into a 64 bits float.

// func Float[I Bytes]() Parser[I, float64] {
// 	digitsParser := TakeWhileOneOf[I]([]rune("0123456789")...)
// 	minusParser := Char[I]('-')
// 	dotParser := Char[I]('.')

// 	return func(input I) Result[float64, I] {
// 		var negative bool

// 		minusresult := minusParser(input)
// 		if result.Err == nil {
// 			negative = true
// 		}

// 		result = digitsParser(result.Remaining)
// 		// result = Expect(digitsParser, "digits")(result.Remaining)
// 		// if result.Err != nil {
// 		// 	return result
// 		// }

// 		parsed, ok := result.Output.(string)
// 		if !ok {
// 			err := fmt.Errorf("failed parsing floating point value; " +
// 				"reason: converting Float() parser result's output to string failed",
// 			)
// 			return Failure(NewFatalError(input, err, "float"), input)
// 		}
// 		if resultTest := dotParser(result.Remaining); resultTest.Err == nil {
// 			if resultTest = digitsParser(resultTest.Remaining); resultTest.Err == nil {
// 				parsed = parsed + "." + resultTest.Output.(string)
// 				result = resultTest
// 			}
// 		}

// 		floatingPointValue, err := strconv.ParseFloat(parsed, 64)
// 		if err != nil {
// 			err = fmt.Errorf("failed to parse '%v' as float; reason: %w", parsed, err)
// 			return Failure(NewFatalError(input, err), input)
// 		}

// 		if negative {
// 			floatingPointValue = -floatingPointValue
// 		}

// 		result.Output = floatingPointValue

// 		return result
// 	}
// }
