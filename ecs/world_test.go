package ecs

import "testing"

type lifecycleComponent struct {
	Component
	Value int
}

func TestDestroyEntityClearsPooledComponentReference(t *testing.T) {
	w := NewWorld()
	component := SpawnComponent[*lifecycleComponent](w)
	component.Value = 1
	entity := SpawnEmptyEntity(w, component)

	w.GetCommands().DestroyEntity(entity).Execute()
	if got := GetComponent[*lifecycleComponent](entity); got != nil {
		t.Fatalf("destroyed entity retained pooled component %p", got)
	}

	reused := SpawnComponent[*lifecycleComponent](w)
	if reused != component {
		t.Fatal("test requires the component pool to reuse the component")
	}
	reused.Value = 2
	SpawnEmptyEntity(w, reused)
	if got := GetComponent[*lifecycleComponent](entity); got != nil {
		t.Fatalf("destroyed entity observed reused component: %+v", got)
	}
}

func TestShutdownDestroysEntities(t *testing.T) {
	w := NewWorld()
	entity := SpawnEmptyEntity(w, SpawnComponent[*lifecycleComponent](w))
	w.Shutdown()
	if len(entity.GetComponentContainer()) != 0 {
		t.Fatal("Shutdown left components attached to an entity")
	}
}

type lifecycleResource struct {
	Component
	destroyed *int
}

func (r *lifecycleResource) Destroy() {
	(*r.destroyed)++
}

func TestResourcesDestroyOnReplaceRemoveAndShutdown(t *testing.T) {
	w := NewWorld()
	destroyed := 0
	first := &lifecycleResource{destroyed: &destroyed}
	second := &lifecycleResource{destroyed: &destroyed}

	w.GetCommands().SetResource(first).SetResource(first)
	if destroyed != 0 {
		t.Fatal("setting the same resource destroyed it")
	}
	w.GetCommands().SetResource(second)
	if destroyed != 1 {
		t.Fatalf("replacing a resource destroyed %d resources, want 1", destroyed)
	}
	w.GetCommands().RemoveResource(&lifecycleResource{})
	if destroyed != 2 {
		t.Fatalf("removing a resource destroyed %d resources, want 2", destroyed)
	}
	w.GetCommands().SetResource(&lifecycleResource{destroyed: &destroyed})
	w.Shutdown()
	if destroyed != 3 {
		t.Fatalf("shutdown destroyed %d resources, want 3", destroyed)
	}
}

func TestSpawnReplacesDuplicateComponentWithoutLeakingPoolEntry(t *testing.T) {
	w := NewWorld()
	first := SpawnComponent[*lifecycleComponent](w)
	second := SpawnComponent[*lifecycleComponent](w)
	entity := SpawnEmptyEntity(w, first, second, second)
	if got := GetComponent[*lifecycleComponent](entity); got != second {
		t.Fatalf("entity component = %p, want last component %p", got, second)
	}
	if reused := SpawnComponent[*lifecycleComponent](w); reused != first {
		t.Fatal("replaced component was not returned to the pool")
	}
}

type commonComponent struct{ Component }
type rareComponent struct{ Component }

type customEntity struct {
	*Entity
	Name string
}

func TestSpawnCustomEntityInitializesEmbeddedEntity(t *testing.T) {
	w := NewWorld()
	entity := SpawnEntity[*customEntity](w, SpawnComponent[*commonComponent](w))
	if entity.Entity == nil || entity.ID() == 0 || !w.GetQuery().Has(entity, &commonComponent{}) {
		t.Fatalf("custom entity was not initialized: %+v", entity)
	}
}

func TestQueryMatchesRegardlessOfComponentOrder(t *testing.T) {
	w := NewWorld()
	for range 3 {
		SpawnEmptyEntity(w, SpawnComponent[*commonComponent](w))
	}
	wanted := SpawnEmptyEntity(w, SpawnComponent[*commonComponent](w), SpawnComponent[*rareComponent](w))

	for _, components := range [][]IComponent{{&commonComponent{}, &rareComponent{}}, {&rareComponent{}, &commonComponent{}}} {
		got := w.GetQuery().Query(components...)
		if len(got) != 1 || got[0] != wanted {
			t.Fatalf("query %T first = %v, want entity %d", components[0], got, wanted.ID())
		}
	}
}
