package ecs

type IEventData[T any] interface {
	Has() bool
	Get() T
	Set(T)
	Clear()
}

type EventData[T any] struct {
	IEventData[T]
	data T
}

func NewEventData[T any]() *EventData[T] {
	return &EventData[T]{}
}

type Events[T any] struct {
	data   *EventData[T]
	reader *EventReader[T]
	writer *EventWriter[T]
}

func NewEvents[T any]() *Events[T] {
	data := NewEventData[T]()
	return &Events[T]{
		data:   data,
		reader: NewEventReader[T](data),
		writer: NewEventWriter[T](data),
	}
}

type EventReader[T any] struct {
	data IEventData[T]
}

func NewEventReader[T any](data IEventData[T]) *EventReader[T] {
	return &EventReader[T]{
		data: data,
	}
}

func (e *EventReader[T]) Has() bool {
	return e.data.Has()
}

func (e *EventReader[T]) Get() interface{} {
	return e.data.Get()
}

type EventWriter[T any] struct {
	data IEventData[T]
}

func NewEventWriter[T any](data IEventData[T]) *EventWriter[T] {
	return &EventWriter[T]{
		data: data,
	}
}

func (e *EventWriter[T]) Send(data T) {
	e.data.Set(data)
}
