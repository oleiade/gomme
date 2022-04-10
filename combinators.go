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
package gomme

import (
	"fmt"
)

// Parser is the common signature of a parser function
type Parser func(input []rune) Result

// Result represents the result of a parser given an input
type Result struct {
	Output    interface{}
	Err       *Error
	Remaining []rune
}

// Success creates a Result with a output set from
// the result of a successful parsing.
func Success(output interface{}, remaining []rune) Result {
	return Result{output, nil, remaining}
}

// Failure creates a Result with an error set from
// the result of a failed parsing.
func Failure(err *Error, input []rune) Result {
	return Result{nil, err, input}
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

// DiscardAll effectively applies a parser and discards its result (Output),
// effectively only returning the remaining.
func DiscardAll(parser Parser) Parser {
	return func(input []rune) Result {
		res := parser(input)

		for res.Err == nil {
			res = parser(res.Remaining)
		}

		res.Output = nil
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

			if res.Output != nil {
				results = append(results, res.Output)
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

		return Success(result.Output, result.Remaining)
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

		return Success(result.Output, suffixResult.Remaining)
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
