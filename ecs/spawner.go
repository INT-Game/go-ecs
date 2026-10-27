package ecs

import "reflect"

func SpawnEmptyEntity(w IWorld, components ...IComponent) IEntity {
	entity := NewEntity(w)
	w.GetCommands().doSpawn(entity, components...)
	return entity
}

func SpawnEntity[T IEntity](w IWorld, components ...IComponent) T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Ptr {
		panic("ecs: entity type must be a pointer")
	}
	e := reflect.New(t.Elem()).Interface().(T)

	setValue := reflect.ValueOf(e)
	setValue = setValue.Elem()
	baseEntityType := reflect.TypeOf((*Entity)(nil))

	found := false
	for i := 0; i < setValue.NumField(); i++ {
		field := setValue.Field(i)
		fieldType := setValue.Type().Field(i)

		if fieldType.Anonymous && fieldType.Type == baseEntityType {
			entity := NewEntity(w)
			field.Set(reflect.ValueOf(entity))
			found = true
			break
		}
	}
	if !found {
		panic("ecs: custom entity must embed *ecs.Entity")
	}

	w.GetCommands().doSpawn(e, components...)
	return e
}

func SpawnComponent[T IComponent](w IWorld) T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	componentId := ComponentId(w.GetCompId(t))
	if _, ok := w.GetComponentMap()[componentId]; !ok {
		w.GetComponentMap()[componentId] = NewComponentInfo[T](w)
	}
	componentInfo := w.GetComponentMap()[componentId]
	return componentInfo.CreateComponent().(T)
}
