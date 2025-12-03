package ecs

import (
	"reflect"
)

type IdentityGetter struct {
	idIncr uint64
	idMap  map[reflect.Type]uint64
}

func NewIdentityGetter() *IdentityGetter {
	return &IdentityGetter{
		idMap: make(map[reflect.Type]uint64),
	}
}

func (ig *IdentityGetter) GetID(t reflect.Type) uint64 {
	if id, ok := ig.idMap[t]; ok {
		return id
	}

	ig.idIncr += 1
	ig.idMap[t] = ig.idIncr
	return ig.idIncr
}
