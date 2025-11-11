package ecs

import "reflect"

type Resources struct {
	w *World
}

func NewResources(w *World) *Resources {
	return &Resources{
		w: w,
	}
}

func (r *Resources) Has(resource IComponent) bool {
	resourceId := ResIdGetter.GetID(reflect.TypeOf(resource))
	if _, ok := r.w.resourceMap[ComponentId(resourceId)]; ok && r.w.resourceMap[ComponentId(resourceId)].resource != nil {
		return true
	}
	return false
}

// Get 获取指定类型的资源
func (r *Resources) Get(resource IComponent) (IComponent, bool) {
	resourceId := ResIdGetter.GetID(reflect.TypeOf(resource))
	if resourceInfo, ok := r.w.resourceMap[ComponentId(resourceId)]; ok && resourceInfo.resource != nil {
		return resourceInfo.resource, true
	}
	return nil, false
}

// GetResource 添加泛型获取方法，使用更方便
func GetResource[T IComponent](r *Resources) (T, bool) {
	var zero T
	t := reflect.TypeOf(zero)
	if t == nil {
		t = reflect.TypeOf((*T)(nil)).Elem()
	}

	resourceId := ResIdGetter.GetID(t)
	if resourceInfo, ok := r.w.resourceMap[ComponentId(resourceId)]; ok && resourceInfo.resource != nil {
		if res, ok := resourceInfo.resource.(T); ok {
			return res, true
		}
	}
	return zero, false
}
