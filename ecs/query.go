package ecs

import "reflect"

type Query struct {
	w *World
}

func NewQuery(w *World) *Query {
	return &Query{
		w: w,
	}
}

func (q *Query) Query(components ...IComponent) []IEntity {
	entities := make([]IEntity, 0)

	if len(components) == 0 {
		return entities
	}

	componentIDs := make([]ComponentId, len(components))
	baseIndex := -1
	var density []uint64
	for i, component := range components {
		componentID := ComponentId(q.w.GetCompId(reflect.TypeOf(component)))
		componentInfo, ok := q.w.componentMap[componentID]
		if !ok {
			return entities
		}
		componentIDs[i] = componentID
		candidate := componentInfo.Density()
		if baseIndex == -1 || len(candidate) < len(density) {
			baseIndex = i
			density = candidate
		}
	}

	entities = make([]IEntity, 0, len(density))
	for _, id := range density {
		entity, ok := q.w.entities[EntityId(id)]
		if !ok {
			continue
		}
		container := entity.GetComponentContainer()
		matched := true
		for i, componentID := range componentIDs {
			if i == baseIndex {
				continue
			}
			if _, ok := container[componentID]; !ok {
				matched = false
				break
			}
		}
		if matched {
			entities = append(entities, entity)
		}
	}

	return entities
}

// Has 判断实体是否包含指定组件
func (q *Query) Has(e IEntity, c IComponent) bool {
	componentId := q.w.GetCompId(reflect.TypeOf(c))
	if _, ok := e.GetComponentContainer()[ComponentId(componentId)]; ok {
		return true
	}

	return false
}

// Contains 判断实体是否包含所有指定组件
func (q *Query) Contains(e IEntity, components ...IComponent) bool {
	for _, c := range components {
		if !q.Has(e, c) {
			return false
		}
	}
	return true
}

// Get 获取实体的指定组件
func (q *Query) Get(e IEntity, c IComponent) (IComponent, bool) {
	componentId := q.w.GetCompId(reflect.TypeOf(c))
	component, ok := e.GetComponentContainer()[ComponentId(componentId)]
	return component, ok
}
