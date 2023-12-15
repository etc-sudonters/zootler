package entity

// responsible for the total administration of a population of models
type Pool interface {
	Queryable
	Factory
}

// responsible for creation of models
type Factory interface {
	Create() (View, error)
}

// responsible for looking either individual models or creating a subset of the
// population that matches the provided selectors
type Queryable interface {
	// return a subset of the population that matches the provided filter
	Query(f any) ([]View, error)
}
