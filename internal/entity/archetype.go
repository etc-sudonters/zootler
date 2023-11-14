package entity

// applies one or more components to an entity, useful for grouping components
// commonly used together
type Archetype interface {
	Apply(entity View) error
}
