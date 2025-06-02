package utils

type Float interface {
	~float32 | ~float64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Integer interface {
	Unsigned | Signed
}

type Number interface {
	Float | Integer
}

type Range[T Number] struct {
	Start T
	End T
}

func NewRange[T Number](start, end T) Range[T] {
	return Range[T]{
		Start: start,
		End: end,
	}
}

func (r Range[T]) Length() T {
	return r.End - r.Start
}