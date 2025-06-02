package utils

type Pair[A, B any] struct {
	First A
	Second B
}

func NewPair[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{ a, b }
}
