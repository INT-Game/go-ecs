package main

import (
	"fmt"

	"github.com/INT-Game/go-ecs/ecs"
)

type Position struct {
	ecs.Component
	X, Y float64
}

func (p *Position) Init() {
	*p = Position{Component: p.Component}
}

type Velocity struct {
	ecs.Component
	X, Y float64
}

func (v *Velocity) Init() {
	*v = Velocity{Component: v.Component}
}

type GameTime struct {
	ecs.Component
	DeltaSeconds float64
}

type MovedEvent struct {
	EntityID uint64
	X, Y     float64
}

type MovementSystem struct {
	*ecs.System
	resources *ecs.Resources
	moved     *ecs.EventWriter[MovedEvent]
}

func NewMovementSystem(w *ecs.World, resources *ecs.Resources, moved *ecs.EventWriter[MovedEvent]) *MovementSystem {
	return &MovementSystem{
		System:    ecs.NewSystem(w, &Position{}, &Velocity{}),
		resources: resources,
		moved:     moved,
	}
}

func (s *MovementSystem) Update() {
	gameTime, ok := ecs.GetResource[*GameTime](s.resources)
	if !ok {
		return
	}
	s.RangeEntities(func(entity ecs.IEntity) {
		position := ecs.GetComponent[*Position](entity)
		velocity := ecs.GetComponent[*Velocity](entity)
		position.X += velocity.X * gameTime.DeltaSeconds
		position.Y += velocity.Y * gameTime.DeltaSeconds
		s.moved.Send(MovedEvent{EntityID: entity.ID(), X: position.X, Y: position.Y})
	})
}

type ReportSystem struct {
	*ecs.System
	moved *ecs.EventReader[MovedEvent]
}

func NewReportSystem(w *ecs.World, moved *ecs.EventReader[MovedEvent]) *ReportSystem {
	return &ReportSystem{System: ecs.NewSystem(w), moved: moved}
}

func (s *ReportSystem) Update() {
	for _, event := range s.moved.Get() {
		fmt.Printf("entity=%d position=(%.1f, %.1f)\n", event.EntityID, event.X, event.Y)
	}
}

func main() {
	w := ecs.NewWorld()
	defer w.Shutdown()

	resources := ecs.NewResources(w)
	w.GetCommands().SetResource(&GameTime{DeltaSeconds: 1})

	movedEvents := ecs.RegisterEvents[MovedEvent](w)
	w.AddUpdateSystem(NewMovementSystem(w, resources, movedEvents.Writer()))
	w.AddUpdateSystem(NewReportSystem(w, movedEvents.Reader()))
	w.Startup()

	position := ecs.SpawnComponent[*Position](w)
	position.X, position.Y = 0, 0
	velocity := ecs.SpawnComponent[*Velocity](w)
	velocity.X, velocity.Y = 1, 0.5
	entity := ecs.SpawnEmptyEntity(w, position, velocity)

	for tick := 1; tick <= 3; tick++ {
		fmt.Printf("tick %d\n", tick)
		w.Update()
		w.GetCommands().Execute()
	}

	entity.RemoveComponents(&Velocity{})
	fmt.Printf("moving entities after removing Velocity: %d\n", len(w.GetQuery().Query(&Position{}, &Velocity{})))

	w.GetCommands().DestroyEntity(entity).Execute()
	fmt.Printf("entities after destroy: %d\n", len(w.GetEntities()))
}
