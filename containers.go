package gomme

// PairContainer allows returning a pair of results from a parser.
type PairContainer[L, R any] struct {
	Left  L
	Right R
}

// NewPairContainer instantiates a new Pair
func NewPairContainer[L, R any](left L, right R) *PairContainer[L, R] {
	return &PairContainer[L, R]{
		Left:  left,
		Right: right,
	}
}
