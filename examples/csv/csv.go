// Package csv implements a parser for CSV files.
//
// It is a simple, incomplete, example of how to use the gomme
// parser combinator library to build a parser targetting the
// format described in [RFC4180].
//
// [RFC4180]: https://tools.ietf.org/html/rfc4180
package csv

import "github.com/oleiade/gomme"

func ParseCSV(input string) ([][]string, error) {
	parser := gomme.SeparatedList1(
		gomme.SeparatedList1(
			gomme.Alternative(
				gomme.Alphanumeric1[string](),
				gomme.Delimited(gomme.Char[string]('"'), gomme.Alphanumeric1[string](), gomme.Char[string]('"')),
			),
			gomme.Char[string](','),
		),
		gomme.CRLF[string](),
	)

	result := parser(input)
	if result.Err != nil {
		return nil, result.Err
	}

	return result.Output, nil
}
