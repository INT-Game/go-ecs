package ecs

import "reflect"

func SpawnEmptyEntity(w IWorld, components ...IComponent) IEntity {
	entity := NewEntity(w)
	w.GetCommands().doSpawn(entity, components...)
	return entity
}

func SpawnEntity[T IEntity](w IWorld, components ...IComponent) T {
	var e T
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() == reflect.Ptr {
		v := reflect.New(t.Elem())
		e = v.Interface().(T)
	} else {
		e = reflect.New(t).Interface().(T)
	}

	setValue := reflect.ValueOf(e)
	if setValue.Kind() == reflect.Ptr {
		setValue = setValue.Elem()
	}

	for i := 0; i < setValue.NumField(); i++ {
		field := setValue.Field(i)
		fieldType := setValue.Type().Field(i)

		// Entity
		if fieldType.Name == "Entity" && fieldType.Type.String() == "*ecs.Entity" {
			entity := NewEntity(w)
			field.Set(reflect.ValueOf(entity))
			break
		}
	}

	w.GetCommands().doSpawn(e, components...)
	return e
}

func SpawnComponent[T IComponent](w IWorld) T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	componentId := ComponentId(CompIdGetter.GetID(t))
	if _, ok := w.GetComponentMap()[componentId]; !ok {
		w.GetComponentMap()[componentId] = NewComponentInfo[T]()
	}
	componentInfo := w.GetComponentMap()[componentId]
	return componentInfo.CreateComponent().(T)
}
