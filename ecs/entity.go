package ecs2

import "sync/atomic"

var (
	eidIncr uint64
)

type IEntity interface {
	Identifier
	GetComponentContainer() ComponentContainer
}

type Entity struct {
	IEntity
	id                 uint64
	componentContainer ComponentContainer
}

func NewEntity() *Entity {
	return &Entity{
		id:                 atomic.AddUint64(&eidIncr, 1),
		componentContainer: make(ComponentContainer),
	}
}

func NewEntities(amount int) []*Entity {
	entities := make([]*Entity, amount)

	lastID := atomic.AddUint64(&eidIncr, uint64(amount))
	for i := 0; i < amount; i++ {
		// 完全初始化每个实体
		entities[i] = &Entity{
			id:                 lastID - uint64(amount) + uint64(i) + 1,
			componentContainer: make(ComponentContainer),
		}
	}

	return entities
}

func (e *Entity) ID() uint64 {
	return e.id
}

func (e *Entity) GetComponentContainer() ComponentContainer {
	return e.componentContainer
}
