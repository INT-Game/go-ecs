package sparse_set

const EMPTY = -1 // -1 indicates empty slot

type SparseSet[T int | uint64] struct {
	density  Density[T]
	sparse   [][]int
	pageSize int
}

func NewSparseSet[T int | uint64](pageSize int) *SparseSet[T] {
	return &SparseSet[T]{
		density:  NewDensity[T](),
		sparse:   make([][]int, 0),
		pageSize: pageSize,
	}
}

func (s *SparseSet[T]) Density() Density[T] {
	return s.density
}

func (s *SparseSet[T]) Add(t T) {
	if s.Contains(t) {
		return
	}

	s.density.PushBack(t)
	s.ensurePage(t)
	s.setDensityIndex(t, len(s.density)-1)
}

func (s *SparseSet[T]) Remove(t T) (ok bool) {
	if !s.Contains(t) {
		return
	}

	i := s.getDensityIndex(t)
	if i == len(s.density)-1 {
		s.setDensityIndex(t, EMPTY)
		s.density.PopBack()
	} else {
		var last T
		if last, ok = s.density.Back(); !ok {
			return
		}
		s.setDensityIndex(last, i)

		s.density.Swap(i, len(s.density)-1)

		s.setDensityIndex(t, EMPTY)
		s.density.PopBack()
	}

	return true
}

func (s *SparseSet[T]) Contains(t T) bool {
	p := s.page(t)
	o := s.offset(t)
	return p < len(s.sparse) && o < len(s.sparse[p]) && s.sparse[p][o] > EMPTY
}

func (s *SparseSet[T]) Clear() {
	s.density = make([]T, 0)
	s.sparse = make([][]int, 0)
}

func (s *SparseSet[T]) page(t T) int {
	return int(t) / s.pageSize
}

func (s *SparseSet[T]) offset(t T) int {
	return int(t) % s.pageSize
}

func (s *SparseSet[T]) getDensityIndex(t T) int {
	return s.sparse[s.page(t)][s.offset(t)]
}

func (s *SparseSet[T]) setDensityIndex(t T, index int) {
	s.sparse[s.page(t)][s.offset(t)] = index
}

func (s *SparseSet[T]) ensurePage(t T) {
	p := s.page(t)
	if p >= len(s.sparse) {
		for i := len(s.sparse); i <= p; i++ {
			newPage := make([]int, s.pageSize)
			for j := range newPage {
				newPage[j] = EMPTY
			}
			s.sparse = append(s.sparse, newPage)
		}
	}
}
