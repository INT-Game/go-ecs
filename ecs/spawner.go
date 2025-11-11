package ecs

import "reflect"

func SpawnEmptyEntity(w *World, components ...IComponent) IEntity {
	entity := NewEntity(w)
	w.Commands.doSpawn(entity, components...)
	return entity
}

func SpawnEntity[T IEntity](w *World, components ...IComponent) T {
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

	w.Commands.doSpawn(e, components...)
	return e
}

func SpawnComponent[T IComponent](w *World) T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	componentId := ComponentId(CompIdGetter.GetID(t))
	if _, ok := w.componentMap[componentId]; !ok {
		w.componentMap[componentId] = NewComponentInfo[T]()
	}
	componentInfo := w.componentMap[componentId]
	return componentInfo.CreateComponent().(T)
}
