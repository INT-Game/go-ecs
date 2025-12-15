package ecs

type ISystem interface {
	GetWorld() *World
	StartUp()
	Update(entity IEntity)
	rangeEntities(fn func(entity IEntity))
}

type System struct {
	ISystem
	World     *World
	Commands  *Commands
	Query     *Query
	queryList []IComponent
}

func NewSystem(w *World, queryList ...IComponent) *System {
	return &System{
		World:     w,
		Commands:  w.commands,
		Query:     w.query,
		queryList: queryList,
	}
}

func (s *System) GetWorld() *World {
	return s.World
}

func (s *System) StartUp() {

}

func (s *System) Update(_ IEntity) {

}

func (s *System) rangeEntities(fn func(entity IEntity)) {
	for _, entity := range s.Query.Query(s.queryList...) {
		fn(entity)
	}
}
