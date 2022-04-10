package gomme

// Tag parses a provided candidate string.
// Given a "tag" to parse, it will try to consume an exact
// match from the input.
func Tag(tag string) Parser {
	termRunes := []rune(tag)

	return func(input []rune) Result {
		if len(input) < len(termRunes) {
			return Failure(NewError(input, tag), input)
		}

		for i, c := range termRunes {
			if input[i] != c {
				return Failure(NewError(input, tag), input)
			}
		}

		return Success(tag, input[len(termRunes):])
	}
}
