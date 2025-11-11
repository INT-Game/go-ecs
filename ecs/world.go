package ecs

type EntityId uint64
type ComponentId uint64

type IWorld interface {
}

type World struct {
	IWorld
	Commands        *Commands
	Query           *Query
	resourceMap     map[ComponentId]*ResourceInfo
	componentMap    map[ComponentId]IComponentInfo
	entities        map[EntityId]IEntity
	destroyEntities []IEntity
	startUpSystems  []ISystem
	updateSystems   []ISystem
}

func NewWorld() *World {
	w := &World{
		resourceMap:    make(map[ComponentId]*ResourceInfo),
		componentMap:   make(map[ComponentId]IComponentInfo),
		entities:       make(map[EntityId]IEntity),
		startUpSystems: make([]ISystem, 0),
		updateSystems:  make([]ISystem, 0),
	}

	w.Commands = NewCommands(w)
	w.Query = NewQuery(w)

	return w
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
	for componentId, component := range entity.GetComponentContainer() {
		componentInfo := w.componentMap[componentId]
		componentInfo.DestroyComponent(component)
		componentInfo.RemoveEntity(entity)
	}
	delete(w.entities, EntityId(entity.ID()))
}

func (w *World) Startup() {
	for _, system := range w.startUpSystems {
		system.StartUp()
	}
}

func (w *World) Update() {
	for _, system := range w.updateSystems {
		system.Update()
	}
}

func (w *World) Shutdown() {
	w.resourceMap = make(map[ComponentId]*ResourceInfo)
	w.componentMap = make(map[ComponentId]IComponentInfo)
	w.entities = make(map[EntityId]IEntity)
	w.startUpSystems = make([]ISystem, 0)
	w.updateSystems = make([]ISystem, 0)
}
