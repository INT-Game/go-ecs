package main

import (
	"fmt"
	"github.com/INT-Game/go-ecs/ecs"
)

type NameComponent struct {
	ecs.Component
	Name string
}

type IDComponent struct {
	ecs.Component
	Id string
}

type Timer struct {
	ecs.Component
}

type NameSystem struct {
	ecs.System
}

func NewNameSystem(commands *ecs.Commands, query *ecs.Query) *NameSystem {
	return &NameSystem{
		System: *ecs.NewSystem(commands, query),
	}
}

func (s *NameSystem) Update() {
	entities := s.Query.Query(&NameComponent{})
	for _, entity := range entities {
		comp, ok := s.Query.Get(entity, &NameComponent{})
		if ok {
			fmt.Println(comp.(*NameComponent).Name)
		}
	}
}

type IdSystem struct {
	ecs.System
}

func NewIdSystem(commands *ecs.Commands, query *ecs.Query) *IdSystem {
	return &IdSystem{
		System: *ecs.NewSystem(commands, query),
	}
}

func (s *IdSystem) Update() {
	entities := s.Query.Query(&IDComponent{})
	for _, entity := range entities {
		comp, ok := s.Query.Get(entity, &IDComponent{})
		if ok {
			fmt.Println(comp.(*IDComponent).Id)
		}
	}
}

func main() {
	w := ecs.NewWorld()
	commands := ecs.NewCommands(w)
	query := ecs.NewQuery(w)

	w.AddUpdateSystem(NewNameSystem(commands, query))
	w.AddUpdateSystem(NewIdSystem(commands, query))

	nameComponent := ecs.CreateComponent[*NameComponent](w)
	nameComponent.Name = "TestNameComponent"

	idComponent := ecs.CreateComponent[*IDComponent](w)
	idComponent.Id = "TestIDComponent"

	commands.Spawn(nameComponent, idComponent)
	commands.Spawn(nameComponent)
	entityA := commands.SpawnAndReturnEntity(nameComponent)
	w.Update()

	fmt.Println("================================")

	commands.DestroyEntity(entityA)
	w.Update()
}
