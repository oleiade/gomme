package gomme

// PairContainer allows returning a pair of results from a parser.
type PairContainer[Left, Right any] struct {
	Left  Left
	Right Right
}

// NewPairContainer instantiates a new Pair
func NewPairContainer[Left, Right any](left Left, right Right) *PairContainer[Left, Right] {
	return &PairContainer[Left, Right]{
		Left:  left,
		Right: right,
	}
}
