// A somewhat limited json parser
//
// TODO: json strings don't support escaping, or spaces, or unicode characters
package json

import "github.com/oleiade/gomme"

func ParseJSON(input string) (JSONValue, error) {
	res := Value()(input)

	return res.Output, res.Err
}

type JSONValue struct {
	Kind   JSONKind
	Null   JSONNull
	Bool   JSONBool
	String JSONString
	Number JSONNumber
	Object JSONObject
	Array  []*JSONValue
}

type JSONKind string

const (
	JSONNullKind   JSONKind = "null"
	JSONBoolKind   JSONKind = "bool"
	JSONStringKind JSONKind = "string"
	JSONNumberKind JSONKind = "number"
	JSONObjectKind JSONKind = "object"
	JSONArrayKind  JSONKind = "array"
)

func Value[O JSONValue]() gomme.Parser[string, JSONValue] {
	return gomme.Preceded(
		separator(),
		gomme.Alternative(
			// Objects
			gomme.Map(Object(), func(input JSONObject) (JSONValue, error) {
				return JSONValue{Kind: JSONObjectKind, Object: input}, nil
			}),

			// Strings
			gomme.Map(String(), func(input JSONString) (JSONValue, error) {
				return JSONValue{Kind: JSONStringKind, String: input}, nil
			}),

			// Numbers
			gomme.Map(Number(), func(input JSONNumber) (JSONValue, error) {
				return JSONValue{Kind: JSONNumberKind, Number: input}, nil
			}),

			// Boolean
			gomme.Map(Boolean(), func(input JSONBool) (JSONValue, error) {
				return JSONValue{Kind: JSONBoolKind, Bool: input}, nil
			}),

			// Null
			gomme.Assign(JSONValue{Kind: JSONNullKind, Null: JSONNull{}}, Null()),
		),
	)
}

type JSONNull struct{}

func Null() gomme.Parser[string, JSONNull] {
	return gomme.Assign(JSONNull(struct{}{}), gomme.Token[string]("null"))
}

func nullValue(input JSONNull) (JSONValue, error) {
	return JSONValue{Kind: JSONNullKind, Null: JSONNull{}}, nil
}

type JSONBool bool

func Boolean() gomme.Parser[string, JSONBool] {
	return gomme.Alternative(
		gomme.Assign(JSONBool(true), gomme.Token[string]("true")),
		gomme.Assign(JSONBool(false), gomme.Token[string]("false")),
	)
}

type JSONString string

func String() gomme.Parser[string, JSONString] {
	return gomme.Map(
		gomme.Delimited(gomme.Char[string]('"'), parseString(), gomme.Char[string]('"')),
		func(input string) (JSONString, error) {
			return JSONString(input), nil
		},
	)
}

type JSONNumber float64

func Number() gomme.Parser[string, JSONNumber] {
	return gomme.Map(
		gomme.Number[string](),
		func(input float64) (JSONNumber, error) {
			return JSONNumber(input), nil
		},
	)
}

type JSONObject map[JSONString]JSONValue

func Object() gomme.Parser[string, JSONObject] {
	return gomme.Delimited(
		gomme.Token[string]("{"),
		gomme.Map(gomme.Token[string]("abc:123"), func(input string) (JSONObject, error) {
			obj := make(map[JSONString]JSONValue)
			obj[JSONString("abc")] = JSONValue{Kind: JSONNumberKind, Number: JSONNumber(123)}
			return JSONObject(obj), nil
		}),
		gomme.Token[string]("}"),
	)
}

type JSONArray []JSONValue

type keyvalue struct {
	Key   string
	Value JSONValue
}

func separator() gomme.Parser[string, string] {
	return gomme.Whitespace0[string]()
}

func parseString() gomme.Parser[string, string] {
	return gomme.Alphanumeric0[string]()
}
