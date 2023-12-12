<p align="center"><img src="logo.png" alt="motus logo"/></p>
<h1 align="center">A parser combinator library for Go</h1>

<p align="center">
    <a href="https://choosealicense.com/licenses/mit/"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License"></a>
    <a href="https://github.com/oleiade/gomme/actions/workflows/go.yml"><img src="https://github.com/oleiade/gomme/actions/workflows/go.yml/badge.svg" alt="Build Status"></a>
    <a href="https://pkg.go.dev/github.com/oleiade/gomme#pkg-types"><img src="https://pkg.go.dev/badge/github.com/oleiade/gomme#pkg-types.svg" alt="Go Documentation"></a>
    <a href="https://goreportcard.com/report/github.com/oleiade/gomme"><img src="https://goreportcard.com/badge/github.com/oleiade/gomme" alt="Go Report Card"></a>
    <a href="https://img.shields.io/github/go-mod/go-version/oleiade/gomme" alt="Go Version">
</p>

Gomme is a library that simplifies building parsers in Go.

Inspired by Rust's renowned `nom` crate, Gomme provides a developer-friendly toolkit that allows you to quickly and easily create reliable parsers for both textual and binary formats.

With the power of Go's newly introduced Generics, Gomme gives you the flexibility to design your own parsers while ensuring optimal compile-time type safety. Whether you're a seasoned developer or just starting out, Gomme is designed to make the process of building parsers efficient, enjoyable, and less intimidating.

## Table of content

<!-- toc -->
- [Table of content](#table-of-content)
- [Getting started](#getting-started)
- [Why Gomme?](#why-gomme)
- [Examples](#examples)
- [Documentation](#documentation)
- [Table of content](#table-of-content-1)
- [Documentation](#documentation-1)
- [Installation](#installation)
- [Guide](#guide)
  - [List of combinators](#list-of-combinators)
    - [Base combinators](#base-combinators)
    - [Bytes combinators](#bytes-combinators)
    - [Character combinators](#character-combinators)
    - [Combinators for Sequences](#combinators-for-sequences)
    - [Combinators for Applying Parsers Many Times](#combinators-for-applying-parsers-many-times)
    - [Combinators for Choices](#combinators-for-choices)
- [Installation](#installation-1)
- [Frequently asked questions](#frequently-asked-questions)
  - [Q: What are parser combinators?](#q-what-are-parser-combinators)
  - [Q: Why would I use parser combinators instead of a specific parser?](#q-why-would-i-use-parser-combinators-instead-of-a-specific-parser)
  - [Q: Where can I learn more about parser combinators?](#q-where-can-i-learn-more-about-parser-combinators)
- [Acknowledgements](#acknowledgements)
- [Authors](#authors)


## Getting started

Here's how to quickly parse [hexadecimal color codes](https://developer.mozilla.org/en-US/docs/Web/CSS/color) using Gomme:

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

It's as simple as that! Feel free to explore more in the [examples](examples/) directory.

## Why Gomme?

While it's true that learning parser combinators might initially seem daunting, their power, flexibility, and efficiency make them an invaluable tool for parsing textual and binary formats. We've created Gomme with a focus on making this learning curve as smooth as possible, providing clear documentation and a wide array of examples.

Once you get the hang of it, you'll find that Gomme's parser combinators are intuitive, adaptable, and perfect for quickly building parsers for various formats. They're easy to test and maintain, and they can help you create parsers that are as fast as their hand-written counterparts.

## Examples

See Gomme in action with these handy examples:
- [Parsing a simple CSV file](./examples/csv)
- [Parsing Redis' RESP protocol](./examples/redis)
- [Parsing hexadecimal color codes](./examples/hexcolor)

## Documentation

For more detailled information, refer to the official [documentation](https://pkg.go.dev/github.com/oleiade/gomme).
## Table of content

## Documentation

[Documentation](https://pkg.go.dev/github.com/oleiade/gomme)

## Installation

```bash
go get github.com/oleiade/gomme
```

## Guide

In this guide, we provide a detailed overview of the various combinators available in Gomme. Combinators are fundamental building blocks in parser construction, each designed for a specific task. By combining them, you can create complex parsers suited to your specific needs. For each combinator, we've provided a brief description and a usage example. Let's explore!

### List of combinators

#### Base combinators

| Combinator                                                           | Description                                                                                                                                                                                                             | Example                                                            |
| :------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------- |
| [`Map`](https://pkg.go.dev/github.com/oleiade/gomme#Map)             | Applies a function to the result of the provided parser, allowing you to transform the parser's result. | `Map(Digit1(), func(s string)int { return 123 })`                  |
| [`Optional`](https://pkg.go.dev/github.com/oleiade/gomme#Optional)   | Makes a parser optional. If unsuccessful, the parser returns a nil `Result.Output`.Output`.                                                                                                                         | `Optional(CRLF())`                                                 |
| [`Peek`](https://pkg.go.dev/github.com/oleiade/gomme#Peek)           | Applies the provided parser without consuming the input.                                                                                                                                                               |                                                                    |
| [`Recognize`](https://pkg.go.dev/github.com/oleiade/gomme#Recognize) | Returns the consumed input as the produced value when the provided parser is successful.                                                                                                                              | `Recognize(SeparatedPair(Token("key"), Char(':'), Token("value"))` |
| [`Assign`](https://pkg.go.dev/github.com/oleiade/gomme#Assign)       | Returns the assigned value when the provided parser is successful.                                                                                                                                                   | `Assign(true, Token("true"))`                                      |

#### Bytes combinators

| Combinator                                                               | Description                                                                                                                                                                                                        | Example                               |
| :----------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------ |
| [`Take`](https://pkg.go.dev/github.com/oleiade/gomme#Take)               | Parses the first N elements of the input.                                                                                                                                                                               | `Take(5)`                             |
| [`TakeUntil`](https://pkg.go.dev/github.com/oleiade/gomme#TakeUntil)     | Parses the input until the provided parser argument succeeds.                                                                                                                                                     | `TakeUntil(CRLF()))`                  |
| [`TakeWhileMN`](https://pkg.go.dev/github.com/oleiade/gomme#TakeWhileMN) | Parses the longest input slice fitting the length expectation (m <= input length <= n) and matching the predicate. The parser argument is a function taking a `rune` as input and returning a `bool`. | `TakeWhileMN(2, 6, gomme.isHexDigit)` |
| [`Token`](https://pkg.go.dev/github.com/oleiade/gomme#Token)             | Recognizes a specific pattern. Compares the input with the token's argument and returns the matching part.                                                                                                   | `Token("tolkien")`                    |

#### Character combinators

| Combinator | Description | Example |
| :--- | :--- | :--- |
| [`Char`](https://pkg.go.dev/github.com/oleiade/gomme#Char) | Parses a single instance of a provided character. | `Char('$')` |
| [`AnyChar`](https://pkg.go.dev/github.com/oleiade/gomme#AnyChar) | Parses a single instance of any character. | `AnyChar()` |
| [`Alpha0`](https://pkg.go.dev/github.com/oleiade/gomme#Alpha0) | Parses zero or more alphabetical ASCII characters (case insensitive). | `Alpha0()` |
| [`Alpha1`](https://pkg.go.dev/github.com/oleiade/gomme#Alpha1) | Parses one or more alphabetical ASCII characters (case insensitive). | `Alpha1()` |
| [`Alphanumeric0`](https://pkg.go.dev/github.com/oleiade/gomme#Alphanumeric0) | Parses zero or more alphabetical and numerical ASCII characters (case insensitive). | `Alphanumeric0()` |
| [`Alphanumeric1`](https://pkg.go.dev/github.com/oleiade/gomme#Alphanumeric1) | Parses one or more alphabetical and numerical ASCII characters (case insensitive). | `Alphanumeric1()` |
| [`Digit0`](https://pkg.go.dev/github.com/oleiade/gomme#Digit0) | Parses zero or more numerical ASCII characters: 0-9. | `Digit0()` |
| [`Digit1`](https://pkg.go.dev/github.com/oleiade/gomme#Digit1) | Parses one or more numerical ASCII characters: 0-9. | `Digit1()` |
| [`HexDigit0`](https://pkg.go.dev/github.com/oleiade/gomme#HexDigit0) | Parses zero or more hexadecimal ASCII characters (case insensitive). | `HexDigit0()` |
| [`HexDigit1`](https://pkg.go.dev/github.com/oleiade/gomme#HexDigit1) | Parses one or more hexadecimal ASCII characters (case insensitive). | `HexDigit1()` |
| [`Whitespace0`](https://pkg.go.dev/github.com/oleiade/gomme#Whitespace0) | Parses zero or more whitespace ASCII characters: space, tab, carriage return, line feed. | `Whitespace0()` |
| [`Whitespace1`](https://pkg.go.dev/github.com/oleiade/gomme#Whitespace1) | Parses one or more whitespace ASCII characters: space, tab, carriage return, line feed. | `Whitespace1()` |
| [`LF`](https://pkg.go.dev/github.com/oleiade/gomme#LF) | Parses a single new line character '\n'. | `LF()` |
| [`CRLF`](https://pkg.go.dev/github.com/oleiade/gomme#CRLF) | Parses a '\r\n' string. | `CRLF()` |
| [`OneOf`](https://pkg.go.dev/github.com/oleiade/gomme#OneOf) | Parses one of the provided characters. Equivalent to using `Alternative` over a series of `Char` parsers. | `OneOf('a', 'b' , 'c')` |
| [`Satisfy`](https://pkg.go.dev/github.com/oleiade/gomme#Satisfy) | Parses a single character, asserting that it matches the provided predicate. The predicate function takes a `rune` as input and returns a `bool`. `Satisfy` is useful for building custom character matchers. | `Satisfy(func(c rune)bool { return c == '{' || c == '[' })` |
| [`Space`](https://pkg.go.dev/github.com/oleiade/gomme#Space) | Parses a single space character ' '. | `Space()` |
| [`Tab`](https://pkg.go.dev/github.com/oleiade/gomme#Tab) | Parses a single tab character '\t'. | `Tab()` |
| [`Int64`](https://pkg.go.dev/github.com/oleiade/gomme#Int64) | Parses an `int64` from its textual representation. | `Int64()` |
| [`Int8`](https://pkg.go.dev/github.com/oleiade/gomme#Int8) | Parses an `int8` from its textual representation. | `Int8()` |
| [`UInt8`](https://pkg.go.dev/github.com/oleiade/gomme#UInt8) | Parses a `uint8` from its textual representation. | `UInt8()` |

#### Combinators for Sequences

| Combinator | Description | Example |
| :--- | :--- | :--- |
| [`Preceded`](https://pkg.go.dev/github.com/oleiade/gomme#Preceded) | Applies the prefix parser and discards its result. It then applies the main parser and returns its result. It discards the prefix value. It proves useful when looking for data prefixed with a pattern. For instance, when parsing a value, prefixed with its name. | `Preceded(Token("name:"), Alpha1())` |
| [`Terminated`](https://pkg.go.dev/github.com/oleiade/gomme#Terminated) | Applies the main parser, followed by the suffix parser whom it discards the result of, and returns the result of the main parser. Note that if the suffix parser fails, the whole operation fails, regardless of the result of the main parser. It proves useful when looking for suffixed data while not interested in retaining the suffix value itself. For instance, when parsing a value followed by a control character. | `Terminated(Digit1(), LF())` |
| [`Delimited`](https://pkg.go.dev/github.com/oleiade/gomme#Delimited) | Applies the prefix parser, the main parser, followed by the suffix parser, discards the result of both the prefix and suffix parsers, and returns the result of the main parser. Note that if any of the prefix or suffix parsers fail, the whole operation fails, regardless of the result of the main parser. It proves useful when looking for data surrounded by patterns helping them identify it without retaining its value. For instance, when parsing a value, prefixed by its name and followed by a control character. | `Delimited(Tag("name:"), Digit1(), LF())` |
| [`Pair`](https://pkg.go.dev/github.com/oleiade/gomme#Pair) | Applies two parsers in a row and returns a pair container holding both their result values. | `Pair(Alpha1(), Tag("cm"))` |
| [`SeparatedPair`](https://pkg.go.dev/github.com/oleiade/gomme#SeparatedPair) | Applies a left parser, a separator parser, and a right parser discards the result of the separator parser, and returns the result of the left and right parsers as a pair container holding the result values. | `SeparatedPair(Alpha1(), Tag(":"), Alpha1())` |
| [`Sequence`](https://pkg.go.dev/github.com/oleiade/gomme#Sequence) | Applies a sequence of parsers sharing the same signature. If any of the provided parsers fail, the whole operation fails. | `Sequence(SeparatedPair(Tag("name"), Char(':'), Alpha1()), SeparatedPair(Tag("height"), Char(':'), Digit1()))` |

#### Combinators for Applying Parsers Many Times

| Combinator | Description | Example |
| :--- | :--- | :--- |
| [`Count`](https://pkg.go.dev/github.com/oleiade/gomme#Count) | Applies the provided parser `count` times. If the parser fails before it can be applied `count` times, the operation fails. It proves useful whenever one needs to parse the same pattern many times in a row. | `Count(3, OneOf('a', 'b', 'c'))` |
| [`Many0`](https://pkg.go.dev/github.com/oleiade/gomme#Many0) | Keeps applying the provided parser until it fails and returns a slice of all the results. Specifically, if the parser fails to match, `Many0` still succeeds, returning an empty slice of results. It proves useful when trying to consume a repeated pattern, regardless of whether there's any match, like when trying to parse any number of whitespaces in a row. | `Many0(Char(' '))` |
| [`Many1`](https://pkg.go.dev/github.com/oleiade/gomme#Many1) | Keeps applying the provided parser until it fails and returns a slice of all the results. If the parser fails to match at least once, `Many1` fails. It proves useful when trying to consume a repeated pattern, like any number of whitespaces in a row, ensuring that it appears at least once. | `Many1(LF())` |
| [`SeparatedList0`](https://pkg.go.dev/github.com/oleiade/gomme#SeparatedList0) |  |  |
| [`SeparatedList1`](https://pkg.go.dev/github.com/oleiade/gomme#SeparatedList1) |  |  |

#### Combinators for Choices

| Combinator | Description | Example |
| :--- | :--- | :--- |
| [`Alternative`](https://pkg.go.dev/github.com/oleiade/gomme#Alternative) | Tests a list of parsers, one by one, until one succeeds. Note that all parsers must share the same signature (`Parser[I, O]`). | `Alternative(Token("abc"), Token("123"))` |


## Installation

Add the library to your Go project with the following command:

```bash
    go get github.com/oleiade/gomme@latest
```

## Frequently asked questions

### Q: What are parser combinators?

**A**: Parser combinators offer a new way of building parsers. Instead of writing a complex parser that analyzes an entire format, you create small, simple parsers that handle the smallest units of the format. These small parsers can then be combined to build more complex parsers. It's a bit like using building blocks to construct whatever structure you want.

### Q: Why would I use parser combinators instead of a specific parser?

**A**: Parser combinators are incredibly flexible and intuitive. Once you're familiar with them, they enable you to quickly create, maintain, and modify parsers. They offer you a high degree of freedom in designing your parser and how it's used.

### Q: Where can I learn more about parser combinators?

A: Here are some resources we recommend:
- [You could have invented parser combinators](https://theorangeduck.com/page/you-could-have-invented-parser-combinators)
- [Functional Parsing](https://www.youtube.com/watch?v=dDtZLm7HIJs)
- [Building a Mapping Language in Go with Parser Combinators](https://www.youtube.com/watch?v=JiViND-bpmw)

## Acknowledgements

We can frankly take close to zero credit for this library, apart from work put into assembling the already existing elements of theory and implementation into a single autonomous project.

We've stood on the shoulders of giants to create Gomme. The library draws heavily on the extensive theoretical work done in the parser combinators space, and we owe a huge debt to Rust's [nom](https://github.com/Geal/nom) and [benthos'](https://github.com/benthosdev/benthos) blob lang implementation. Our goal was to consolidate these diverse elements into a single, easy-to-use Go library.
## Authors

- [@oleiade](https://github.com/oleiade)
