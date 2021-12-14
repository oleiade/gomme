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
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChar(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Char('(')

	// Act
	result := parser([]rune("(foo"))

	// Arrange
	assert.Equal(t, "(", result.Payload)
	assert.Equal(t, "foo", string(result.Remaining))
}

func TestCharFailsOnNotFoundChar(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Char('(')

	// Act
	result := parser([]rune("*foo"))

	// Arrange
	assert.NotNil(t, result.Err)
	assert.Nil(t, result.Err.Err) // A parsing error embeddeds no underlying errors
	assert.Equal(t, "*foo", string(result.Err.Input))
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, "(", result.Err.Expected[0])
}

func BenchmarkChar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Char('-')([]rune("-"))
	}
}

func TestDigit(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Digit()

	// Act
	results := make([]Result, 10)
	for i := 0; i < 10; i++ {
		results[i] = parser([]rune(strconv.Itoa(i)))
	}

	// Assert
	for i, result := range results {
		assert.Nil(t, result.Err, "result should not hold any error")
		assert.Equal(t, int32(i), result.Payload, "result payload should be the matching numerical character")
		assert.Equal(t, "", string(result.Remaining), "result remaining should be empty")
	}
}

func TestDigitFailsOnOutOfBoundValues(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Digit()

	// Act
	tooLowResult := parser([]rune("/"))  // ASCII character 47 < '0' (character 48)
	tooHighResult := parser([]rune(":")) // ASCII character 58 > '9' (character 57)

	// Assert
	assert.NotNil(t, tooLowResult.Err, "result should hold an error")
	assert.NotNil(t, tooHighResult.Err, "result should hold an error")
	assert.Equal(t, "/", string(tooLowResult.Remaining), "result should return the input as remaining")
	assert.Equal(t, ":", string(tooHighResult.Remaining), "result should return the input as remaining")
}

func BenchmarkDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Digit()([]rune("9"))
	}
}

func TestAlphaLowercaseAlphabetical(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Alpha()

	// Act
	alpha := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	results := make([]Result, len(alpha))
	for i, a := range alpha {
		results[i] = parser([]rune{a})
	}

	// Assert
	for i, result := range results {
		assert.Nil(t, result.Err)
		assert.Equal(t, rune(alpha[i]), result.Payload, "result payload should be the matching alphabetical character")
		assert.Equal(t, "", string(result.Remaining), "result remaining should be empty")
	}
}

func TestAlphaFailsOnOutOfBoundValues(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Alpha()

	// Act
	lowerBoundResult := parser([]rune("@"))              // ASCII character 64 < 'A' (character 65)
	intermediateLowerBoundResult := parser([]rune("["))  // ASCII character 91 > 'Z' (character 90)
	intermediateHigherBoundResult := parser([]rune("`")) // ASCII character 96 < 'a' (character 97)
	higherBoundResult := parser([]rune("{"))             // ASCII character 123 > 'z' (character 122)

	// Assert
	assert.NotNil(t, lowerBoundResult.Err, "result should hold an error")
	assert.NotNil(t, intermediateLowerBoundResult.Err, "result should hold an error")
	assert.NotNil(t, intermediateHigherBoundResult.Err, "result should hold an error")
	assert.NotNil(t, higherBoundResult.Err, "result should hold an error")
	assert.Equal(t, "@", string(lowerBoundResult.Remaining), "result should return the input as remaining")
	assert.Equal(t, "[", string(intermediateLowerBoundResult.Remaining), "result should return the input as remaining")
	assert.Equal(t, "`", string(intermediateHigherBoundResult.Remaining), "result should return the input as remaining")
	assert.Equal(t, "{", string(higherBoundResult.Remaining), "result should return the input as remaining")
}

func BenchmarkAlpha(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Alpha()([]rune("z"))
	}
}

func TestLF(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := LF()

	// Act
	result := parser([]rune("\n"))
	failingResult := parser([]rune("\r\n"))

	// assert
	assert.Nil(t, result.Err, "result shouldn't hold an error")
	assert.Equal(t, '\n', result.Payload, "result payload should be the \\n character")
	assert.Equal(t, "", string(result.Remaining), "result remaining should be empty")
	assert.NotNil(t, failingResult.Err, "result should hold an error")
	assert.Equal(t, nil, failingResult.Payload, "result's payload should be empty")
	assert.Equal(t, "\r\n", string(failingResult.Remaining), "result's remaining should contain the input")
}

func BenchmarkLF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LF()([]rune("\n"))
	}
}

func TestCR(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := CR()

	// Act
	result := parser([]rune("\r"))
	failingResult := parser([]rune("\n"))

	// assert
	assert.Nil(t, result.Err, "result shouldn't hold an error")
	assert.Equal(t, '\r', result.Payload, "result payload should be the \\r character")
	assert.Equal(t, "", string(result.Remaining), "result remaining should be empty")
	assert.NotNil(t, failingResult.Err, "result should hold an error")
	assert.Equal(t, nil, failingResult.Payload, "result's payload should be empty")
	assert.Equal(t, "\n", string(failingResult.Remaining), "result's remaining should contain the input")
}

func BenchmarkCR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CR()([]rune("\r"))
	}
}

func TestCRLF(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := CRLF()

	// Act
	result := parser([]rune("\r\n"))
	failingResult := parser([]rune("\r"))

	// assert
	assert.Nil(t, result.Err, "result shouldn't hold an error")
	assert.Equal(t, "\r\n", result.Payload, "result payload should be the \\r\\n string")
	assert.Equal(t, "", string(result.Remaining), "result remaining should be empty")
	assert.NotNil(t, failingResult.Err, "result should hold an error")
	assert.Equal(t, nil, failingResult.Payload, "result's payload should be empty")
	assert.Equal(t, "\r", string(failingResult.Remaining), "result's remaining should contain the input")
}

func BenchmarkCRLF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CRLF()([]rune("\r\n"))
	}
}

func TestNewline(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Newline()

	// Act
	lfResult := parser([]rune("\n"))
	crlfResult := parser([]rune("\r\n"))
	failingResult := parser([]rune("\r"))

	// assert
	assert.Nil(t, lfResult.Err, "result shouldn't hold an error")
	assert.Equal(t, "\n", lfResult.Payload, "result payload should be the \\r\\n string")
	assert.Equal(t, "", string(lfResult.Remaining), "result remaining should be empty")
	assert.Nil(t, crlfResult.Err, "result shouldn't hold an error")
	assert.Equal(t, "\r\n", crlfResult.Payload, "result payload should be the \\r\\n string")
	assert.Equal(t, "", string(crlfResult.Remaining), "result remaining should be empty")
	assert.NotNil(t, failingResult.Err, "result should hold an error")
	assert.Equal(t, nil, failingResult.Payload, "result's payload should be empty")
	assert.Equal(t, "\r", string(failingResult.Remaining), "result's remaining should contain the input")
}

func BenchmarkNewline_LF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Newline()([]rune("\n"))
	}
}

func BenchmarkNewline_CRLF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Newline()([]rune("\r\n"))
	}
}

func TestWhileOneOf(t *testing.T) {
	t.Parallel()

	result := TakeWhileOneOf([]rune("0123456789")...)([]rune("123abc"))
	assert.Equal(t, "123", result.Payload)
	assert.Equal(t, "abc", string(result.Remaining))
}

func BenchmarkTakeWhileOneOf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TakeWhileOneOf([]rune("0123")...)([]rune(strings.Repeat("0123", 1024)))
	}
}

func TestOptional(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Optional(Char('-'))

	// Act
	presentResult := parser([]rune("-123"))
	absentResult := parser([]rune("123"))

	// Assert
	assert.Equal(t, "-", presentResult.Payload)
	assert.Equal(t, "123", string(presentResult.Remaining))
	assert.Nil(t, presentResult.Err)

	assert.Equal(t, nil, absentResult.Payload)
	assert.Equal(t, "123", string(absentResult.Remaining))
	assert.Nil(t, absentResult.Err)
}

func BenchmarkOptional(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Optional(Char('-'))([]rune("-123"))
	}
}

func TestFloatPositiveWithDecimal(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Float()

	// Act
	result := parser([]rune("123.456"))

	// Assert
	assert.Equal(t, float64(123.456), result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestFloatNegativeWithDecimal(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Float()

	// Act
	result := parser([]rune("-123.456"))

	// Assert
	assert.Equal(t, float64(-123.456), result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestFloatPositiveWithoutDecimal(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Float()

	// Act
	result := parser([]rune("123"))

	// Assert
	assert.Equal(t, float64(123), result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestFloatNegativeWithoutDecimal(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Float()

	// Act
	result := parser([]rune("-123"))

	// Assert
	assert.Equal(t, float64(-123), result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestFloatInvalidNumberFormat(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Float()

	// Act
	result := parser([]rune("foo.123"))

	// Assert
	assert.NotNil(t, result.Err)
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, "digits", result.Err.Expected[0])
	assert.Equal(t, "foo.123", string(result.Err.Input))
}

func BenchmarkFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float()([]rune("1234567890.0987654321"))
	}
}

func TestTag(t *testing.T) {
	t.Parallel()

	result := Tag("foo")([]rune("foo bar"))
	assert.Equal(t, "foo", result.Payload)
	assert.Equal(t, " bar", string(result.Remaining))
}

func BenchmarkTag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tag("foo")([]rune("foobar"))
	}
}

func TestAlternative(t *testing.T) {
	t.Parallel()

	result := Alternative(Tag("foo"), Tag("bar"), Tag("baz"))([]rune("bar hello"))
	assert.Equal(t, "bar", result.Payload)
	assert.Equal(t, " hello", string(result.Remaining))
}

func BenchmarkAlternative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Alternative(Tag("foo"), Tag("bar"), Tag("baz"))([]rune("baz world bonjour"))
	}
}

func TestSpace(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Space()

	// Act
	spaceResult := parser([]rune(" foo"))

	// Assert
	assert.Equal(t, ' ', spaceResult.Payload)
	assert.Equal(t, "foo", string(spaceResult.Remaining))
}

func TestTab(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Tab()

	// Act
	spaceResult := parser([]rune("\tfoo"))

	// Assert
	assert.Equal(t, '\t', spaceResult.Payload)
	assert.Equal(t, "foo", string(spaceResult.Remaining))
}

func TestWhitespaceSingleWhitespace(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Whitespace()

	// Act
	spaceResult := parser([]rune(" "))

	// Assert
	assert.Equal(t, " ", spaceResult.Payload)
	assert.Equal(t, "", string(spaceResult.Remaining))
}

func TestWhitespaceMultipleWhitespaces(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Whitespace()

	// Act
	spaceResult := parser([]rune("   "))

	// Assert
	assert.Equal(t, "   ", spaceResult.Payload)
	assert.Equal(t, "", string(spaceResult.Remaining))
}

func BenchmarkWhitespace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Whitespace()([]rune("     "))
	}
}

func TestDiscardAll(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := DiscardAll(Whitespace())

	// Act
	result := parser([]rune("  "))

	// Assert
	assert.Equal(t, nil, result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestDiscardAllSequence(t *testing.T) {
	t.Parallel()

	// Arrange
	whitespace := DiscardAll(Whitespace())
	parser := Sequence(
		Char('a'),
		whitespace,
		Char('b'),
	)

	// Act
	result := parser([]rune("a  b"))

	// Assert
	assert.Equal(t, []interface{}{"a", "b"}, result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func BenchmarkDiscardAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DiscardAll(Whitespace())([]rune("   foo   "))
	}
}

func TestSequence(t *testing.T) {
	t.Parallel()

	result := Sequence(Tag("foo"), Char(' '), Tag("bar"))([]rune("foo bar"))
	assert.Equal(t, []interface{}{"foo", " ", "bar"}, result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestSequenceWithFailingSubparser(t *testing.T) {
	t.Parallel()

	result := Sequence(Tag("foo"), Char(' '), Tag("bar"))([]rune("foobar"))

	assert.NotNil(t, result.Err)
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, " ", result.Err.Expected[0])
	assert.Equal(t, "bar", string(result.Err.Input))
}

func BenchmarkSequence(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sequence(Tag("foo"), Char(' '), Tag("bar"))([]rune("foo bar"))
	}
}

func TestPreceded(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Preceded(Char('('), Tag("foo"))

	// Act
	result := parser([]rune("(foo"))

	// Assert
	assert.Nil(t, result.Err)
	assert.Equal(t, "foo", result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestPrecededFailsOnMissingDelimiterFromInput(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Preceded(Char('('), Tag("foo"))

	// Act
	result := parser([]rune("foo"))

	// Assert
	assert.NotNil(t, result.Err)
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, "(", result.Err.Expected[0])
	assert.Equal(t, "foo", string(result.Err.Input))
}

func TestPrecededFailsOnPresentDelimiterButFailingSuccessorParser(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Preceded(Char('('), Tag("foo"))

	// Act
	result := parser([]rune("(bar"))

	// Assert
	assert.NotNil(t, result.Err)
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, "foo", result.Err.Expected[0])
	assert.Equal(t, "bar", string(result.Err.Input))
}

func BenchmarkPreceded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Preceded(Char('('), Tag("foo"))([]rune("(foo)"))
	}
}

func TestTerminated(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Terminated(Tag("foo"), Char(')'))

	// Act
	result := parser([]rune("foo)"))

	// Assert
	assert.Nil(t, result.Err)
	assert.Equal(t, "foo", result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestTerminatedFailsOnPresentDelimiterButFailingSuccessorParser(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Terminated(Tag("foo"), Char(')'))

	// Act
	result := parser([]rune("bar)"))

	// Assert
	assert.NotNil(t, result.Err)
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, "foo", result.Err.Expected[0])
	assert.Equal(t, "bar)", string(result.Err.Input))
}

func TestTerminatedFailsOnMissingDelimiterFromInput(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Terminated(Tag("foo"), Char(')'))

	// Act
	result := parser([]rune("foo"))

	// Assert
	assert.NotNil(t, result.Err)
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, ")", result.Err.Expected[0])
	assert.Equal(t, "", string(result.Err.Input))
}

func BenchmarkTerminated(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Terminated(Tag("foo"), Char(')'))([]rune("foo)"))
	}
}

func TestDelimited(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Delimited(Char('('), Tag("foo"), Char(')'))

	// Act
	result := parser([]rune("(foo)"))

	// Assert
	assert.Nil(t, result.Err)
	assert.Equal(t, "foo", result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestDelimitedFailsOnMissingPrefix(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Delimited(Char('('), Tag("foo"), Char(')'))

	// Act
	result := parser([]rune("foo)"))

	// Assert
	assert.NotNil(t, result.Err)
	assert.Equal(t, "foo)", string(result.Remaining))
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, "(", result.Err.Expected[0])
	assert.Equal(t, "foo)", string(result.Err.Input))
}

func TestDelimitedFailsOnMissingMainParser(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Delimited(Char('('), Tag("foo"), Char(')'))

	// Act
	result := parser([]rune("()"))

	// Assert
	assert.NotNil(t, result.Err)
	assert.Equal(t, "()", string(result.Remaining))
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, "foo", result.Err.Expected[0])
	assert.Equal(t, ")", string(result.Err.Input))
}

func TestDelimitedFailsOnMissingSuffix(t *testing.T) {
	t.Parallel()

	// Arrange
	parser := Delimited(Char('('), Tag("foo"), Char(')'))

	// Act
	result := parser([]rune("(foo"))

	// Assert
	assert.NotNil(t, result.Err)
	assert.Equal(t, "(foo", string(result.Remaining))
	assert.Len(t, result.Err.Expected, 1)
	assert.Equal(t, ")", result.Err.Expected[0])
	assert.Equal(t, "", string(result.Err.Input))
}

func BenchmarkDelimited(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Delimited(Char('('), Tag("foo"), Char(')'))([]rune("(foo)"))
	}
}
