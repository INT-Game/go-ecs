package sparse_set

import "math"

const EMPTY = -1 // -1 indicates empty slot

// SparseSet 使用混合存储策略：
// - 对于 <= MaxInt32 的值，使用 int32 数组存储（节省内存）
// - 对于 > MaxInt32 的值，使用 map 存储（处理极端情况）
type SparseSet[T int32 | int64] struct {
	density   Density[T]
	sparse    [][]T               // 使用 int32 存储索引，节省内存
	sparseMap map[uint64]struct{} // 处理超过 int32 范围的值
	pageSize  int
}

func NewSparseSet[T int32 | int64](pageSize int) *SparseSet[T] {
	return &SparseSet[T]{
		density:   NewDensity[T](),
		sparse:    make([][]T, 0),
		sparseMap: make(map[uint64]struct{}),
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

	if s.shouldUseMap(t) {
		s.sparseMap[uint64(t)] = struct{}{}
	} else {
		s.density.PushBack(t)
		densityIdx := len(s.density) - 1
		s.ensurePage(t)
		s.setDensityIndex(t, T(densityIdx))
	}
}

func (s *SparseSet[T]) Remove(t T) (ok bool) {
	if !s.Contains(t) {
		return
	}

	if s.shouldUseMap(t) {
		delete(s.sparseMap, uint64(t))
		return true
	}

	i := s.getDensityIndex(t)
	if i == T(len(s.density)-1) {
		s.deleteDensityIndex(t)
		s.density.PopBack()
	} else {
		var last T
		if last, ok = s.density.Back(); !ok {
			return
		}
		
		s.setDensityIndex(last, i)
		s.density.Swap(int(i), len(s.density)-1)

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
	return p < len(s.sparse) && o < len(s.sparse[p]) && s.sparse[p][o] > EMPTY
}

func (s *SparseSet[T]) Clear() {
	s.density = make([]T, 0)
	s.sparse = make([][]T, 0)
	s.sparseMap = make(map[uint64]struct{})
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

func (s *SparseSet[T]) getDensityIndex(t T) T {
	return s.sparse[s.page(t)][s.offset(t)]
}

func (s *SparseSet[T]) setDensityIndex(t, index T) {
	s.sparse[s.page(t)][s.offset(t)] = index
}

func (s *SparseSet[T]) deleteDensityIndex(t T) {
	s.setDensityIndex(t, EMPTY)
}

func (s *SparseSet[T]) ensurePage(t T) {
	p := s.page(t)
	if p >= len(s.sparse) {
		for i := len(s.sparse); i <= p; i++ {
			newPage := make([]T, s.pageSize)
			for j := range newPage {
				newPage[j] = EMPTY
			}
			s.sparse = append(s.sparse, newPage)
		}
	}
}
