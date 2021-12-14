/*
* Copyright (c) 2020 Ashley Jeffs
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in
* all copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
* THE SOFTWARE.
 */
/*
*
* k6 - a next-generation load testing tool
* Copyright (C) 2021 Load Impact
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU Affero General Public License as
* published by the Free Software Foundation, either version 3 of the
* License, or (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU Affero General Public License for more details.
*
* You should have received a copy of the GNU Affero General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*
 */

// Package combinators implements a minimalistic parser combinators library.
//
// N.B: The code in this package is mostly either copied, or very inspired by
// Jeffhail Benthos bloblang's parser combinator code: https://tinyurl.com/6duh2yfe.
// Some of the APIs were modifed, some functions and methods have been added,
// but all the credit goes to Jeff for this combinator library.
package combinators

import (
	"fmt"
	"strconv"
)

// Parser is the common signature of a parser function
type Parser func(input []rune) Result

// Result represents the result of a parser given an input
type Result struct {
	Payload   interface{}
	Err       *Error
	Remaining []rune
}

// Success creates a Result with a payload set from
// the result of a successful parsing.
func Success(payload interface{}, remaining []rune) Result {
	return Result{payload, nil, remaining}
}

// Failure creates a Result with an error set from
// the result of a failed parsing.
func Failure(err *Error, input []rune) Result {
	return Result{nil, err, input}
}

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
		if len(input) != 2 || (input[0] != '\r' || input[1] != '\n') {
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

		switch payload := res.Payload.(type) {
		case rune:
			return Success(string(payload), res.Remaining)
		default:
			return Success(res.Payload, res.Remaining)
		}
	}
}

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

// Float parses a sequence of numerical characters into a float64.
// The '.' character is used as the optional decimal delimiter. Any
// number without a decimal part will still be parsed as a float64.
//
// N.B: it is not the parser's role to make sure the floating point
// number you're attempting to parse fits into a 64 bits float.
func Float() Parser {
	digitsParser := TakeWhileOneOf([]rune("0123456789")...)
	minusParser := Char('-')
	dotParser := Char('.')

	return func(input []rune) Result {
		var negative bool

		result := minusParser(input)
		if result.Err == nil {
			negative = true
		}

		result = Expect(digitsParser, "digits")(result.Remaining)
		if result.Err != nil {
			return result
		}

		parsed, ok := result.Payload.(string)
		if !ok {
			err := fmt.Errorf("failed parsing floating point value; " +
				"reason: converting Float() parser result's payload to string failed",
			)
			return Failure(NewFatalError(input, err, "float"), input)
		}
		if resultTest := dotParser(result.Remaining); resultTest.Err == nil {
			if resultTest = digitsParser(resultTest.Remaining); resultTest.Err == nil {
				parsed = parsed + "." + resultTest.Payload.(string)
				result = resultTest
			}
		}

		floatingPointValue, err := strconv.ParseFloat(parsed, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse '%v' as float; reason: %w", parsed, err)
			return Failure(NewFatalError(input, err), input)
		}

		if negative {
			floatingPointValue = -floatingPointValue
		}

		result.Payload = floatingPointValue

		return result
	}
}

// TakeWhileOneOf parses any number of characters present in the
// provided collection of runes.
func TakeWhileOneOf(collection ...rune) Parser {
	index := make(map[rune]struct{}, len(collection))

	for _, r := range collection {
		index[r] = struct{}{}
	}

	expected := fmt.Sprintf("chars(%v)", string(collection))

	return func(input []rune) Result {
		if len(input) == 0 {
			return Failure(NewError(input, expected), input)
		}

		pos := 0
		for ; pos < len(input); pos++ {
			_, exists := index[input[pos]]
			if !exists {
				if pos == 0 {
					return Failure(NewError(input, expected), input)
				}

				break
			}
		}

		return Success(string(input[:pos]), input[pos:])
	}
}

// Optional applies a an optional child parser. Will return nil
// if not successful.
//
// N.B: unless a FatalError is encountered, Optional will ignore
// any parsing failures and errors.
func Optional(p Parser) Parser {
	return func(input []rune) Result {
		result := p(input)
		if result.Err != nil && !result.Err.IsFatal() {
			result.Err = nil
		}

		return result
	}
}

// Alternative applies a list of parsers one by one until one succeeds.
func Alternative(parsers ...Parser) Parser {
	return func(input []rune) Result {
		var err *Error

		for _, p := range parsers {
			res := p(input)
			if res.Err == nil || res.Err.IsFatal() {
				return res
			}

			if err == nil || len(err.Input) > len(res.Err.Input) {
				err = res.Err
			} else if len(err.Input) == len(res.Err.Input) {
				err.Add(res.Err)
			}
		}

		return Failure(err, input)
	}
}

// Expect applies a parser and, if an error is returned, the list of expected
// candidates is replaced with the given strings. This is useful for providing
// better context to users.
func Expect(p Parser, expected ...string) Parser {
	return func(input []rune) Result {
		res := p(input)
		if res.Err != nil && !res.Err.IsFatal() {
			res.Err.Expected = expected
		}

		return res
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

// DiscardAll effectively applies a parser and discards its result (Payload),
// effectively only returning the remaining.
func DiscardAll(parser Parser) Parser {
	return func(input []rune) Result {
		res := parser(input)

		for res.Err == nil {
			res = parser(res.Remaining)
		}

		res.Payload = nil
		res.Err = nil

		return res
	}
}

// Sequence applies a sequence of parsers and returns either a
// slice of results or an error if any parser fails.
func Sequence(parsers ...Parser) Parser {
	return func(input []rune) Result {
		results := make([]interface{}, 0, len(parsers))
		res := Result{Remaining: input}

		for _, parser := range parsers {
			if res = parser(res.Remaining); res.Err != nil {
				return Failure(res.Err, input)
			}

			if res.Payload != nil {
				results = append(results, res.Payload)
			}
		}

		return Success(results, res.Remaining)
	}
}

// Preceded parses and discards a result from the prefix parser. It
// then parses a result from the main parser and returns its result.
//
// Preceded is effectively equivalent to applying DiscardAll(prefix),
// and then applying the main parser.
func Preceded(prefix, parser Parser) Parser {
	return func(input []rune) Result {
		prefixResult := prefix(input)
		if prefixResult.Err != nil {
			return prefixResult
		}

		result := parser(prefixResult.Remaining)
		if result.Err != nil {
			return Failure(result.Err, input)
		}

		return Success(result.Payload, result.Remaining)
	}
}

// Terminated parses a result from the main parser, it then
// parses the result from the suffix parser and discards it; only
// returning the result of the main parser.
func Terminated(parser, suffix Parser) Parser {
	return func(input []rune) Result {
		result := parser(input)
		if result.Err != nil {
			return result
		}

		suffixResult := suffix(result.Remaining)
		if suffixResult.Err != nil {
			return Failure(suffixResult.Err, input)
		}

		return Success(result.Payload, suffixResult.Remaining)
	}
}

// Delimited parses and discards the result from the prefix parser, then
// parses the result of the main parser, and finally parses and discards
// the result of the suffix parser.
func Delimited(prefix, parser, suffix Parser) Parser {
	return func(input []rune) Result {
		return Terminated(Preceded(prefix, parser), suffix)(input)
	}
}
