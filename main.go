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

func NewNameSystem(w *ecs.World) *NameSystem {
	return &NameSystem{
		System: *ecs.NewSystem(w),
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

func NewIdSystem(w *ecs.World) *IdSystem {
	return &IdSystem{
		System: *ecs.NewSystem(w),
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

	w.AddUpdateSystem(NewNameSystem(w))
	w.AddUpdateSystem(NewIdSystem(w))

	nameComponent := ecs.SpawnComponent[*NameComponent](w)
	nameComponent.Name = "TestNameComponent"

	idComponent := ecs.SpawnComponent[*IDComponent](w)
	idComponent.Id = "TestIDComponent"

	entityA := ecs.SpawnEmptyEntity(w, nameComponent)
	w.Update()

	fmt.Println("================================")

	w.Commands.DestroyEntity(entityA)
	w.Update()
}
