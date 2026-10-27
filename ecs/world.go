package ecs

import (
	"reflect"
	"sync/atomic"
)

type EntityId uint64
type ComponentId uint64

type IWorld interface {
	IncrEntId() uint64
	GetResId(t reflect.Type) uint64
	GetCompId(t reflect.Type) uint64
	GetCommands() *Commands
	GetQuery() *Query
	GetComponentMap() map[ComponentId]IComponentInfo
	GetEntities() map[EntityId]IEntity
}

type World struct {
	entId        uint64
	resIdGetter  *IdentityGetter
	compIdGetter *IdentityGetter

	commands        *Commands
	query           *Query
	resourceMap     map[ComponentId]*ResourceInfo
	componentMap    map[ComponentId]IComponentInfo
	entities        map[EntityId]IEntity
	destroyEntities []IEntity
	startUpSystems  []ISystem
	updateSystems   []ISystem
	events          []eventQueue
}

type eventQueue interface {
	Advance()
	Clear()
}

func NewWorld() *World {
	w := &World{
		resIdGetter:  NewIdentityGetter(),
		compIdGetter: NewIdentityGetter(),

		resourceMap:    make(map[ComponentId]*ResourceInfo),
		componentMap:   make(map[ComponentId]IComponentInfo),
		entities:       make(map[EntityId]IEntity),
		startUpSystems: make([]ISystem, 0),
		updateSystems:  make([]ISystem, 0),
	}

	w.commands = NewCommands(w)
	w.query = NewQuery(w)

	return w
}

func (w *World) IncrEntId() uint64 {
	return atomic.AddUint64(&w.entId, 1)
}

func (w *World) GetResId(t reflect.Type) uint64 {
	return w.resIdGetter.GetID(t)
}

func (w *World) GetCompId(t reflect.Type) uint64 {
	return w.compIdGetter.GetID(t)
}

func (w *World) GetCommands() *Commands {
	return w.commands
}

func (w *World) GetQuery() *Query {
	return w.query
}

func (w *World) GetComponentMap() map[ComponentId]IComponentInfo {
	return w.componentMap
}
func (w *World) GetEntities() map[EntityId]IEntity {
	return w.entities
}

func (w *World) AddStartUpSystem(startUpSystem ISystem) *World {
	w.startUpSystems = append(w.startUpSystems, startUpSystem)
	return w
}

func (w *World) AddUpdateSystem(updateSystem ISystem) *World {
	w.updateSystems = append(w.updateSystems, updateSystem)
	return w
}

func (w *World) destroy(entity IEntity) {
	if entity == nil || entity.GetEcsWorld() != w {
		return
	}
	stored, ok := w.entities[EntityId(entity.ID())]
	if !ok {
		return
	}
	container := stored.GetComponentContainer()
	for componentId, component := range container {
		componentInfo := w.componentMap[componentId]
		if componentInfo != nil {
			componentInfo.DestroyComponent(component)
			componentInfo.RemoveEntity(stored)
		}
		delete(container, componentId)
	}
	delete(w.entities, EntityId(stored.ID()))
}

func (w *World) Startup() {
	for _, system := range w.startUpSystems {
		system.StartUp()
	}
}

func (w *World) Update() {
	for _, events := range w.events {
		events.Advance()
	}
	for _, system := range w.updateSystems {
		system.Update()
	}
}

func (w *World) Shutdown() {
	for _, entity := range w.entities {
		w.destroy(entity)
	}
	for _, resource := range w.resourceMap {
		if resource.resource != nil {
			resource.resource.Destroy()
			resource.resource = nil
		}
	}
	for _, events := range w.events {
		events.Clear()
	}
	w.resourceMap = make(map[ComponentId]*ResourceInfo)
	w.componentMap = make(map[ComponentId]IComponentInfo)
	w.entities = make(map[EntityId]IEntity)
	w.destroyEntities = nil
	w.startUpSystems = make([]ISystem, 0)
	w.updateSystems = make([]ISystem, 0)
	w.events = nil
}
