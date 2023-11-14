package entity

type Archetype interface {
	Apply(entity View) error
}
