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
	Null   *JSONNull
	Bool   *JSONBool
	String *JSONString
	Number *JSONNumber
	Object *JSONObject
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
			gomme.Assign(JSONValue{Kind: JSONNullKind, Null: &JSONNull{}}, Null()),
			// gomme.Map(Null(), nullValue),
			// gomme.Map(Boolean(), booleanValue),
			// gomme.Map(String(), stringValue),
		),
	)
}

type JSONNull struct{}

func Null() gomme.Parser[string, JSONNull] {
	return gomme.Assign(JSONNull(struct{}{}), gomme.Token[string]("null"))
}

func nullValue(input JSONNull) (JSONValue, error) {
	return JSONValue{Kind: JSONNullKind, Null: &JSONNull{}}, nil
}

type JSONBool bool

func Boolean() gomme.Parser[string, JSONBool] {
	return gomme.Alternative(
		gomme.Assign(JSONBool(true), gomme.Token[string]("true")),
		gomme.Assign(JSONBool(false), gomme.Token[string]("false")),
	)
}

func booleanValue(input JSONBool) (JSONValue, error) {
	return JSONValue{Kind: JSONBoolKind, Bool: &input}, nil
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

func stringValue(input JSONString) (JSONValue, error) {
	return JSONValue{Kind: JSONStringKind, String: &input}, nil
}

type JSONNumber float64
type JSONObject map[string]JSONValue
type JSONArray []JSONValue

type keyvalue struct {
	Key   string
	Value JSONValue
}

// func KeyValue() gomme.Parser[string, keyvalue] {
// 	gomme.SeparatedPair(
// 		gomme.Preceded(separated, String()),
// 		gomme.Preceded(separated, )
// 	)
// }

func separator() gomme.Parser[string, string] {
	return gomme.TakeWhileOneOf[string]('\t', '\r', '\n', ' ')
}

func parseString() gomme.Parser[string, string] {
	return gomme.Alphanumeric0[string]()
}
