package dsa


type Stack[T any] interface {
	Push(val T) 
	Pop() T
	Peek() T
	Empty() bool
	Size() int
}

// using golang slices
type SliceStack[T any] struct {
	arr []T
}

func NewSliceStack[T any]() *SliceStack[T] {
	return &SliceStack[T]{
		arr: make([]T, 0),
	}
}

func (s *SliceStack[T]) Push(val T) {
	s.arr = append(s.arr, val)
}

func (s *SliceStack[T]) Peek() T {
	return s.arr[len(s.arr)-1]
}

func (s *SliceStack[T]) Pop() T {
	val := s.Peek()
	s.arr = s.arr[:len(s.arr)-1]
	return val
}

func (s *SliceStack[T]) Empty() bool {
	return len(s.arr) == 0
}

func (s *SliceStack[T]) Size() int {
	return len(s.arr)
}