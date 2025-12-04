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

	component := components[0]
	remains := components[1:]

	componentId := ComponentId(q.w.GetCompId(reflect.TypeOf(component)))
	if componentInfo, ok := q.w.componentMap[componentId]; ok {
		density := componentInfo.Density()
		for i := 0; i < len(density); i++ {
			var entity IEntity
			entityId := EntityId(density[i])
			if entity, ok = q.w.entities[entityId]; ok {
				if q.doQueryRemains(entity, remains...) {
					entities = append(entities, entity)
				}
			}
		}
	}

	return entities
}

func (q *Query) doQueryRemains(e IEntity, components ...IComponent) bool {
	if len(components) == 0 {
		return true
	}

	component := components[0]
	remains := components[1:]

	componentId := q.w.GetCompId(reflect.TypeOf(component))
	if _, ok := e.GetComponentContainer()[ComponentId(componentId)]; ok {
		return q.doQueryRemains(e, remains...)
	}

	return false
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
	if !q.Has(e, c) {
		return nil, false
	}

	componentId := q.w.GetCompId(reflect.TypeOf(c))
	if component, ok := e.GetComponentContainer()[ComponentId(componentId)]; ok {
		return component, true
	}

	return nil, false
}
