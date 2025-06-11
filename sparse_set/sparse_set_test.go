package sparse_set

import (
	"testing"
)

func TestSparseSet_Add_Contains_Remove(t *testing.T) {
	s := NewSparseSet[int](8)
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
	s := NewSparseSet[int](4)
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Clear()
	if s.Contains(1) || s.Contains(2) || s.Contains(3) {
		t.Errorf("Clear failed: set should be empty")
	}
}

func TestSparseSet_DuplicateAdd(t *testing.T) {
	s := NewSparseSet[int](4)
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
