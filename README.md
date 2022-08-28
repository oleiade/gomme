[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

# Gomme

Gomme is a parser combinator library for the Go programming language.
It provides a toolkit for developers to build reliable, fast, flexible, and easy-to-develop and maintain parsers
for both textual and binary formats. It extensively uses the recent introduction of Generics in the Go programming
language to offer flexibility in how combinators can be mixed and matched to produce the desired output while
providing as much compile-time type safety as possible.

## Usage/Examples

FIXME: replace with something simpler and more explicit?

```golang
type HexColor struct {
    red   uint8
    green uint8
    blue  uint8
}

func ParseHexColor(input string) (HexColor, error) {
    prefixRes := gomme.Token("#")(input)
    if prefixRes.Err != nil {
        return HexColor{}, prefixRes.Err
    }

    input = prefixRes.Remaining
    parser := gomme.Map(
        gomme.Count(HexColorComponent(), 3),
        func(components []uint8) (HexColor, error) {
            return HexColor{components[0], components[1], components[2]}, nil
        },
    )

    res := parser(input)
    if res.Err != nil {
        return HexColor{}, res.Err
    }

    return res.Output, nil
}

func HexColorComponent() gomme.Parser[string, uint8] {
    return func(input string) gomme.Result[uint8, string] {
        return gomme.Map(
            gomme.TakeWhileMN(2, 2, gomme.IsHexDigit),
            fromHex,
        )(input)
    }
}

func fromHex(input string) (uint8, error) {
    res, err := strconv.ParseInt(input, 16, 16)
    if err != nil {
        return 0, err
    }

    return uint8(res), nil
}
```

## Installation

Add the library to your Go project with the following command:

```bash
    go get github.com/oleiade/gomme@latest
```

## Documentation

[Documentation](https://linktodocumentation)

## FAQ

#### What are parser combinators?

TODO
Answer 1

#### Why would I want to use parser combinators, and not write my own specific parser?

TODO
Answer 2

#### How fast are parser combinators?

TODO
Answer 3

## Acknowledgements

We can frankly take close to zero credit for this library, apart for the work put into assembling the already existing elements of theory and implementation into a single autonomous project.

This library relies heavily on the whole theorical work done in the parser combinators space. From the implementation side of things, it was specifically started with the intention to have something similar to Rust's incredible [nom](https://github.com/Geal/nom) library in Go. This project was made possible by the pre-existing implementation of some parser combinators in [benthos'](https://github.com/benthosdev/benthos) blob lang implementation. Although the end-result is somewhat different from it, this project wouldn't have been possible without this pre-existing resource as a guiding example.

## Authors

- [@oleiade](https://github.com/oleiade)
