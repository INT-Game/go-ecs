package array

type Array[T comparable] []T

func New[T comparable]() Array[T] {
	return make(Array[T], 0)
}

func (a *Array[T]) PopBack() {
	if len(*a) > 0 {
		*a = (*a)[:len(*a)-1]
	}
}

func (a *Array[T]) PushBack(t T) {
	*a = append(*a, t)
}

func (a *Array[T]) Len() int {
	return len(*a)
}

func (a *Array[T]) Swap(i, j int) {
	if i < 0 || i >= len(*a) || j < 0 || j >= len(*a) {
		return
	}
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
}

func (a *Array[T]) Back() T {
	if len(*a) == 0 {
		var zero T
		return zero
	}
	return (*a)[len(*a)-1]
}

func (a *Array[T]) Clear() {
	*a = make(Array[T], 0)
}

func (a *Array[T]) Empty() bool {
	return len(*a) == 0
}

func (a *Array[T]) Contain(t T) (index int, ok bool) {
	for i, v := range *a {
		if v == t {
			return i, true
		}
	}

	return
}
