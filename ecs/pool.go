package ecs

import (
	"reflect"

	"github.com/INT-Game/go-ecs/array"
)

type Pool[T IComponent] struct {
	w         IWorld
	instances array.Array[IComponent]
	caches    array.Array[IComponent]
}

func NewPool[T IComponent](w IWorld) *Pool[T] {
	return &Pool[T]{
		w:         w,
		instances: array.New[IComponent](),
		caches:    array.New[IComponent](),
	}
}

func (p *Pool[T]) Create() T {
	if !p.caches.Empty() {
		component := p.caches.Back()
		p.caches.PopBack()
		component.Init()
		p.instances.PushBack(component)
	} else {
		componentId := p.w.GetCompId(reflect.TypeOf((*T)(nil)).Elem())
		component := p.doCreate()
		component.SetID(componentId)
		component.Init()
		p.instances.PushBack(component)
	}

	return p.instances.Back().(T)
}

func (p *Pool[T]) doCreate() T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Ptr {
		panic("ecs: component type must be a pointer")
	}
	return reflect.New(t.Elem()).Interface().(T)
}

func (p *Pool[T]) Destroy(elem IComponent) {
	if i, ok := p.instances.Contain(elem); ok {
		p.caches.PushBack(elem)
		p.instances.Swap(i, p.instances.Len()-1)
		p.instances.PopBack()

		elem.Destroy()
	}
}
