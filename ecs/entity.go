package ecs

import (
	"reflect"
	"sync/atomic"
)

var (
	eidIncr uint64
)

type IEntity interface {
	Identifier
	GetComponentContainer() ComponentContainer
	AddComponents(components ...IComponent)
	RemoveComponents(components ...IComponent)
}

type Entity struct {
	IEntity
	w                  IWorld
	id                 uint64
	componentContainer ComponentContainer
}

func NewEntity(w IWorld) *Entity {
	return &Entity{
		w:                  w,
		id:                 atomic.AddUint64(&eidIncr, 1),
		componentContainer: make(ComponentContainer),
	}
}

func (e *Entity) ID() uint64 {
	return e.id
}

func (e *Entity) GetComponentContainer() ComponentContainer {
	return e.componentContainer
}

func (e *Entity) AddComponents(components ...IComponent) {
	if _, ok := e.w.GetEntities()[EntityId(e.ID())]; !ok {
		e.w.GetEntities()[EntityId(e.ID())] = e
	}

	for _, component := range components {
		componentId := ComponentId(CompIdGetter.GetID(reflect.TypeOf(component)))
		componentInfo, ok := e.w.GetComponentMap()[componentId]
		if !ok {
			return
		}

		if target, exists := e.componentContainer[componentId]; exists {
			componentInfo.DestroyComponent(target)
			delete(e.componentContainer, componentId)
			componentInfo.RemoveEntity(e)
		}

		e.componentContainer[componentId] = component
		componentInfo.AddEntity(e)
	}
}

func (e *Entity) RemoveComponents(components ...IComponent) {
	for _, component := range components {
		componentId := ComponentId(CompIdGetter.GetID(reflect.TypeOf(component)))
		componentInfo, ok := e.w.GetComponentMap()[componentId]
		if !ok {
			return
		}

		if target, exists := e.componentContainer[componentId]; exists {
			componentInfo.DestroyComponent(target)
			delete(e.componentContainer, componentId)
			componentInfo.RemoveEntity(e)
		}
	}
}

func GetComponent[T IComponent](e IEntity) T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	componentId := CompIdGetter.GetID(t)
	component, ok := e.GetComponentContainer()[ComponentId(componentId)]
	if !ok {
		var zero T
		return zero
	}
	return component.(T)
}
