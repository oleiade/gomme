# Gomme

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![Build Status](https://github.com/oleiade/gomme/actions/workflows/go.yml/badge.svg)](https://github.com/oleiade/gomme/actions/workflows/go.yml)
[![Go Documentation](https://pkg.go.dev/badge/github.com/oleiade/gomme#pkg-types.svg)](https://pkg.go.dev/github.com/oleiade/gomme#pkg-types)
[![Go Report Card](https://goreportcard.com/badge/github.com/oleiade/gomme)](https://goreportcard.com/report/github.com/oleiade/gomme)
![Go Version](https://img.shields.io/github/go-mod/go-version/oleiade/gomme)

Gomme is a parser combinator library for the Go programming language.

It provides a toolkit for developers to build reliable, fast, flexible, easy-to-develop, and support parsers
for both textual and binary formats. It extensively uses the recent introduction of Generics in the Go programming
language. Those afford the library to offer flexibility in how users can mix and match combinators to produce the desired output while
providing as much compile-time type safety as possible.

## Table of content

<!-- toc -->- [Gomme](#gomme)
  - [Table of content](#table-of-content)
  - [Example](#example)
  - [Why would you want to use Gomme](#why-would-you-want-to-use-gomme)
  - [Table of content](#table-of-content-1)
  - [Documentation](#documentation)
    - [How to](#how-to)
    - [List of combinators](#list-of-combinators)
      - [Base combinators](#base-combinators)
      - [Bytes combinators](#bytes-combinators)
      - [Character combinators](#character-combinators)
      - [Combinators for sequences](#combinators-for-sequences)
      - [Combinators for applying parsers many times](#combinators-for-applying-parsers-many-times)
      - [Combinators for choices](#combinators-for-choices)
  - [Installation](#installation)
  - [Frequently asked questions](#frequently-asked-questions)
      - [What are parser combinators](#what-are-parser-combinators)
      - [Why would one want to use parser combinators and not write a specific parser](#why-would-one-want-to-use-parser-combinators-and-not-write-a-specific-parser)
      - [Where can I read/watch about parser combinators](#where-can-i-readwatch-about-parser-combinators)
  - [Acknowledgements](#acknowledgements)
  - [Authors](#authors)

- [Example](#example)
- [Why would you want to use Gomme?](#why-would-you-want-to-use-gomme)
- [Documentation](#documentation)
- [Installation](#installation)
- [FAQ](#faq)
- [Acknowledgements](#acknowledgements)
- [Authors](#authors)


## Example

How to parse [hexadecimal color codes](https://developer.mozilla.org/en-US/docs/Web/CSS/color), using the Gomme library:

```golang
// RGBColor stores the three bytes describing a color in the RGB space.
type RGBColor struct {
    red   uint8
    green uint8
    blue  uint8
}

// ParseRGBColor creates a new RGBColor from a hexadecimal color string.
// The string must be a six-digit hexadecimal number, prefixed with a "#".
func ParseRGBColor(input string) (RGBColor, error) {
    parser := gomme.Preceded(
        gomme.Token[string]("#"),
        gomme.Map(
            gomme.Count(HexColorComponent(), 3),
            func(components []uint8) (RGBColor, error) {
                return RGBColor{components[0], components[1], components[2]}, nil
            },
        ),
    )

    result := parser(input)
    if result.Err != nil {
        return RGBColor{}, result.Err
    }

    return result.Output, nil
}

// HexColorComponent produces a parser that parses a single hex color component,
// which is a two-digit hexadecimal number.
func HexColorComponent() gomme.Parser[string, uint8] {
    return func(input string) gomme.Result[uint8, string] {
        return gomme.Map(
            gomme.TakeWhileMN[string](2, 2, gomme.IsHexDigit),
            fromHex,
        )(input)
    }
}

// fromHex converts two digits hexadecimal numbers to their decimal value.
func fromHex(input string) (uint8, error) {
    res, err := strconv.ParseInt(input, 16, 16)
    if err != nil {
        return 0, err
    }

    return uint8(res), nil
}

```

Find more usage examples in the [examples](./examples) directory:
- [Parsing a simple CSV file](./examples/csv)
- [Parsing Redis' RESP protocol](./examples/redis)
- [Parsing hexadecimal color codes](./examples/hexcolor)

## Why would you want to use Gomme

Parser combinators arguably come with a steep learning curve, but they're a potent tool for parsing textual and binary formats. The benefits of parser combinators outweigh the cost of learning them, hence this project. We believe most of the cost of learning them comes from the lack of good documentation and examples. As a result, this project also emphasizes trying to offer comprehensive documentation and a large set of examples.

In practice and with experience, parser combinators prove intuitive and flexible and can be used to build parsers for various formats relatively quickly. They're also straightforward to test and can be used to create easy-to-maintain and extend parsers. Parser combinators also tend to be very fast and can be used to build parsers turning out as fast as hand-written ones.
## Table of content

## Documentation

[Documentation](https://pkg.go.dev/github.com/oleiade/gomme)

### How to

### List of combinators

#### Base combinators

| Combinator                                                           | Description                                                                                                                                                                                                             | Example                                                            |
| :------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------- |
| [`Map`](https://pkg.go.dev/github.com/oleiade/gomme#Map)             | Map applies a function to the result of the provided parser. Use `Map` to apply transformations to a parser result. | `Map(Digit1(), func(s string)int { return 123 })`                  |
| [`Optional`](https://pkg.go.dev/github.com/oleiade/gomme#Optional)   | Makes a parser optional. If unsuccessful, the parser returns a nil `Result.Output`.                                                                                                                         | `Optional(CRLF())`                                                 |
| [`Peek`](https://pkg.go.dev/github.com/oleiade/gomme#Peek)           | Applies the provided parser without consuming the input.                                                                                                                                                                |                                                                    |
| [`Recognize`](https://pkg.go.dev/github.com/oleiade/gomme#Recognize) | Returns the consumed input as the produced value when the provided parser is successful.                                                                                                                                | `Recognize(SeparatedPair(Token("key"), Char(':'), Token("value"))` |
| [`Assign`](https://pkg.go.dev/github.com/oleiade/gomme#Assign)       | Returns the assigned value, when the provided parser is successful.                                                                                                                                                     | `Assign(true, Token("true"))`                                      |

#### Bytes combinators

| Combinator                                                               | Description                                                                                                                                                                                                        | Example                               |
| :----------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------ |
| [`Take`](https://pkg.go.dev/github.com/oleiade/gomme#Take)               | Parses the N first element of input.                                                                                                                                                                               | `Take(5)`                             |
| [`TakeUntil`](https://pkg.go.dev/github.com/oleiade/gomme#TakeUntil)     | Parses the input until the provided parser argument succeeds.                                                                                                                                                      | `TakeUntil(CRLF()))`                  |
| [`TakeWhileMN`](https://pkg.go.dev/github.com/oleiade/gomme#TakeWhileMN) | Parses the longest input slice fitting the length expectation (m <= input length <= n) and matching the predicate. Its parser argument is a function taking a rune as input and returning a `bool`. | `TakeWhileMN(2, 6, gomme.isHexDigit)` |
| [`Token`](https://pkg.go.dev/github.com/oleiade/gomme#Token)             | Recognizes a specific patter. Compares the input with the token's argument and returns the matching part.                                                                                                   | `Token("tolkien")`                    |

#### Character combinators

| Combinator                                                                   | Description                                                                                                                                                                                                                    | Example                                              |
| :--------------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------- |
| [`Char`](https://pkg.go.dev/github.com/oleiade/gomme#Char)                   | Parses a single instance of a provided character.                                                                                                                                                                              | `Char('$)`                                           |
| [`AnyChar`](https://pkg.go.dev/github.com/oleiade/gomme#AnyChar)             | Parses a single instance of any character.                                                                                                                                                                                     | `AnyChar()`                                          |
| [`Alpha0`](https://pkg.go.dev/github.com/oleiade/gomme#Alpha0)               | Parses zero or more alphabetical ASCII characters (case insensitive).                                                                                                                                                          | `Alpha0()`                                           |
| [`Alpha1`](https://pkg.go.dev/github.com/oleiade/gomme#Alpha1)               | Parses one or more alphabetical ASCII characters (case insensitive).                                                                                                                                                           | `Alpha1()`                                           |
| [`Alphanumeric0`](https://pkg.go.dev/github.com/oleiade/gomme#Alphanumeric0) | Parses zero or more alphabetical and numerical ASCII characters (case insensitive).                                                                                                                                            | `Alphanumeric0()`                                    |
| [`Alphanumeric1`](https://pkg.go.dev/github.com/oleiade/gomme#Alphanumeric1) | Parses one or more alphabetical and numerical ASCII characters (case insensitive).                                                                                                                                              | `Alphanumeric1()`                                    |
| [`Digit0`](https://pkg.go.dev/github.com/oleiade/gomme#Digit0)               | Parses zero or more numerical ASCII characters: *0-9*.                                                                                                                                                                           | `Digit0()`                                           |
| [`Digit1`](https://pkg.go.dev/github.com/oleiade/gomme#Digit1)               | Parses one or more numerical ASCII characters: *0-9*.                                                                                                                                                                            | `Digit1()`                                           |
| [`HexDigit0`](https://pkg.go.dev/github.com/oleiade/gomme#HexDigit0)         | Parses zero or more hexadecimal ASCII characters (case insensitive).                                                                                                                                                           | `HexDigit0()`                                        |
| [`HexDigit1`](https://pkg.go.dev/github.com/oleiade/gomme#HexDigit1)         | Parses one or more hexadecimal ASCII characters (case insensitive).                                                                                                                                                            | `HexDigit1()`                                        |
| [`Whitespace0`](https://pkg.go.dev/github.com/oleiade/gomme#Whitespace0)     | Parses zero or more whitespace ASCII characters: *space, tab, carriage return, line feed. | `Whitespace0()` |
| [`Whitespace1`](https://pkg.go.dev/github.com/oleiade/gomme#Whitespace1)     | Parses one or more whitespace ASCII characters: *space, tab, carriage return, line feed. | `Whitespace1()` |
| [`LF`](https://pkg.go.dev/github.com/oleiade/gomme#LF)                       | Parses a single new line character '\n'.                                                                                                                                                                                       | `LF()`                                               |
| [`CRLF`](https://pkg.go.dev/github.com/oleiade/gomme#CRLF)                   | Parses a '\r\n' string.                                                                                                                                                                                                        | `CRLF()`                                             |
| [`OneOf`](https://pkg.go.dev/github.com/oleiade/gomme#OneOf)                 | Parses one of the provided characters. Same as using `Alternative` over a bunch of `Char` parsers.                                                                                                                 | `OneOf('a', 'b' , 'c')`                              |
| [`Satisfy`](https://pkg.go.dev/github.com/oleiade/gomme#Satisfy)             | Parses a single character, and asserts that it matches the provided predicate. The predicate function takes a `rune` as input and returns a `bool`. `Satisfy` proves useful to build custom character matchers. | `Satisfy(func(c rune)bool { c == '{' || c == '[' })` |
| [`Space`](https://pkg.go.dev/github.com/oleiade/gomme#Space)                 | Parses a single space character ' '.                                                                                                                                                                                           | `Space()`                                            |
| [`Tab`](https://pkg.go.dev/github.com/oleiade/gomme#Tab)                     | Parses a single tab character '\t'.                                                                                                                                                                                            | `Tab()`                                              |
| [`Int64`](https://pkg.go.dev/github.com/oleiade/gomme#Int64)                 | Parses an `int64` from its textual representation.                                                                                                                                                                             | `Int64()`                                            |
| [`Int8`](https://pkg.go.dev/github.com/oleiade/gomme#Int8)                   | Parses an `int8` from its textual representation.                                                                                                                                                                              | `Int8()`                                             |
| [`UInt8`](https://pkg.go.dev/github.com/oleiade/gomme#UInt8)                 | Parses a `uint8` from its textual representation.                                                                                                                                                                              | `UInt8()`                                            |

#### Combinators for sequences

| Combinator                                                                   | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | Example                                                                                                        |
| :--------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------- |
| [`Preceded`](https://pkg.go.dev/github.com/oleiade/gomme#Preceded)           | Applies the prefix parser and discards its result. It then applies the main parser and returns its result. It discards the prefix value. It proves useful when looking for data prefixed with a pattern. For instance, when parsing a value, prefixed with its name.                                                                                                                                                                                                                                                  | `Preceded(Token("name:"), Alpha1())`                                                                           |
| [`Terminated`](https://pkg.go.dev/github.com/oleiade/gomme#Terminated)       | Applies the main parser, followed by the suffix parser whom it discards the result of, and returns the result of the main parser. Note that if the suffix parser fails, the whole operation fails, regardless of the result of the main parser. It proves useful when looking for suffixed data while not interested in retaining the suffix value itself. For instance, when parsing a value followed by a control character.                                                                                                 | `Terminated(Digit1(), LF())`                                                                                   |
| [`Delimited`](https://pkg.go.dev/github.com/oleiade/gomme#Delimited)         | Applies the prefix parser, the main parser, followed by the suffix parser, discards the result of both the prefix and suffix parsers, and returns the result of the main parser. Note that if any of the prefix or suffix parsers fail, the whole operation fails, regardless of the result of the main parser. It proves useful when looking for data surrounded by patterns helping them identify it without retaining its value. For instance, when parsing a value, prefixed by its name and followed by a control character. | `Delimited(Tag("name:"), Digit1(), LF())`                                                                      |
| [`Pair`](https://pkg.go.dev/github.com/oleiade/gomme#Pair)                   | Applies two parsers in a row and returns a pair container holding both their result values.                                                                                                                                                                                                                                                                                                                                                                                                                                                              | `Pair(Alpha1(), Tag("cm"))`                                                                                    |
| [`SeparatedPair`](https://pkg.go.dev/github.com/oleiade/gomme#SeparatedPair) | Applies a left parser, a separator parser, and a right parser discards the result of the separator parser, and returns the result of the left and right parsers as a pair container holding the result values.                                                                                                                                                                                                                                                                                                                                           | `SeparatedPair(Alpha1(), Tag(":"), Alpha1())`                                                                  |
| [`Sequence`](https://pkg.go.dev/github.com/oleiade/gomme#Sequence)           | Applies a sequence of parsers sharing the same signature. If any of the provided parsers fail, the whole operation fails.                                                                                                                                                                                                                                                                                                                                                                                                                                | `Sequence(SeparatedPair(Tag("name"), Char(':'), Alpha1()), SeparatedPair(Tag("height"), Char(':'), Digit1()))` |

#### Combinators for applying parsers many times

| Combinator                                                                     | Description                                                                                                                                                                                                                                                                                                                                                                   | Example                         |
| :----------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------ |
| [`Count`](https://pkg.go.dev/github.com/oleiade/gomme#Count)                   | Applies the provided parser `count` times. If the parser fails before it can be applied `count` times, the operation fails. It proves useful whenever one needs to parse the same pattern many times in a row.                                                                                                                                        | `Count(3, OneOf('a', 'b', 'c')` |
| [`Many0`](https://pkg.go.dev/github.com/oleiade/gomme#Many0)                   | Keeps applying the provided parser until it fails and returns a slice of all the results. Specifically, if the parser fails to match, `Many0` still succeeds, returning an empty slice of results. It proves useful when trying to consume a repeated pattern, regardless of whether there's any match, like when trying to parse any number of whitespaces in a row. | `Many0(Char(' '))`              |
| [`Many1`](https://pkg.go.dev/github.com/oleiade/gomme#Many1)                   | Keeps applying the provided parser until it fails and returns a slice of all the results. If the parser fails to match at least once, `Many1` fails. It proves useful when trying to consume a repeated pattern, like any number of whitespaces in a row, ensuring that it appears at least once.                                                                     | `Many1(LF())`                   |
| [`SeparatedList0`](https://pkg.go.dev/github.com/oleiade/gomme#SeparatedList0) |                                                                                                                                                                                                                                                                                                                                                                               |                                 |
| [`SeparatedList1`](https://pkg.go.dev/github.com/oleiade/gomme#SeparatedList1) |                                                                                                                                                                                                                                                                                                                                                                               |                                 |

#### Combinators for choices

| Combinator    | Description                                                                                                                   | Example                                   |
| :------------ | :---------------------------------------------------------------------------------------------------------------------------- | :---------------------------------------- |
| [`Alternative`](https://pkg.go.dev/github.com/oleiade/gomme#Alternative) | Tests a list of parsers, one by one, until one succeeds. Note that all parsers must share the same signature (`Parser[I, O]`). | `Alternative(Token("abc"), Token("123"))` |

## Installation

Add the library to your Go project with the following command:

```bash
    go get github.com/oleiade/gomme@latest
```

## Frequently asked questions

#### What are parser combinators

Parser combinators are a programming paradigm for building parsers. They adopt a functional programming approach to parsing instead of the hand-written or generated parser. They are based on the idea of composing parsers together to build more complex parsers. Thus, Instead of writing a parser parsing a whole format (by analyzing and branching based on each character of your input) one would write a set of parsers parsing the smallest possible unit of the format, and composing them together to build more complex parsers.

A key concept to understand is that parser combinators are not parsers themselves but rather a toolkit that allows you to build parsers. This is why parser combinators are often referred to as a "parser building toolkit." Parser combinators generally are functions producing other functions ingesting some input byte by byte based on some predicate and returning a result. The result is a structure containing the output of the parser, the remaining part (once the combinator's predicate is not matched anymore, it stops and returns both what it "consumed" and what was left of the input), and an error if the parser failed to parse the input. The parser's output is the result of the parsing process and can be of any type. The error is a Go error and can be used to provide more information about the parsing failure.

#### Why would one want to use parser combinators and not write a specific parser

Parser combinators are very flexible, and once you get a good hang of them, they'll allow you to write parsers that are very easy to maintain, modify and extend very easily and very fast. They are also allegedly quite intuitive and descriptive of what the underlying data format they parse looks like. Because they're essentially a bunch of functions, generating other functions, composed in various ways depending on the need, they afford you much freedom in building your specific parser and how you want to use it.
#### Where can I read/watch about parser combinators

We recommend the following resources:
- [You could have invented parser combinators](https://theorangeduck.com/page/you-could-have-invented-parser-combinators)
- [Functional Parsing](https://www.youtube.com/watch?v=dDtZLm7HIJs)
- [Building a Mapping Language in Go with Parser Combinators](https://www.youtube.com/watch?v=JiViND-bpmw)

## Acknowledgements

We can frankly take close to zero credit for this library, apart from work put into assembling the already existing elements of theory and implementation into a single autonomous project.

This library relies heavily on the whole theoretical work done in the parser combinators space. From the implementation side of things, it was specifically started to have something similar to Rust's incredible [nom](https://github.com/Geal/nom) library in Go. This project was made possible by the pre-existing implementation of some parser combinators in [benthos'](https://github.com/benthosdev/benthos) blob lang implementation. Although the result is somewhat different, this project wouldn't have been possible without this pre-existing resource as a guiding example.

## Authors

- [@oleiade](https://github.com/oleiade)
