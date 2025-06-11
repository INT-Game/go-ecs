package ecs

import (
	"mmo-server/sparse_set"
)

type ComponentContainer map[ComponentId]IComponent

type ResourceInfo struct {
	resource    IComponent
	createFunc  func()
	destroyFunc func()
}

func NewResourceInfo(createFunc func(), destroyFunc func()) *ResourceInfo {
	return &ResourceInfo{
		createFunc:  createFunc,
		destroyFunc: destroyFunc,
	}
}

type IComponentInfo interface {
	AddEntity(e IEntity)
	RemoveEntity(e IEntity)
	CreateComponent() IComponent
	DestroyComponent(elem IComponent)
	Density() []uint64
}

type ComponentInfo[T IComponent] struct {
	IComponentInfo
	pool      *Pool[T]
	sparseSet *sparse_set.SparseSet[uint64]
}

func NewComponentInfo[T IComponent]() *ComponentInfo[T] {
	return &ComponentInfo[T]{
		pool:      NewPool[T](),
		sparseSet: sparse_set.NewSparseSet[uint64](32),
	}
}

func (c *ComponentInfo[T]) AddEntity(e IEntity) {
	c.sparseSet.Add(e.ID())
}

func (c *ComponentInfo[T]) RemoveEntity(e IEntity) {
	c.sparseSet.Remove(e.ID())
}

func (c *ComponentInfo[T]) CreateComponent() IComponent {
	return c.pool.Create()
}

func (c *ComponentInfo[T]) DestroyComponent(elem IComponent) {
	c.pool.Destroy(elem)
}

func (c *ComponentInfo[T]) Density() []uint64 {
	return c.sparseSet.Density()
}

type IComponent interface {
	Identifier
	IdentifierSetter
	Destroy()
}

type IComparableComponent interface {
	comparable
	IComponent
}

type Component struct {
	IComponent
	id uint64
}

func (c *Component) SetID(id uint64) {
	c.id = id
}

func (c *Component) ID() uint64 {
	return c.id
}

func (c *Component) Destroy() {

}
