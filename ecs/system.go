package ecs2

type ISystem interface {
	StartUp()
	Update()
}

type System struct {
	ISystem
	Commands *Commands
	Query    *Query
}

func NewSystem(commands *Commands, query *Query) *System {
	return &System{
		Commands: commands,
		Query:    query,
	}
}

func (s *System) StartUp() {

}

func (s *System) Update() {

}
