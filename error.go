package gomme

import (
	"fmt"
	"strings"
)

// Error represents a parsing error. It holds the input that was being parsed,
// the parsers that were tried, and the error that was produced.
type Error[Input Bytes] struct {
	Input    Input
	Err      error
	Expected []string
}

// NewError produces a new Error from the provided input and names of
// parsers expected to succeed.
func NewError[Input Bytes](input Input, expected ...string) *Error[Input] {
	return &Error[Input]{Input: input, Expected: expected}
}

// Error returns a human readable error string.
func (e *Error[Input]) Error() string {
	return fmt.Sprintf("expected %v", strings.Join(e.Expected, ", "))
}

// IsFatal returns true if the error is fatal.
func (e *Error[Input]) IsFatal() bool {
	return e.Err != nil
}
