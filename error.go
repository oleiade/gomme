package gomme

import (
	"fmt"
	"strings"
)

type GenericError[Input Bytes] struct {
	Input    Input
	Err      error
	Expected []string
}

func NewGenericError[Input Bytes](input Input, expected ...string) *GenericError[Input] {
	return &GenericError[Input]{Input: input, Expected: expected}
}

// Error returns a human readable error string.
func (e *GenericError[Input]) Error() string {
	return fmt.Sprintf("expected %v", strings.Join(e.Expected, ", "))
}

func (e *GenericError[Input]) IsFatal() bool {
	return e.Err != nil
}
