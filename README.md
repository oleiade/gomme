# Gomme

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![Build Status](https://github.com/oleiade/gomme/actions/workflows/go.yml/badge.svg)](https://github.com/oleiade/gomme/actions/workflows/go.yml)
[![Go Documentation](https://pkg.go.dev/badge/github.com/oleiade/gomme#pkg-types.svg)](https://pkg.go.dev/github.com/oleiade/gomme#pkg-types)
[![Go Report Card](https://goreportcard.com/badge/github.com/oleiade/gomme)](https://goreportcard.com/report/github.com/oleiade/gomme)
![Go Version](https://img.shields.io/github/go-mod/go-version/oleiade/gomme)

Gomme is a parser combinator library for the Go programming language.

It provides a toolkit for developers to build reliable, fast, flexible, and easy-to-develop and maintain parsers
for both textual and binary formats. It extensively uses the recent introduction of Generics in the Go programming
language to offer flexibility in how combinators can be mixed and matched to produce the desired output while
providing as much compile-time type safety as possible.

## Why would you want to use Gomme?

Parser combinators arguably come with a steep learning curve, but they are a potent tool for parsing textual and binary formats. We believe that the benefits of parser combinators outweigh the cost of learning them, and that's why we built Gomme. Our intuition is that most of the cost of learning them is due to the lack of good documentation and examples, and that's why we are trying to provide comprehensive documentation and a large set of examples.

In practice, we have found that parser combinators are intuitive and flexible and can be used to build parsers for various formats. They are also straightforward to test and can be used to create very easy-to-maintain and extend parsers. We have also found that parser combinators are very fast and can be used to build parsers that can turn out as fast as hand-written.
## Table of Content

<!-- toc -->

- [Example](#example)
- [Documentation](#documentation)
- [Installation](#installation)
- [FAQ](#faq)
- [Acknowledgements](#acknowledgements)
- [Authors](#authors)

## Example

Here's an example of how to parse [hexadecimal color codes](https://developer.mozilla.org/en-US/docs/Web/CSS/color), using the Gomme library:

```golang
// RGBColor stores the three bytes describing a color in the RGB space.
type RGBColor struct {
    red   uint8
    green uint8
    blue  uint8
}

// ParseRGBColor creates a new RGBColor from a hexadecimal color string.
// The string must be a six digit hexadecimal number, prefixed with a "#".
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
// which is a two digit hexadecimal number.
func HexColorComponent() gomme.Parser[string, uint8] {
    return func(input string) gomme.Result[uint8, string] {
        return gomme.Map(
            gomme.TakeWhileMN[string](2, 2, gomme.IsHexDigit),
            fromHex,
        )(input)
    }
}

// fromHex converts a two digits hexadecimal number to its decimal value.
func fromHex(input string) (uint8, error) {
    res, err := strconv.ParseInt(input, 16, 16)
    if err != nil {
        return 0, err
    }

    return uint8(res), nil
}

```

More examples can be found in the [examples](./examples) directory:
- [Parsing a simple CSV file](./examples/csv)
- [Parsing Redis' RESP protocol](./examples/redis)
- [Parsing hexadecimal color codes](./examples/hexcolor)


## Documentation

[Documentation](https://pkg.go.dev/github.com/oleiade/gomme)



## Installation

Add the library to your Go project with the following command:

```bash
    go get github.com/oleiade/gomme@latest
```

## FAQ

#### What are parser combinators?

Parser combinators are a programming paradigm for building parsers. As opposed to the hand-written or generated parser, they adopt a functional programming approach to parsing and are based on the idea of composing parsers together to build more complex parsers. What that means in practice is that instead of writing a parser that parses a whole format, by analyzing and branching based on each character of your input, you write a set of parsers that parse the smallest possible unit of the format and then compose them together to build more complex parsers.

A key concept to understand is that parser combinators are not themselves parsers but rather a toolkit that allows you to build parsers. This is why parser combinators are often referred to as a "parser building toolkit." Parser combinators generally are functions producing other functions ingesting some input byte by byte based on some predicate and returning a result. The result is a structure containing the output of the parser, the remaining part (once the combinator's predicate is not matched anymore, it stops and returns both what it "consumed" and what was left of the input), and an error if the parser failed to parse the input. The output of the parser is the result of the parsing process and can be of any type. The error is a Go error and can be used to provide more information about the parsing failure.

#### Why would I want to use parser combinators and not write my specific parser?

Parser combinators are very flexible, and once you get a good hang of them, they'll allow you to write parsers that are very easy to maintain, modify and extend very easily and very fast. They are also allegedly quite intuitive and descriptive of what the underlying data format they parse looks like. Because they're essentially a bunch of functions, generating other functions, composed in various ways depending on the need, they afford you much freedom in how you want to build your specific parser and how you want to use it.
#### Where can I read/watch about Parser Combinators?

We recommend the following resources:
- [You could have invented parser combinators](https://theorangeduck.com/page/you-could-have-invented-parser-combinators)
- [Functional Parsing](https://www.youtube.com/watch?v=dDtZLm7HIJs)
- [Building a Mapping Language in Go with Parser Combinators](https://www.youtube.com/watch?v=JiViND-bpmw)

## Acknowledgements

We can frankly take close to zero credit for this library, apart from work put into assembling the already existing elements of theory and implementation into a single autonomous project.

This library relies heavily on the whole theoretical work done in the parser combinators space. From the implementation side of things, it was specifically started to have something similar to Rust's incredible [nom](https://github.com/Geal/nom) library in Go. This project was made possible by the pre-existing implementation of some parser combinators in [benthos'](https://github.com/benthosdev/benthos) blob lang implementation. Although the result is somewhat different, this project wouldn't have been possible without this pre-existing resource as a guiding example.

## Authors

- [@oleiade](https://github.com/oleiade)
