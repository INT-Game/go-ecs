package sparse_set

import "math"

const EMPTY uint32 = 0

// SparseSet stores density indexes as index+1 so the uint32 zero value means empty.
// Large values use a map to avoid allocating an impractically large page table.
type SparseSet[T uint32 | uint64] struct {
	density   Density[T]
	sparse    [][]uint32
	sparseMap map[uint64]uint32
	pageSize  int
}

func NewSparseSet[T uint32 | uint64](pageSize int) *SparseSet[T] {
	if pageSize <= 0 {
		panic("sparse_set: page size must be positive")
	}
	return &SparseSet[T]{
		density:   NewDensity[T](),
		sparse:    make([][]uint32, 0),
		sparseMap: make(map[uint64]uint32),
		pageSize:  pageSize,
	}
}

func (s *SparseSet[T]) Density() Density[T] {
	return s.density
}

func (s *SparseSet[T]) Add(t T) {
	if s.Contains(t) {
		return
	}
	if uint64(len(s.density)) >= uint64(math.MaxUint32) {
		panic("sparse_set: uint32 index capacity exceeded")
	}

	s.density.PushBack(t)
	densityIdx := len(s.density) - 1

	s.ensurePage(t)
	s.setSparseIndex(t, densityIdx)
}

func (s *SparseSet[T]) Remove(t T) (ok bool) {
	if !s.Contains(t) {
		return
	}

	i := s.getDensityIndex(t)
	if i == len(s.density)-1 {
		s.deleteDensityIndex(t)
		s.density.PopBack()
	} else {
		var last T
		if last, ok = s.density.Back(); !ok {
			return
		}

		// 更新被交换的最后一个元素的索引
		s.setSparseIndex(last, i)
		s.density.Swap(i, len(s.density)-1)

		s.deleteDensityIndex(t)
		s.density.PopBack()
	}

	return true
}

func (s *SparseSet[T]) Contains(t T) bool {
	if s.shouldUseMap(t) {
		_, exists := s.sparseMap[uint64(t)]
		return exists
	}

	p := s.page(t)
	o := s.offset(t)
	return p < len(s.sparse) && o < len(s.sparse[p]) && s.sparse[p][o] != EMPTY
}

func (s *SparseSet[T]) Clear() {
	s.density = make([]T, 0)
	s.sparse = make([][]uint32, 0)
	s.sparseMap = make(map[uint64]uint32)
}

// shouldUseMap 判断是否应该使用 map 存储
func (s *SparseSet[T]) shouldUseMap(t T) bool {
	return uint64(t) > math.MaxInt32
}

func (s *SparseSet[T]) page(t T) int {
	return int(t) / s.pageSize
}

func (s *SparseSet[T]) offset(t T) int {
	return int(t) % s.pageSize
}

func (s *SparseSet[T]) getDensityIndex(t T) int {
	if s.shouldUseMap(t) {
		return int(s.sparseMap[uint64(t)] - 1)
	}
	return int(s.sparse[s.page(t)][s.offset(t)] - 1)
}

func (s *SparseSet[T]) setSparseIndex(t T, index int) {
	encoded := uint32(index) + 1
	if s.shouldUseMap(t) {
		s.sparseMap[uint64(t)] = encoded
		return
	}
	s.sparse[s.page(t)][s.offset(t)] = encoded
}

func (s *SparseSet[T]) deleteDensityIndex(t T) {
	if s.shouldUseMap(t) {
		delete(s.sparseMap, uint64(t))
		return
	}
	s.sparse[s.page(t)][s.offset(t)] = EMPTY
}

func (s *SparseSet[T]) ensurePage(t T) {
	if s.shouldUseMap(t) {
		return
	}

	p := s.page(t)
	if p >= len(s.sparse) {
		for i := len(s.sparse); i <= p; i++ {
			s.sparse = append(s.sparse, make([]uint32, s.pageSize))
		}
	}
}
