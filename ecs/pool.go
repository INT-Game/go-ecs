package ecs2

import (
	"mmo-server/array"
	"reflect"
)

type Pool[T IComponent] struct {
	instances array.Array[IComponent]
	caches    array.Array[IComponent]
}

func NewPool[T IComponent]() *Pool[T] {
	return &Pool[T]{
		instances: array.New[IComponent](),
		caches:    array.New[IComponent](),
	}
}

func (p *Pool[T]) Create() T {
	if !p.caches.Empty() {
		p.instances.PushBack(p.caches.Back())
		p.caches.PopBack()
	} else {
		componentId := CompIdGetter.GetID(reflect.TypeOf((*T)(nil)).Elem())
		component := p.doCreate()
		component.SetID(componentId)
		p.instances.PushBack(component)
	}

	return p.instances.Back().(T)
}

func (p *Pool[T]) doCreate() T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() == reflect.Ptr {
		v := reflect.New(t.Elem())
		return v.Interface().(T)
	}
	return reflect.New(t).Interface().(T)
}

func (p *Pool[T]) Destroy(elem IComponent) {
	if i, ok := p.instances.Contain(elem); ok {
		p.caches.PushBack(elem)
		p.instances.Swap(i, p.instances.Len()-1)
		p.instances.PopBack()

		// 调用销毁函数进行资源清理
		if elem.Destroy != nil {
			elem.Destroy()
		}
	}
}
