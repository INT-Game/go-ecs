package sparse_set

type Density[T uint32 | uint64] []T

func NewDensity[T uint32 | uint64]() Density[T] {
	return make(Density[T], 0)
}

func (d *Density[T]) PopBack() {
	if len(*d) > 0 {
		*d = (*d)[:len(*d)-1]
	}
}

func (d *Density[T]) PushBack(t T) {
	*d = append(*d, t)
}

func (d *Density[T]) Len() int {
	return len(*d)
}

func (d *Density[T]) Swap(i, j int) {
	if i < 0 || i >= len(*d) || j < 0 || j >= len(*d) {
		return
	}
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

func (d *Density[T]) Back() (T, bool) {
	if len(*d) == 0 {
		var zero T
		return zero, false
	}
	return (*d)[len(*d)-1], true
}
