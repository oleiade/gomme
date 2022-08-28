package main

import (
	"fmt"
	"strconv"

	"github.com/oleiade/gomme"
)

func main() {
	fmt.Println(ParseHexColor("#ff00ff"))
}

type HexColor struct {
	red   uint8
	green uint8
	blue  uint8
}

func ParseHexColor(input string) (HexColor, error) {
	prefixRes := gomme.Token[string]("#")(input)
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
			gomme.TakeWhileMN[string](2, 2, gomme.IsHexDigit),
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
