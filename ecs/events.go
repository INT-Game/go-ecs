package ecs

// Events is a typed, double-buffered event queue. Events sent during one
// World update become readable at the start of the next update.
type Events[T any] struct {
	current []T
	pending []T
	reader  EventReader[T]
	writer  EventWriter[T]
}

func NewEvents[T any]() *Events[T] {
	events := &Events[T]{}
	events.reader.events = events
	events.writer.events = events
	return events
}

// RegisterEvents creates an event queue managed by w.
func RegisterEvents[T any](w *World) *Events[T] {
	events := NewEvents[T]()
	w.events = append(w.events, events)
	return events
}

func (e *Events[T]) Reader() *EventReader[T] {
	return &e.reader
}

func (e *Events[T]) Writer() *EventWriter[T] {
	return &e.writer
}

// Advance expires the current batch and exposes the pending batch.
func (e *Events[T]) Advance() {
	clear(e.current)
	e.current, e.pending = e.pending, e.current[:0]
}

func (e *Events[T]) Clear() {
	clear(e.current)
	clear(e.pending)
	e.current = e.current[:0]
	e.pending = e.pending[:0]
}

type EventReader[T any] struct {
	events *Events[T]
}

func (r *EventReader[T]) Has() bool {
	return len(r.events.current) > 0
}

// Get returns the current batch. Treat the returned slice as read-only; it is
// valid until the next Advance or Clear.
func (r *EventReader[T]) Get() []T {
	return r.events.current
}

type EventWriter[T any] struct {
	events *Events[T]
}

func (w *EventWriter[T]) Send(event T) {
	w.events.pending = append(w.events.pending, event)
}
