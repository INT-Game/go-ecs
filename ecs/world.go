package ecs

type EntityId uint64
type ComponentId uint64

type IWorld interface {
	GetCommands() *Commands
	GetQuery() *Query
	GetComponentMap() map[ComponentId]IComponentInfo
	GetEntities() map[EntityId]IEntity
}

type World struct {
	IWorld
	commands        *Commands
	query           *Query
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

	w.commands = NewCommands(w)
	w.query = NewQuery(w)

	return w
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
