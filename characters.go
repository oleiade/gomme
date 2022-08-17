package gomme

// Char parses a single character and matches it with
// a provided candidate.
func Char(character rune) Parser {
	return func(input []rune) Result {
		if len(input) == 0 || input[0] != character {
			return Failure(NewError(input, string(character)), input)
		}

		return Success(string(character), input[1:])
	}
}

// AnyChar parses any single character.
func AnyChar() Parser {
	return func(input []rune) Result {
		if len(input) == 0 {
			return Failure(NewError(input, "any character"), input)
		}

		return Success(input[0], input[1:])
	}
}
// Digit parses a single numerical character: 0-9.
func Digit() Parser {
	return func(input []rune) Result {
		if len(input) == 0 || (input[0] < '0' || input[0] > '9') {
			return Failure(NewError(input, "digit"), input)
		}

		// Considering runes are numerical (int32) representations
		// of UTF-8 characters, we need to subtract the actual rune
		// we find from the '0' rune in order to convert its actual
		// numerical value and store it in an int.
		return Success(input[0]-'0', input[1:])
	}
}

// Alpha parses a single lowercase and uppercase alphabetic character: a-z, A-Z
func Alpha() Parser {
	return func(input []rune) Result {
		if len(input) == 0 ||
			(input[0] < 'a' || input[0] > 'z') &&
				(input[0] < 'A' || input[0] > 'Z') {
			return Failure(NewError(input, "alpha"), input)
		}

		// Considering runes are numerical (int32) representations
		// of UTF-8 characters, we need to subtract the actual rune
		// we find from the '0' rune in order to convert its actual
		// numerical value and store it in an int.
		return Success(input[0], input[1:])
	}
}

// LF parses a line feed `\n` character.
func LF() Parser {
	return func(input []rune) Result {
		if len(input) == 0 || input[0] != '\n' {
			return Failure(NewError(input, "line feed ('\\n')"), input)
		}

		return Success(input[0], input[1:])
	}
}

// CR parses a carriage return `\r` character.
func CR() Parser {
	return func(input []rune) Result {
		if len(input) == 0 || input[0] != '\r' {
			return Failure(NewError(input, "carriage return ('\\r')"), input)
		}

		return Success(input[0], input[1:])
	}
}

// CRLF parses the string `\r\n`.
func CRLF() Parser {
	return func(input []rune) Result {
		if len(input) < 2 || (input[0] != '\r' || input[1] != '\n') {
			return Failure(NewError(input, "CRLF ('\\r\\n')"), input)
		}

		return Success(string(input[:2]), input[2:])
	}
}

// Newline parses a newline symbol: either LF (`\n`) or CRLF (`\r\n`).
func Newline() Parser {
	parser := Expect(Alternative(LF(), CRLF()), "new line")

	return func(input []rune) Result {
		res := parser(input)
		if res.Err != nil {
			return res
		}

		switch output := res.Output.(type) {
		case rune:
			return Success(string(output), res.Remaining)
		default:
			return Success(res.Output, res.Remaining)
		}
	}
}

// Space parses a space character.
func Space() Parser {
	return func(input []rune) Result {
		if len(input) == 0 || input[0] != ' ' {
			return Failure(NewError(input, "space"), input)
		}

		return Success(input[0], input[1:])
	}
}

// Tab parses a tab character.
func Tab() Parser {
	return func(input []rune) Result {
		if len(input) == 0 || input[0] != '\t' {
			return Failure(NewError(input, "tab"), input)
		}

		return Success(input[0], input[1:])
	}
}

// Whitespace parses any number of space or tab characters.
func Whitespace() Parser {
	return Expect(TakeWhileOneOf([]rune(" ")...), "whitespace")
}
