package ecs

import (
	"slices"
	"testing"
)

type testEvent struct {
	Value int
}

func TestEvents_WorldAdvancesOrderedBatches(t *testing.T) {
	w := NewWorld()
	events := RegisterEvents[testEvent](w)
	reader := events.Reader()
	writer := events.Writer()

	writer.Send(testEvent{Value: 1})
	writer.Send(testEvent{Value: 2})
	if reader.Has() {
		t.Fatal("pending events must not be visible before update")
	}

	w.Update()
	if got := reader.Get(); !slices.Equal(got, []testEvent{{Value: 1}, {Value: 2}}) {
		t.Fatalf("first batch = %v", got)
	}

	writer.Send(testEvent{Value: 3})
	if got := reader.Get(); !slices.Equal(got, []testEvent{{Value: 1}, {Value: 2}}) {
		t.Fatalf("writes changed the current batch: %v", got)
	}

	w.Update()
	if got := reader.Get(); !slices.Equal(got, []testEvent{{Value: 3}}) {
		t.Fatalf("second batch = %v", got)
	}
	w.Update()
	if reader.Has() {
		t.Fatalf("expired batch remains visible: %v", reader.Get())
	}
}

func TestEvents_MultipleReadersSeeSameBatch(t *testing.T) {
	events := NewEvents[int]()
	first := events.Reader()
	second := &EventReader[int]{events: events}
	events.Writer().Send(7)
	events.Advance()
	if !slices.Equal(first.Get(), []int{7}) || !slices.Equal(second.Get(), []int{7}) {
		t.Fatal("readers did not observe the same event batch")
	}
	events.Clear()
	if first.Has() || second.Has() {
		t.Fatal("Clear did not remove all events")
	}
}
