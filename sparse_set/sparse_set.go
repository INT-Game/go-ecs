package sparse_set

import "math"

const EMPTY = -1 // -1 indicates empty slot

// SparseSet 使用混合存储策略：
// - 对于 <= MaxInt32 的值，使用 int32 数组存储（节省内存）
// - 对于 > MaxInt32 的值，使用 map 存储（处理极端情况）
type SparseSet[T uint32 | uint64] struct {
	density   Density[T]
	sparse    [][]int32    // 使用 int32 存储索引，节省内存
	sparseMap map[uint64]T // 处理超过 int32 范围的值
	pageSize  int
}

func NewSparseSet[T uint32 | uint64](pageSize int) *SparseSet[T] {
	return &SparseSet[T]{
		density:   NewDensity[T](),
		sparse:    make([][]int32, 0),
		sparseMap: make(map[uint64]T),
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
	return p < len(s.sparse) && o < len(s.sparse[p]) && s.sparse[p][o] > EMPTY
}

func (s *SparseSet[T]) Clear() {
	s.density = make([]T, 0)
	s.sparse = make([][]int32, 0)
	s.sparseMap = make(map[uint64]T)
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
		return int(s.sparseMap[uint64(t)])
	}
	return int(s.sparse[s.page(t)][s.offset(t)])
}

func (s *SparseSet[T]) setSparseIndex(t T, index int) {
	if s.shouldUseMap(t) {
		s.sparseMap[uint64(t)] = T(index)
		return
	}
	s.sparse[s.page(t)][s.offset(t)] = int32(index)
}

func (s *SparseSet[T]) deleteDensityIndex(t T) {
	if s.shouldUseMap(t) {
		delete(s.sparseMap, uint64(t))
		return
	}
	s.setSparseIndex(t, EMPTY)
}

func (s *SparseSet[T]) ensurePage(t T) {
	if s.shouldUseMap(t) {
		return
	}

	p := s.page(t)
	if p >= len(s.sparse) {
		for i := len(s.sparse); i <= p; i++ {
			newPage := make([]int32, s.pageSize)
			for j := range newPage {
				newPage[j] = EMPTY
			}
			s.sparse = append(s.sparse, newPage)
		}
	}
}
