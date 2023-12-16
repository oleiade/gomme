package main

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/oleiade/gomme"
)

//go:embed test.json
var testJSON string

func main() {
	result := parseJSON(testJSON)
	if result.Err != nil {
		log.Fatal(result.Err)
		return
	}

	fmt.Println(result.Output)
}

type (
	// JSONValue represents any value that can be encountered in
	// JSON, including complex types like objects and arrays.
	JSONValue interface{}

	// JSONString represents a JSON string value.
	JSONString string

	// JSONNumber represents a JSON number value, which internally is treated as float64.
	JSONNumber float64

	// JSONObject represents a JSON object, which is a collection of key-value pairs.
	JSONObject map[string]JSONValue

	// JSONArray represents a JSON array, which is a list of JSON values.
	JSONArray []JSONValue

	// JSONBool represents a JSON boolean value.
	JSONBool bool

	// JSONNull represents the JSON null value.
	JSONNull struct{}
)

// parseJSON is a convenience function to start parsing JSON from the given input string.
func parseJSON(input string) gomme.Result[JSONValue, string] {
	return parseValue(input)
}

// parseValue is a parser that attempts to parse different types of
// JSON values (object, array, string, etc.).
func parseValue(input string) gomme.Result[JSONValue, string] {
	return gomme.Alternative(
		parseObject,
		parseArray,
		parseString,
		parseNumber,
		parseTrue,
		parseFalse,
		parseNull,
	)(input)
}

// parseObject parses a JSON object, which starts and ends with
// curly braces and contains key-value pairs.
func parseObject(input string) gomme.Result[JSONValue, string] {
	return gomme.Map(
		gomme.Delimited[string, rune, map[string]JSONValue, rune](
			gomme.Char[string]('{'),
			gomme.Optional[string, map[string]JSONValue](
				gomme.Preceded(
					ws(),
					gomme.Terminated[string, map[string]JSONValue](
						parseMembers,
						ws(),
					),
				),
			),
			gomme.Char[string]('}'),
		),
		func(members map[string]JSONValue) (JSONValue, error) {
			return JSONObject(members), nil
		},
	)(input)
}

// Ensure parseObject is a Parser[string, JSONValue]
var _ gomme.Parser[string, JSONValue] = parseObject

// parseArray parses a JSON array, which starts and ends with
// square brackets and contains a list of values.
func parseArray(input string) gomme.Result[JSONValue, string] {
	return gomme.Map(
		gomme.Delimited[string, rune, []JSONValue, rune](
			gomme.Char[string]('['),
			gomme.Alternative(
				parseElements,
				gomme.Map(ws(), func(s string) ([]JSONValue, error) { return []JSONValue{}, nil }),
			),
			gomme.Char[string](']'),
		),
		func(elements []JSONValue) (JSONValue, error) {
			return JSONArray(elements), nil
		},
	)(input)
}

// Ensure parseArray is a Parser[string, JSONValue]
var _ gomme.Parser[string, JSONValue] = parseArray

func parseElement(input string) gomme.Result[JSONValue, string] {
	return gomme.Map(
		gomme.Delimited[string](ws(), parseValue, ws()),
		func(v JSONValue) (JSONValue, error) { return v, nil },
	)(input)
}

// Ensure parseElement is a Parser[string, JSONValue]
var _ gomme.Parser[string, JSONValue] = parseElement

// parseNumber parses a JSON number.
func parseNumber(input string) gomme.Result[JSONValue, string] {
	return gomme.Map[string](
		gomme.Sequence(
			gomme.Map(integer(), func(i int) (string, error) { return strconv.Itoa(i), nil }),
			gomme.Optional(fraction()),
			gomme.Optional(exponent()),
		),
		func(parts []string) (JSONValue, error) {
			// Construct the float string from parts
			var floatStr string

			// Integer part
			floatStr += parts[0]

			// Fraction part
			if parts[1] != "" {
				fractionPart, err := strconv.Atoi(parts[1])
				if err != nil {
					return 0, err
				}

				if fractionPart != 0 {
					floatStr += fmt.Sprintf(".%d", fractionPart)
				}
			}

			// Exponent part
			if parts[2] != "" {
				floatStr += fmt.Sprintf("e%s", parts[2])
			}

			f, err := strconv.ParseFloat(floatStr, 64)
			if err != nil {
				return JSONNumber(0.0), err
			}

			return JSONNumber(f), nil
		},
	)(input)
}

// Ensure parseNumber is a Parser[string, JSONValue]
var _ gomme.Parser[string, JSONValue] = parseNumber

// parseString parses a JSON string.
func parseString(input string) gomme.Result[JSONValue, string] {
	return gomme.Map(
		stringParser(),
		func(s string) (JSONValue, error) {
			return JSONString(s), nil
		},
	)(input)
}

// Ensure parseString is a Parser[string, JSONValue]
var _ gomme.Parser[string, JSONValue] = parseString

// parseFalse parses the JSON boolean value 'false'.
func parseFalse(input string) gomme.Result[JSONValue, string] {
	return gomme.Map(
		gomme.Token[string]("false"),
		func(_ string) (JSONValue, error) { return JSONBool(false), nil },
	)(input)
}

// Ensure parseFalse is a Parser[string, JSONValue]
var _ gomme.Parser[string, JSONValue] = parseFalse

// parseTrue parses the JSON boolean value 'true'.
func parseTrue(input string) gomme.Result[JSONValue, string] {
	return gomme.Map(
		gomme.Token[string]("true"),
		func(_ string) (JSONValue, error) { return JSONBool(true), nil },
	)(input)
}

// Ensure parseTrue is a Parser[string, JSONValue]
var _ gomme.Parser[string, JSONValue] = parseTrue

// parseNull parses the JSON 'null' value.
func parseNull(input string) gomme.Result[JSONValue, string] {
	return gomme.Map(
		gomme.Token[string]("null"),
		func(_ string) (JSONValue, error) { return nil, nil },
	)(input)
}

// Ensure parseNull is a Parser[string, JSONValue]
var _ gomme.Parser[string, JSONValue] = parseNull

// parseElements parses the elements of a JSON array.
func parseElements(input string) gomme.Result[[]JSONValue, string] {
	return gomme.Map(
		gomme.SeparatedList0[string](
			parseElement,
			gomme.Token[string](","),
		),
		func(elems []JSONValue) ([]JSONValue, error) {
			return elems, nil
		},
	)(input)
}

// Ensure parseElements is a Parser[string, []JSONValue]
var _ gomme.Parser[string, []JSONValue] = parseElements

// parseElement parses a single element of a JSON array.
func parseMembers(input string) gomme.Result[map[string]JSONValue, string] {
	return gomme.Map(
		gomme.SeparatedList0[string](
			parseMember,
			gomme.Token[string](","),
		),
		func(kvs []kv) (map[string]JSONValue, error) {
			obj := make(JSONObject)
			for _, kv := range kvs {
				obj[kv.key] = kv.value
			}
			return obj, nil
		},
	)(input)
}

// Ensure parseMembers is a Parser[string, map[string]JSONValue]
var _ gomme.Parser[string, map[string]JSONValue] = parseMembers

// parseMember parses a single member (key-value pair) of a JSON object.
func parseMember(input string) gomme.Result[kv, string] {
	return member()(input)
}

// Ensure parseMember is a Parser[string, kv]
var _ gomme.Parser[string, kv] = parseMember

// member creates a parser for a single key-value pair in a JSON object.
//
// It expects a string followed by a colon and then a JSON value.
// The result is a kv struct with the parsed key and value.
func member() gomme.Parser[string, kv] {
	mapFunc := func(p gomme.PairContainer[string, JSONValue]) (kv, error) {
		return kv{p.Left, p.Right}, nil
	}

	return gomme.Map(
		gomme.SeparatedPair[string](
			gomme.Delimited(ws(), stringParser(), ws()),
			gomme.Token[string](":"),
			element(),
		),
		mapFunc,
	)
}

// element creates a parser for a single element in a JSON array.
//
// It wraps the element with optional whitespace on either side.
func element() gomme.Parser[string, JSONValue] {
	return gomme.Map(
		gomme.Delimited(ws(), parseValue, ws()),
		func(v JSONValue) (JSONValue, error) { return v, nil },
	)
}

// kv is a struct representing a key-value pair in a JSON object.
//
// 'key' holds the string key, and 'value' holds the corresponding JSON value.
type kv struct {
	key   string
	value JSONValue
}

// stringParser creates a parser for a JSON string.
//
// It expects a sequence of characters enclosed in double quotes.
func stringParser() gomme.Parser[string, string] {
	return gomme.Delimited[string, rune, string, rune](
		gomme.Char[string]('"'),
		characters(),
		gomme.Char[string]('"'),
	)
}

// integer creates a parser for a JSON number's integer part.
//
// It handles negative and positive integers including zero.
func integer() gomme.Parser[string, int] {
	return gomme.Alternative(
		// "-" onenine digits
		gomme.Preceded(
			gomme.Token[string]("-"),
			gomme.Map(
				gomme.Pair(onenine(), digits()),
				func(p gomme.PairContainer[string, string]) (int, error) {
					return strconv.Atoi(p.Left + p.Right)
				},
			),
		),

		// onenine digits
		gomme.Map(
			gomme.Pair(onenine(), digits()),
			func(p gomme.PairContainer[string, string]) (int, error) {
				return strconv.Atoi(p.Left + p.Right)
			},
		),

		// "-" digit
		gomme.Preceded(
			gomme.Token[string]("-"),
			gomme.Map(
				digit(),
				strconv.Atoi,
			),
		),

		// digit
		gomme.Map(digit(), strconv.Atoi),
	)
}

// digits creates a parser for a sequence of digits.
//
// It concatenates the sequence into a single string.
func digits() gomme.Parser[string, string] {
	return gomme.Map(gomme.Many1(digit()), func(digits []string) (string, error) {
		return strings.Join(digits, ""), nil
	})
}

// digit creates a parser for a single digit.
//
// It distinguishes between '0' and non-zero digits.
func digit() gomme.Parser[string, string] {
	return gomme.Alternative(
		gomme.Token[string]("0"),
		onenine(),
	)
}

// onenine creates a parser for digits from 1 to 9.
func onenine() gomme.Parser[string, string] {
	return gomme.Alternative(
		gomme.Token[string]("1"),
		gomme.Token[string]("2"),
		gomme.Token[string]("3"),
		gomme.Token[string]("4"),
		gomme.Token[string]("5"),
		gomme.Token[string]("6"),
		gomme.Token[string]("7"),
		gomme.Token[string]("8"),
		gomme.Token[string]("9"),
	)
}

// fraction creates a parser for the fractional part of a JSON number.
//
// It expects a dot followed by at least one digit.
func fraction() gomme.Parser[string, string] {
	return gomme.Preceded(
		gomme.Token[string]("."),
		gomme.Digit1[string](),
	)
}

// exponent creates a parser for the exponent part of a JSON number.
//
// It handles the exponent sign and the exponent digits.
func exponent() gomme.Parser[string, string] {
	return gomme.Preceded(
		gomme.Token[string]("e"),
		gomme.Map(
			gomme.Pair(sign(), digits()),
			func(p gomme.PairContainer[string, string]) (string, error) {
				return p.Left + p.Right, nil
			},
		),
	)
}

// sign creates a parser for the sign part of a number's exponent.
//
// It can parse both positive ('+') and negative ('-') signs.
func sign() gomme.Parser[string, string] {
	return gomme.Optional(
		gomme.Alternative[string, string](
			gomme.Token[string]("-"),
			gomme.Token[string]("+"),
		),
	)
}

// characters creates a parser for a sequence of JSON string characters.
//
// It handles regular characters and escaped sequences.
func characters() gomme.Parser[string, string] {
	return gomme.Optional(
		gomme.Map(
			gomme.Many1[string, rune](character()),
			func(chars []rune) (string, error) {
				return string(chars), nil
			},
		),
	)
}

// character creates a parser for a single JSON string character.
//
// It distinguishes between regular characters and escape sequences.
func character() gomme.Parser[string, rune] {
	return gomme.Alternative(
		// normal character
		gomme.Satisfy[string](func(c rune) bool {
			return c != '"' && c != '\\' && c >= 0x20 && c <= 0x10FFFF
		}),

		// escape
		escape(),
	)
}

// escape creates a parser for escaped characters in a JSON string.
//
// It handles common escape sequences like '\n', '\t', etc., and unicode escapes.
func escape() gomme.Parser[string, rune] {
	mapFunc := func(chars []rune) (rune, error) {
		// chars[0] will always be '\\'
		switch chars[1] {
		case '"':
			return '"', nil
		case '\\':
			return '\\', nil
		case '/':
			return '/', nil
		case 'b':
			return '\b', nil
		case 'f':
			return '\f', nil
		case 'n':
			return '\n', nil
		case 'r':
			return '\r', nil
		case 't':
			return '\t', nil
		default: // for unicode escapes
			return chars[1], nil
		}
	}

	return gomme.Map(
		gomme.Sequence(
			gomme.Char[string]('\\'),
			gomme.Alternative(
				gomme.Char[string]('"'),
				gomme.Char[string]('\\'),
				gomme.Char[string]('/'),
				gomme.Char[string]('b'),
				gomme.Char[string]('f'),
				gomme.Char[string]('n'),
				gomme.Char[string]('r'),
				gomme.Char[string]('t'),
				unicodeEscape(),
			),
		),
		mapFunc,
	)
}

// unicodeEscape creates a parser for a unicode escape sequence in a JSON string.
//
// It expects a sequence starting with 'u' followed by four hexadecimal digits and
// converts them to the corresponding rune.
func unicodeEscape() gomme.Parser[string, rune] {
	mapFunc := func(chars []rune) (rune, error) {
		// chars[0] will always be 'u'
		hex := string(chars[1:5])
		codePoint, err := strconv.ParseInt(hex, 16, 32)
		if err != nil {
			return 0, err
		}
		return rune(codePoint), nil
	}

	return gomme.Map(
		gomme.Sequence(
			gomme.Char[string]('u'),
			hex(),
			hex(),
			hex(),
			hex(),
		),
		mapFunc,
	)
}

// hex creates a parser for a single hexadecimal digit.
//
// It can parse digits ('0'-'9') as well as
// letters ('a'-'f', 'A'-'F') used in hexadecimal numbers.
func hex() gomme.Parser[string, rune] {
	return gomme.Satisfy[string](func(r rune) bool {
		return ('0' <= r && r <= '9') || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
	})
}

// ws creates a parser for whitespace in JSON.
//
// It can handle spaces, tabs, newlines, and carriage returns.
// The parser accumulates all whitespace characters and returns them as a single string.
func ws() gomme.Parser[string, string] {
	parser := gomme.Many0(
		gomme.Satisfy[string](func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n' || r == '\r'
		}),
	)

	mapFunc := func(runes []rune) (string, error) {
		return string(runes), nil
	}

	return gomme.Map(parser, mapFunc)
}
