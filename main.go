package main

import (
	"fmt"
	ecs2 "github.com/INT-Game/go-ecs/ecs"
)

type NameComponent struct {
	ecs2.Component
	Name string
}

type IDComponent struct {
	ecs2.Component
	Id string
}

type Timer struct {
	ecs2.Component
}

type NameSystem struct {
	ecs2.System
}

func NewNameSystem(commands *ecs2.Commands, query *ecs2.Query) *NameSystem {
	return &NameSystem{
		System: *ecs2.NewSystem(commands, query),
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
	ecs2.System
}

func NewIdSystem(commands *ecs2.Commands, query *ecs2.Query) *IdSystem {
	return &IdSystem{
		System: *ecs2.NewSystem(commands, query),
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
	w := ecs2.NewWorld()
	commands := ecs2.NewCommands(w)
	query := ecs2.NewQuery(w)

	w.AddUpdateSystem(NewNameSystem(commands, query))
	w.AddUpdateSystem(NewIdSystem(commands, query))

	nameComponent := ecs2.CreateComponent[*NameComponent](w)
	nameComponent.Name = "TestNameComponent"

	idComponent := ecs2.CreateComponent[*IDComponent](w)
	idComponent.Id = "TestIDComponent"

	commands.Spawn(nameComponent, idComponent)
	commands.Spawn(nameComponent)
	entityA := commands.SpawnAndReturnEntity(nameComponent)
	w.Update()

	fmt.Println("================================")

	commands.DestroyEntity(entityA)
	w.Update()
}
