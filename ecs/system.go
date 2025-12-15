package ecs

type ISystem interface {
	GetWorld() *World
	StartUp()
	Update()
	UpdateEntity(entity IEntity)
	RangeEntities(fn func(entity IEntity))
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

func (s *System) Update() {
	s.RangeEntities(func(entity IEntity) {
		s.UpdateEntity(entity)
	})
}

func (s *System) UpdateEntity(_ IEntity) {

}

func (s *System) RangeEntities(fn func(entity IEntity)) {
	for _, entity := range s.Query.Query(s.queryList...) {
		fn(entity)
	}
}
