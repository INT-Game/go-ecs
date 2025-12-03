package ecs

import "reflect"

type Commands struct {
	w *World
}

func NewCommands(w *World) *Commands {
	return &Commands{
		w: w,
	}
}

func (c *Commands) doSpawn(entity IEntity, components ...IComponent) {
	if len(components) == 0 {
		return
	}

	component := components[0]
	remains := components[1:]

	// 设置组件的ID
	component.SetID(CompIdGetter.GetID(reflect.TypeOf(component)))
	componentId := ComponentId(component.ID())

	if _, ok := c.w.componentMap[componentId]; !ok {
		return
	}

	// 根据组件和实体的映射关系
	componentInfo := c.w.componentMap[componentId]
	componentInfo.AddEntity(entity)

	// 建立实体和组件的映射关系
	if _, ok := c.w.entities[EntityId(entity.ID())]; !ok {
		c.w.entities[EntityId(entity.ID())] = entity
	}

	c.w.entities[EntityId(entity.ID())].GetComponentContainer()[componentId] = component
	c.doSpawn(entity, remains...)
}

func (c *Commands) DestroyEntity(entity IEntity) *Commands {
	c.w.destroyEntities = append(c.w.destroyEntities, entity)
	return c
}

func (c *Commands) Execute() {
	for _, entity := range c.w.destroyEntities {
		c.w.destroy(entity)
	}
	c.w.destroyEntities = c.w.destroyEntities[:0]
}

func (c *Commands) SetResource(component IComponent) *Commands {
	resourceId := ResIdGetter.GetID(reflect.TypeOf(component))
	if _, ok := c.w.resourceMap[ComponentId(resourceId)]; !ok {
		c.w.resourceMap[ComponentId(resourceId)] = NewResourceInfo(func() {}, func() {})
	}
	c.w.resourceMap[ComponentId(resourceId)].resource = component
	return c
}

func (c *Commands) RemoveResource(component IComponent) *Commands {
	resourceId := ResIdGetter.GetID(reflect.TypeOf(component))
	if resourceInfo, ok := c.w.resourceMap[ComponentId(resourceId)]; ok {
		resourceInfo.destroyFunc()
		resourceInfo.resource = nil
		delete(c.w.resourceMap, ComponentId(resourceId))
	}
	return c
}
