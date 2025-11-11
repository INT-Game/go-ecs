package ecs

type ISystem interface {
	StartUp()
	Update()
}

type System struct {
	ISystem
	Commands *Commands
	Query    *Query
}

func NewSystem(w *World) *System {
	return &System{
		Commands: w.commands,
		Query:    w.query,
	}
}

func (s *System) StartUp() {

}

func (s *System) Update() {

}
