package sparse_set

import (
	"math"
	"testing"
)

func TestSparseSet_Add_Contains_Remove(t *testing.T) {
	s := NewSparseSet[uint32](8)
	s.Add(10)
	s.Add(20)
	s.Add(30)

	if !s.Contains(10) || !s.Contains(20) || !s.Contains(30) {
		t.Errorf("Add or Contains failed")
	}

	s.Remove(20)
	if s.Contains(20) {
		t.Errorf("DestroyEntity failed: 20 should not be contained")
	}

	s.Remove(10)
	if s.Contains(10) {
		t.Errorf("DestroyEntity failed: 10 should not be contained")
	}

	s.Remove(30)
	if s.Contains(30) {
		t.Errorf("DestroyEntity failed: 30 should not be contained")
	}
}

func TestSparseSet_Clear(t *testing.T) {
	s := NewSparseSet[uint32](4)
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Clear()
	if s.Contains(1) || s.Contains(2) || s.Contains(3) {
		t.Errorf("Clear failed: set should be empty")
	}
}

func TestSparseSet_DuplicateAdd(t *testing.T) {
	s := NewSparseSet[uint32](4)
	s.Add(5)
	s.Add(5)
	count := 0
	for _, v := range s.density {
		if v == 5 {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Duplicate Add failed: expected 1, got %d", count)
	}
}

func TestSparseSet_ZeroMeansEmptyAndIndexesStartAtOne(t *testing.T) {
	s := NewSparseSet[uint32](4)
	s.Add(0)
	if !s.Contains(0) {
		t.Fatal("zero must remain a valid set value")
	}
	if got := s.sparse[0][0]; got != 1 {
		t.Fatalf("first density index encoding = %d, want 1", got)
	}
	s.Remove(0)
	if got := s.sparse[0][0]; got != EMPTY {
		t.Fatalf("removed sparse slot = %d, want EMPTY", got)
	}
}

func TestSparseSet_Uint32IndexBeyondMaxInt32(t *testing.T) {
	s := NewSparseSet[uint32](4)
	s.ensurePage(1)
	want := int64(math.MaxInt32) + 1
	s.setSparseIndex(1, int(want))
	if got := int64(s.getDensityIndex(1)); got != want {
		t.Fatalf("density index = %d, want %d", got, want)
	}
}

func TestSparseSet_LargeValueUsesMap(t *testing.T) {
	s := NewSparseSet[uint64](4)
	value := uint64(math.MaxInt32) + 1
	s.Add(value)
	if !s.Contains(value) || !s.Remove(value) || s.Contains(value) {
		t.Fatal("large value map lifecycle failed")
	}
}

func TestSparseSet_RejectsInvalidPageSize(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewSparseSet must reject a non-positive page size")
		}
	}()
	NewSparseSet[uint32](0)
}
