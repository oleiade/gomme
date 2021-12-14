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

package combinators

import (
	"fmt"
	"strings"
)

// Error represents an that has occurred whilst attempting to apply a
// parser function to a given input. A slice of abstract names should
// be provided outlining tokens or characters that were expected and not
// found at the input in order to provide a useful error message.
type Error struct {
	Input    []rune
	Err      error
	Expected []string
}

// NewError creates a parser error from the input and a list of
// expected token. This is a "passive" error, indicating this particular
// parser did not succeed, but that others parsers should be tried if
// applicable.
func NewError(input []rune, expected ...string) *Error {
	return &Error{
		input,
		nil,
		expected,
	}
}

// NewFatalError creates a parser error from the input and wraps the underlying
// error indicating this parse succeeded only partially, but as a requirement
// was not meant, the parsed input is to be considered invalid, and all parsing
// should stop.
func NewFatalError(input []rune, err error, expected ...string) *Error {
	return &Error{Input: input, Err: err, Expected: expected}
}

// Error returns a human readable error string.
func (e *Error) Error() string {
	return fmt.Sprintf("expected %v", strings.Join(e.Expected, ", "))
}

// Unwrap returns the underlying fatal error (or nil).
func (e *Error) Unwrap() error {
	return e.Err
}

// ErrorAtChar returns a human readable error string including the
// character position of the error.
func (e *Error) ErrorAtChar(fullInput []rune) string {
	char := len(fullInput) - len(e.Input)
	return fmt.Sprintf("char at position %v, %v", char+1, e.Error())
}

// IsFatal indicates this parser error should be considered
// fatal, and thus the sibling parser candidates should not
// be tried.
func (e *Error) IsFatal() bool {
	return e.Err != nil
}

// Add context from another error into this one.
func (e *Error) Add(from *Error) {
	e.Expected = append(e.Expected, from.Expected...)
	if e.Err == nil {
		e.Err = from.Err
	}
}
