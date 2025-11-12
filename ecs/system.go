package ecs

type ISystem interface {
	StartUp()
	Update()
}

type System struct {
	ISystem
	World    *World
	Commands *Commands
	Query    *Query
}

func NewSystem(w *World) *System {
	return &System{
		World:    w,
		Commands: w.commands,
		Query:    w.query,
	}
}

func (s *System) StartUp() {

}

func (s *System) Update() {

}
