package ecs

type Identifier interface {
	ID() uint64
}

type IdentifierSetter interface {
	SetID(id uint64)
}

type IdentifierSlice []Identifier

func (s IdentifierSlice) Len() int {
	return len(s)
}

func (s IdentifierSlice) Less(i, j int) bool {
	return s[i].ID() < s[j].ID()
}

func (s IdentifierSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
